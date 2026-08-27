# External-surface — configuring and publishing

How a template/wrapper user sets up and runs a publication. The whole point is that
a user never learns the internals or edits YAML by hand — they talk to the
coordinator in plain language, and the `cc-publish` skill does the rest. This is the
material that ships in the getting-started docs when the scope is released.

## Prerequisite (one-time, at the host)

Connect the provider at the **host** first — the ClickUp / Slack / Jira / … MCP or
credentials. Config files are credential-free (INV-SEC-01), so the token lives at the
host, never in the workspace. If the provider is not reachable, `cc-publish` reports
`host-blocked` and stops.

## Configuring is create-on-first-use

There is no separate "configure" step. The first time a user asks to publish, if no
matching publication exists, `cc-publish` creates one — writing
`publication/<name>/config.yaml` and asking only for what it needs (which
list/channel, and any authoring `instructions`). After that the publication is reused.

**First-time setup prompts** — any of these:

- *"Publish plan 0100 to ClickUp."* → no `plans-clickup` publication yet, so it asks
  which ClickUp list, creates `publication/plans-clickup/config.yaml`, then publishes.
- *"Set up publishing our plans to the Engineering list in ClickUp."* → creates the
  publication, no publish yet.
- *"Open a Slack discussion for plan 0100's open questions in #eng."* → creates
  `publication/thread-slack/config.yaml`, then posts the `[thread]` parent and one
  reply per question.
- *"…and write it in Bahasa Indonesia, everyday tone, keep technical terms in
  English."* → recorded as the publication's `instructions`.

The user may name the publication or let the convention name it `<subject>-<provider>`
(`plans-clickup`, `plans-github`, `thread-slack`). Several can coexist.

## Triggering (every time after)

- *"Publish plan 0100 to ClickUp."* / *"Push 0100 to our tracker."*
- *"Re-publish 0100."* — updates the same items in place; no duplicates.
- *"Open the discussion thread for 0100's open questions."*
- Slash-command form the host may surface: `/cc-publish 0100`.

Users say **"publish"** for this and **"deliver / open a PR"** for git — the two
never share a word, so intent is unambiguous.

## Instructions (optional)

One free-text field guides how the agent authors, for any kind:

- **Language** — *"write in Bahasa Indonesia."* Default is the plan's language.
- **Tone and terms** — *"everyday conversational tone; keep technical terms in
  English."*
- **Optional field enrichment** — *"estimate a time estimate per task and a target
  date range,"* which the agent fills as best-effort estimates on the tracker's own
  fields.

Instructions only shape wording and optional field values; they never change what is
published or leak workspace internals.

## What gets created (visible to the user)

```
publication/plans-clickup/config.yaml            # what this publication is (credential-free)
publication/plans-clickup/published/0100-….yaml  # the record, after publishing
```

The user can open or hand-tweak `config.yaml`, but never has to — the conversation
writes it. Everything a publication reads or writes lives under `publication/`;
nothing is written under `plans/`.
