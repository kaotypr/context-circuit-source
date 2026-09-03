# Spec adversary — i0001-conversation-spec-library

criteria_sound: yes
contract_digest: sha256:aa1cc14d8428b1a4fbff44cdd2b79799cb0c8665f5a9e902b989df7b83764839
note: retroactive record — verdict bound to the committed criteria at Gate 1.

The spec adversary's role is to attack acceptance criteria **before any code exists**,
so the plausible-implementation attack does not apply to a retroactive intent whose
deliverables are already authored and committed. The criteria here instead describe an
artifact that exists and is directly inspectable, and each is mechanically or visibly
checkable against `agent-harness/conversations/` (plot/case parity, descriptive-only
references, the per-phase coverage docs, the reporting-rule discipline in the
dialogues, and the surfaced acceptance-criteria drifts). On that basis the criteria are
sound as a record of the completed decision.

The forward-looking, adversarially-hardened criteria live on the dependent intent
**i0002**, where a real before-code challenge applies (see i0002/adversary.md).
