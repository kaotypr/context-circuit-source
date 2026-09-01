# Spec adversary role

The spec adversary is an independent actor that attacks an intent's acceptance
criteria **before any code exists** (Context Circuit v1.0, Mechanism 1,
INV-INTENT-01). It applies the same "author ≠ checker" discipline the product
already trusts for the verifier, pointed at a different input: the criteria, not
the implementation.

It receives only the intent's `contract.yaml` — goal, non-goals, constraints,
acceptance criteria, and scope. It sees **no** implementation, because it runs
before the worker. It is read-only over the intent contract; it proposes, it
never edits the contract itself, and it never approves the intent (the human
decides at the gate).

Its task, recorded in the intent's `adversary.md`:

- **Find a way to satisfy every criterion and still be wrong.** Construct the
  most plausible implementation that passes all acceptance criteria yet violates
  the goal, a constraint, or a non-goal.
- **Name the missing paths.** List the edge / error / security / concurrency /
  data-loss cases the criteria do not cover.

Its output is a list of findings, each `{severity, statement, suggested
criterion}`, plus a verdict `criteria_sound: yes | needs-work`. Findings become
new or revised criteria, or explicit open questions, **before** the human
approves. Criteria that cannot survive the attack are rewritten or logged, so the
downstream verification rigor certifies the right thing rather than faithfully
certifying the wrong one.

Independence: the adversary must be a distinct child from whatever authored the
criteria; otherwise it is self-review of the spec. Which model or host runs it is
bounded `host_evidence` and authorizes nothing (INV-HOST-01).

Its depth scales with the intent's consequence tier (INV-ASSURE-01): light or a
single inline pass at Explore, a full battery at Critical. It is measured by how
often its findings actually change a criterion; a low change-rate means it is
fabricating implausible cases and is miscalibrated, and it must be tuned toward
real, reachable failure rather than theatre.

If the host cannot create the adversary child, the coordinator reports
`host-blocked` and surfaces the criteria to the human unchallenged rather than
faking the pass — an unchallenged criteria set is reported as such, never as
"checked."
