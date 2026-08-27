# External-surface — the `thread` kind

Continues from [design.md](./design.md); the config and record homes it uses are in
[configuration-and-records.md](./configuration-and-records.md). This file specifies
the `thread` kind: opening a **discussion thread** in a chat system (Slack first)
from a plan's **open questions**, so a team can resolve them where they already talk.

## What it does

On `/cc-publish <plan-id>` against a `kind: thread` publication (for example
`thread-slack`), the `cc-publish` skill reads the plan's open questions and posts a
**parent message** plus **one threaded reply per question**. It is one-way, manual,
and self-contained; humans discuss in the thread, and nothing they write flows back.

## Source: the plan's open questions

A plan records what is unresolved as **open questions** (INV-PLAN-04 — missing or
contradictory information becomes an explicit open question, assumption, or risk).
The `thread` kind reads those and nothing else its `config.yaml` does not name
(`reads: [plans]`). It publishes questions — not tasks, not status — so it is a
genuinely different kind from `plan`, on the same backbone.

## Mapping to the chat thread

| Context Circuit | chat |
| --- | --- |
| the plan (its open questions) | a **parent message** that opens the thread |
| each open question | **one threaded reply**, carrying that question's **full description** |

- **Parent message** — headline `[thread] [<plan-number>] <plan title> — open
  questions`, then one line of context (why these need a decision). The `[thread]`
  prefix marks the kind; `[<plan-number>]` is the one allowed cross-reference.
- **One reply per question, fully described.** Each open question is a single reply
  that states the question *and* the context needed to answer it — not a terse
  one-liner — so a reader can engage with just that message. One question, one
  message; discussion happens in that message's own replies.

## External text is self-contained (INV-EXTERNAL-03)

Every message reads as an ordinary team discussion to someone who has never heard of
Context Circuit: no workspace file, path, or internal id, and no internal mechanism.
The plan id in the parent headline is the only cross-reference.

## Language

The thread is written in the publication's configured `language` (see
[configuration-and-records.md](./configuration-and-records.md#language)). `cc-publish`
authors the parent and every reply in that language, whatever language the plan is
written in.

## One-way and non-authoritative

The thread is a place to *ask*, not a mirror of state. Human replies are the point,
and they never return to the workspace (INV-PLAN-01, INV-EXTERNAL-02). Statuses,
task structure, and plan text are not reflected here — only the questions.

## Idempotent re-run (different from the `plan` kind)

A discussion must never be clobbered, so the `thread` kind does **not** re-sync like
the `plan` kind:

- The record keeps the parent `ts` and each owned question's reply `ts`.
- On re-run, `cc-publish` may **edit its own** question messages in place (a question
  was reworded) and **append a reply** for a newly-added question. It **never**
  touches, edits, or deletes a human's reply, and it does not delete messages —
  a question that is no longer open is left in place (optionally a short "resolved"
  note is added as a new reply), so the discussion's history stays intact.

## Record

The thread record is a different shape from the plan record — it is keyed by the
thread and its questions, not by tasks — so the `thread` kind has its own record
schema (`publication-thread-record.yaml`), stored like every publication record
under the publication's `published/` folder:

```yaml
# publication/thread-slack/published/0100-csv-export-reports.yaml
plan: 0100-csv-export-reports
provider: slack
thread:
  channel: <channel-id>
  parent_ts: "…"
  url: "https://…"
questions:
  - { ref: scope,         reply_ts: "…" }
  - { ref: large-reports, reply_ts: "…" }
  - { ref: formatting,    reply_ts: "…" }
```

`ref` is a stable per-question key (derived from the question, not a workspace id)
so a re-run can find and edit the right message.

## Provider wrinkles

- **Slack** — the parent message `ts` is the thread root; each reply sets
  `thread_ts` to it. Editing a message on re-run uses the message update API;
  deleting is avoided entirely (the kind never deletes messages).
- Other chat providers (Teams, …) use the same shape: a root message plus
  per-question replies, keyed by whatever id the provider returns.

## Worked trace

1. **Configure.** `publication/thread-slack/config.yaml`: `kind: thread`,
   `provider: slack`, `reads: [plans]`, `target_ref.channel_id`, and a `language`.
2. **Publish.** *"open a discussion thread for plan 0100's open questions."*
   `cc-publish` posts the `[thread] [0100] …` parent and one fully-described reply
   per open question, then writes the record to
   `publication/thread-slack/published/0100-csv-export-reports.yaml`.
3. **Discuss.** The team replies in the thread; those replies are theirs and are
   never read back into the workspace.
4. **Re-run after a question is reworded.** `cc-publish` edits only that question's
   own message in place; every human reply and every other message is untouched.
