---
name: audit-cc-skills
description: Audit the shipped cc-* skills in this source checkout against the CLI, the workspace structure, the entry instruction, and each other, and report findings for a pull request. Use when a change touches product/skills/, when asked to audit or review the shipped skills, or before agent tests for a skill change. Maintainer tool; never shipped; edits nothing.
---

# Audit the shipped cc-* skills

Produce a findings report on the shipped skills under `product/skills/`. This
skill is maintainer tooling for the source checkout. It is not a generated
workspace skill, appears in no manifest, and ships nowhere.

## Rules

- Read files and run the CLI. Edit nothing in the checkout: report findings and
  suggest fixes, but do not apply them. A fix to a shipped skill is a product
  change with its own design case and agent tests.
- Build everything in a fresh temporary directory (`mktemp -d`), never in the
  checkout and never over existing output.
- Start no coding host and call no model API. The audit costs one coding
  session, and a clean audit does not replace agent tests: only runs show whether
  a host follows a skill.
- Audit the 13 shipped `cc-*` skills only, not the maintainer skills in
  `.agents/skills/`.
- Judge each skill against the code and files in this checkout, never against an
  earlier version of the skills.
- Claude Code and Cursor trigger on a skill's `route.md` description and Codex
  triggers on its `SKILL.md` description. They may differ on purpose. A
  difference is not a finding; check only that each description still fits what
  the skill does.

## Set up

From the repository root, with `T=$(mktemp -d)`:

1. Build the CLI with the version the template pins, so the workspace does not
   warn about a mismatch:
   `go build -trimpath -ldflags "-X main.version=$(cat CLI_VERSION)" -o "$T/cli" ./cmd/context-circuit`.
2. Make a throwaway workspace: `git init -q -b main "$T/ws"`, then
   `"$T/cli" --workspace "$T/ws" init --name Acme --purpose "Invoicing" --member maya --member-name Maya`,
   then `"$T/cli" --workspace "$T/ws" repo init --id api --path "$T/ws/repositories/api" --base main`.
3. Capture `"$T/cli" help` and `help <topic>` for every topic it lists, and the
   `--json` output of these:

   ```sh
   "$T/cli" --workspace "$T/ws" --json check
   "$T/cli" --workspace "$T/ws" --json status
   "$T/cli" --workspace "$T/ws" --json record create --kind intent --slug sample --title "Sample intent"
   "$T/cli" --workspace "$T/ws" --json record create --kind plan --slug sample --title "Sample plan" --repo api
   ```

   On this setup `check` exits 1 with one issue, because the workspace is a Git
   checkout that `workspace connect` was never run on; its JSON is still the
   evidence. These outputs are the evidence for C1 and C3.
4. List what each skill names: commands, flags, JSON fields, next-step strings,
   file paths, and other skills. Check each against step 3, the Go code in
   `internal/`, and `scripts/release-manifest.txt`.

## Checks

For each skill and each check, find the evidence or record a finding.

General checks:

- **G1** Frontmatter is valid YAML, `name` equals the folder name, and
  `description` is present.
- **G2** The description says what the skill does and when to use it, specifically
  enough to choose it over a neighbouring skill.
- **G3** The body is focused imperative instruction. Report its size in bytes and
  approximate tokens (bytes divided by 4); rarely needed detail belongs in a doc
  the skill names.
- **G4** Every file a skill names exists at that path in the built workspace.
- **G5** Bundled scripts (under `cc-cli/scripts/`) run and are documented.
- **G6** Examples are complete and correct as written.
- **G7** Nothing is duplicated across files where the copies could drift.
- **G8** No vague trigger, rambling, or dead reference.

Context Circuit checks:

- **C1 The CLI is truthful.** Every command, flag, JSON field, and next-step
  string a skill names exists in the CLI built from this commit and means what the
  skill says. Existing Go tests cover commands and flags only.
- **C2 Launcher discipline.** Every call goes through the workspace launcher
  (`.context-circuit/cli`, or `.context-circuit/cli.ps1` where Windows is meant); no bare command;
  no `CC_CLI_VERSION=local` except where `cc-cli` explains it.
- **C3 Records and folders.** Paths, frontmatter fields, and headings match
  `internal/workspace/records.go` and the real layout: intent and plan paths,
  plan folders and `required_files`, per-member id files, members, bands and
  roles, and `*.local.yaml`.
- **C4 Consistent with `product/AGENTS.md.in`.** No skill weakens, widens, or
  restates a gate differently: approval versus planning, commit authorization,
  push, pull request and merge, attribution, `sources/` passive, borrowed
  knowledge read-only, the setup gate.
- **C5 Consistent across skills.** The same rule reads the same way everywhere
  (for example what follows approval and completion, and that running a plan does
  not edit it); hand-offs name the right skill.
- **C6 Per-host routing.** Each skill has a `route.md` mapped to the Claude Code
  and Cursor skill paths in `scripts/release-manifest.txt` and a `SKILL.md`
  mapped for Codex; names match the folder; host-specific text (the Codex exec
  cell, the Claude Code Agent tool, the Cursor Task tool, the agent definition
  folders) matches what `agent setup` and `agent dispatch` produce.
- **C7 Neighbouring triggers.** The descriptions tell `cc-direct` from
  `cc-intent`, `cc-plan` from `cc-stacked`, `cc-complete` from `cc-knowledge`, and
  `cc-review` from ordinary review, on every host.
- **C8 Nothing is stale.** No retired terms, outdated commands or paths, or
  hard-coded version that should read "or later"; no default `--reuse`;
  `planning_required` only as the sign of an old CLI.
- **C9 Agrees with this checkout's `context/`.** Notes on the same subject say the same. When they
  disagree, the code decides which side is wrong; report the disagreement.
- **C10 Language rules.** Record headings and field names stay English;
  repository text follows the repository; prose follows the member.
- **C11 Named docs exist and agree.** Every shipped doc a skill names exists and
  does not say something different.
- **C12 Reachability.** `product/AGENTS.md.in` names every skill, and a skill
  names the skills it hands off to.

## Report

Produce one report that pastes into a pull request as it is:

1. A findings table, one row per finding: skill, check id, severity (`wrong`,
   `stale`, `inconsistent`, `unclear`), `file:line`, shipped (yes or no), evidence (the CLI output or
   code location that contradicts the text), and a suggested fix. Report nothing
   without evidence; if a suspected problem turns out to be fine, say so in the
   summary instead of the table.
2. A per-skill summary: the checks that passed, the findings, and the size from
   G3.
3. Candidate Go tests: findings a test could check mechanically, so they stop
   being re-checked by hand. List them; do not write them.

A finding in a file under `product/` or `template/` is shipped. Findings never change anything
by themselves: hand the report to the maintainer.
