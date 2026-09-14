#!/bin/sh
# Exercise publication guards and local commit/tag behavior only in fixtures.
# Release assembly itself is tested by check-release.sh; this fixture reuses its
# binary to export a seed instead of compiling six targets a second time.
set -eu
[ "$#" -eq 1 ] || { printf 'usage: check-publication.sh <native-binary>\n' >&2; exit 2; }
cc_fixture_binary=$1
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM
export GIT_CONFIG_GLOBAL="$work/gitconfig"
export GIT_CONFIG_NOSYSTEM=1
printf '[commit]\n  gpgsign = false\n[tag]\n  gpgsign = false\n' > "$GIT_CONFIG_GLOBAL"
export MAINTAINER_NAME=Fixture MAINTAINER_EMAIL=fixture@example.invalid
export CC_FIXTURE_BINARY="$cc_fixture_binary"
mkdir -p "$work/source/scripts" "$work/source/release/requests" "$work/template"
cp "$source_root/scripts/publish-template.sh" "$work/source/scripts/"
printf 'destination_ref: main\n' > "$work/source/release/binding.yaml"
printf '%s\n' --- 'version: 2.0.0-test' --- 'Fixture release notes.' > "$work/source/release/requests/2.0.0-test.md"
cat > "$work/source/scripts/release-artifact.sh" <<'ASSEMBLER'
#!/bin/sh
set -eu
out=$2
version=$3
mkdir -p "$out"
"$CC_FIXTURE_BINARY" template export --path "$out/context-circuit-$version" >/dev/null
printf '%s\n' "${version#v}" > "$out/context-circuit-$version/.context-circuit/VERSION"
(cd "$out/context-circuit-$version" && tar -cf - .) | gzip -n > "$out/context-circuit-$version.tar.gz"
if [ -n "${CC_FIXTURE_DIRTY_TARGET:-}" ]; then printf 'Edit during build\n' > "$CC_FIXTURE_DIRTY_TARGET/keep.local"; fi
ASSEMBLER

git -C "$work/template" init -q -b main
printf '# Old seed\n' > "$work/template/README.md"
printf '*.local\n' > "$work/template/.gitignore"
git -C "$work/template" add .
git -C "$work/template" -c user.name=Fixture -c user.email=fixture@example.invalid commit -q -m 'test: seed publication fixture'
original=$(git -C "$work/template" rev-parse HEAD)
publish() { sh "$work/source/scripts/publish-template.sh" 2.0.0-test "$work/template"; }
# A tracked edit must remain untouched.
printf 'Keep my edit\n' >> "$work/template/README.md"
if publish > "$work/dirty.log" 2>&1; then printf 'FAIL: published dirty checkout\n' >&2; exit 1; fi
[ "$(git -C "$work/template" rev-parse HEAD)" = "$original" ]
[ "$(tail -n 1 "$work/template/README.md")" = 'Keep my edit' ]
git -C "$work/template" restore README.md
# Ignored data must also prevent full-tree replacement.
printf 'Keep local data\n' > "$work/template/keep.local"
if publish > "$work/ignored.log" 2>&1; then printf 'FAIL: published over ignored data\n' >&2; exit 1; fi
[ "$(cat "$work/template/keep.local")" = 'Keep local data' ]
rm "$work/template/keep.local"
# Simulate a user editing the destination while assembly is in progress.
export CC_FIXTURE_DIRTY_TARGET="$work/template"
if publish > "$work/concurrent.log" 2>&1; then printf 'FAIL: replaced an edit made during assembly\n' >&2; exit 1; fi
[ "$(cat "$work/template/keep.local")" = 'Edit during build' ]
[ "$(git -C "$work/template" rev-parse HEAD)" = "$original" ]
unset CC_FIXTURE_DIRTY_TARGET
rm "$work/template/keep.local"
# Only discard this test's generated output so the success case can run.
rm -rf "$work/source/dist"
publish > "$work/published.log"
[ "$(git -C "$work/template" log -1 --format='%B')" = 'chore(release): publish v2.0.0-test' ]
[ "$(cat "$work/template/.context-circuit/VERSION")" = 2.0.0-test ]
git -C "$work/template" rev-parse --verify 'refs/tags/v2.0.0-test^{commit}' >/dev/null
[ -z "$(git -C "$work/template" status --porcelain --untracked-files=all)" ]
[ -f "$work/source/dist/v2.0.0-test/context-circuit-v2.0.0-test.tar.gz" ]
if publish > "$work/repeated.log" 2>&1; then printf 'FAIL: republished an existing tag\n' >&2; exit 1; fi
printf 'PASS: publication preserves user files and creates a clean versioned fixture\n'
