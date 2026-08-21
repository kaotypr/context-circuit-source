# Opus 4.8 Plan v3 — Optimize the real token cost, prove it on a live run

Status: **draft — human review and approval required before execution.** No
merge, publication, deployment, or lifecycle status change is authorized.

Scope: the Context Circuit wrapper source and the workspace template it ships —
instructions, routing, docs, tests, release, and the live-workspace upgrade path.

Goals, unchanged: **less context**, **one clear route**, **one owner per rule**,
**zero change to safety semantics** — now measured against the *actual* outcome,
not a byte proxy.

---

## 1. Thesis, and what this adds beyond the converged design

The two prior drafts converged: measure first, ship the certain token win first,
freeze behavior with normalized fixtures, make every abstraction pay rent, repair
identity manifest-first, treat NL-routing as a reported metric. **v3 takes all of
that as settled** and wins on the three things both converged plans still get
wrong:

- **A. Optimize where the tokens actually are.** Both plans budget a *static*
  one-plan snapshot. But a plan's real cost is a **lifecycle**: root entry + K
  child sessions + R resumes, each re-reading wrapper context. The multiplier —
  not the snapshot — is the lever. v3 budgets the lifecycle and cuts the
  re-read, not just the read.

- **B. Prove the goal, don't proxy it.** Both validate against byte budgets and
  self-authored fixtures. The stated goal is *fewer tokens an agent consumes* and
  *clearer routing* — so v3 adds a **ground-truth harness** that counts real
  input tokens on a live session-entry + run-plan, and a **Tier-0-only routing-
  clarity metric** that couples the two goals: if compression breaks routing, it
  fails. A byte budget that doesn't move real tokens is not accepted.

- **C. Protect the real users.** This is a *shipped template*; workspaces are
  already instantiated and may be mid-plan. Rewriting instruction files changes
  what a deployed workspace reads on its next template pull. Neither plan designs
  that upgrade. v3 defines a **wrapper-version compatibility contract** so a live,
  mid-plan workspace adopts the new wrapper without corrupting runtime state.

Everything else below is the converged design, stated compactly; the novelty is
A, B, C and the sharper rent rule in §2.

---

## 2. Rules every task obeys

1. **Same safety, less context.** Gates, lifecycle values, leases, worktrees,
   recovery, and independent verification are unchanged.
2. **Semantic equivalence, tested.** Fixtures compare normalized route + safety
   *decisions*, never wording. Each fixture is classified `preserve` (must match
   the frozen baseline) or `intentional-correction` (a human-approved change to
   current behavior, so a golden test never freezes a current bug).
3. **Every abstraction pays rent — falsifiably, and confirmed on real tokens.**
   A new always-read artifact must remove **≥ 2×** its own bytes from Tier 0. A
   heavier representation (routing YAML, invariant schema, `template/` move,
   host binding) ships **only if** it cuts its affected lifecycle profile by
   **≥ 10%** *measured by the §4B ground-truth harness, not by byte proxy*, **or**
   makes a material safety property testable that the lighter form cannot. Each
   candidate is decided independently; none are bundled.
4. **Ship the certain win before the speculative one.** The largest reduction
   lands in the first edit phase, with no new substrate; every phase is
   independently reviewable and revertable.

---

## 3. Measured baseline (this checkout) — static, then lifecycle

Static profiles (bytes are the portable gate; tokens reported when a tokenizer
exists):

| Profile | Words | Bytes |
| --- | ---: | ---: |
| Tier-0 / every-session set | 3,068 | 22.3 KB |
| Full root entry (spine + deep) | 8,874 | 64.5 KB |
| Writer wrapper context | ~4,075 | ~30.4 KB |
| Verifier wrapper context | ~4,033 | ~30.1 KB |
| All 15 skill bodies | 7,252 | 51.1 KB |

Duplication: writer-child rule in **12 files**, single-plan entry in 7,
no-global-current-session in 4. Identity defect: the repo is the product yet
`workspace.yaml`/`PROJECT.md` call it "uninitialized."

**Lifecycle cost (the number that matters), modeled per plan run:**

```
plan_tokens ≈ root_entry
            + Σ_children (child_bootstrap + delegated_context)
            + Σ_resumes  (re_entry + handoff_reread)
```

A 3-task plan with a writer child, a verifier child, and one resume re-reads the
wrapper spine and role context **4–5 times**. Static compression helps once;
cutting the re-read helps every time. §4A targets the multiplier.

---

## 4. The three differentiators

### 4A. Lifecycle token model — cut the re-read, not just the read

- Budget the **whole-plan lifecycle**, not a single snapshot (hard limits in §6).
- **Handoffs carry resolved context, not pointers to re-derive.** A child packet
  and a resume record include the already-selected invariant IDs, the chosen
  context set, frozen SHAs, and acceptance criteria — so a child or resume reads
  its bounded packet instead of re-running root orientation. Children still
  inspect *primary* evidence; packets never copy large source summaries.
