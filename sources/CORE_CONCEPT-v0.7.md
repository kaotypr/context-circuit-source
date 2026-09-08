# Context Circuit — Core Concept

> A maintainer-facing conceptual map of what Context Circuit **is**, how it
> **works**, and how the same repository serves two purposes at once: as the
> product **source** and as the shippable **template**.
>
> This document is descriptive, not authoritative. The binding rules live in
> `wrapper/contracts/invariants.yaml` (one rule, one owner) and the conversational
> contract lives in `wrapper/adapters/WORKFLOW.md`. Where this document and a
> contract disagree, the contract wins.

---

## 1. What Context Circuit is

Context Circuit is a **universal project workspace for AI-assisted work**. You
talk to an agent in ordinary language; the agent holds durable **Product
Knowledge**, turns requests into **grounded plans**, and executes those plans
**safely across one or more Git repositories** — with a hard separation between
*thinking about a change*, *approving it*, *making it*, *proving it*, and
*shipping it*.

The whole design exists to answer one question safely:

> "An AI made changes across my repositories. How do I know exactly what it did,
> that a human approved it, that something independent checked it, and that
> nothing shipped without me asking?"

Its answer is a set of **non-negotiable separations**:

```mermaid
flowchart LR
    A[Inspect] -->|distinct| B[Mutate knowledge]
    C[Approve] -->|distinct| D[Execute]
    E[Change repositories] -->|distinct| F[Deliver / ship]
    G[Verify] -->|distinct| H[Complete]

    classDef sep fill:#eef,stroke:#557,stroke-width:1px,color:#1f2937;
    class A,B,C,D,E,F,G,H sep;
```

Each arrow is a boundary a human crosses on purpose. Nothing on the right side
happens because something on the left side "looked good."

---

## 2. The dual nature: source **and** template

This one repository (`context-circuit-source`) is both the factory and a working
example of what the factory produces.

```mermaid
flowchart TB
    subgraph SRC["context-circuit-source (this repo)"]
        direction TB
        W["wrapper/ — shipped runtime, contracts, schemas, adapters"]
        AG[".agents/skills/ + agents/ — skill packets & role deltas"]
        DOC["docs/ — shipped guides & templates"]
        T["template/ — blank mutable seed"]
        MAINT["context/ · plans/ · sources/ · test/ · scripts/<br/>(maintainer-only, never shipped)"]
    end

    SRC -->|release assembly<br/>scripts/| TMPL

    subgraph TMPL["context-circuit-template (released artifact)"]
        direction TB
        RW["wrapper/ (copied)"]
        RA[".agents/skills/ + agents/ (copied)"]
        RD["docs/ (copied)"]
        RT["template/ contents → become the blank workspace"]
    end

    TMPL -->|user instantiates| WS["A live universal project workspace"]

    classDef ship fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    classDef seed fill:#fff8e1,stroke:#aa4,color:#1f2937;
    classDef never fill:#fdecea,stroke:#c44,color:#1f2937;
    class W,AG,DOC,RW,RA,RD ship;
    class T,RT seed;
    class MAINT never;
```

Two hats, one checkout:

| Hat | Meaning | Where it lives |
| --- | --- | --- |
| **Source** | The maintainer material that *produces* the product: design, tests, release scripts, the maintainer's own Product Knowledge and plans. | `sources/`, `test/`, `scripts/`, `context/`, `plans/` |
| **Template** | The clean, uninitialized workspace a user actually receives and fills in. | `template/` (seed) + the shipped `wrapper/`, `.agents/`, `agents/`, `docs/` |

The **release boundary** (`wrapper/manifest.yaml`) makes this precise:

