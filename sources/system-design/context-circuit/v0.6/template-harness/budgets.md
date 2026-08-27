# Template harness v0.6 — budgets and units

## The units a budget is stated in

A budget is only meaningful once its unit is fixed. v0.6 pins them:

| Field | Unit | Why |
| --- | --- | --- |
| `max_tokens` | **generated (output) tokens** for the action | The clean "work produced" measure. It excludes cache-read/context volume, which is an artifact of conversation length, not of how much the action actually did. |
| `max_turns` | **conversational** turns (human↔coordinator exchanges) | Matches the intended 1–6 magnitude and what a maintainer reasons about. |
| `max_agent_turns` *(optional)* | agent-loop turns (`num_turns`) | Bounds internal loop/tool churn (tens per plan); optional because it is the noisiest number. |
| `max_cost_usd` *(optional)* | dollars (`total_cost_usd`) | The fullest aggregate, incl. sub-agents; useful to bound real spend. |

`context_peak` (max cache-read input tokens) is **always tracked and reported**
but is a **diagnostic, not a budget** — it tells you how large the context window
grew, which is informative for trend-watching but too conversation-dependent to
bound.

## Case `budgets` schema (delta on grader.md §"Case grader schema")

```yaml
budgets:                               # dimension D (soft), keyed by action
  <action>:
    max_turns: N            # conversational turns
    max_tokens: M           # generated (output) tokens
    max_agent_turns: K      # optional — agent-loop turns
    max_cost_usd: C         # optional — dollars
```

`max_turns` and `max_tokens` stay required (backward compatible); the two optional
fields are new. Absent optional fields are simply not compared.

## Recalibration method

Budgets are **targets**, not gates, so they should sit at comfortable headroom
over a known-good run rather than at the edge:

1. Run the case once green; read its `$CC_TELEMETRY` ledger.
2. Set `max_tokens` ≈ 1.5–2× the observed `output_tokens`.
3. Set `max_turns` to the observed conversational turns (usually the case's own
   turn count), and, if used, `max_agent_turns` ≈ 1.5× observed `agent_turns`.
4. If used, `max_cost_usd` ≈ 1.5–2× observed `cost_usd`.
5. Re-run; confirm the ledger reports WITHIN.

### Worked recalibration for the two v0.6 product cases

From measured runs (see [evaluation.md](./evaluation.md) for the full ledgers):

| Case / action | observed output | observed conv / agent turns | observed cost | recalibrated budget |
| --- | --- | --- | --- | --- |
| 11 / execute-plan | ~9.7k | 2 / ~26 | ~$2.5 | `max_turns: 3, max_tokens: 20000, max_agent_turns: 45, max_cost_usd: 5` |
| 10 / execute-plan | ~37k | 2 / ~91 | ~$8.6 | `max_turns: 3, max_tokens: 80000, max_agent_turns: 140, max_cost_usd: 18` |

The old values (`max_tokens` 120k/400k, `max_turns` 5/6) are replaced: the token
budgets tighten toward real output, and `max_turns` drops to the true
conversational count with the agent-loop churn moved to `max_agent_turns`.

## What is not changing

- Dimension D stays **soft, warning-only**. Recalibration makes the numbers
  meaningful; it does not make them a gate.
- The other cases (01–09) keep their existing `budgets`; they may be recalibrated
  opportunistically the first time each is run green under the new ledger, but no
  case is required to change to keep the suite green.