- **Tier-0 is read once per session, never per task.** Task iteration inside an
  approved plan must not reload the spine.
- Acceptance: a 3-task plan's modeled lifecycle tokens fall by ≥ the static
  reduction **plus** a measured re-read saving, both reported by §4B.

### 4B. Ground-truth harness — measure the actual goal

- `test/token-truth.sh`: runs a real **session-entry** and a real **run-plan**
  against a fixture workspace using an available agent/model, and records **actual
  input tokens consumed**, before vs. after, for root / writer / verifier / full
  lifecycle. Byte budgets are gates on the artifact; this harness is the gate on
  the *outcome*. A change that hits its byte budget but does not move real tokens
  is rejected.
- **Tier-0-only routing-clarity metric:** run the route decision on ≥ 50 fixtures
  **using only the Tier-0 spine** (deep reads disabled). This directly measures
  "one clear route": if the compressed spine can't route correctly on its own,
  the metric drops and the phase fails — coupling the token goal to the routing
  goal so we can't win one by sacrificing the other.
- **Independent-reader route check:** the `reviewer` role (a separate context,
  blind to the routing table) resolves a sample of fixtures; disagreement flags a
  route that is ambiguous to a fresh agent, not just to the author.

### 4C. Live-workspace upgrade — don't break shipped instances

- Add a `wrapper_version` to the wrapper catalog and record it in runtime session
  records at creation.
- Define the **compatibility contract**: which wrapper changes are backward-
  compatible for a workspace mid-plan (instruction/doc compaction, added
  invariant IDs) vs. breaking (renamed record fields, changed route enums), and
  the required migration note for a breaking change.
- A workspace whose live `.runtime/` was created under an older wrapper version
  must either keep working or receive an explicit, human-gated migration — never
  silent corruption.
- Acceptance: a fixture workspace instantiated on the *current* wrapper, with an
  in-flight plan and live lease, upgrades to the new wrapper and resumes its plan
  with no runtime-state loss.

---

## 5. Target model (converged — stated compactly)

- **Two tiers.** Tier-0 spine ≤ 8 KB (`AGENTS.md` safety/precedence/source
  boundary; `WORKFLOW.md` lifecycle + the one router; `workspace.yaml`;
  `INDEX.md` navigation-only; `WORKSPACE.md`+`PROJECT.md` identity). Tier-1 =
  deep contracts, skills, roles, loaded only by the route's named context set.
- **One route decision** with distinct `intent / phase / eligibility /
  authorization / capability / context_set / reason_codes / human_gate`.
  `cc-session-entry` evaluates; `cc-whats-next` renders it and owns no
  conditions. Representation stays a canonical table until §6's cost gate
  promotes it.
- **One owner per rule** via a flat `docs/invariants.md` (ID → owner → canonical
  text); other docs reference IDs; the registry is not Tier 0. Tests assert
  invariants, not prose.
- **Identity** repaired manifest-first (source describes itself; ship-vs-
  maintainer paths explicit); `template/` move only if cost-gated.

---

## 6. Work packages (value-first; each states its halting gate)

| Profile | Hard | Stretch | Min. reduction |
| --- | ---: | ---: | ---: |
| Tier-0 every-session | 8 KB | 6 KB | 64% |
| Full root entry | 20 KB | 16 KB | 74% |
| Writer wrapper | 12 KB | 10 KB | 61% |
| Verifier wrapper | 11 KB | 9 KB | 63% |
| One-plan **lifecycle** total | 47 KB static + measured re-read saving | 39 KB | ≥70% static **and** measured live-token drop |

**V3-001 — Freeze static + lifecycle + build both harnesses.** Record static and
lifecycle profiles; capture ≥ 50 normalized route decisions with safety outcome,
each `preserve`/`intentional-correction`; stand up `test/context-budget.sh`,
`test/token-truth.sh` (§4B), and the Tier-0-only clarity metric; inventory
duplicated rules and prose-pinning tests.
*Gate:* baselines reproducible; every fixture has a golden decision + class;
harnesses run; no normative edits.

**V3-002 — Ship Tier-0 compaction + ownership (the certain win).** `invariants.md`
+ IDs; rewrite spine to ≤ 8 KB; `INDEX.md` navigation-only; one routing table in
`WORKFLOW.md`; enforce 8 KB + the 2× rule.
*Gate:* spine ≤ 8 KB; writer-child + single-plan rules each one owner; `preserve`
fixtures unchanged; **Tier-0-only clarity metric ≥ baseline** (compression didn't
hurt routing); real-token drop reported.

