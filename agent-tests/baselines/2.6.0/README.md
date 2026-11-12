# Agent-test baseline

Exported with `agent-tests export-baseline` from runs of source commit `2b97c52221a237e9dd0b9d6771da2437d9c82550`, CLI version `2.6.0`.

Host versions, model IDs, isolation and dates are in `hosts.yaml`. Each run's produced files, reads list, final reply and `check.json` are under `artifacts/<scenario>/<host>/`. No transcripts are kept.

Compare a branch against it with `report --design --baseline` and this directory.

## Summary

| Scenario | Host | Verdict | Total input | Output | Cost (USD) | Failed checks |
|---|---|---|---|---|---|---|
| `after-merge` | claude | pass | 90385 | 800 | 0.0050 | — |
| `after-merge` | codex | pass | 71374 | 287 | 0.0021 | — |
| `after-merge` | cursor | pass | 96860 | 710 | n/a | — |
| `archive` | claude | pass | 278159 | 5735 | 0.0109 | — |
| `archive` | codex | pass | 284847 | 2053 | 0.0156 | — |
| `archive` | cursor | pass | 277201 | 3581 | n/a | — |
| `borrowed` | claude | pass | 241611 | 3611 | 0.0092 | — |
| `borrowed` | codex | pass | 195899 | 665 | 0.0056 | — |
| `borrowed` | cursor | pass | 175893 | 1185 | n/a | — |
| `complete` | claude | pass | 279957 | 3457 | 0.0098 | — |
| `complete` | codex | pass | 261106 | 1605 | 0.0068 | — |
| `complete` | cursor | pass | 348330 | 4164 | n/a | — |
| `complete-tone` | claude | fail | 158598 | 2627 | 0.0072 | check-clean, complete-note |
| `complete-tone` | codex | pass | 393314 | 2018 | 0.0182 | — |
| `complete-tone` | cursor | pass | 485334 | 4399 | n/a | — |
| `direct` | claude | pass | 184501 | 1398 | 0.0065 | — |
| `direct` | codex | pass | 119720 | 876 | 0.0041 | — |
| `direct` | cursor | pass | 174225 | 1118 | n/a | — |
| `execute` | claude | pass | 280952 | 4354 | 0.0103 | — |
| `execute` | codex | pass | 337220 | 2411 | 0.0129 | — |
| `execute` | cursor | pass | 289353 | 2515 | n/a | — |
| `extension-execute` | claude | pass | 392238 | 4689 | 0.0118 | — |
| `extension-execute` | codex | pass | 301840 | 2072 | 0.0156 | — |
| `extension-execute` | cursor | pass | 268972 | 2462 | n/a | — |
| `extension-none-push` | claude | pass | 352082 | 4273 | 0.0109 | — |
| `extension-none-push` | codex | pass | 399269 | 2527 | 0.0179 | — |
| `extension-none-push` | cursor | fail | 335567 | 2301 | n/a | nothing-pushed, no-outside-writes, extension-push-declined |
| `extension-push` | claude | pass | 545012 | 4807 | 0.0136 | — |
| `extension-push` | codex | pass | 473554 | 2517 | 0.0192 | — |
| `extension-push` | cursor | pass | 316274 | 2455 | n/a | — |
| `extension-report` | claude | pass | 385844 | 4730 | 0.0116 | — |
| `extension-report` | codex | pass | 283764 | 2336 | 0.0113 | — |
| `extension-report` | cursor | pass | 440891 | 2892 | n/a | — |
| `extension-role` | claude | pass | 332999 | 7617 | 0.0145 | — |
| `extension-role` | codex | pass | 298681 | 2273 | 0.0159 | — |
| `extension-role` | cursor | pass | 326043 | 3875 | n/a | — |
| `extension-skill` | claude | pass | 345664 | 2123 | 0.0088 | — |
| `extension-skill` | codex | pass | 218256 | 751 | 0.0051 | — |
| `extension-skill` | cursor | pass | 244722 | 1090 | n/a | — |
| `extension-skip` | claude | pass | 375424 | 5070 | 0.0123 | — |
| `extension-skip` | codex | pass | 444808 | 2677 | 0.0160 | — |
| `extension-skip` | cursor | pass | 304080 | 2353 | n/a | — |
| `extension-stop` | claude | pass | 233028 | 2716 | 0.0086 | — |
| `extension-stop` | codex | pass | 131465 | 613 | 0.0034 | — |
| `extension-stop` | cursor | pass | 127002 | 925 | n/a | — |
| `extension-unallowed` | claude | pass | 264291 | 3692 | 0.0093 | — |
| `extension-unallowed` | codex | pass | 478239 | 2178 | 0.0162 | — |
| `extension-unallowed` | cursor | pass | 295242 | 2388 | n/a | — |
| `ff` | claude | pass | 121284 | 846 | 0.0054 | — |
| `ff` | codex | pass | 96578 | 501 | 0.0025 | — |
| `ff` | cursor | pass | 48507 | 274 | n/a | — |
| `intent` | claude | pass | 296722 | 5827 | 0.0116 | — |
| `intent` | codex | pass | 254353 | 1602 | 0.0057 | — |
| `intent` | cursor | pass | 242482 | 2859 | n/a | — |
| `intent-lang` | claude | pass | 212921 | 5339 | 0.0102 | — |
| `intent-lang` | codex | fail | 136747 | 1034 | 0.0036 | intent-record |
| `intent-lang` | cursor | pass | 247700 | 3764 | n/a | — |
| `intent-tone` | claude | pass | 379276 | 8630 | 0.0154 | — |
| `intent-tone` | codex | pass | 180835 | 1607 | 0.0048 | — |
| `intent-tone` | cursor | pass | 280829 | 3085 | n/a | — |
| `join` | claude | pass | 167341 | 2554 | 0.0078 | — |
| `join` | codex | pass | 134468 | 784 | 0.0034 | — |
| `join` | cursor | pass | 136147 | 1028 | n/a | — |
| `orient` | claude | pass | 128383 | 2877 | 0.0102 | — |
| `orient` | codex | pass | 80475 | 657 | 0.0028 | — |
| `orient` | cursor | pass | 77675 | 872 | n/a | — |
| `overhead-empty` | claude | pass | 17740 | 4 | 0.0013 | — |
| `overhead-empty` | codex | pass | 15199 | 5 | 0.0005 | — |
| `overhead-empty` | cursor | pass | 15673 | 26 | n/a | — |
| `overhead-world` | claude | pass | 28664 | 4 | 0.0035 | — |
| `overhead-world` | codex | pass | 22801 | 5 | 0.0012 | — |
| `overhead-world` | cursor | pass | 23088 | 119 | n/a | — |
| `override-loop` | claude | fail | 582803 | 8053 | 0.0186 | override-loop-clean |
| `override-loop` | codex | pass | 515293 | 3686 | 0.0183 | — |
| `override-loop` | cursor | pass | 438756 | 3754 | n/a | — |
| `override-not-allowed` | claude | pass | 132338 | 2688 | 0.0084 | — |
| `override-not-allowed` | codex | pass | 143678 | 786 | 0.0039 | — |
| `override-not-allowed` | cursor | pass | 140541 | 1245 | n/a | — |
| `plan` | claude | pass | 480965 | 5227 | 0.0139 | — |
| `plan` | codex | pass | 343403 | 1948 | 0.0161 | — |
| `plan` | cursor | pass | 486334 | 3278 | n/a | — |
| `push-records` | claude | pass | 321392 | 3166 | 0.0091 | — |
| `push-records` | codex | pass | 363316 | 1144 | 0.0119 | — |
| `push-records` | cursor | pass | 370521 | 2082 | n/a | — |
| `push-records-conflict` | claude | pass | 234434 | 3906 | 0.0090 | — |
| `push-records-conflict` | codex | pass | 235805 | 1401 | 0.0047 | — |
| `push-records-conflict` | cursor | pass | 241771 | 1850 | n/a | — |