```mermaid
flowchart LR
    subgraph shipped["shipped (replaceable on upgrade)"]
        s1[wrapper/]
        s2[.agents/skills/]
        s3[agents/]
        s4[docs/]
    end
    subgraph seed["mutable_seed (copied once)"]
        m1[template/]
    end
    subgraph never["never_ship"]
        n1[.git/ · .runtime/]
        n2[context/ · plans/ · sources/]
        n3[test/ · scripts/ · repositories/]
    end

    classDef ship fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    classDef seed fill:#fff8e1,stroke:#aa4,color:#1f2937;
    classDef never fill:#fdecea,stroke:#c44,color:#1f2937;
    class s1,s2,s3,s4 ship;
    class m1 seed;
    class n1,n2,n3 never;
```

And the **upgrade boundary** protects the user: an upgrade may replace only the
shipped machinery; it must always preserve what the user created.

```mermaid
flowchart LR
    U["Upgrade a workspace"] --> R["Replace: wrapper/, .agents/skills/, agents/, docs/"]
    U --> P["Preserve: workspace.yaml, context/, sources/, plans/,<br/>repositories.local.yaml, repositories/, .runtime/"]

    classDef repl fill:#e3f2fd,stroke:#37a,color:#1f2937;
    classDef pres fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    class R repl;
    class P pres;
```

> **Why it matters:** the machinery evolves; the user's knowledge, plans, and
> repository bindings never get clobbered by an upgrade.

---

## 3. The mental model: one workspace, many repositories

A workspace is not a repository. It is an **agent-oriented coordination layer**
that sits *above* one or more Git repositories and holds the knowledge and plans
that span them.

```mermaid
flowchart TB
    WS["Universal Project Workspace"]
    WS --> K["Product Knowledge<br/>(context/ — accepted facts)"]
    WS --> PL["Plans<br/>(plans/ — intent + human status)"]
    WS --> RT[".runtime/<br/>(execution evidence, host-local)"]

    WS -.binds.-> R1["repo A<br/>(anchor: develop)"]
    WS -.binds.-> R2["repo B<br/>(anchor: main)"]
    WS -.binds.-> R3["workspace root as repo '.'<br/>(optional)"]

    classDef ws fill:#ede7f6,stroke:#63c,color:#1f2937;
    classDef repo fill:#e0f7fa,stroke:#087,color:#1f2937;
    class WS,K,PL,RT ws;
    class R1,R2,R3 repo;
```

Key identity rules (owned by `INV-REPO-*`):

- **Portable identity** (`workspace.yaml`) holds only a logical key, an optional
  credential-free URL, and an optional `default_branch`. **No machine paths, no
  credentials.**
- **Local binding** (`repositories.local.yaml`, gitignored) holds the
  host-specific path and the **`anchor_branch`** — the required execution base
  and default pull-request target.
- `default_branch` is only clone/setup guidance; it is **never** silently used as
  the execution base or PR target.

---

## 4. Layered architecture

Context Circuit is built as strict layers, each with a narrow job. Higher layers
carry *intelligence*; lower layers carry *determinism*.

```mermaid
flowchart TB
    subgraph L5["Conversation (host transport)"]
        HOST["Codex CLI · Claude Code · Cursor Agent<br/>= thin transports only"]
    end
    subgraph L4["Roles"]
        CO["Coordinator (root agent)"]
        WR["Worker (bounded writer)"]
        VE["Verifier (independent, read-only)"]
    end
    subgraph L3["Skill packets"]
        SK[".agents/skills/*/SKILL.md<br/>read-as-procedure"]
    end
    subgraph L2["Contracts"]
        INV["invariants.yaml (one rule, one owner)"]
        SCH["schemas/*.yaml (record shapes)"]
    end
    subgraph L1["Deterministic runtime"]
        ENG["wrapper/runtime/engine.sh<br/>host-neutral POSIX library"]
    end
    subgraph L0["State on disk"]
        STATE["workspace.yaml · context/ · plans/ · .runtime/ · Git"]
    end

    HOST --> CO
    CO --> SK
    WR --> SK
    VE --> SK
    SK --> INV
    SK --> ENG
    CO --> ENG
    ENG --> STATE
    INV -.governs.-> SK
    SCH -.shapes.-> STATE

    classDef host fill:#fff3e0,stroke:#e80,color:#1f2937;
    classDef role fill:#ede7f6,stroke:#63c,color:#1f2937;
    classDef skill fill:#e1f5fe,stroke:#08a,color:#1f2937;
    classDef contract fill:#f3e5f5,stroke:#93c,color:#1f2937;
    classDef rt fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    classDef state fill:#eceff1,stroke:#678,color:#1f2937;
    class HOST host;
    class CO,WR,VE role;
    class SK skill;
    class INV,SCH contract;
    class ENG rt;
    class STATE state;
```

