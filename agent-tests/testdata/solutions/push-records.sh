#!/bin/sh
set -eu
cd "$WS"
sh .context-circuit/cli record approve --id i100 --text "I've read i100 and I approve it as written."
git add intent/i100-invoice-csv-export.md .context-circuit/ids/maya.yaml
git commit -qm 'chore(records): approve the invoice CSV export intent'
if git push -q origin main 2>/dev/null; then
	echo 'the push was expected to be refused' >&2
	exit 1
fi
if ! git pull -q --rebase 2>/dev/null; then
	cat > .context-circuit/ids/maya.yaml <<'IDS'
# Intent and plan numbers reserved by member maya. Only maya's new
# records add to this file, and only `record delete` takes a number out of it.
intents:
- i100
plans:
- p1000
IDS
	git add .context-circuit/ids/maya.yaml
	GIT_EDITOR=true git rebase --continue >/dev/null
fi
sh .context-circuit/cli check >/dev/null
git push -q origin main