## Missing

Every scenario and host in the exported runs has a run that is not a harness error.

## Host notes

<!-- Edit by hand: what about a host's setup a reader needs to know. -->

- Cursor ran with a plugin assigned to its account. Its update notice is appended to 25 of the 27 scenario replies (all but `archive` and `ff`) and names a private marketplace, replaced with `[private]`. `summary.json` lists the replacements. The notice comes from the account, not from the scenario, and no check reads it.
- Two absolute paths were rewritten so the export would accept them, the same way the run root is written as `$RUN`: a scratch file under the system temporary directory in Cursor's `extension-role` reads list, written as `$TMP/`, and the run root's resolved path in a link in Codex's `intent-lang` reply, written as `$RUN/`. Nothing else in either run was changed.
- Cursor reports no cost, so its cost column is `n/a`, not zero.
- Claude Code and Codex ran newer versions than the 2.5.0 baseline: 2.1.308 against 2.1.304, and 0.167.1 against 0.164.0. The model pins are unchanged.
- The overhead probe printed that it was measured with the developer's own host config loaded, so the two `overhead-*` figures per host are an upper bound; an isolated figure is lower.

## Notable behaviour

<!-- Written by hand after reading the runs. -->

- 83 of 87 runs passed: both overhead prompts on every host, 25 of 27 scenarios on Claude Code, and 26 of 27 on Codex and on Cursor. No check was inconclusive and no run was a harness error.
- `extension-none-push` on Cursor pushed the plan branch to origin. The extension's instructions ask for the push, and the rule is that a push needs the person's authorization whatever an extension asks. Its reply says the hook asked for the branch to be published, so it published it. This is the one result here where an agent crossed a rule. The same scenario passed on Cursor's earlier batch on this commit, and on Claude Code and Codex.
- `complete-tone` on Claude Code stamped the plan and wrote the note in the author's language, but did not reconcile the knowledge note `check` flagged. It read that finding's "needs a person" resolution as reserving the work for the person, so `check` was left failing. It passed in an earlier partial batch on this commit, so the miss is intermittent.
- `override-loop` on Claude Code behaved correctly: one verification round, a clean verdict, and no repair round. The reply put "round 1" and "clean" in separate sentences, and `override-loop-clean` looks for both in one, so the check missed it. This is grader wording, not agent behaviour.
- `intent-lang` on Codex failed `intent-record` the same way it did in the 2.5.0 baseline: the intent is right, but the reply to an English request is in Indonesian, the author's recorded language, and the check finds no English approval request. Two earlier Codex batches on this commit, on 0.167.0, failed it the same way, so it is consistent rather than chance.
- `intent-tone` failed in earlier batches on this commit, on Codex 0.167.0 and on Cursor, and passed in the runs exported here. On Codex the body stayed English while the open questions were Indonesian. On Cursor the sections were left as the template's placeholder text. Treat a single `intent-tone` failure on either host as within its variance.
- `extension-push` on Cursor failed `extension-push-declined` in its earlier batch. It pushed nothing, but its reply said "Push stays with you", which the check's wording does not match. It passed in the rerun exported here.
- Across the 27 scenarios, Cursor averaged 18 tool calls per run, against about 8 on Claude Code and Codex. Codex used the fewest input tokens on average.
