# Intention — i0003

_Status: draft, waiting for your approval._

## Intention

Use a shorter, three-digit number for every intent identifier.

The existing intents and every place that refers to them should move together,
so the workspace has one consistent naming style. The rules that create, check,
approve, archive, restore, and connect intents should understand the new style
for all future work too.

## Expectations

- Intent identifiers use three digits after the initial `i`, such as `i001`.
- Existing intent folders and records are renamed consistently, without losing
  their history, status, approval, or relationships.
- Every current reference to an existing intent uses its new identifier.
- New intents, including archived intents, are allocated and validated with the
  three-digit form.
- The change does not alter what an intent means or how approval, planning,
  execution, verification, completion, or delivery work.
- Checks prove that the old four-digit form and stale old references cannot remain
  active in the workspace.

## The plans

1. **Change the naming rule**
   _After this:_ the workspace consistently defines and validates the three-digit
   intent form.
2. **Migrate the existing records**
   _After this:_ current intent folders, records, indexes, relationships, and
   references use the new identifiers together.
3. **Check the whole workspace**
   _After this:_ allocation, lifecycle operations, authorization, and repository
   acceptance checks prove the migration is complete without changing lifecycle
   behavior.

## How carefully this is checked

**Standard**

This touches the identifier contract and runtime lifecycle surfaces, so an
independent verifier should check both the migration and the unchanged lifecycle
behavior.

## Open questions

None.
