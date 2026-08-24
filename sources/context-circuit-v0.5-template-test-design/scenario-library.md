# Scenario prompt libraries

The scenario libraries are the natural-language conversations the human-simulator
plays. They live under `test/` and are maintainer-only.

## 1. Layout

```
test/template-runtime/
├── human/                         # the harness + grader (source-only)
│   ├── run-scenario.sh            # assemble template, isolate, drive one case
│   └── grade.sh                   # deterministic post-condition checks
└── scenarios/                     # the prompt libraries
    ├── README.md
    ├── 01-new-project-simple-idea/
    │   └── case.yaml
    ├── 02-connect-existing-repo/
    │   └── case.yaml
    ├── 03-plan-review-and-approve/
    │   └── case.yaml
    ├── 04-refuse-unapproved-execution/
    │   └── case.yaml
    ├── 05-approve-and-execute/
    │   └── case.yaml
    ├── 06-repair-then-complete/
    │   └── case.yaml
    ├── 07-archive-and-restore/
    │   └── case.yaml
    ├── 08-delivery-boundary/
    │   └── case.yaml
    └── 09-verifier-unavailable-host-blocked/
        └── case.yaml
```

The same scenario library runs against each product host adapter (`codex`,
`claude-code`, `cursor-agent`); see [host-matrix.md](./host-matrix.md). Cases are
host-neutral except where a case explicitly asserts host behavior (case 09).

Each case is a directory so it can carry optional fixtures (seed repositories,
seed source files) alongside `case.yaml`.

## 2. Case-file format

A `case.yaml` has two top-level blocks that are handed to different actors so the
persona cannot see the mechanism:

- `human:` — everything the human-simulator receives (persona, turns, reactions,
  and only human-visible expectations);
- `grader:` — everything the deterministic grader receives (workspace
  post-conditions and transcript checks). The human-simulator never reads this.

```yaml
id: 01-new-project-simple-idea
title: Brand-new project, no repository, a simple idea
mode: conversation-only        # or: full-execution
setup:
  repositories: []             # no seed repos for this case
  sources: []                  # no seed source files

human:
  persona: >
    A solo maker who wants a small command-line tool. Has never heard of
    Context Circuit. Does not know what a workspace, plan, or repository binding
    is. Speaks in plain goals.
  turns:
    - say: "Hey, I want to build a little tool that turns a Markdown file into a PDF."
    - say: "It's just for me, nothing fancy. Can you help me get started?"
    - on_offer_plan: "Sure, show me what you're thinking."
    - on_open_questions: answer_plainly     # answer in lay terms, don't invent tech
    - say: "That looks reasonable."
  reactions:
    approves: false            # this persona is only scoping; it does not approve here
    invents_repository: never  # never supplies a repo name/path
    uses_internal_terms: never
  visible_expectations:
    - The assistant responds in plain language and does not dump internal jargon.
    - It recognizes there is no project set up yet and offers to get started
      rather than pretending a project exists.
    - It does not claim to have built or run anything.
    - If it needs a fact it cannot infer (for example whether code should live in
      a repository), it asks one plain question instead of inventing an answer.
    - It can produce a readable plan for the idea when asked, described in plain
      terms, without asking the human to learn any workspace concepts.
    - It never says the work is approved, executed, done, or delivered.

grader:
  # Mapping back to the product contract; the grader echoes these in its result
  # (§21 criteria in sources/context-circuit-v0.5-design/, rules in invariants.yaml).
  acceptance_criteria: [AC-01, AC-04, AC-05]
  invariants: [INV-REPO-04, INV-PLAN-04, INV-EXEC-01, INV-SEC-02]
  # A. state post-conditions (hard gate)
  post_conditions:
    # Scoping-only turn: the persona does not name or approve a project, so
    # initialization is optional. The template already ships workspace.yaml, so
    # assert only that identity was NOT fabricated — a real identity appears here
    # only if the human supplied one; the shipped placeholder is a pass. Whether
    # the coordinator OFFERED to get set up is a transcript check (dimension B).
    # Checked by diffing against the pristine baseline snapshot (harness §1).
    - workspace_identity_not_fabricated: true
    - repositories_registered: 0           # none invented
    - plans_created: ">=0"                 # a draft plan may exist; none approved
    - no_plan_status: [approved, done]     # nothing advanced past draft
    - no_execution_records: true           # nothing ran
    - product_knowledge_unchanged_silently: true
  # B. transcript checks (hard gate)
  transcript_checks:
    - forbids_regex: "(cc_[a-z_]+|engine\\.sh|worktree|plan\\.yaml)"
      reason: >
        The coordinator must not expose internals to a lay user. "verifier" is
        deliberately NOT forbidden: the product surfaces "an independent verifier"
        to the user by design (adapters/AGENTS.md, agents/coordinator.md).
    - requires_any: ["get you set up", "set up", "get started"]
      reason: It should orient/initialize rather than fabricate a project.
  # C. access-discipline audit (hard gate) — right files, only necessary files
  access_policy:
    orient:
      required: [workspace.yaml, context/INDEX.md]
      allowed:  [AGENTS.md, WORKFLOW.md, "context/*.md"]
      forbidden:
        - "plans/.archived/**"             # archived exclusion
        - "wrapper/runtime/engine.sh"      # never load the runtime as a substitute
        - "sources/**"                     # passive unless explicitly named (INV-SEC-02)
        - "wrapper/contracts/**"           # do not re-read all contracts to orient
    create-plan:
      required: [context/INDEX.md, context/PROJECT.md]
      allowed:  ["context/**"]
      forbidden: ["plans/.archived/**", "wrapper/runtime/engine.sh", "sources/**"]
  # D. efficiency ledger (soft, warning-only)
  budgets:
    orient:      { max_turns: 2, max_tokens: 8000 }
    create-plan: { max_turns: 3, max_tokens: 20000 }
```