The single most important architectural line:

> **Intelligence lives in roles and skills. Determinism lives in the runtime.**
> The runtime never interprets knowledge, writes plans, routes conversation, or
> ships anything. Roles never bypass the runtime's safe state operations.

---

## 5. Authority order

When instructions could conflict, this fixed order decides who wins
(owned by `wrapper/adapters/WORKFLOW.md`):

```mermaid
flowchart TB
    A1["1 · Host & system instructions"] --> A2["2 · AGENTS.md safety spine + WORKFLOW.md"]
    A2 --> A3["3 · workspace.yaml identity + invariants.yaml ownership"]
    A3 --> A4["4 · Product Knowledge (context/)"]
    A4 --> A5["5 · Plan files (plans/&lt;id&gt;/)"]
    A5 --> A6["6 · Runtime evidence (.runtime/)"]

    A6 -. "never overrides" .-> A5

    classDef top fill:#fdecea,stroke:#c44,color:#1f2937;
    classDef mid fill:#fff8e1,stroke:#aa4,color:#1f2937;
    classDef low fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    class A1,A2 top;
    class A3,A4 mid;
    class A5,A6 low;
```

The last edge is deliberate and load-bearing: **runtime execution evidence never
overrides human plan status.** A plan is `done` only because a human said so
(§10), never because a machine record implies it.

---

## 6. Invariants: one rule, one owner

### In plain terms

An **invariant** is a promise the system always keeps — a rule about how Context
Circuit must behave, no matter which host, which repository, or which agent is
acting. "Only an approved plan may execute." "A pull request targets the anchor
branch, never the default branch." "Never silently clean up failed work." Each
such promise gets a **stable name** (`INV-<AREA>-NN`, e.g. `INV-EXEC-01`) and is
written down in exactly **one place**.

Think of `wrapper/contracts/invariants.yaml` as the project's **constitution plus
a table of contents**: it states every rule once and records which single file is
responsible for it. It is *not* program code that runs — it is the authoritative
index that people and agents consult to answer "what is the rule here, and who
owns it?"

### What problem it solves

Context Circuit's behavior is spread across many files: role descriptions, skill
packets, schemas, the runtime, and host adapters. Without a single owner per rule,
two files would eventually describe the *same* rule slightly differently and drift
apart — "parallel policy." The invariant model forbids that:

> **One rule has one owner. Every other file may *reference* the rule by its ID,
> but must never re-state or re-define it.**

```mermaid
flowchart TB
    subgraph BAD["❌ Without invariants — parallel policy drifts"]
        b1["writer.md: 'commit like this…'"]
        b2["cc-execute: 'commit like that…'"]
        b3["adapter: 'and also this…'"]
        b1 -.slowly disagree.- b2
        b2 -.slowly disagree.- b3
    end
    subgraph GOOD["✅ With invariants — one owner, others cite it"]
        INV["INV-COMMIT-01<br/>(defined once in invariants.yaml)"]
        g1["writer.md"] -->|cites ID| INV
        g2["cc-execute"] -->|cites ID| INV
        g3["adapter"] -->|cites ID| INV
    end

    classDef bad fill:#fdecea,stroke:#c44,color:#1f2937;
    classDef good fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    classDef owner fill:#f3e5f5,stroke:#93c,stroke-width:2px,color:#1f2937;
    class b1,b2,b3 bad;
    class g1,g2,g3 good;
    class INV owner;
```

