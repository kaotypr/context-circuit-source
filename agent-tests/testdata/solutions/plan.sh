#!/bin/sh
set -eu
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
sh "$WS/.context-circuit/cli" record create --kind plan --slug invoice-csv-export --title 'Invoice CSV export' --intent i100 --repo api
sh "$here/../../fixtures/bin/setbody.sh" "$WS/plans/p1000-invoice-csv-export/plan.md" "$here/../../fixtures/records/plan-csv.md"
