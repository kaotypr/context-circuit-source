# How the Walkthrough Maps to Context Circuit

This guide describes the expected behavior when the prompts are used in an
instantiated Context Circuit workspace. It does not run the lifecycle in this
maintainer source repository.

## 1. Initialize the workspace

The first prompt creates the minimum workspace and records only the supplied
product name and purpose. It establishes portable identity, empty knowledge and
lifecycle indexes, roster support, contracts, and private runtime support. It
also records that the project is greenfield. It does not create or connect
product repositories because repository creation remains a separate action.

## 2. Bind the repositories

The second prompt explicitly asks for two new repositories. The workspace should
create empty `invoice-tracker-api` and `invoice-tracker-web` Git repositories in
its normal repository location, establish `main`, create the initial base commit
required for later isolated execution, record their portable identities and
host-local bindings, and validate both checkouts.

No remote repository is created and nothing is pushed. Those are separate
delivery actions that require explicit destinations later.

## 3. Configure child-agent roles

No configuration prompt is required for a first-time user. The workspace uses
the host defaults for planner, worker, and verifier children. If the user asks
the optional plain-language question, the coordinator explains how planning,
building, and independent checking will be handled without requiring the user
to choose model IDs or reasoning levels.

An experienced user may configure host-local role settings later. Those settings
affect cost and capability only; they do not change authorization, consequence
tier, or verifier independence.

## 4. Establish the human identity

The fourth prompt selects or adds the human in the shared roster, allocates
non-overlapping intent and plan number bands when needed, and records the local
machine identity. All allocation details remain internal; the human only needs
to say who they are.

## 5. Build living Product Knowledge

No separate gathering action is useful before the system design exists. There
is no code, repository documentation, or previous product behavior to inspect.
The initial idea remains human-supplied evidence for authoring the design, while
unknown tax, numbering, currency, privacy, retention, hosting, and technology
choices remain explicit rather than being guessed.

Gathering becomes useful after the whole-system design has been written. The
later question, “How can we implement all of this designed system?”, names that
design as the bounded source. Before deriving implementation increments, the
coordinator retrieves any existing knowledge, synthesizes the durable decisions
from the design into the relevant `context/` units, and updates
`context/INDEX.md` in the same operation. Raw design prose is not copied
wholesale, and there is no separate knowledge-approval gate.

## Design the whole system

“Lets design the whole system for this idea” asks for a product-level system
design, not an implementation plan. The design belongs under
`sources/system-design/` using a product and grouping index, then concern-based
scope folders. Each scope starts with an orienting `README.md` and a normative
`design.md`; implementation-depth concern files are added only when needed.

The design should cover the complete product shape across client management,
invoice lifecycle and numbering, calculations and rounding, PDF documents,
manual payment records, overdue state, reminders, identity and privacy,
dashboard behavior, data ownership, operational concerns, and rollout. Scopes
are divided by product concern rather than by Git repository. Important flows
and relationships should use embedded Mermaid diagrams where they improve
understanding.

This design is passive source material. Creating it does not approve an outcome,
write Product Knowledge, create lifecycle plans, or authorize implementation.

## Explain how the whole design can be implemented

“How can we implement all of this designed system?” asks the coordinator to
explain the route from the whole-product design to controlled increments. The
named design is first gathered into living Product Knowledge, then decomposed
into separate outcome-level intents. Each approved intent can produce one or
more repository-grounded, single-repository plans with explicit dependencies.

The response should propose a product-oriented sequence, such as foundations,
client and business profiles, invoice lifecycle and calculations, PDF output,
payment and overdue tracking, dashboard experience, reminders, and operational
readiness. It is an implementation roadmap, not authority to approve or execute
all work at once. The human reviews the sequence and starts with one bounded
outcome.

## 6. Draft the product decision

The sixth numbered prompt selects the first bounded outcome from the proposed
implementation path and creates one draft intent containing its desired outcome,
non-goals, constraints, acceptance criteria, consequence tier, and open
questions. Its optional detail is organized by product concerns rather than
repositories. No target code is read while the intent is drafted.

Because invoice data is financially sensitive, uncertainty should not lower the
assurance tier. Important decisions such as numbering uniqueness, rounding,
time zones, taxes, edits after issue, payment allocation, and PDF immutability
must remain visible for the human to resolve before approval.

## 7. Approve and derive grounded plans

The seventh prompt is the first human gate. Approval freezes what correct means.
One planner child per repository then inspects the initialized repository and
its guidance. Because the repositories are new, the grounded findings should
say so plainly and derive feasible single-repository plans for scaffolding and
implementation with real intended paths, tasks, risks, dependencies, and
runnable checks.

The planners derive the actual order from the approved outcome, Product
Knowledge, repository state, and repository guidance. API foundations, data
rules, and contracts will likely precede dependent web flows. Any required
technology or architecture choice that changes the product decision returns to
the human instead of being silently settled in a plan. Approval itself does not
execute work.

## 8. Execute one plan and verify it

The eighth prompt starts the first ready Standard or Critical plan. The runtime
checks authorization and the clean base, then creates an isolated execution
branch and worktree. One worker implements, tests, and commits the plan before a
separate read-only verifier checks the exact candidate.

Failures return evidence to the same worker while attempts remain. Missing
independent verification produces a block. Every outcome preserves the work and
does not imply delivery or completion.

## 9. Run the remaining plan graph

The ninth prompt schedules all remaining authorized plans from current evidence.
Verified work is not repeated. Independent plans may overlap when dependencies,
path reservations, Git bases, and host capacity allow it. Each plan retains its
own worker, verifier, candidate, and outcome.

Cross-repository dependencies order API and web work without pretending they
share one Git history. A failure holds its descendants while unrelated work may
continue. The stack can resume after repair or resolution of a blocker.

## 10. Deliver explicit external effects

Because this began as a greenfield local project, the delivery prompts first
authorize creating private remote repositories at an exact GitHub account or
organization, then pushing the reviewed branches, and finally opening pull
requests. Each action is explicit and its result is reviewed before the next.

For pull requests, the coordinator groups verified plans by covering repository
tip, binds human acceptance to current candidates, checks base drift, and opens
the required pull requests against `main`. At least one pull request is required
per repository.

Merge and release remain separate prompts because each creates a distinct
external effect. Drift requires bringing the work onto the current base and
independently checking the new candidate. Delivery does not mark plans done or
update Product Knowledge.

## 11. Complete and reconcile durable truth

The eleventh prompt explicitly marks the release plans done and reconciles what
was actually implemented into living Product Knowledge. The update describes
current client, invoice, PDF, payment-status, dashboard, and ownership behavior,
while the retrieval index is updated in the same operation.

Commits, branches, runtime logs, and agent reports stay out of Product Knowledge.
Missing evidence is reported as knowledge debt rather than converted into an
assumed product fact.

## 12. Publish a self-contained external view

The preview prompt prepares a manual, export-only view of the completed work and
makes no provider write. After review, the publish prompt sends a self-contained
representation to the named target.

Stable publication records allow later runs to update existing parent and child
items, create only missing items, and skip unchanged fields. External tracker
state never becomes authority over workspace intents, plans, execution, or
Product Knowledge.

## Expected end state

After all twelve sections are completed deliberately:

- the workspace has a stable product identity and two exact repository bindings;
- Product Knowledge describes the implemented invoice behavior and constraints;
- the approved intent and grounded plans preserve the decision and reasoning;
- every implemented candidate has independent evidence appropriate to its tier;
- Git and release actions are separately authorized and recorded;
- plan completion reflects explicit human status authority; and
- the external tracker holds an idempotent view rather than a second source of
  truth.