### What a rule connects to: the owner map

The file has two parts. `invariants:` lists the rules; `owners:` maps each
*concern* to the one file that actually implements/enforces it. The invariant
**names** the promise; the **owner file** is where the behavior truly lives.

```mermaid
flowchart LR
    RULE["INV-COMMIT-01<br/>'commits follow Conventional Commits,<br/>no AI attribution'"]
    RULE -->|owner map points to| OWN["invariants.yaml<br/>(owns commit_convention)"]
    OWN -->|enforced/expressed in| ENF["writer.md commits · cc-execute steps ·<br/>proven by test/ semantic suite"]
    REF["writer.md · cc-execute · adapters/*.md"] -.->|"reference by ID<br/>(never redefine)"| RULE

    classDef rule fill:#f3e5f5,stroke:#93c,stroke-width:2px,color:#1f2937;
    classDef own fill:#e1f5fe,stroke:#08a,color:#1f2937;
    classDef enf fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    classDef ref fill:#eceff1,stroke:#678,color:#1f2937;
    class RULE rule;
    class OWN own;
    class ENF enf;
    class REF ref;
```

Different concerns have different owners — a schema file owns record shapes, the
runtime owns locking, a role file owns worker behavior. The `owners:` block is the
lookup table from concern → owning file.

### Who reads the invariants, and when

The invariants file is consulted, not executed. Three kinds of reader use it:

```mermaid
flowchart TB
    INV["wrapper/contracts/invariants.yaml"]

    A["🤖 Agents (coordinator/worker/verifier)<br/>at decision time"] -->|"'what is the rule, and which file owns it?'<br/>→ follow the owner, never invent parallel policy"| INV
    M["🛠 Maintainers<br/>when changing behavior"] -->|"'find the canonical owner before editing,<br/>then update the fixture that proves it'"| INV
    T["✅ Semantic acceptance tests (test/)<br/>on every meaningful change"] -->|"prove each invariant still holds"| INV

    classDef inv fill:#f3e5f5,stroke:#93c,stroke-width:2px,color:#1f2937;
    classDef reader fill:#ede7f6,stroke:#63c,color:#1f2937;
    class INV inv;
    class A,M,T reader;
```

- **Agents** consult it to find the *one* authoritative source for a rule so they
  act consistently and don't quietly write a competing policy into a skill or role
  file (`AGENTS.md`: "do not add parallel policy").
- **Maintainers** consult it to locate a rule's owner *before* changing it, then
  update the test that proves it (`WORKFLOW.md`: "identify its canonical owner").
- **Tests** consult it as the checklist of promises the semantic acceptance suite
  must keep proving true.

### Worked example: `INV-COMMIT-01`

> *Rule:* every commit an agent authors follows Conventional Commits and carries
> **no** AI attribution. *Owner:* `invariants.yaml` (`commit_convention`).

```mermaid
flowchart LR
    W["Worker writes a commit"] -->|"must obey"| INV["INV-COMMIT-01"]
    SK["cc-execute skill"] -->|"cites for repair commits"| INV
    AD["adapters/AGENTS.md + CLAUDE.md"] -->|"restate + explicitly cite<br/>'this restates INV-COMMIT-01'"| INV
    TST["semantic suite"] -->|"asserts commit shape<br/>+ no attribution trailer"| INV

    classDef owner fill:#f3e5f5,stroke:#93c,stroke-width:2px,color:#1f2937;
    classDef node fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    class INV owner;
    class W,SK,AD,TST node;
```

The adapter files even say so out loud — they note their commit text "restates
INV-COMMIT-01, which owns the rule." That sentence is the invariant model working
as designed: mention the effect, point at the one owner, never fork it.

### The full catalog

