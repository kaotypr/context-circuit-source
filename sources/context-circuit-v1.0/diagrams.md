# Diagrams

Mermaid views of the full flow, the roles-per-tier, the object states, and the edge
cases, each with a short walkthrough. These render inline (the `cc-system-design`
convention). Node labels avoid parentheses so they render on every viewer.

## 1 — Core planned flow (Standard / Critical)

```mermaid
flowchart TD
    R[Request] --> I[cc-intent drafts INTENT.md and contract.yaml, no code read]
    I --> G1{{Gate 1: human approves the intent}}
    G1 -->|correct me| I
    G1 -->|approved, contract_digest frozen| DISC[Discovery spawns: one read-only child per repo, parallel]
    DISC -->|intent itself is wrong| I
    DISC --> MANI[Discovery reports a manifest: files, call-sites, risks, executable done-checks, tier signal]
    MANI --> ENV{envelope check against discovery findings}
    ENV -->|EXCEEDS| RG[Hold and re-gate to human]
    RG --> G1
    ENV -->|WITHIN| P[cc-plan creates plan.yaml and tasks from the manifest]
    P --> ENV2{envelope check against the plan}
    ENV2 -->|EXCEEDS| RG
    ENV2 -->|WITHIN| REV[Human may informally review the plan - no gate]
    REV --> TRIG[Human triggers execution]
    TRIG --> EX[Execute: one worker in an isolated worktree]
    EX --> CAND[Candidate = digest of commits + bases + contract_digest]
    CAND --> T{tier}
    T -->|Standard or Critical| VER[Independent verifier checks the candidate]
    T -->|Explore| SUP[Human-supervised, no verifier]
    VER --> AC[Human accepts the candidate]
    SUP --> AC
    AC --> G2{{Gate 2: human authorizes delivery}}
    G2 --> DEL[Pull request, merge, push]
    DEL --> DEBT[Reconciliation debt emitted]
    DEBT --> DONE[Done inferred, explicit at Critical]
```

**Walkthrough.** A request becomes an **intent** — the decision, with its criteria and
scope — drafted with **no code read**. **Gate 1** is the human approving that plain
intent; approval both confirms the coordinator understood correctly and freezes the
criteria as a digest. Approval **spawns discovery**: one read-only child per
repository, in parallel, reading the real code and **reporting back a manifest** —
file/call-site map, concrete risks, a tier signal, and the executable checks that prove
the outcome criteria. The **envelope check** runs at two points — against discovery's
findings, and again against each plan — and is the only thing standing between
derivation and execution: stay inside the approved scope and work proceeds; step
outside and it is held and re-gated to the human. Discovery can also **kick back to the
intent** if the goal itself turns out to be wrong. From the manifest, the coordinator
**creates the plan(s)** (no plan gate); the human may informally review them, and
**execution is a separate, human-triggered action** — nothing runs until asked.
Execution is unchanged from today — one worker, isolated worktree. Its result becomes a
**candidate**: a fingerprint of the exact commits + bases + frozen criteria. The
**tier** decides what checks the candidate gets — an independent verifier at
Standard/Critical, human supervision at Explore. The human accepts the candidate, then
**Gate 2** authorizes the irreversible delivery. Two human gates total (the two
double-bordered nodes); everything else is mechanical.

## 2 — Explore and promotion

```mermaid
flowchart TD
    R[Request: small, live change] --> EX[Explore session: coordinator + one worker]
    EX -->|no intent, no plan, no candidate| WRK[Worker writes in the cc-pair branch and worktree]
    WRK --> Q{real work worth keeping?}
    Q -->|no| ENDS[Human-supervised output, never verified]
    Q -->|yes, PROMOTE| ATT[Create intent + criteria, raise tier]
    ATT --> DISC[Discovery spawns on approval, reads the pairing-branch code]
    DISC --> MADE[Plan file created from the manifest, candidate formed]
    MADE --> STD[Enter the Standard or Critical pipeline from Execute onward]
```

