#!/bin/sh
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
sh "$WS/.context-circuit/cli" record create --kind intent --slug invoice-csv-export --title 'Invoice CSV export'
sh "$here/../../fixtures/bin/setbody.sh" "$WS/intent/i100-invoice-csv-export.md" "$here/files/intent-id.md"
