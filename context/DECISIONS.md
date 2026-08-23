# Decisions

## 2026-08-21 — wrapper/template split

Decision: keep shipped wrapper ownership under `wrapper/` and the blank mutable
seed under `template/`.

Rationale: source identity must not be confused with an instantiated workspace.
Consequence: release assembly overlays root adapters and template state while
preserving no user data.

## 2026-08-21 — separate lifecycle gates

Decision: review, approval, execution, completion, delivery, archive, takeover,
and cleanup remain separate human actions.

Rationale: eligibility is not authorization and runtime evidence is not Done.
Consequence: the router emits an exact gate or read-only recommendation.