These are all the invariant families and the concern each one guards — the "table
of contents" of promises:

```mermaid
mindmap
  root((Invariants))
    Knowledge
      INV-KNOWLEDGE-01 retrieval-first index
      INV-KNOWLEDGE-02 no silent context accept
    Plan
      INV-PLAN-01 plan.yaml owns status
      INV-PLAN-02 explicit repo + path scope
      INV-PLAN-03 stable never-reused ids
      INV-PLAN-04 grounded, no invented decisions
    Approve / Execute
      INV-APPROVE-01 conversational gate
      INV-EXEC-01 only approved executes
      INV-EXEC-02 one worker, one execution
      INV-EXEC-03 deterministic branch + worktree
      INV-EXEC-04 commit before verify, no amend
    Commit
      INV-COMMIT-01 Conventional Commits + no AI attribution
    Verify / Repair
      INV-VERIFY-01 independent read-only
      INV-VERIFY-02 no verifier means blocked
      INV-REPAIR-01 three-failure limit
    Preserve / Complete / Archive
      INV-PRESERVE-01 never silently clean up
      INV-COMPLETE-01 human marks done, only if verified
      INV-COMPLETE-02 completion record + reconcile
      INV-ARCHIVE-01 exact move, no status change
      INV-ARCHIVE-02 archived is out of context
    Repo / Deliver
      INV-REPO-01..04 portable identity, fail-closed binding
      INV-DELIVER-01 PR targets anchor_branch
      INV-DELIVER-02 block, never infer a remote
    Platform
      INV-SEC-01/02 paths, credentials, passive sources
      INV-RUNTIME-01/02 host-neutral, atomic records
      INV-OWN-01 one-writer lock
      INV-HOST-01 host evidence never authorizes
      INV-SKILL-01 skills are read-as-procedure packets
```

---

## 7. The three roles

Exactly three roles exist. They map cleanly onto whatever a host offers (root
session, task/subagent), but the boundaries never move.

```mermaid
flowchart TB
    subgraph coord["Coordinator — the root conversation"]
        C1["orient · gather context · plan · review"]
        C2["interpret approval · completion"]
        C3["archive / restore · delivery discussion"]
        C4["launch worker & verifier · report in plain language"]
    end

    subgraph worker["Worker — single bounded writer"]
        W1["reads immutable plan snapshot"]
        W2["writes only inside assigned worktrees & declared paths"]
        W3["commits each repo · writes handoff"]
        W4["repairs with NEW commits (never amends)"]
    end

    subgraph verifier["Verifier — independent, read-only"]
        V1["inspects latest commit per repo"]
        V2["replays acceptance/verification evidence"]
        V3["reports passed / failed / blocked"]
        V4["never edits, never repairs, never self-verifies"]
    end

    coord -->|"one execution brief"| worker
    coord -->|"latest revisions"| verifier
    verifier -->|"failure evidence"| worker

    classDef c fill:#ede7f6,stroke:#63c,color:#1f2937;
    classDef w fill:#e3f2fd,stroke:#37a,color:#1f2937;
    classDef v fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    class C1,C2,C3,C4 c;
    class W1,W2,W3,W4 w;
    class V1,V2,V3,V4 v;
```

Hard separations between them:

- The **coordinator** never writes code and never self-verifies.
- The **worker** never approves, never completes, never marks its own work
  verified, never merges/pushes/publishes.
- The **verifier** is a *separate actor*, strictly read-only w.r.t. product
  files. **If the host cannot create an independent verifier, the execution is
  `host-blocked` — nobody self-verifies as a substitute** (`INV-VERIFY-02`).

---

## 8. The full lifecycle

A request flows through explicit, separately-authorized stages. Each stage is a
human decision or an independent check — never an inference.

