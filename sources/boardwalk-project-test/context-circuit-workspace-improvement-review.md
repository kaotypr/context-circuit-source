# Context Circuit workspace improvement review

This is a passive review artifact for the Boardwalk workspace. It records
recommendations for improving Context Circuit itself; it is not an accepted
Boardwalk context summary, an implementation plan, or authorization to modify
the wrapper.

Primary evidence:

- `sources/boardwalk-conversation.txt`
- Agent conversation `dd45100b-2c2c-42f4-9658-63134ac38153`

## Executive summary

The conversation shows that Context Circuit's lifecycle model is valuable:

- workspace identity was accepted before other writes;
- a plan was drafted before implementation;
- approval and execution were separate actions;
- a writer and an independent verifier had distinct roles;
- push, deployment, cleanup, and final status remained separately gated.

The main weakness is implementation overhead. Safety policy was repeatedly
reconstructed through broad file searches, duplicated reads, handwritten
runtime records, and large child-agent prompts. This increased latency and
created opportunities for state drift, ambiguous effects, and verification
claims that were stronger than the evidence.

Recommended priorities:

1. Enforce workspace and context integrity before execution.
2. Make entry probes and path boundaries machine-checkable.
3. Generate runtime records atomically through the engine.
4. Match verification evidence to the actual promised behavior.
5. Clarify gate effects and reduce redundant confirmations.
6. Make validated delegation packets the only source of child scope.

## Findings and recommendations

### 1. High: the workspace wrapper was not versioned

Execution preflight explicitly observed that the wrapper had no Git repository,
but execution continued and created a nested Boardwalk repository.

That leaves `workspace.yaml`, `context/`, plans, and other durable coordination
state outside the version history that owns the project workspace. The product
repository can then be clean and committed while its accepted intent and
workspace context remain untracked.

Recommendations:

- Require a versioned wrapper before executing a product plan.
- Provide a separate, explicit wrapper-bootstrap action for a new workspace.
- Keep `/repositories/` ignored by the wrapper while versioning workspace
  identity, accepted context, plans, and wrapper configuration.
- Block execution when relevant wrapper ownership or dirty-state checks fail.
- Distinguish wrapper Git state from nested implementation-repository Git state.

### 2. High: identity projections became inconsistent

The Boardwalk repository was registered in `workspace.yaml` and
`context/WORKSPACE.md`, but `context/PROJECT.md` continued to report no
registered repositories. `context/INDEX.md` also retained wording written for
the pre-initialization state because changing it was outside the earlier gate's
declared effects.

The narrow write boundary was respected, but Tier 0 now contains conflicting
summaries.

Initialization also said roles had to be confirmed. The user's confirmation
omitted roles, and the agent silently accepted the proposed `none` default.

Recommendations:

- Define identity summaries as projections of `workspace.yaml`.
- Update every identity-owned projection through one atomic engine operation.
- Add `cc_validate_workspace_projections` and run it during entry and preflight.
- Have gate cards enumerate all required projection changes.
- Treat omitted confirmation fields as explicit displayed defaults or reject
  the confirmation as incomplete; do not decide this afterward.

### 3. High: context loading exceeded the bounded-probe design

The initialization and planning phases repeatedly read the same entry, gate,
schema, context, and documentation files. They also used broad searches over
the wrapper and documentation trees.

Planning inspected sibling workspaces, including
`claude-context-circuit-v0.5.0` and `tembiter-workspace`, to find plan examples.
There was no focused scope-expansion decision before those reads.

This conflicts with the intended single-probe design and weakens project
boundaries, provenance, and confidentiality.

Recommendations:

- Make Stage A emit an exact path allowlist and byte budget.
- Require loaders to reject paths outside the selected context set.
- Add an explicit scope-expansion route with a human-visible reason.
- Ship canonical plan, task, gate, and runtime fixtures so agents do not need
  sibling projects as examples.
- Cache evidence by digest and reload only changed references.
- Add transcript regression tests for redundant Tier 0 reads, broad globs, and
  cross-workspace access.

### 4. High: runtime records were manually assembled

Session, receipt, delegation, lease, completion, and handoff records were
created through many individual writes. In the richer agent transcript, one
receipt temporarily contained a placeholder delegation digest and byte count,
then was patched after the delegation was written.

The verifier validation command and verifier launch were submitted together.
That makes the intended ordering—validate first, launch second—dependent on tool
scheduling instead of an engine precondition.

Recommendations:

- Add engine operations such as:
  - `cc_create_session`
  - `cc_create_receipt`
  - `cc_create_delegation`
  - `cc_start_child`
  - `cc_record_completion`
- Generate timestamps, versions, byte counts, and digests in the engine.
- Write records to temporary files, validate them, then publish atomically.
- Reject placeholder or pending revisions in active receipts.
- Validate the complete runtime ownership graph before child launch.
- Require a successful validation token or receipt as a launch precondition.

### 5. High: verification evidence was below the promised behavior

