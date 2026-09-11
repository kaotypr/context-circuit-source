# Plan 0045 — Complete real Codex, Claude Code, and Cursor session integration

**Intent:** i030-integrate-real-host-sessions  
**Repository:** agent-harness  
**Tier:** Standard  
**Status:** draft

## Objective

Implement the approved real-host outcome afresh from registered `agent-harness`
main. Correctly launch and parse genuine Codex, Claude Code, and Cursor sessions,
stream only redacted bounded evidence into complete independent per-host runs,
and require strict authenticated execution—not probing—to prove success.

## Grounding

The executable base is clean registered main revision
`3f25ccb187ba320acde306e227d334845d67b392`. It provides scenario resolution,
immutable artifact staging, synthetic Git fixtures, disposable run worlds,
manifests, finalization, retention, cleanup, and reproduction, but no native host
layer.

Plans 0043 and 0044 remain failed diagnostic history, not implementation
dependencies. The final 0044 candidate at
`c37591b414584cf512f6ef21ec2a7ca0e67105a1` proved how little command-string
smoke tests establish: all three driver launches failed while a probe-only live
test passed. Independent direct checks established the working entry points:
Codex requires `--skip-git-repo-check`, Claude stream JSON requires `--verbose`,
and Cursor Agent requires `--trust`. These facts seed new tests and adapters; no
failed commit is reused.

## Decomposition

This is one plan with four ordered tasks because the decisive properties cross
module boundaries: parsers populate evidence, grading trusts it, manifests
persist it, cleanup depends on outcome and process ownership, and live acceptance
must exercise that same public path. Splitting those into separately verifiable
plans could allow another probe/parser shell to appear complete before the
integrated lifecycle is real.

1. Establish bounded streaming, pre-persistence redaction, typed evidence,
   precedence, and owned-process semantics.
2. Implement corrected protocol-specific adapters and adversarial tests.
3. Connect complete per-host worlds, manifests, reports, aggregation, retention,
   cleanup, and reproduction.
4. Prove the public path with meaningful deterministic tests and direct
   authenticated native sessions.

Canonical paths, dependencies, acceptance statements, and runnable commands are
in `plan.yaml`.

## Risks

- Native schemas vary by installed version; translators must preserve unknown
  redacted events while requiring proven root and terminal continuity.
- Structured output and stderr may carry secrets; redaction and buffer limits
  apply before every retained representation, including errors and partial lines.
- Process-group cleanup can harm unrelated sessions if ownership is vague;
  sentinel and descendant tests cover every exit path.
- Live proof consumes authenticated provider time. Checks use synthetic repos,
  minimal read-only prompts, five-minute per-host limits, and supported usage
  controls, but every explicitly requested host must really run.

## Verification outcome required

The deterministic suite rejects probe-only, replayed, mock, synthesized,
malformed, contradictory, unredacted, or incomplete evidence. Both strict live
commands then launch Codex, Claude Code, and Cursor and report three exact passes,
each backed by its own fresh root, native terminal, redacted durable stream,
typed manifest evidence, report, and cleanup decision. Any requested-host skip
or unavailable/missing proof fails with retained diagnostics.
