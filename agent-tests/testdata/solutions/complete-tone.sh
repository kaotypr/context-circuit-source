#!/bin/sh
set -eu
sh "$WS/.context-circuit/cli" record complete --id p1000 --text 'CSV export is on main and the invoices note now mentions CSV.'
cat > "$WS/context/domains/invoices.md" <<'EOF'
# Invoices

An invoice is voided only while it is a draft; once issued it is credited with a credit note. Invoices export as a plain-text summary or as CSV.

Owner:

- `api@internal/invoice/` — voiding, credit notes, the plain-text summary, and CSV.
EOF
today=$(date -u +%Y-%m-%d)
tmp=$TMPDIR/invoices-index
awk -v today="$today" '
/domains\/invoices.md/ {
	sub(/reviewed 2026-01-15/, "reviewed " today)
	sub(/credit, summary/, "credit, summary, csv")
}
{ print }
' "$WS/context/INDEX.md" > "$tmp"
mv "$tmp" "$WS/context/INDEX.md"
