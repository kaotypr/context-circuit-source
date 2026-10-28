#!/bin/sh
# Draft a product's release request from the change files under .changelog/.
#
# Each pull request that changes something users see adds one change file, so a
# note is written with its change, reviewed in the same diff, and carried along
# when a fix is cherry-picked onto a release branch. At release time this script
# turns the files on the release branch into release/requests/<product>/<version>.md,
# and the maintainer edits that draft in the release-prep pull request. The
# request stays what Publish reads, so the published text is in git and was
# reviewed before it shipped.
#
# No tag range is needed. A stable release removes the files it used, so the
# files left on the branch are exactly what no stable release has shipped yet,
# and a candidate's draft covers everything since the last stable release.
#
# Usage:
#   release-notes.sh <cli|template> <version> [--stdout] [--force]
#   release-notes.sh --validate <file>...
#
# --stdout prints the draft and writes or removes nothing. --force overwrites an
# existing request. --validate parses each change file and exits nonzero naming
# the first one that does not parse; check-change-files.sh uses it, so the format
# is read in one place.
#
# It reads only the repository: the pull request link and the fallback
# contributor come from the subject GitHub gives a merge commit, which merge
# commits only guarantees. Nothing here calls an API or needs a token, so both
# forges and the offline check run it the same way.
set -eu

fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
note() { printf 'note: %s\n' "$1" >&2; }

# parse <file>: sets product, kind, upgrade, authors, and writes the body to
# "$work/body". Fails naming the file when it does not parse.
parse() {
  file=$1
  [ -f "$file" ] || fail "$file: no such change file"
  case "${file##*/}" in
    *[!A-Za-z0-9._-]*|.*) fail "$file: a change file name uses only letters, digits, '.', '_' and '-'" ;;
    *.md) ;;
    *) fail "$file: a change file ends in .md" ;;
  esac
  tr -d '\r' < "$file" > "$work/file"
  [ "$(sed -n 1p "$work/file")" = --- ] || fail "$file: does not open with a --- frontmatter line"
  awk 'NR>1 && /^---[[:space:]]*$/{print NR; exit}' "$work/file" > "$work/close"
  close=$(cat "$work/close")
  [ -n "$close" ] || fail "$file: frontmatter is not closed with ---"
  sed -n "2,$((close - 1))p" "$work/file" > "$work/front"
  sed -n "$((close + 1)),\$p" "$work/file" | awk '
    { lines[NR]=$0 }
    END {
      first=0; last=0
      for (i=1; i<=NR; i++) if (lines[i] ~ /[^[:space:]]/) { if (!first) first=i; last=i }
      if (first) for (i=first; i<=last; i++) print lines[i]
    }' > "$work/body"
  [ -s "$work/body" ] || fail "$file: has no entry text after its frontmatter"
  product='' kind='' upgrade='' authors='' prs=''
  while IFS= read -r line || [ -n "$line" ]; do
    case "$line" in ''|'#'*) continue ;; esac
    key=${line%%:*}
    [ "$key" != "$line" ] || fail "$file: frontmatter line is not key: value: $line"
    value=$(printf '%s' "${line#*:}" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' \
      -e 's/^"\(.*\)"$/\1/' -e "s/^'\\(.*\\)'\$/\\1/")
    case "$key" in
      product) product=$value ;;
      kind) kind=$value ;;
      upgrade) upgrade=$value ;;
      authors) authors=$value ;;
      # Spaces between numbers collapse, and a value of spaces alone is no
      # value, so it cannot replace the link a merge would give.
      pr) prs=$(printf '%s' "$value" | tr -s ' \t' '  ' | sed -e 's/^ //' -e 's/ $//') ;;
      *) fail "$file: unknown frontmatter key: $key (expected product, kind, upgrade, authors, pr)" ;;
    esac
  done < "$work/front"
  case "$product" in
    cli|template|both) ;;
    '') fail "$file: no product (expected cli, template, or both)" ;;
    *) fail "$file: unknown product: $product (expected cli, template, or both)" ;;
  esac
  case "$kind" in
    added|fixed|changed|breaking) ;;
    '') fail "$file: no kind (expected added, fixed, changed, or breaking)" ;;
    *) fail "$file: unknown kind: $kind (expected added, fixed, changed, or breaking)" ;;
  esac
  for n in $prs; do
    case "$n" in
      *[!0-9]*|0*) fail "$file: pr is not a list of pull request numbers: $prs" ;;
    esac
  done
  for handle in $authors; do
    case "$handle" in
      @*) ;;
      *) fail "$file: author is not a GitHub handle starting with @: $handle" ;;
    esac
  done
}

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM

if [ "${1:-}" = --validate ]; then
  shift
  [ "$#" -gt 0 ] || fail 'usage: release-notes.sh --validate <file>...'
  for f in "$@"; do parse "$f"; done
  exit 0
fi

target=${1:-}
version=${2:-}
case "$target" in
  cli|template) ;;
  *) fail 'usage: release-notes.sh <cli|template> <version> [--stdout] [--force]' ;;
esac
case "$version" in ''|v*|*[!A-Za-z0-9.+-]*|.*|-*) fail "invalid version: $version" ;; esac
shift 2
to_stdout=false
force=false
for arg in "$@"; do
  case "$arg" in
    --stdout) to_stdout=true ;;
    --force) force=true ;;
    *) fail "unknown option: $arg" ;;
  esac
done

source_root=$(git rev-parse --show-toplevel 2>/dev/null) || fail 'not a git source checkout'
cd "$source_root"