```mermaid
stateDiagram-v2
    [*] --> Orient
    Orient --> Context: "gather context about X"
    Context --> Orient: proposal recorded (not auto-accepted)

    Orient --> Draft: "create a plan for F"
    Context --> Draft
    Draft --> Draft: "review plan X" (status-preserving)

    Draft --> Approved: "approve plan X" (explicit gate)
    Approved --> Draft: edits reopen discussion

    Approved --> Executing: "execute plan X"
    Executing --> Verifying: worker commits each repo
    Verifying --> Repairing: verifier FAILED
    Repairing --> Verifying: new worker attempt
    Repairing --> Stopped: 3rd failure (evidence preserved)
    Verifying --> Verified: verifier PASSED
    Verifying --> Blocked: no independent verifier

    Verified --> Done: "mark X complete" (human, only if verified)
    Done --> Reconcile: record completion + reconcile knowledge

    Verified --> Delivered: "open a pull request" (separate action)
    Done --> Delivered

    Draft --> Archived: "archive plan X" (no status change)
    Archived --> Draft: "restore plan X"

    Stopped --> [*]
    Blocked --> [*]
    Reconcile --> [*]
    Delivered --> [*]
```

Notice what is **not** an edge: verification does not flow to `Done`; execution
does not flow to `Delivered`; a passing verifier does not complete a plan. Those
gaps are the product.

---

## 9. Conversation → action

The coordinator maps ordinary language onto exactly one contract action. It never
infers a consequential action from "okay" or "looks good."

| You say | Action | Side effect? |
| --- | --- | --- |
| What is this workspace? | Read-only orientation | none |
| Gather context about X. | Propose a context update (with provenance) | proposal only |
| Connect / clone / init repo R. | Register & bind; clone/init only if asked | repo binding |
| Create a plan for F. | Draft a grounded plan | plan draft |
| Review plan X. | Non-executing discussion | may edit draft |
| Approve plan X. | Explicit gate: draft → approved | status change |
| Execute plan X. | Execute **only if approved** | branches, commits |
| Approve plan X and execute it. | Two sequential explicit actions | both |
| What happened with X? | Summarize execution evidence | none |
| Repair the failed X. | Another worker attempt (if allowed) | new commits |
| Mark X complete. | Human completion, **only if verified** | done + reconcile |
| Archive / restore plan X. | Move out of / into active area | no status change |
| Open a pull request for X. | Delivery: source = execution branch, target = anchor branch | PR |

---

## 10. Execution internals (hidden, but knowable)

Execution is deterministic and isolated. The user sees plain language; the
mechanism underneath is fixed by `INV-EXEC-*` and `INV-OWN-01`.

```mermaid
sequenceDiagram
    participant U as Human
    participant C as Coordinator
    participant E as Runtime (engine.sh)
    participant W as Worker
    participant V as Verifier

    U->>C: "execute plan 0001"
    C->>E: execution-begin
    Note over E: validate approval + bindings<br/>capture anchor tips (reject dirty)<br/>acquire one-writer lock<br/>snapshot plan<br/>create cc/{plan}/{repo} branch + worktree per repo
    E-->>C: execution brief + worktrees

    C->>W: launch ONE worker (brief + worktrees)
    W->>W: implement tasks in order, commit each repo to Git
    W-->>C: commits + handoff (a claim, not evidence)
    C->>E: capture commits + handoff into .runtime/

    C->>V: prepare read-only access, launch independent verifier
    V->>E: record own result (runtime enforces read-only)
    V-->>C: passed / failed / blocked

    alt failed
        C->>W: same worker, failure evidence → NEW commit
        Note over C,V: repeat until pass — counter increments each rejection, stops at 3
    else passed
        C-->>U: "verified — nothing marked complete yet"
    end
```

Guarantees baked into this flow:

- **The anchor checkout is never written during execution** — all work happens in
  isolated worktrees off captured anchor tips (`INV-EXEC-03`).
- **Commit before verify; repairs are new commits, never amendments**
  (`INV-EXEC-04`) — history never hides a repair attempt.