**Walkthrough.** Explore is the fast path (today's `cc-pair`) — the coordinator plus
**one worker**, live human supervision, and **nothing recorded** beyond a working copy:
no intent, no plan, no candidate. Most quick fixes end there, honestly labeled
"human-supervised, not verified." The interesting arrow is **PROMOTE**: the moment the
human decides the work is real, an intent and its criteria are created and approved,
the tier rises, discovery spawns and reads the code, and *only now* does a plan file
and a candidate exist. So Explore
is the one path that creates no plan — until it is promoted, at which point it joins the
Diagram 1 flow from Execute onward. This is the ramp that replaces the old cliff between
pairing and plans.

## 3 — Roles and spawning by tier

```mermaid
flowchart TD
    CO[Coordinator - root session, never writes] --> T{tier}
    T -->|Explore| E1[Spawn: 1 worker only, no discovery child]
    T -->|Standard| S1[On approval: discovery, 1 per repo, proportionate depth]
    T -->|Standard| S2[At execution: worker]
    T -->|Standard| S3[Per candidate: independent verifier]
    T -->|Critical| C1[On approval: discovery, 1 per repo, exhaustive depth + completeness proofs]
    T -->|Critical| C2[Worker]
    T -->|Critical| C3[Independent verifier + explicit human completion]
```

**Walkthrough.** There are four structural roles, and they are not all spawned every
time. The **coordinator** is the root session you talk to — it never writes code. The
discovery, worker, and verifier are spawned children whose presence the **tier
decides**: Explore spawns just **one worker** (no discovery child, no verifier — at
Explore the human reads the code live alongside the agent); Standard and Critical
spawn **discovery** automatically on approval (one child per repository, in parallel)
and add the **independent verifier** (once per candidate, after code). They never all
run at once — discovery, when spawned, happens right after approval, before any
worktree exists; the worker during execution; the verifier after each candidate.
Critical differs from Standard mainly by discovery's depth — exhaustive, with
completeness proofs required — and an explicit human completion instead of an
inferred one.

## 4 — Object states

```mermaid
stateDiagram-v2
    [*] --> draft
    draft --> approved: human approves the plain intent
    approved --> archived: after delivery
    note right of approved
        approval freezes contract_digest, spawns discovery
    end note
```

```mermaid
stateDiagram-v2
    [*] --> derived
    derived --> regated: exceeds envelope
    regated --> derived: intent widened or plan narrowed
    derived --> executing: within envelope
    executing --> candidate
    candidate --> verified: Standard or Critical
    candidate --> supervised: Explore
    verified --> accepted
    supervised --> accepted
    accepted --> delivered
    delivered --> done
    verified --> voided: new commit or criteria change
    voided --> candidate: recompute and re-verify
```

**Walkthrough.** The **intent** (first diagram) has a tiny lifecycle: `draft` →
`approved` (which freezes the criteria digest) → `archived` after the change ships. The
**plan/candidate** (second diagram) is richer. A derived plan either executes (inside
the envelope) or is `regated` (outside it) and returns once the intent is widened or the
plan narrowed. Execution produces a `candidate`, which is `verified` at Standard/Critical
or `supervised` at Explore, then `accepted`, `delivered`, and finally `done`. The
load-bearing edge is `verified → voided`: **any new commit or criteria change voids the
candidate**, forcing a recompute and re-verify. That single transition is what makes "it
passed earlier" impossible — evidence cannot outlive the exact code and criteria it saw.

## 5 — Edge cases

```mermaid
flowchart LR
    C1[New commit or criteria change] --> C1b[Candidate changes, prior evidence and acceptance VOID, re-verify and re-accept]
    C2[Discovery findings or a plan reach outside intent scope] --> C2b[Envelope EXCEEDS, hold, re-gate to human]
    C3[Change set integration will not build] --> C3b[BASE_UNBUILDABLE, blocked not a worker failure, human splits or reorders]
    C4[Delivered work not reconciled] --> C4b[Knowledge debt, next grounding blocks at Std or Crit and warns at Explore]
    C5[Base drifted at delivery] --> C5b[Rebase, new candidate, re-verify and re-accept]
    C6[Host cannot spawn a verifier] --> C6b[host-blocked, read-only, never self-verify]
    C7[Worker fails verification 3 times] --> C7b[Failure limit reached, execution stops, work preserved]
```

**Walkthrough.** Each row is a mechanical branch, not a human gate — the design's safety
comes from these firing reliably rather than from asking a human at every step. C1 and C5
are the candidate rule doing its job (a change voids evidence). C2 is crown jewel 1 (the
envelope) pulling the human back in when scope drifts. C3 keeps the honest multi-repo
stance — combined work that will not build is blocked, never faked. C4 is the closed
knowledge loop refusing to let reconciliation be forgotten. C6 and C7 are preserved from
today (host-blocked never self-verifies; the three-failure cap stops and preserves).
Every one of them **fails safe**: it either re-checks, re-gates, blocks, or preserves —
none of them proceeds on a guess.

## 6 — Sources, intent, plan, knowledge

```mermaid
flowchart LR
    SRC[sources: raw evidence and system-design docs] -->|grounds| INT[intent: the approved decision]
    PK[context: durable Product Knowledge] -->|grounds| INT
    INT -->|approval spawns| DISC[discovery: reads the real code, reports a manifest]
    DISC -->|coordinator creates| PLAN[plans: task breakdown]
    PLAN --> CHANGE[executed and delivered change]
    CHANGE -->|reconciliation debt| RECON[knowledge proposals]
    RECON -->|human accepts| PK
```

**Walkthrough.** This is the information flow around a change. `sources/` — raw evidence
and, for larger work, the multi-topic **system-design** docs — is passive material that
*grounds* an intent but carries no authority. Durable **Product Knowledge** (`context/`)
also grounds the intent. The **intent** is where authority sits: its approval spawns
**discovery**, which reads the real code and reports a manifest the coordinator uses to
create the **plan(s)**, which become an executed, delivered change. Delivery raises
**reconciliation debt**, which produces knowledge **proposals** — and only a **human
acceptance** folds
them back into Product Knowledge. The loop is closed but never automatic: the arrows into
`context` always pass through a human. This is also the layering that answers
"multi-topic big picture": it lives in `sources/system-design/` on the far left and
spawns one intent per topic.
