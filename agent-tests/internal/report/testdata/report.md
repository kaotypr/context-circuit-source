## Agent tests

**Change kind:** always-loaded, skill-text
**Surface diff:** 2 workspace files and 1 CLI output differ from base aaaaaaaaaaaa
**Scenarios:** required by the trigger map: `intent`, `plan` · asked for by the reviewer: — · one-off: —
**Before:** committed baseline `2.5.0`

| Host | Host version | Model | Effort | Isolation |
|---|---|---|---|---|
| claude | 2.1.304 | `claude-opus-5-5` | default | empty config dir, no MCP |
| codex | 0.162.1 | `gpt-6-sol` | medium | empty CODEX_HOME, no user config |
| cursor | 2026.10.22-27dcf5f | `auto` | default | empty home |

| Scenario | Host | Branch | Before | Total input tokens (branch / before / change) | Output tokens | Cost per run (branch / before) | Cost basis |
|---|---|---|---|---|---|---|---|
| borrowed | claude | inconclusive | inconclusive | 80k / 80k / +0% | 2k / 2k | $0.20 / $0.20 | reported |
| complete | claude | pass (2 runs) | pass | 110k / 100k / +10% (mean of 2 branch runs) | 3k / 2k | $0.50 / $0.40 | reported |
| direct | cursor | pass | pass | n/a / n/a | n/a / n/a | n/a / n/a | none |
| execute | claude | pass | pass | 120k / 100k / +20% | 1k / 1k | $0.30 / $0.25 | reported |
| intent | claude | pass | pass | 212k / 198k / +7% | 4.1k / 3.9k | $0.61 / $0.57 | reported |
| plan | claude | fail | pass | 880k / 860k / +2% | 12k / 12k | $2.10 / $2.02 | reported |

**Total cost:** $4.21 (list-price equivalent; runs without a cost are omitted) · **Hosts with no token report:** cursor
**Harness errors (excluded, rerun):** 1 · **Inconclusive:** 1

### Warnings (explain each below; the reviewer decides)
- [ ] `execute` on claude: total input tokens above +15% (+20%)
- [ ] `plan` on claude fails on the branch (nothing-pushed, plan-record)
- [ ] `plan` on claude: a lower result than before (fail, before pass)

### Gate checks
- [x] U1 launcher only · [ ] U2 nothing pushed · [x] U3 no commit outside allowed branch
- [x] U4 no attribution · [x] U5 check clean at finish · [x] U6 `sources/` untouched · [x] U7 no writes outside workspace

### Notable behaviour
<!-- What the agent did differently that the design case does not already cover, and any failure.
     Cite the run id; quote at most a few transcript lines. -->

### Evidence
- Report and before/after: pasted above
- Cleaned zip of all runs: bundle.zip
  It holds every run's JSON line, stream and produced files. The runner cleans it: home and run-root
  paths become `$HOME`/`$RUN`, the username is removed, and auth files are dropped. If it finds a
  key-shaped string, it refuses to build the zip.

### Reviewer checklist
Design case
- [ ] Present and before the evidence; light form only where no behaviour change is intended
- [ ] Larger change: links the issue where the reviewer agreed the design case
- [ ] The problem is shown with evidence, not asserted
- [ ] Before/after use the artifact the trigger map names, from the same scenario, fixture, host and model
- [ ] Every "better" claim is visible in the after artifact or a check, with its run count
- [ ] The catch gives a size and a payer for tokens, compliance, existing workspaces and customized templates
- [ ] At least one alternative, including leaving it as it is

Agent tests
- [ ] Surface diff and scenarios match the trigger map; host-specific files ran on that host
- [ ] "Before" is a matching baseline, or base runs on the same machine
- [ ] Every warning is explained, and I accept the explanation (or asked for changes)
- [ ] Opened the zip for any failed or flagged run, and for the before/after pair
- [ ] More runs, scenarios or hosts asked for? (write which, or "none"; ask for more runs when a result is borderline)
- [ ] Hosts the contributor lacked: run by the reviewer (write which, or "none")