- **At most one active writer**, enforced by an atomic exclusive-create lock; a
  competing writer gets read-only or blocked, never a silent steal (`INV-OWN-01`).

### The repair loop and the three-failure limit

```mermaid
flowchart TB
    S["worker attempt (implementation counts as attempt #1)"] --> VF{verifier}
    VF -->|passed| OK["verified evidence"]
    VF -->|failed| INC["failure counter += 1"]
    INC --> CHK{counter < 3?}
    CHK -->|yes| RE["same worker, new commit(s)"] --> VF
    CHK -->|no| STOP["execution STOPS<br/>all evidence preserved"]
    VF -->|blocked| HB["host-blocked<br/>never self-verify"]

    classDef ok fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    classDef bad fill:#fdecea,stroke:#c44,color:#1f2937;
    class OK ok;
    class STOP,HB bad;
```

---

## 11. The runtime engine: a deterministic boundary

`wrapper/runtime/engine.sh` is a small, host-neutral, POSIX-`sh` deterministic
library. It is the *only* thing allowed to mutate workspace/Git/execution state,
and it is deliberately unintelligent.

```mermaid
flowchart LR
    subgraph owns["engine.sh OWNS (deterministic)"]
        o1["safe path/id checks · atomic writes · digests"]
        o2["workspace + binding validation"]
        o3["anchor validation · branch/worktree prep"]
        o4["plan structure + approval-state validation"]
        o5["archive/restore exact moves"]
        o6["execution/attempt records · commit capture"]
        o7["verifier-result + read-only enforcement"]
        o8["three-failure counter · one-writer lock"]
        o9["completion eligibility + record · recovery inspect"]
    end
    subgraph not["engine.sh does NOT own"]
        n1["provider child launch · model prompts/SDKs"]
        n2["Product Knowledge interpretation"]
        n3["plan-writing intelligence"]
        n4["conversational routing policy"]
        n5["confirmation cards / tokens"]
        n6["product test semantics"]
        n7["auto PR/merge/push/publish/deploy/cleanup/completion"]
    end

    classDef ok fill:#e8f5e9,stroke:#4a4,color:#1f2937;
    classDef no fill:#fdecea,stroke:#c44,color:#1f2937;
    class o1,o2,o3,o4,o5,o6,o7,o8,o9 ok;
    class n1,n2,n3,n4,n5,n6,n7 no;
```

> A coordinator **invokes** the engine as a tool; it never *reads the engine's
> implementation* as context. The execution brief is enough to act.

---

## 12. Skills as read-as-procedure packets

Product behavior is packaged as skills at `.agents/skills/<name>/SKILL.md`. A
skill is a **procedure the coordinator reads and follows** — resolved by path, not
a separate authority or route (`INV-SKILL-01`).

```mermaid
flowchart LR
    cc-workspace["cc-workspace<br/>orient · register/clone/init repos"]
    cc-plan["cc-plan<br/>create / review grounded plan"]
    cc-execute["cc-execute<br/>approve · execute · repair"]
    cc-verify["cc-verify<br/>run independent verifier"]
    cc-complete["cc-complete<br/>human completion + reconcile"]
    cc-deliver["cc-deliver<br/>PR / merge / push"]
    cc-archive["cc-archive<br/>archive / restore"]

    cc-workspace --> cc-plan --> cc-execute --> cc-verify --> cc-complete
    cc-execute -.-> cc-deliver
    cc-plan -.-> cc-archive

    classDef s fill:#e1f5fe,stroke:#08a,color:#1f2937;
    class cc-workspace,cc-plan,cc-execute,cc-verify,cc-complete,cc-deliver,cc-archive s;
```

A host may *additionally* surface a skill as a slash command (e.g. Claude Code
symlinks under `.claude/skills/`), but that bridge is a host-local convenience
that grants no route or authority the read-as-procedure path doesn't already
carry. The shipped workspace never contains a host-specific skill directory.

