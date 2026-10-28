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
if git pull -q --rebase >/dev/null 2>&1; then
	echo 'the pull was expected to stop on a conflict in i100' >&2
	exit 1
fi
