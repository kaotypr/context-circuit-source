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
demonstrates:
  acceptance_criteria: [AC-03, AC-04, AC-06, AC-07, AC-08, AC-14, AC-16, AC-22, AC-35]
  invariants: [INV-INTENT-01, INV-APPROVE-01, INV-EXEC-01, INV-VERIFY-01, INV-CANDIDATE-01, INV-DELIVER-01, INV-COMPLETE-01, INV-COMPLETE-02, INV-KNOWLEDGE-02]
hidden:
  - internal file names, paths, digests (cand-…, contract_digest)
  - branch/worktree names (cc/…, refs/cc-…), model names, tier vocabulary as jargon
decision_points:
  - id: intent-presented-with-adversary
    when: the human states a goal
    coordinator_must: >
      present the goal, non-goals, and scope as one readable thing; fold in the
      spec adversary's finding in plain language; state the risk level plainly
      ("independent check") without exposing the tier label as jargon; ask for a
      single approval decision — NOT a plan approval.
  - id: gate1-approve
    when: the human approves the intent
    coordinator_must: >
      treat this as the one upstream decision, freeze it, derive the plan
      automatically, and NOT ask for a second plan approval.
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
```

## Dialogue

> **H:** I want a command that exports my notes to a Markdown file.
>
> **C:** Here's what I understand you want: **a command that writes your notes out
> to a Markdown file you choose.** Out of scope: changing how notes are stored, and
> any network/upload behavior. It touches the notes project only.
>
> I had this checked adversarially first. One thing came back: it didn't say what
> happens if the output file already exists — an implementation could pass every
> check and still silently overwrite your file. I've added that as a criterion:
> never overwrite; report and stop.
>
> This is ordinary, reversible, well-covered work, so it gets an independent check.
> Approve this and I'll build it — or tell me what to change.
> `[decision_point: intent-presented-with-adversary]`
>
> **H:** That's right — approve it.
>
> **C:** Building it now. `[decision_point: gate1-approve]`
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
