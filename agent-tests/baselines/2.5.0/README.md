# Agent-test baseline

Exported with `agent-tests export-baseline` from runs of source commit `650c7e0909d4d08a7fb56ca924b1678a41f29706`, CLI version `2.5.0-rc.1`.

Host versions, model IDs, isolation and dates are in `hosts.yaml`. Each run's produced files, reads list, final reply and `check.json` are under `artifacts/<scenario>/<host>/`. No transcripts are kept.

Compare a branch against it with `report --design --baseline` and this directory.

## Summary

| Scenario | Host | Verdict | Total input | Output | Cost (USD) | Failed checks |
|---|---|---|---|---|---|---|
| `archive` | claude | pass | 433854 | 5555 | 0.0130 | — |
| `archive` | codex | pass | 231622 | 1210 | 0.0070 | — |
| `archive` | cursor | pass | 254207 | 2212 | n/a | — |
| `borrowed` | claude | pass | 165618 | 3242 | 0.0080 | — |
| `borrowed` | codex | pass | 158765 | 817 | 0.0037 | — |
| `borrowed` | cursor | pass | 273669 | 1787 | n/a | — |
| `complete` | claude | pass | 296725 | 4002 | 0.0100 | — |
| `complete` | codex | pass | 261384 | 1991 | 0.0075 | — |
| `complete` | cursor | pass | 260537 | 2984 | n/a | — |
| `direct` | claude | pass | 175968 | 1403 | 0.0061 | — |
| `direct` | codex | pass | 135800 | 729 | 0.0029 | — |
| `direct` | cursor | pass | 144613 | 1149 | n/a | — |
| `execute` | claude | pass | 338080 | 4675 | 0.0108 | — |
| `execute` | codex | pass | 180740 | 1652 | 0.0042 | — |
| `execute` | cursor | pass | 221140 | 2605 | n/a | — |
| `intent` | claude | pass | 221137 | 5316 | 0.0107 | — |
| `intent` | codex | pass | 165425 | 752 | 0.0038 | — |
| `intent` | cursor | pass | 243911 | 2274 | n/a | — |
| `intent-lang` | claude | pass | 217343 | 5501 | 0.0105 | — |
| `intent-lang` | codex | fail | 132070 | 765 | 0.0034 | intent-record |
| `intent-lang` | cursor | pass | 195731 | 2371 | n/a | — |
| `join` | claude | pass | 182586 | 2701 | 0.0091 | — |
| `join` | codex | pass | 153489 | 1024 | 0.0049 | — |
| `join` | cursor | pass | 105228 | 1233 | n/a | — |
| `orient` | claude | pass | 147145 | 3004 | 0.0086 | — |
| `orient` | codex | pass | 75181 | 407 | 0.0033 | — |
| `orient` | cursor | pass | 71725 | 736 | n/a | — |
| `overhead-empty` | claude | pass | 17757 | 4 | 0.0013 | — |
| `overhead-empty` | codex | pass | 14586 | 5 | 0.0004 | — |
| `overhead-empty` | cursor | pass | 15524 | 65 | n/a | — |
| `overhead-world` | claude | pass | 27700 | 4 | 0.0033 | — |
| `overhead-world` | codex | pass | 21459 | 5 | 0.0011 | — |
| `overhead-world` | cursor | pass | 22434 | 109 | n/a | — |
| `plan` | claude | pass | 341771 | 5270 | 0.0121 | — |
| `plan` | codex | pass | 270468 | 2060 | 0.0078 | — |
| `plan` | cursor | pass | 310904 | 3014 | n/a | — |
| `push-records` | claude | pass | 302277 | 3009 | 0.0092 | — |
| `push-records` | codex | pass | 252955 | 1614 | 0.0063 | — |
| `push-records` | cursor | pass | 344308 | 2031 | n/a | — |
| `push-records-conflict` | claude | pass | 215114 | 2516 | 0.0075 | — |
| `push-records-conflict` | codex | pass | 256516 | 1580 | 0.0053 | — |
| `push-records-conflict` | cursor | pass | 278411 | 1787 | n/a | — |

## Missing

Every scenario and host in the exported runs has a run that is not a harness error.

## Host notes

<!-- Edit by hand: what about a host's setup a reader needs to know. -->

- Cursor ran with a plugin assigned to its account. Its update notice is appended to eight replies (`direct`, `execute`, `intent`, `intent-lang`, `join`, `orient`, `plan`, `push-records`) and names a private marketplace, replaced with `[private]`. `summary.json` lists the replacements. The notice comes from the account, not from the scenario, and no check reads it.
- Cursor reports no cost, so its cost column is `n/a`, not zero.
- The runs measured `release/2.5` while its version files still read `2.5.0-rc.1`. The release that follows changes only the version string the CLI and the template carry; the product and template files are the ones measured.

## Notable behaviour

<!-- Written by hand after reading the runs. -->

- 41 of 42 runs passed: both overhead prompts on every host, all twelve scenarios on Claude Code and Cursor, and eleven of twelve on Codex. No check was inconclusive and no run was a harness error.
- `intent-lang` on Codex failed `intent-record`. The intent itself is right: English headings, an Indonesian body, `created_by: maya`, one numbered open question. But the reply to an English request was written in Indonesian, the author's recorded language, and the check looks for an English approval request, so it found none. The rule is to talk in the language the person uses; the record's language does not set the reply's. An earlier batch on the commit before this one failed the same way, so it is consistent rather than chance.
- In that earlier batch, `borrowed` on Codex once replied without saying the fix goes through the handbook's own review. It passed here. Treat a single `borrowed` failure on Codex as within its variance.
- Across the twelve scenarios, Cursor averaged 16 tool calls per run, against about 7 on Claude Code and Codex. Codex used the fewest input tokens in eight of the twelve.