**V3-003 — Unify routing + compile context sets + amortized handoffs.** Sole
evaluator + read-only view; per-route required/conditional reads with budgets;
split runtime guidance by concern; **handoffs carry resolved context (§4A)**.
*Gate:* one router; `preserve` fixtures semantically identical; approved
corrections match review; recommendation-only never authorizes; malformed state
visibly blocks; every route within budget.

**V3-004 — Compact skills + roles + bound child packets.** Trigger + unique
procedure only; owner references replace policy copies; measure root/writer/
verifier/**lifecycle** after each change.
*Gate:* all profile budgets met; skills ≤ 700 w (prefer ≤ 450); no lost invariant
coverage; **lifecycle re-read saving measured by token-truth harness.**

**V3-005 — Identity repair + live-workspace upgrade (§4C).** Source identity +
ship-vs-maintainer manifest; `wrapper_version` + compatibility contract; upgrade
fixture (instantiated current wrapper, in-flight plan + live lease → new wrapper
→ resume).
*Gate:* source/artifact identities correct; artifact clean + host-compatible; a
mid-plan workspace upgrades and resumes with **zero runtime-state loss**.

**V3-006 — Cost-gated spikes + final A/B + handoff.** Independent per-candidate
decisions on routing YAML, `template/` move, host binding — adopt only on the §2
rule *confirmed by real tokens*. Independent read-only verifier reproduces every
number; migration + rollback guide.
*Gate:* all hard budgets met; **static ≥70% AND live-token lifecycle drop
confirmed**; fixtures 100% (`preserve`); clarity metric held; NL metric reported
(gate only after two stable runs); zero safety regression; artifact matches
manifest.

Batches: (001–002) freeze + win, (003–004) router + amortization, (005–006)
upgrade + cost gates. Batching authorizes no commit, push, PR, merge, or publish.

---

## 7. Acceptance criteria (outcome-based)

1. All static budgets in §6 met in **both** source and staged template.
2. **Real input-token lifecycle cost** (token-truth harness) drops ≥ the static
   reduction plus a measured re-read saving — the byte win is confirmed on a live
   run, not assumed.
3. **Tier-0-only routing-clarity ≥ baseline** — compression did not degrade
   routing; independent-reader check shows no new ambiguity.
4. One route decision; distinct intent/eligibility/authorization/phase; ≥ 50
   fixtures 100% on `preserve`, approved on `intentional-correction`.
5. One owner per invariant (writer-child = 1); tests assert invariants, not prose.
6. Identity unambiguous; **a mid-plan instantiated workspace upgrades to the new
   wrapper and resumes with no runtime-state loss.**
7. Heavier machinery ships only with §2 evidence confirmed on real tokens.

---

## 8. Risks & controls

| Risk | Control |
| --- | --- |
| Golden tests freeze a current bug | `preserve` vs `intentional-correction` classification (§2.2) |
| Byte win doesn't move real tokens | token-truth harness is an acceptance gate, not a proxy |
| Compression clears budgets but breaks routing | Tier-0-only clarity metric coupled to each phase gate |
| Fixtures encode the author's assumptions | independent-reader route check by the reviewer role |
| New substrate eats the savings | 2× + ≥10% rent rules, confirmed on real tokens, decided per candidate |
| Shipped workspaces break on template pull | `wrapper_version` + compatibility contract + upgrade fixture |
| Registry becomes extra always-read context | store ownership not rule text; keep it out of Tier 0 |
| Lifecycle amortization drops evidence a child needs | packet-completeness fixture; children still read primary evidence |

---

## 9. Non-goals

No change to gates, lifecycle values, exclusive ownership, worktrees, verifier
independence, archive, recovery, or cleanup. No scheduler, daemon, DB, package
manager, or user-facing command layer. No mandatory YAML contract tree,
`template/` move, or host binding without §6 evidence. No `sources/` scan, no
renaming of `AGENTS.md`/`WORKFLOW.md`/`workspace.yaml`, no bypass for
"small" changes, no merge or publication.

---

## 10. Open decisions

1. Confirm §6 static budgets **and** that a live-token drop is a hard gate (not
   just bytes).
2. Confirm the 2× / ≥10%-on-real-tokens rent rules.
3. Is the ground-truth harness required for local acceptance, release, or both?
   (It needs an available agent/model in CI; if none, it runs on-demand and the
   byte budgets + clarity metric are the CI gate.)
4. Confirm identity is manifest-first with a `wrapper_version` compatibility
   contract; `template/` move only if §6 warrants.
5. Three review batches, or one PR?

---

## 11. Why v3 wins

The converged plan compresses the artifact and proves it against byte proxies.
v3 compresses the artifact **and proves the reduction on real consumed tokens,
across the whole plan lifecycle where the tokens actually accumulate, without
degrading routing, and without breaking the workspaces users have already
deployed.** It optimizes the goal the user actually stated — an agent using less
context and routing more clearly — rather than the proxy that is easy to measure.
Same discipline as the converged design; measured against the real outcome.
