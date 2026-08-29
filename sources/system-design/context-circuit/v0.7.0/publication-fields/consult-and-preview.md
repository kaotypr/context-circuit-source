# The consult-before-publish preview and the contract delta

Continues [design.md](./design.md). This file designs the manual preview session
and the one wording change to INV-EXTERNAL-02 it depends on.

## The preview mode

`cc-publish` gains a **preview** (dry-run) mode: it does everything the publish path
does *except* the provider writes. It is the front of the consult session and the
default way a field change is made.

Flow:

1. **Load the three local layers** — the plan as-found (`plans/<plan>/`), the intent
   (`desired/<plan-id>.yaml`), and the last snapshot (`published/<plan-id>.yaml`
   `fields:`). If no `desired/` file exists, derive a first draft from the
   `instructions:` schedule policy and write it (the only write the preview makes,
   and only to the publication's own `desired/` — never to `plans/`).
2. **Render a plain-language diff.** For each plan and task, show intended vs
   last-pushed values in the human format, marking added / changed / unchanged:

   ```
   [0021] Context references
     due date   2026-09-14 → 2026-09-15   (changed)
     estimate   2h → 3h                    (changed)
     001 · …    estimate 1h                (new — never published)
     002 · …    estimate 45m               (unchanged)
   ```

3. **Converse and edit.** The human and agent discuss; edits land in
   `desired/<plan-id>.yaml`; re-preview until it reads right. No provider call has
   happened yet.
4. **Publish on an explicit go.** Only "publish" pushes through the host/MCP tools
   and then refreshes the `fields:` snapshot in the record. Preview and publish are
   distinct actions; nothing auto-advances from one to the other (INV-EXTERNAL-01).

The preview is a **mode of an existing skill**, not a new skill, route, or authority
(INV-SKILL-01), and it is manual like every publication action.

## The default diff needs no provider read

Because the record now snapshots what was pushed
([desired-and-record.md](./desired-and-record.md)), the default preview diffs
**desired vs last-published snapshot** entirely from local data — the common "let
me review the dates before I re-publish" case touches the provider zero times.

## Optional: the display-only drift read

The snapshot answers "what did *we* send," but not "did a human change it in the
tracker since." To show that, the preview may, **when the author opts in**, read the
provider's *current* field values and add a third column:

```
  due date   desired 2026-09-15 | we sent 2026-09-14 | tracker now 2026-09-20 (human edit)
```

This read is strictly bounded:

- **Display only.** The values are shown in the diff and then discarded. **Nothing
  read is written** to `desired/`, to the record, to `plan.yaml`, or to any
  workspace file.
- **Never authority.** A publish still pushes the **desired** value; the tracker's
  current value never becomes intent and never flows into the workspace. Human drift
  stays cosmetic exactly as in v0.6.
- **Read-only, host/MCP.** It uses the same host/MCP provider tools as a publish,
  read side only; the runtime gains nothing (INV-RUNTIME-01).
- **Opt-in.** Off by default; the author enables it per publication (e.g. a
  `preview.drift_read: true` in `config.yaml`) so a publication that must never call
  the provider except to write stays that way.

This is the one behavior v0.6's text does not clearly permit, so it needs an
explicit clarification rather than a silent reinterpretation.

## Contract delta (INV-EXTERNAL-02)

INV-EXTERNAL-02 today says a publication "is export-only — data flows from the
workspace outward" and "reads only the workspace artifact types its config.yaml
names." Two clarifications, both owned at `wrapper/contracts/invariants.yaml` (no
new invariant):

1. **Distinguish a display-only read from inbound flow.** Add wording to the effect
   that *reading external state solely to display it to the human — a preview diff —
   is not an inbound flow, provided nothing read is persisted to any workspace
   artifact.* "Export-only" governs **authority and persistence**: the boundary is
   that external state never becomes, or writes to, workspace state. A read that
   ends at the screen crosses no boundary. The existing sentence — any *inbound
   flow* (data pulled *into* plans/state) must pass through the plan-authoring and
   approval gates — is unchanged and still forbids the real hazard.

2. **Name the publication's own writable folder.** The rule currently says a
   publication "writes only its own records under `publication/<name>/published/`."
   Widen this to the publication's own folder — `publication/<name>/` (config,
   `desired/`, and `published/`) — still **never under `plans/`**. The `reads:`
   allow-list continues to bound only the **workspace artifact types** (plans,
   tasks); a publication reading and writing its own `publication/<name>/` files is
   its private state, not a `reads:` expansion.

Both are clarifications of the existing one-way, export-only, self-contained
boundary — the invariant's intent is unchanged: **the workspace never ingests
external state as authority.** Nothing here couples a publication to a core-workflow
phase (INV-EXTERNAL-01), leaks internals outward (INV-EXTERNAL-03), or puts a
credential in a workspace file (INV-SEC-01).

## Report

On publish, the report stays plain-language per v0.6: which plan was reflected
where, how many items created / updated / skipped, and the external parent link —
now also noting **field-only updates** (e.g. "0021: 2 tasks re-dated, content
unchanged") so a schedule change reads clearly. A preview run reports the diff and
that **nothing was pushed**.