# A candidate is a version whose pre-release label is one of the reserved words
# release-version.sh refuses as a variant label. A variant's own -<label>.N is a
# release of that variant and is treated as stable.
candidate=false
case "$version" in
  *-*)
    label=${version#*-}
    label=${label%%.*}
    for word in alpha beta rc dev test pre preview snapshot; do
      [ "$label" != "$word" ] || candidate=true
    done
    ;;
esac

request="release/requests/$target/$version.md"
if [ "$to_stdout" = false ] && [ -e "$request" ] && [ "$force" = false ]; then
  fail "$request already exists; edit it, or pass --force to draft over it"
fi

# One record per used file: "<sort key>\t<kind>\t<file>\t<link>\t<owner>".
: > "$work/entries"
: > "$work/used"
for f in .changelog/*.md; do
  [ -f "$f" ] || continue
  [ "${f##*/}" != README.md ] || continue
  parse "$f"
  case "$product" in "$target"|both) ;; *) continue ;; esac
  subject=$(git log --first-parent --diff-filter=A --format=%s -- "$f" |
    grep -E '^Merge pull request #[0-9]+ from ' | head -n 1 || true)
  number='' owner=''
  if [ -n "$subject" ]; then
    number=$(printf '%s\n' "$subject" | sed -E 's/^Merge pull request #([0-9]+) from .*/\1/')
    owner=$(printf '%s\n' "$subject" | sed -E 's/^Merge pull request #[0-9]+ from ([^/ ]+).*/\1/')
  elif [ -z "$prs" ]; then
    note "$f reached this branch without a pull request merge, so its entry has no link"
  fi
  # A file that names its pull requests links those: a change made by several
  # pull requests, an entry added after the fact, or a cherry-pick pointing at
  # the original fix. The merge still supplies the fallback contributor.
  [ -z "$prs" ] || number=$prs
  link='' key=999999999
  for n in $number; do
    case "$target" in
      cli) one="#$n" ;;
      template) one="kaotypr/context-circuit-source#$n" ;;
    esac
    if [ -z "$link" ]; then link=$one; key=$(printf '%09d' "$n"); else link="$link, $one"; fi
  done
  {
    tr '\n' ' ' < "$work/body" | sed -e 's/[[:space:]]*$//'
    [ -z "$link" ] || printf ' (%s)' "$link"
    printf '\n'
  } > "$work/text.$kind.$key.${f##*/}"
  if [ -n "$upgrade" ]; then
    {
      printf '%s' "$upgrade"
      [ -z "$link" ] || printf ' (%s)' "$link"
      printf '\n'
    } > "$work/text.upgrade.$key.${f##*/}"
  fi
  if [ -n "$authors" ]; then
    for handle in $authors; do printf '%s\n' "$handle"; done >> "$work/authors"
  elif [ -n "$owner" ]; then
    printf '@%s\n' "$owner" >> "$work/authors"
  fi
  printf '%s\t%s\n' "$f" "$product" >> "$work/used"
done

section() {
  heading=$1
  kind_key=$2
  set -- "$work"/text."$kind_key".*
  [ -e "$1" ] || return 0
  printf '\n## %s\n' "$heading"
  for entry in $(printf '%s\n' "$@" | LC_ALL=C sort); do
    printf '\n'
    cat "$entry"
  done
}

{
  printf '%s\n' ---
  if [ "$candidate" = true ]; then
    printf 'level: prerelease\n'
  else
    printf 'level: stable\n'
  fi
  printf 'reason: One sentence for the person cutting the release; no script reads it.\n'
  if [ "$target" = template ]; then
    if [ "$candidate" = true ]; then
      printf 'changelog: skip\n'
    else
      printf 'changelog: keep\n'
    fi
  fi
  printf '%s\n\n' ---
  case "$target" in
    cli) printf '# Context Circuit CLI %s\n' "$version" ;;
    template) printf '# v%s\n' "$version" ;;
  esac
  printf '\n## Upgrade notes\n'
  set -- "$work"/text.upgrade.*
  if [ -e "$1" ]; then
    for entry in $(printf '%s\n' "$@" | LC_ALL=C sort); do
      printf '\n'
      cat "$entry"
    done
  else
    printf '\nNone.\n'
  fi
  section Breaking breaking
  section Added added
  section Fixed fixed
  section Changed changed
  if [ -s "$work/authors" ]; then
    contributors=$(grep -v '\[bot\]$' "$work/authors" | awk '!seen[$0]++' | paste -sd, - | sed 's/,/, /g' || true)
    [ -z "$contributors" ] || printf '\n## Contributors\n\n%s\n' "$contributors"
  fi
} > "$work/draft"

if [ "$to_stdout" = true ]; then
  cat "$work/draft"
  exit 0
fi

mkdir -p "release/requests/$target"
cp "$work/draft" "$request"
printf 'drafted: %s\n' "$request"

[ "$candidate" = false ] || exit 0
# A stable release removes what it shipped. A file naming both products is
# rewritten to name the other one, so that product's next stable draft still
# finds it.
while IFS="$(printf '\t')" read -r f product; do
  if [ "$product" = both ]; then
    case "$target" in cli) other=template ;; template) other=cli ;; esac
    awk -v other="$other" '
      NR==1 { print; next }
      !closed && /^---[[:space:]]*\r?$/ { closed=1; print; next }
      !closed && /^product:/ { print "product: " other; next }
      { print }' "$f" > "$work/rewrite"
    cp "$work/rewrite" "$f"
    printf 'kept for %s: %s\n' "$other" "$f"
  else
    git rm -q -- "$f" 2>/dev/null || rm -f -- "$f"
    printf 'removed: %s\n' "$f"
  fi
done < "$work/used"
