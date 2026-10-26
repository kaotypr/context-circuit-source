#!/bin/sh
# sh setbody.sh RECORD_FILE BODY_FILE
# Keeps the record's first YAML frontmatter block and replaces everything
# after it with the body file.
set -eu
if [ "$#" -ne 2 ]; then
	echo "usage: sh setbody.sh RECORD_FILE BODY_FILE" >&2
	exit 2
fi
record=$1
body=$2
tmp=$record.setbody
awk 'BEGIN { state = 0 }
state == 0 && $0 == "---" { state = 1; print; next }
state == 1 && $0 == "---" { state = 2; print; next }
state < 2 { print }
' "$record" > "$tmp"
cat "$body" >> "$tmp"
mv "$tmp" "$record"
