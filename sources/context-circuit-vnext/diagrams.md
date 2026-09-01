# Diagrams

Mermaid views of the full flow, the roles-per-tier, the object states, and the edge
cases. These render inline (the `cc-system-design` convention). Node labels avoid
parentheses so they render on every viewer.

## 1 — Core planned flow (Standard / Critical)

Where the plan is created, and the two human gates. This is the path for anything
above Explore.

```mermaid
flowchart TD
    R[Request] --> I[cc-intent: INTENT.md and contract.yaml]
    I --> ADV[Spec adversary attacks the criteria]
    ADV --> G1{{Gate 1: human approves the intent}}
    G1 -->|approved, contract_digest frozen| P[cc-plan derives plan.yaml and tasks]
    P --> ENV{envelope check}
    ENV -->|EXCEEDS| RG[Hold and re-gate to human]
    RG --> G1
    ENV -->|WITHIN| EX[Execute: one worker in an isolated worktree]
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

## 2 — Explore and promotion

Explore creates no intent, no plan, no candidate — only the pairing session. A plan
first appears at **promotion**.

```mermaid
flowchart TD
    R[Request: small, live change] --> EX[Explore session: coordinator + one worker]
    EX -->|no intent, no plan, no candidate| WRK[Worker writes in the cc-pair branch and worktree]
    WRK --> Q{real work worth keeping?}
    Q -->|no| ENDS[Human-supervised output, never verified]
    Q -->|yes, PROMOTE| ATT[Create intent + criteria, run adversary, raise tier]
    ATT --> MADE[Plan file created, candidate formed]
    MADE --> STD[Enter the Standard or Critical pipeline from Execute onward]
```

## 3 — Roles and spawning by tier

The coordinator is the root; the other roles are spawned children, and how many
depends on the tier.

```mermaid
flowchart TD
    CO[Coordinator - root session, never writes] --> T{tier}
    T -->|Explore| E1[Spawn: 1 worker]
    T -->|Standard| S1[Spawn at intent time: spec adversary]
    T -->|Standard| S2[Spawn at execution: worker]
    T -->|Standard| S3[Spawn per candidate: independent verifier]
    T -->|Critical| C1[Spawn: spec adversary - full]
    T -->|Critical| C2[Spawn: worker]
    T -->|Critical| C3[Spawn: independent verifier + explicit human completion]
    E1 -.no verifier, no adversary.-> NOTE1[ ]
```

## 4 — Object states

The intent and the plan each have a lifecycle; the candidate has none — it is an
identity that is recomputed and can go void.

```mermaid
stateDiagram-v2
    [*] --> draft
    draft --> approved: adversary passed and human approves
    approved --> archived: after delivery
    note right of approved
        approval freezes contract_digest
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

## 5 — Edge cases

The branches that make the two-gate model safe. Each is a mechanical check, not a
human gate.

```mermaid
flowchart LR
    C1[New commit or criteria change] --> C1b[Candidate changes, prior evidence and acceptance VOID, re-verify and re-accept]
    C2[Plan touches repo or path outside intent scope] --> C2b[Envelope EXCEEDS, hold, re-gate to human]
    C3[Change set integration will not build] --> C3b[BASE_UNBUILDABLE, blocked not a worker failure, human splits or reorders]
    C4[Delivered work not reconciled] --> C4b[Knowledge debt, next grounding blocks at Std or Crit and warns at Explore]
    C5[Anchor drifted at delivery] --> C5b[Rebase, new candidate, re-verify and re-accept]
    C6[Host cannot spawn a verifier] --> C6b[host-blocked, read-only, never self-verify]
    C7[Worker fails verification 3 times] --> C7b[Failure limit reached, execution stops, work preserved]
```

## 6 — Sources, intent, plan, knowledge

How material flows into a decision and back into durable knowledge. `sources/` stays
passive; `intent/` carries the authority; reconciliation feeds `context/`.

```mermaid
flowchart LR
    SRC[sources: raw evidence and system-design docs] -->|grounds| INT[intent: the approved decision]
    PK[context: durable Product Knowledge] -->|grounds| INT
    INT -->|derives| PLAN[plans: task breakdown]
    PLAN --> CHANGE[executed and delivered change]
    CHANGE -->|reconciliation debt| RECON[knowledge proposals]
    RECON -->|human accepts| PK
```
