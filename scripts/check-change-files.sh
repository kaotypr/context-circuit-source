#!/bin/sh
# Check that a pull request carries a change file for each product it changes.
#
# A change users can see is described where it is made: the pull request adds
# .changelog/<slug>.md naming the product, and release-notes.sh drafts the
# release request from those files. A change users cannot see says `None` as
# its Release note instead. Without this check a missing note is found at
# release time, by the maintainer, from memory; with it, the author is asked in
# the pull request that made the change.
#
# Usage: check-change-files.sh <base> <head>
#   PR_BODY  the pull request description (env); may be empty
#
# Which product a path belongs to is decided here, for this check only:
#   template  product/, template/, release/template-repo/
#   cli       cmd/, internal/, assets.go, go.mod, go.sum, release/cli/,
#             scripts/build-cli.sh; Go tests and testdata excluded
# Everything else, VERSION, CLI_VERSION, release/requests/ and .changelog/
# included, belongs to neither, so a release-prep pull request needs nothing.
#
# It reads git and the body it is given, and calls no API, so the workflow
# passes the event's body and the offline check passes a fixture one.
set -eu

fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }

[ "$#" -eq 2 ] || fail 'usage: check-change-files.sh <base> <head>'
base=$1
head=$2
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM

git rev-parse -q --verify "$base^{commit}" >/dev/null || fail "unknown base: $base"
git rev-parse -q --verify "$head^{commit}" >/dev/null || fail "unknown head: $head"

# The Release note section, up to the next "## ", with carriage returns and HTML
# comments removed; `None` in any case, with or without a period, says the
# change is not one users see.
note=$(printf '%s\n' "${PR_BODY:-}" | tr -d '\r' | awk '
  BEGIN { in_section=0; in_comment=0 }
  {
    line=$0; out=""
    while (length(line) > 0) {
      if (in_comment) {
        i=index(line, "-->")
        if (i == 0) { line=""; break }
        line=substr(line, i + 3); in_comment=0
      } else {
        i=index(line, "<!--")
        if (i == 0) { out=out line; line=""; break }
        out=out substr(line, 1, i - 1); line=substr(line, i + 4); in_comment=1
      }
    }
    if (out ~ /^##[[:space:]]/) {
      heading=tolower(out); sub(/^##[[:space:]]+/, "", heading); sub(/[[:space:]]+$/, "", heading)
      in_section=(heading == "release note"); next
    }
    if (in_section) print out
  }' | tr -s ' \t\n' ' ' | sed -e 's/^ //' -e 's/ $//')
none=false
case "$(printf '%s' "$note" | tr '[:upper:]' '[:lower:]')" in
  none|none.) none=true ;;
esac

touched_cli=false
touched_template=false
git diff --name-only "$base...$head" > "$work/paths"
while IFS= read -r path; do
  case "$path" in
    product/*|template/*|release/template-repo/*) touched_template=true ;;
    *_test.go|*/testdata/*) ;;
    cmd/*|internal/*|assets.go|go.mod|go.sum|release/cli/*|scripts/build-cli.sh) touched_cli=true ;;
  esac
done < "$work/paths"

# Every change file the pull request adds or edits must parse, whatever else it
# touches, and together they say which products are described.
covered_cli=false
covered_template=false
git diff --name-only --diff-filter=AM "$base...$head" -- .changelog > "$work/notes"
mkdir "$work/files"
while IFS= read -r path; do
  [ "${path##*/}" != README.md ] || continue
  case "$path" in .changelog/*/*) fail "$path: change files sit directly under .changelog/" ;; esac
  git show "$head:$path" > "$work/files/${path##*/}"
  sh "$source_root/scripts/release-notes.sh" --validate "$work/files/${path##*/}" 2> "$work/err" ||
    fail "$path: $(sed -e 's/^FAIL: //' -e "s|$work/files/||" "$work/err")"
  product=$(tr -d '\r' < "$work/files/${path##*/}" | awk 'NR>1 && /^---/{exit} /^product:/{sub(/^product:[[:space:]]*/, ""); gsub(/["'"'"'[:space:]]/, ""); print}')
  case "$product" in
    cli) covered_cli=true ;;
    template) covered_template=true ;;
    both) covered_cli=true; covered_template=true ;;
  esac
done < "$work/notes"

missing=''
[ "$touched_cli" = false ] || [ "$covered_cli" = true ] || missing="$missing cli"
[ "$touched_template" = false ] || [ "$covered_template" = true ] || missing="$missing template"
if [ -n "$missing" ] && [ "$none" = false ]; then
  products=$(printf '%s' "${missing# }" | sed 's/ / and /')
  fail "this pull request changes $products files without a change file for each. Add .changelog/<slug>.md (see .changelog/README.md), or write None under Release note if users see no change."
fi

if [ -n "$missing" ]; then
  printf 'PASS: the Release note says None\n'
else
  printf 'PASS: every product this pull request changes has a change file\n'
fi
