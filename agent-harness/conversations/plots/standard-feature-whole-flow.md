# Standard feature — the whole flow, one conversation

Proves the lay-user experience end-to-end at Standard tier: the two real gates
(approve the intent, authorize delivery), silent mechanical middle, candidate
acceptance, inferred completion, and the reconciliation follow-up — the flow the
delta files describe in pieces but never play as one live conversation.

## Spec
```yaml
id: standard-feature-whole-flow
title: Standard feature, whole flow, one conversation
status: proposed
runtime_version: ">=1.0.0"
mode: full-execution
driver: cc-test-case
surface: [cc-intent, cc-plan, cc-execute, cc-verify, cc-deliver, cc-complete]
preconditions:
  repositories:
    - id: notes
      dest: notes
      default_branch: main
      branches: [main]
      seed_files: [notes.md]
      connect: main
  state: fresh                     # no intent, no plan yet
persona: >
  A solo maker whose notes project is connected. Has never heard of Context
  Circuit. Speaks in plain goals; makes real decisions but never rubber-stamps.
human_turns:
  - "I want a command that exports my notes to a Markdown file."
  - "That's right — approve it."
  - "Looks right to me."
  - "Go ahead and open it onto main."
  - "Reconcile it."
reactions:
  approves: true
  invents_repository: never
  uses_internal_terms: never
demonstrates:
  acceptance_criteria: [AC-03, AC-04, AC-06, AC-07, AC-08, AC-14, AC-16, AC-22, AC-35]
  invariants: [INV-INTENT-01, INV-APPROVE-01, INV-EXEC-01, INV-VERIFY-01, INV-CANDIDATE-01, INV-DELIVER-01, INV-COMPLETE-01, INV-COMPLETE-02, INV-KNOWLEDGE-02]
hidden:
  - internal file names, paths, digests (cand-…, contract_digest)
  - branch/worktree names (cc/…, refs/cc-…), model names, tier vocabulary as jargon
decision_points:
  - id: intent-presented-with-open-question
    when: the human states a goal
    coordinator_must: >
      present the goal, non-goals, and scope as one readable thing; raise in plain
      language the one thing the ask doesn't settle (overwrite behavior) and say how
      it proposes to handle it, without inventing a silent answer; state the risk level
      plainly ("independent check") without exposing the tier label as jargon; ask for
      a single approval decision — NOT a plan approval.
  - id: gate1-approve
    when: the human approves the intent
    coordinator_must: >
      treat this as the one upstream decision, freeze it, read the real code to plan
      from a grounded basis, derive the plan automatically, and NOT ask for a second
      plan approval.
  - id: between-gates-silent
    when: execution and verification run
    coordinator_must: >
      report (not gate) that the work was built and independently checked against
      the exact change; never call it "done" or "shipped" yet.
  - id: accept-candidate
    when: the human tries it and says it looks right
    coordinator_must: >
      record acceptance of this exact result; still not delivered or complete.
  - id: gate2-deliver
    when: the human authorizes the pull request
    coordinator_must: >
      treat delivery as the second real decision, target the recorded base
      branch, and open it only on the explicit yes.
  - id: inferred-completion
    when: acceptance + delivery both exist
    coordinator_must: >
      infer completion (Standard) rather than asking for a separate
      completion approval.
  - id: reconcile-debt
    when: completion/delivery emits reconciliation debt
    coordinator_must: >
      surface, in plain language, that merged work has not been folded into what
      the project knows, and ask to reconcile now or mark no-update-needed —
      never accept knowledge silently.
reporting_rules: [plain-language, never-overstate-assurance, faithful-failure]
expected_end_state:
  # A FRESH whole-flow case authors its own intent/plan ids live, so it cannot use
  # the id-keyed post-conditions the seeded cases use (e.g. plan_status: 0001-…:done).
  # It needs id-agnostic assertions — a harness capability this *(opt)* plot requires.
  - repositories_registered: 1
  - intents_approved: 1
  - plans_created: 1
  - all_plans_done: true            # the single derived plan completes (inferred, Standard)
  - deliveries_recorded: 1
  - completion_kind: inferred
  - knowledge_debt_pending: ">=1"
  - product_knowledge_unchanged_silently: true
access_discipline:
  intent:
    required: []
    allowed: ["intent/**", "context/**", AGENTS.md, WORKFLOW.md]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**", "wrapper/contracts/**"]
  execute-deliver:
    required: []
    allowed: ["plans/**", "intent/**", ".runtime/**", "context/**"]
    forbidden: ["plans/archive/**", "wrapper/runtime/engine.sh", "sources/**"]
```

## Dialogue

> **H:** I want a command that exports my notes to a Markdown file.
>
> **C:** Here's what I understand you want: **a command that writes your notes out
> to a Markdown file you choose.** Out of scope: changing how notes are stored, and
> any network/upload behavior. It touches the notes project only.
>
> One thing your ask doesn't settle: what happens if the output file already exists.
> I don't want it silently overwriting your file, so I've written the goal to never
> overwrite — on a clash it reports and stops. Say if you'd rather it went another way.
>
> This is ordinary, reversible, well-covered work, so it gets an independent check.
> Approve this and I'll build it — or tell me what to change.
> `[decision_point: intent-presented-with-open-question]`
>
> **H:** That's right — approve it.
>
> **C:** I looked at the real code to lay it out — it's doable — building it now.
> `[decision_point: gate1-approve]`
>
> *(…work happens; the human is not asked to approve a plan…)*
>
> **C:** Built and independently checked — the criteria pass against the exact
> change, including the no-overwrite one. Want to try it and decide whether it
> ships? `[decision_point: between-gates-silent]`
>
> **H:** Looks right to me.
>
> **C:** Noted that you're happy with this result. Nothing has shipped yet — say
> the word and I'll open the pull request. `[decision_point: accept-candidate]`
>
> **H:** Go ahead and open it onto main.
>
> **C:** Opened the pull request onto `main` for the notes project.
> `[decision_point: gate2-deliver]`
>
> That's the change accepted and delivered, so it's finished from your side — I
> didn't need a separate "mark it done" step. `[decision_point: inferred-completion]`
>
> One thing before your next change here: this merged work hasn't been folded into
> what the project knows yet. Want me to reconcile that now, or note it as no
> update needed? `[decision_point: reconcile-debt]`
>
> **H:** Reconcile it.

Human decisions: **approve the intent**, **accept the result**, **authorize
delivery**, **resolve reconciliation** — four real decisions, none a rubber stamp,
and no separate plan-approval or mark-done step.
