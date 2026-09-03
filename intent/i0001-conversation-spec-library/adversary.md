# Spec adversary — i0001-conversation-spec-library

- contract_digest: (unset — draft; the runtime was not run, so no digest is frozen)
- criteria_sound: yes
- note: retroactive record.

The spec adversary's role is to attack acceptance criteria **before any code exists**,
so the plausible-implementation attack does not apply to a retroactive intent whose
deliverables are already authored and committed. The criteria here instead describe an
artifact that exists and is directly inspectable, and each is mechanically or visibly
checkable against `template-harness/conversations/` (plot/case parity, descriptive-only
references, the per-phase coverage docs, the reporting-rule discipline in the
dialogues, and the surfaced acceptance-criteria drifts). On that basis the criteria are
sound as a record of the completed decision.

The forward-looking, adversarially-hardened criteria live on the dependent intent
**i0002**, where a real before-code challenge applies (see i0002/adversary.md).