The plan used a production build to prove that one command serves Boardwalk.
A successful build proves compilation, not that the server starts or the
browser can load the application.

Store and API tests covered persistence and filtering logic, but acceptance
criteria promised user-visible drag-and-drop, card editing, search, filtering,
and export. The final handoff acknowledged that drag and search clicks were not
in the automated suite.

Recommendations:

- Classify criteria by required evidence layer:
  - schema
  - store
  - API
  - process
  - browser
  - human exploratory check
- Reject verification mappings below the required layer.
- Add a process smoke test for the documented development command and health
  endpoint.
- Add browser journeys for board creation, card creation, drag-and-drop
  persistence, card editing, filtering, and JSON export.
- Record host-capability waivers explicitly and leave affected criteria
  unproven.
- Run verifier commands in a disposable checkout, or declare and compare every
  permitted generated output.

### 6. Medium: gate effects were technically separated but easy to misread

The approval action itself only changed plan and task status. Later execution
created the product repository and local commit. This is not necessarily an
authorization violation: the plan disclosed an independently versioned
Boardwalk repository, and execution was separately requested.

The wording is still ambiguous. Statements such as "approval does not change
Git" describe immediate approval-action effects, while users may read them as
limits on the implementation that approval authorizes.

Recommendations:

- Show two distinct sections on every gate:
  - **This confirmed action changes now**
  - **A later authorized execution may change**
- Represent effects with canonical fields such as:
  - `git.init`
  - `git.commit`
  - `git.remote`
  - `delivery.push`
  - `workspace.register_repository`
- Separate local version-control work from external delivery.
- Derive cards from the same effect fields used by plan and delegation
  validation.
- Reject a delegation that introduces effects absent from approved intent.

### 7. Medium: child prompts duplicated canonical intent

The writer prompt repeated product requirements, database schema, HTTP routes,
test cases, stop conditions, and file boundaries already represented in the
plan and delegation.

This creates a competing source of truth. A child can follow prompt details that
were not approved or that later diverge from the canonical plan.

Recommendations:

- Keep launch prompts limited to role, workspace root, delegation path, and
  handoff path.
- Put all scope, effects, acceptance criteria, and stop conditions in the
  validated delegation packet.
- Generate task-specific context manifests.
- Require delegation validation to prove it is a subset of approved intent.
- Have children stop when prompt text and delegation disagree.

### 8. Medium: gate UX repeated unchanged work

The separation between initialization, plan approval, execution, finish,
delivery, and cleanup should remain. It prevented accidental implementation and
external side effects.

The friction came from rereading unchanged files and requiring long exact
confirmation phrases after the action and target were already explicit.

Recommendations:

- Give each card a digest and short, state-bound confirmation token.
- Accept `confirm <token>` instead of a long exact sentence.
- Invalidate the token whenever target state or effects change.
- Reuse digested evidence between adjacent lifecycle actions.
- Show a compact effect diff rather than repeating the entire plan summary.
- Consider one-step confirmation only for low-risk deterministic local actions;
  retain two-step confirmation for external or destructive operations.

## Suggested implementation order

### Priority 0: integrity

1. Require or explicitly bootstrap wrapper version control.
2. Validate identity projections.
3. Distinguish immediate gate effects from later execution effects.

### Priority 1: deterministic enforcement

4. Enforce route path allowlists and byte budgets.
5. Generate runtime records atomically.
6. Make validation a hard precondition for child launch.
7. Validate delegation as a subset of approved intent.

### Priority 2: evidence and usability

8. Add evidence-layer metadata and browser/process verification.
9. Add digested confirmation tokens.
10. Reuse unchanged context and shorten child prompts.

## Proposed success criteria

Context Circuit is improved when:

- an unversioned wrapper cannot silently proceed to implementation;
- all Tier 0 identity summaries agree;
- one entry route loads one declared probe within budget;
- sibling workspace reads require explicit scope expansion;
- runtime records contain no placeholders and are atomic;
- no child starts before its delegation and ownership graph validate;
- browser-facing acceptance criteria have browser or explicit human evidence;
- gate text clearly separates immediate effects from later authorized effects;
- the same lifecycle requires materially fewer file reads and prompt bytes;
- delivery and destructive gates remain as strict as before.

## What should remain unchanged

- Recommendations do not mutate accepted Boardwalk Product Knowledge.
- Plan approval must not start execution.
- Runtime evidence must not change plan status.
- Writers remain bounded to assigned paths.
- Verifiers remain independent and do not repair.
- Push, merge, publication, deployment, archive, takeover, and cleanup remain
  explicit human gates.

## Provenance and limitations

This review uses the human-selected export at
`sources/boardwalk-conversation.txt` and the named agent conversation. The text
export omits payloads for some edit, terminal, and child-agent tool calls, so
those details were checked against the richer conversation transcript.

The review describes changes for Context Circuit's product source. It does not
propose implementing those changes inside this Boardwalk workspace.
