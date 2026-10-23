#!/bin/sh
# Exercise release-version.sh's guard table, and the resolver under a variant
# label, in fixtures only. Every case sets CC_VERSION_SUFFIX itself, so a fork's
# CI variable never decides what a case means, and the resolver case sets FORGE
# itself because GitLab sets it for every job.
set -eu
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }

# The two bases differ so a case cannot pass by reading the wrong file.
tree="$work/source"
mkdir -p "$tree/scripts" "$tree/release"
cp "$source_root/scripts/release-version.sh" "$tree/scripts/"
printf '2.5.0\n' > "$tree/CLI_VERSION"
printf '2.5.3\n' > "$tree/VERSION"
variant="$tree/release/variant.env"

# expect <suffix> <product> <version>
expect() {
  got=$(CC_VERSION_SUFFIX=$1 sh "$tree/scripts/release-version.sh" "$2" 2>"$work/stderr") ||
    fail "release-version.sh $2 with CC_VERSION_SUFFIX='$1' failed: $(cat "$work/stderr")"
  [ "$got" = "$3" ] || fail "release-version.sh $2 with CC_VERSION_SUFFIX='$1' printed '$got', want '$3'"
}
# refuse <case> <suffix> [product...]: fails, prints FAIL: on stderr and nothing on stdout.
refuse() {
  case_name=$1
  suffix=$2
  shift 2
  if got=$(CC_VERSION_SUFFIX=$suffix sh "$tree/scripts/release-version.sh" "$@" 2>"$work/stderr"); then
    fail "accepted $case_name: printed '$got'"
  fi
  [ -z "$got" ] || fail "$case_name printed '$got' on stdout while failing"
  grep -q '^FAIL: ' "$work/stderr" || fail "$case_name did not report FAIL: on stderr"
}
# with_variant <line>...: replaces release/variant.env with exactly these lines.
with_variant() { printf '%s\n' "$@" > "$variant"; }

# Public build: no variable, no file.
rm -f "$variant"
expect '' cli 2.5.0
expect '' template 2.5.3
printf '2.6.0-rc.1\n' > "$tree/CLI_VERSION"
expect '' cli 2.6.0-rc.1
printf '2.5.0\n' > "$tree/CLI_VERSION"

# Variant build, with the lines the parser must skip or ignore.
with_variant '# acme variant' '' 'upstream_tag=cli-v2.5.0' 'counter=3'
expect acme cli 2.5.0-acme.3
expect acme template 2.5.3-acme.3

# The states in between.
with_variant 'counter=3'
refuse 'a variant file without a label' '' cli
rm -f "$variant"
refuse 'a label without a variant file' acme cli
# A bad label names the rule it breaks, file or no file.
refuse 'an invalid label without a variant file' Acme cli
grep -q 'invalid label' "$work/stderr" || fail "a bad label without a variant file did not name the label rule: $(cat "$work/stderr")"

with_variant 'counter=3'
for label in Acme acme1 ac-me ac.me ac+me ac_me abcdefghijklmnopqrstu; do
  refuse "label $label" "$label" cli
  grep -q 'invalid label' "$work/stderr" || fail "label $label did not name the label rule: $(cat "$work/stderr")"
done
for label in alpha beta rc dev test pre preview snapshot; do
  refuse "reserved label $label" "$label" cli
  grep -q 'reserved label' "$work/stderr" || fail "label $label did not name the reserved rule: $(cat "$work/stderr")"
done

with_variant '# no counter' 'upstream_tag=cli-v2.5.0'
refuse 'a missing counter' acme cli
for counter in 0 01 1.0 x ''; do
  with_variant "counter=$counter"
  refuse "counter '$counter'" acme cli
done
with_variant 'counter=1' 'counter=2'
refuse 'a counter set twice' acme cli
printf 'counter=1\r\n' > "$variant"
refuse 'a counter with a carriage return' acme cli
with_variant 'counter=1' 'suffix=acme'
refuse 'a suffix key' acme cli
with_variant 'counter=1' 'foo=1'
refuse 'an unknown key' acme cli

with_variant 'counter=1'
printf '2.6.0-rc.1\n' > "$tree/CLI_VERSION"
refuse 'a pre-release base with a label' acme cli
printf '2.5.0\n' > "$tree/CLI_VERSION"

rm -f "$variant"
refuse 'an unknown product' '' workspace
refuse 'no product' ''
refuse 'two products' '' cli template
rm "$tree/VERSION"
refuse 'a missing version file' '' template
printf '2.5.3\n' > "$tree/VERSION"

# The resolver publishes the composed version, needs its request, and refuses an
# explicit version a variant's files do not compose. The fake gh answers "not
# published", so the request lookup is what decides.
repo="$work/resolver"
mkdir -p "$repo/scripts" "$repo/release/requests/cli" "$work/bin"
git init -q "$repo"
cp "$source_root/scripts/resolve-release-version.sh" "$source_root/scripts/release-version.sh" "$repo/scripts/"
printf '2.5.0\n' > "$repo/CLI_VERSION"
printf 'counter=3\n' > "$repo/release/variant.env"
printf '%s\n' --- 'version: 2.5.0-acme.3' --- 'Fixture notes.' > "$repo/release/requests/cli/2.5.0-acme.3.md"
printf '#!/bin/sh\nexit 1\n' > "$work/bin/gh"
chmod +x "$work/bin/gh"
resolve() {
  (
    cd "$repo"
    unset GITHUB_OUTPUT INPUT_VERSION
    if [ "$#" -gt 0 ]; then INPUT_VERSION=$1; export INPUT_VERSION; fi
    PATH="$work/bin:$PATH" FORGE=github CLI_REPOSITORY=fixture/cli \
      RESOLVE_OUTPUT="$work/resolved" CC_VERSION_SUFFIX=acme \
      sh scripts/resolve-release-version.sh cli
  )
}
: > "$work/resolved"
resolve > "$work/resolve.log" 2>&1 || fail "resolver refused the composed version: $(cat "$work/resolve.log")"
grep -qx 'version=2.5.0-acme.3' "$work/resolved" || fail "resolver did not resolve 2.5.0-acme.3: $(cat "$work/resolved")"
grep -qx 'publish=true' "$work/resolved" || fail "resolver did not publish: $(cat "$work/resolved")"
mv "$repo/release/requests/cli/2.5.0-acme.3.md" "$work/request.held"
if resolve > "$work/resolve.log" 2>&1; then fail 'resolver published without a request for the composed version'; fi
mv "$work/request.held" "$repo/release/requests/cli/2.5.0-acme.3.md"
printf '%s\n' --- 'version: 2.5.0' --- 'Fixture notes.' > "$repo/release/requests/cli/2.5.0.md"
if resolve 2.5.0 > "$work/resolve.log" 2>&1; then fail 'resolver accepted an explicit version under a variant label'; fi
grep -q "is not this variant's 2.5.0-acme.3" "$work/resolve.log" || fail "resolver refused the explicit version for another reason: $(cat "$work/resolve.log")"

printf 'PASS: release-version.sh guard table, and the resolver under a variant label\n'
