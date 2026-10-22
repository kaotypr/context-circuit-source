#!/bin/sh
# Print the version a product builds and publishes as. A public checkout prints
# its version file unchanged. A downstream variant supplies a label in
# CC_VERSION_SUFFIX and a counter in release/variant.env, and gets
# <base>-<label>.<counter>. Anything in between fails rather than guessing.
#
# Every release script that needs a version asks this one, so the version a
# product is built, checked, and published as is decided in one place. One line
# on stdout, without a v; on any problem FAIL: on stderr, exit 1, empty stdout.
set -eu
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
[ "$#" -eq 1 ] || fail 'usage: release-version.sh <cli|template>'
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
case "$1" in
  cli) file=CLI_VERSION ;;
  template) file=VERSION ;;
  *) fail 'usage: release-version.sh <cli|template>' ;;
esac
[ -f "$source_root/$file" ] || fail "missing $file"
base=$(tr -d '[:space:]' < "$source_root/$file")
label=${CC_VERSION_SUFFIX:-}
variant="$source_root/release/variant.env"
if [ -z "$label" ]; then
  [ ! -e "$variant" ] || fail "release/variant.env exists but CC_VERSION_SUFFIX is not set; set it to this fork's label"
  printf '%s\n' "$base"
  exit 0
fi
# The label is checked before the file, so a bad label names the rule it breaks
# whether or not this ref carries release/variant.env.
# Read by a Go test; keep each on one line.
reserved_labels='alpha beta rc dev test pre preview snapshot'
label_pattern='^[a-z]{1,20}$'
printf '%s\n' "$label" | grep -Eq "$label_pattern" || fail "invalid label: $label (1-20 lowercase letters)"
for word in $reserved_labels; do [ "$label" != "$word" ] || fail "reserved label: $label"; done
[ -f "$variant" ] || fail 'CC_VERSION_SUFFIX is set but release/variant.env is missing on this ref'
printf '%s\n' "$base" | grep -Eq '^2\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$' || fail "a variant builds only on a stable base: $base"
# Parsed, never sourced: the file is data, and a stray \r fails the counter loudly.
counter=
counter_seen=
while IFS= read -r line || [ -n "$line" ]; do
  case "$line" in ''|'#'*) continue ;; esac
  key=${line%%=*}
  value=${line#*=}
  case "$key" in
    counter) [ -z "$counter_seen" ] || fail 'counter is set twice in release/variant.env'; counter_seen=1; counter=$value ;;
    upstream_*) ;;
    suffix) fail 'the label lives in CC_VERSION_SUFFIX, not in release/variant.env' ;;
    *) fail "unknown key in release/variant.env: $key" ;;
  esac
done < "$variant"
[ -n "$counter_seen" ] || fail 'counter is missing from release/variant.env'
printf '%s\n' "$counter" | grep -Eq '^[1-9][0-9]*$' || fail "invalid counter in release/variant.env: '$counter'"
printf '%s-%s.%s\n' "$base" "$label" "$counter"
