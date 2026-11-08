# Customization

A workspace's rules are of two kinds. Fixed rules hold whatever an extension
says. Overridable defaults are a closed list of four, and an extension changes
one only where its manifest declares the override and a person allowed it in
this workspace. The CLI enforces that: an override reaches you only through
`extension hooks`, and a dispatch that needs one is accepted only with
`agent dispatch --override NAME` naming an allowed one.

## Fixed rules

These hold in every workspace, and no extension or override changes them:

- Approving an intent is a person's answer to it as written.
- Pushing, pull requests, merges, deployment, deletion, and commits in a bound
  checkout need the person's explicit authorization.
- Planning and execution start on a person's request.
- The plan record is left alone during execution.
- No attribution in commits, pull requests, reviews, or comments.
- No secrets in records.
- Borrowed knowledge stays read-only.

A manifest that names any of them as an override is invalid, and `extension
override` refuses one.

## Overridable defaults

| Override | Core behavior | With the override |
| --- | --- | --- |
| `repair-loop` | No automatic repair loops | Implement, verify, fix, until clean or the round limit |
| `verify-during-execution` | Verification only on a person's request | A named read-only role checks each implementation once, and reports |
| `risk-review` | No review triggered by risk classification | The extension's instructions say when a change warrants review, and the reviewer runs then |
| `delegation` | Delegating is a judgment about size | A named stage always goes to a named role |

`repair-loop` includes the verification it needs, so it doesn't also need
`verify-during-execution`.

A manifest declares the overrides it wants:

```yaml
overrides:
  - default: repair-loop
    instructions: instructions/repair-loop.md   # what counts as clean, what to fix first
    verifier: parity-verifier                   # a read-only role this extension declares, or reviewer
    rounds: 3                                   # 1 to 10
  - default: risk-review
    instructions: instructions/risk-review.md   # when a change warrants review
  - default: delegation
    stage: implement                            # plan, implement, or review
    role: worker
```

A verifier is read-only: the shipped `reviewer`, or a role the extension
declares with `access: read-only`.

## Allowing one

Only a person allows an override, and only on their request:

```sh
sh .context-circuit/cli extension override --id parity --allow repair-loop --rounds 2
sh .context-circuit/cli extension override --id parity --revoke repair-loop
```

The allowance is written under the entry's `overrides` in `extensions.yaml`,
keeping the file's comments. `--rounds` may lower the manifest's rounds and
never raise them. The command refuses a default the manifest doesn't declare.
An allowance is per workspace, and it applies to a stage only where the
extension is on for the member whose agent runs it and its repositories cover
the plan. `status` lists each allowed override with its extension, rounds, and
whether it is on for you. Without an allowance, nothing changes.

## Where an override arrives

`extension hooks` returns an `overrides` block at the point where each one
applies:

```json
"overrides": [
  {"default": "repair-loop", "extension": "parity", "rounds": 2,
   "verifier": "x-parity-parity-verifier",
   "instructions": "extensions/parity/instructions/repair-loop.md"}
]
```

| Override | Point that returns it | Dispatch it unlocks |
| --- | --- | --- |
| `repair-loop` | `after implement` | the verifier, then a worker with `--task` carrying the findings |
| `verify-during-execution` | `after implement` | the verifier, once |
| `risk-review` | `after implement` | the reviewer, when the instructions say the change warrants it |
| `delegation` | `before <stage>` | none: the stage goes to the role, dispatched as usual |

Read the `after implement` result before reporting. Each unlocked dispatch
names the override, and its round where there is one:

```sh
sh .context-circuit/cli --json agent dispatch --host HOST --role x-parity-parity-verifier \
  --plan p0001 --override repair-loop --round 1 --path WORKTREE
sh .context-circuit/cli --json agent dispatch --host HOST --role worker \
  --plan p0001 --task 'FINDINGS' --override repair-loop --round 1 --path WORKTREE
```

Any other use of `--override` is refused, naming what is missing. Without
`--override`, a reviewer still needs `--review-requested`.

## The repair loop

1. The worker implements and commits, and `after implement` returns
   `repair-loop`.
2. Dispatch the verifier with `--override repair-loop --round 1`.
3. Clean: report "clean after N rounds".
4. Findings, with rounds left: dispatch a worker with the findings as `--task`
   and the same round, then the verifier again with the next round.
5. Findings in the last round: stop, and give the findings to the person.

- You run every step. The verifier never dispatches the worker, and the worker
  never dispatches the verifier.
- A round is one verification. Round 1 checks the first implementation, so
  `rounds: 2` allows one fix.
- Each fix is committed on the worktree branch like any worker commit. The plan
  record is not edited: the rounds go in the report.
- A finding outside the plan's outcome is reported, not fixed. A finding that
  bears on the outcome goes to the person, as any question about the outcome
  does.
- You may make a fix yourself when you implemented the plan yourself. The
  verifier is still dispatched, because it must not be the session that wrote
  the code. Where the host can't spawn an agent, don't run the loop: report
  that the override needs a host that can.

## The report

Say which override ran, from which extension, how many rounds, and how it
ended: clean after N rounds, or the findings left when the limit was reached.
Name each fix's commit. `verify-during-execution` and `risk-review` report the
verifier's or reviewer's result with the implementation's.

## What `check` reports

Every finding here needs a person:

- an override in a manifest outside the closed list
- an allowance the manifest doesn't declare, or for a default outside the list
- an allowance on an extension whose source can't be read
- `rounds` above the manifest's value or outside 1 to 10
- a `verifier` that isn't read-only, or names a role nobody declares
- two extensions allowing `delegation` for the same stage with different roles
