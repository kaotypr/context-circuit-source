#!/bin/sh
set -eu
# The installed skill runs its named stage around the QA pass; the scripted
# reply stands in for what vet and the tests printed.
cli="$WS/.context-circuit/cli"
cat "$WS/.agents/skills/x-qa-pass-qa/SKILL.md" >/dev/null
sh "$cli" --json extension hooks --stage qa-pass/qa --when before --plan p1000 >/dev/null
sh "$cli" --json extension hooks --stage qa-pass/qa --when after --plan p1000 >/dev/null
