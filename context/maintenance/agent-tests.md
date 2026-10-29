# Agent tests: what agent behavior is measured

What this repository measures about coding-agent behavior, where, and what the
measurement does not establish.

## What is measured

The product is what an agent in a generated workspace reads, runs, and writes,
so a green native test suite does not show that the product works. The agent
tests run a real coding host against a fresh fixture workspace and grade the
result from what the agent leaves behind: files, Git refs, and the product's own
`check --json`, plus pattern checks on the final reply where a scenario asks
for an answer. A person still reads the items a pattern cannot judge.

Nine scenarios cover orientation, joining a workspace, intent writing in two
languages, planning, execution, completion, a direct change, and borrowed
knowledge. Each names its checks, for example:

- nothing was pushed;
- no stray commit or attribution was written;
- nothing outside the workspace changed;
- a planted canary in material meant to stay passive did not appear in the
  host's output.

A check that cannot be decided because the host exposes no tool output is
reported as inconclusive and never counted as a pass.

Three hosts are supported: Claude Code, Codex, and Cursor. A result is only
meaningful as a statement about one host, one model, and one date: host H, model
M, k of N runs passed, on date D. The 2.5.0 baseline, under
`agent-tests/baselines/2.5.0/`, confirmed the pinned models and host versions:
Claude Code with claude-haiku-5-5, Codex with gpt-6-luna, and Cursor with
grok-4.7-low, one run per scenario per host on 2026-10-29. The validation note
cites its result.

## Why it is not CI

A run starts a coding host with the contributor's own credentials and spends
money on their account. It cannot be a required check: it would charge a person
who did not choose to pay, and the runs skip permission prompts inside a
throwaway world, so they never measure those prompts. Contributors run the tests before
opening a pull request and attach the evidence; the reviewer decides.

## When a change needs it

A change is in scope when an agent in a generated workspace could see different
bytes or different CLI output after it. The deterministic surface comparison
decides that for free, without a host.

- An empty comparison needs nothing more.
- A non-empty one names the scenarios and hosts to run and the before-and-after
  artifact to show.
- The pull request then carries a design case before the test evidence: what
  changes, what is wrong with the current design, what is better, what it costs,
  and the alternatives.
- A larger change agrees its design case in an issue first.
- Warnings (a failed check, a lower result, more than 15% more input tokens)
  are weighed by the reviewer and do not block.

## The reversal, and how it differs from the retired harness

The first version of this product had an acceptance harness that was retired,
and for a while this repository said host behavior was reported as not measured. That
is reversed on purpose. A product whose behavior lives in instructions needs
evidence about those instructions, and saying nothing was measured left changes
to skills and instruction files to be judged by reading alone.

The tests differ from the retired harness in three ways:

- they grade mainly the end state the agent leaves, not its account of it;
- they run in an isolated world with an empty host configuration and an
  environment allowlist;
- contributors run them on demand, and they are not part of a release gate.

They live in their own module under `agent-tests/` and ship nowhere.

## Refresh policy

A baseline is a committed set of runs for one host version and model, recorded
in the same format the report reads. A contributor whose host version and
model match it compares against it; otherwise they run the base and their
branch alternately on the same machine. It is refreshed after each template
release and whenever a pinned model or host version changes. A claim about host behavior stays no stronger
than the runs behind it, so anything not exercised is reported as not exercised.

## What it does not prove

Permission prompts, which a headless run cannot show. Whether a host that was
not run behaves the same. Behavior over a long project, because each scenario is
one short task. Anything about a model other than the pinned one.

Owner:

- `context-circuit-source@agent-tests/README.md` — the contributor guide: how to
  run, what a run touches, what it costs.
- `context-circuit-source@agent-tests/scenarios/` — the twelve scenarios and
  the checks each names.
- `context-circuit-source@agent-tests/hosts.yaml` — the hosts, the model pins, and the credential variable names.
- `context-circuit-source@agent-tests/triggers.yaml` — which changed paths need
  which scenarios, hosts, and artifacts.
- `context-circuit-source@CONTRIBUTING.md` — the rule a pull request follows for
  changes agents can see.
