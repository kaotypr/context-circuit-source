# Template harness v0.6 — evaluation (dimension D delta)

## What grade.sh dimension D does now

v0.5: dimension D printed one INFO line per telemetry row and, in practice, always
"no telemetry … unavailable". v0.6 makes it a real (still soft) comparison:

1. If `$CC_TELEMETRY` is absent or empty → unchanged: report "unavailable"; D is
   not evaluated. (Non-`claude-code` hosts stay in this state.)
2. Otherwise, for each ledger row, look up the action's `budgets`:
   - compare `output_tokens` vs `max_tokens`, `conv_turns` vs `max_turns`, and
     (when present) `agent_turns` vs `max_agent_turns`, `cost_usd` vs `max_cost_usd`;
   - print, per action, **WITHIN** or **OVER** with observed vs budgeted for each
     compared field, plus the reported `context_peak`.
3. **Never change the verdict.** Dimension D remains warning-only: a WITHIN or OVER
   result annotates the run and feeds trend tracking; only A, B, and C's forbidden
   check gate the pass (v0.5 §4.2 / grader.md "Verdict"). The final verdict line is
   unchanged in meaning.

The comparison is a plain numeric check (`awk`/shell), deterministic, no model.

## Report shape

```
--- D. efficiency ledger (soft) ---
[INFO] execute-plan: output 37008/80000 WITHIN · conv-turns 1/3 WITHIN · agent-turns 90/140 WITHIN · cost $8.50/$18 WITHIN · context-peak 5.19M
```

An OVER on any field prints OVER for that field (still `[INFO]`, still soft), so a
maintainer sees drift immediately without the run failing.

## Worked ledgers (measured)

The numbers this scope is calibrated against, read from each case's run result:

### Case 11 — repository grounding (1 plan)

```
action        conv  agent  output  context_peak  cost_usd
execute-plan   2     26     9703    949176        2.46
```
Recalibrated budget: `max_turns 3, max_tokens 20000, max_agent_turns 45, max_cost_usd 5` → all WITHIN.

### Case 10 — run-stack (10 plans)

```
action        conv  agent  output  context_peak  cost_usd
execute-plan   2     91     37441   5297625       8.58
```
Recalibrated budget: `max_turns 3, max_tokens 80000, max_agent_turns 140, max_cost_usd 18` → all WITHIN.

(Output/turns/cost vary run to run; budgets carry headroom so normal variance
stays WITHIN. `context_peak` is reported, never budgeted.)

## Acceptance criteria

The v0.6 template-harness delta is acceptable when the v0.5 harness behavior still
holds and:

1. The `claude-code` driver writes a per-action `$CC_TELEMETRY` ledger derived
   from the runner's own result numbers, with no model call and no change to how
   the coordinator runs.
2. The ledger carries, per action, conversational turns, agent-loop turns,
   generated output tokens, context peak, and cost.
3. `grade.sh` dimension D, when a ledger exists, reports each action WITHIN/OVER
   against its `budgets`, showing observed vs budgeted for every compared field.
4. Dimension D never changes the pass/fail verdict; A, B, and C's forbidden check
   remain the only hard gates.
5. With no telemetry, dimension D degrades to "unavailable" exactly as before.
6. Each budgeted field has a defined unit (`max_tokens` = output tokens,
   `max_turns` = conversational turns, optional `max_agent_turns` / `max_cost_usd`).
7. Cases 10 and 11 are recalibrated from measured actuals and report WITHIN.
8. No product surface changes; the release manifest still ships nothing from
   `test/` or `template-harness/`.
9. The deterministic `sh test/acceptance.sh` suite stays green (this scope does not
   touch it; the harness efficiency ledger is exercised only by the live driver).