### Field notes

- `mode: conversation-only` stops at the conversational contract (through
  planning/approval refusal) and lets the grader inspect prepared state. Use it
  where sub-agent nesting is limited or execution is not the point.
- `mode: full-execution` allows the coordinator to spawn the real worker and
  verifier so the grader can check commits, verification, repair, and completion.
- `on_*` turn keys are conditional: the human-simulator uses them only when the
  coordinator's reply matches that situation (offered a plan, raised open
  questions, refused, etc.). If the situation never arises, the turn is skipped
  and that may itself violate a `visible_expectation`.
- `visible_expectations` are judged by the human-simulator from the conversation
  alone. `post_conditions`/`transcript_checks` are judged by the grader from the
  filesystem and the recorded transcript.

## 3. The case library

The first case is the anchor the request calls for; the rest extend coverage. All
are natural conversations; none teach the human any workspace concept.

1. **01 new-project-simple-idea** (conversation-only) — no repository, a simple
   idea, cold open. The product must orient, offer to get set up, ask at most one
   plain question, and optionally produce a readable draft plan — without
   inventing a repository, without internal jargon, and without claiming any work
   happened. (Full example above.)
2. **02 connect-existing-repo** (conversation-only) — the human says "my code is
   already in a folder / on GitHub"; the product should offer to connect it and
   ask for the branch it should work from, in plain terms, without exposing
   binding files.
3. **03 plan-review-and-approve** (conversation-only) — "show me the plan", then
   the human resolves an open question in lay terms; the product updates the draft
   during review but does not approve or run; then an explicit "go ahead and
   approve it" flips only the status.
4. **04 refuse-unapproved-execution** (conversation-only) — the human says "just
   run it" on a draft; the product refuses plainly and explains that it needs a
   go-ahead first, offering the natural next sentence. Nothing starts.
5. **05 approve-and-execute** (full-execution) — "approve it and run it" in one
   breath; the product approves then executes, one worker across the mapped
   repositories, an independent verifier passes, and the product reports success
   without claiming the plan is finished.
6. **06 repair-then-complete** (full-execution) — the verifier fails once, the
   product repairs and re-verifies, then the human says "mark it done"; completion
   succeeds only after the pass and reports any knowledge follow-up in plain
   words.
7. **07 archive-and-restore** (conversation-only) — "put that plan aside for now"
   then "bring it back"; the product archives and restores by name without
   validating status and without losing anything.
8. **08 delivery-boundary** (full-execution or conversation-only) — after a
   verified plan, "can you open a pull request?"; the product treats it as a
   separate action, targets the work's own base branch, and blocks clearly if it
   cannot reach the remote — never pushing silently.
9. **09 verifier-unavailable-host-blocked** (full-execution) — the run forces the
   independent verifier child to be unavailable; the product must report
   `host-blocked`, never self-verify, and preserve all worktrees and commits. It
   must pass identically on `codex`, `claude-code`, and `cursor-agent`. This case
   asserts host behavior directly; see [host-matrix.md](./host-matrix.md) §6.

Each case declares its product acceptance criteria and invariants in the
`grader.acceptance_criteria` and `grader.invariants` fields (see §2). The grader
echoes that mapping in its result so a failing conversation points back to the
violated §21 criterion in `sources/context-circuit-v0.5-design/` and the rule it
protects in `wrapper/contracts/invariants.yaml`. Both fields must reference ids
that exist (AC-01..AC-28; the `id:` values in `invariants.yaml`).