---

## 13. Host adapters: transports, not deciders

Codex CLI, Claude Code, and Cursor Agent CLI all sit on the **same** instruction
surface (`AGENTS.md` / `WORKFLOW.md`). They are transports. What a host reports
about itself is bounded, provider-neutral **`host_evidence`** and it authorizes
**nothing** (`INV-HOST-01`).

```mermaid
flowchart TB
    subgraph hosts["Hosts (transports)"]
        H1["Codex CLI"]
        H2["Claude Code"]
        H3["Cursor Agent"]
    end
    MAP["Shared surface:<br/>AGENTS.md · WORKFLOW.md · invariants.yaml"]
    ROLES["Roles: coordinator / worker / verifier"]

    H1 --> MAP
    H2 --> MAP
    H3 --> MAP
    MAP --> ROLES

    EV["host_evidence:<br/>identity, version, capability,<br/>permission mode, provider status"]
    EV -. "observation, never authorization" .-> ROLES
    EV -. "forbids credentials, payloads,<br/>transcripts, auth state" .-> EV

    classDef host fill:#fff3e0,stroke:#e80,color:#1f2937;
    classDef map fill:#f3e5f5,stroke:#93c,color:#1f2937;
    classDef role fill:#ede7f6,stroke:#63c,color:#1f2937;
    classDef ev fill:#eceff1,stroke:#678,color:#1f2937;
    class H1,H2,H3 host;
    class MAP map;
    class ROLES role;
    class EV ev;
```

- The **root session** is always the **coordinator**.
- A native **task/subagent** maps only to the single **worker** or the
  independent **verifier**.
- A host permission flag is an observation, not a grant. If child creation is
  unavailable → `host-blocked`, read-only, **never self-verify**.

---

## 14. Safety spine (the short version)

Everything above collapses into a handful of always-on rules:

```mermaid
flowchart TB
    R1["Read first: WORKFLOW.md · workspace.yaml · context/INDEX.md<br/>then only what an action needs"]
    R2["One rule, one owner — no parallel policy in skills/roles"]
    R3["Only an approved plan executes<br/>(explicit gate, no cards/tokens)"]
    R4["One writer, one independent read-only verifier<br/>no verifier ⇒ host-blocked, never self-verify"]
    R5["sources/ is passive: read only exact named files<br/>plans/.archived/ only via explicit restore"]
    R6["Completion, PR, merge, push, publish, deploy, archive, cleanup<br/>= separate explicit human actions"]
    R7["Failed/interrupted work is preserved, never silently cleaned"]
    R8["Credentials stay in host Git config / host agent<br/>never in workspace files or runtime records"]

    R1 --> R2 --> R3 --> R4 --> R5 --> R6 --> R7 --> R8

    classDef spine fill:#fff8e1,stroke:#aa4,color:#1f2937;
    class R1,R2,R3,R4,R5,R6,R7,R8 spine;
```

---

## 15. One-paragraph summary

Context Circuit is a universal, agent-oriented project workspace that layers
durable Product Knowledge, grounded plans, and safe multi-repository execution on
top of one or more Git repositories. It draws hard, human-crossed boundaries
between inspecting and mutating, approving and executing, changing and delivering,
and verifying and completing. Intelligence lives in three roles (coordinator,
worker, verifier) and read-as-procedure skill packets; determinism lives in a
small host-neutral runtime; and every product rule has exactly one owning file so
policy can never quietly fork. This repository is simultaneously the **source**
that assembles the product and the **template** — the clean, uninitialized
workspace a user instantiates — with a release boundary that ships the machinery,
seeds the blank workspace once, and preserves everything the user creates across
upgrades.

---

*Descriptive companion to the contracts. Authoritative rules:
`wrapper/contracts/invariants.yaml`. Conversational contract:
`wrapper/adapters/WORKFLOW.md`. Release/upgrade boundary: `wrapper/manifest.yaml`.*
