#!/bin/sh
# All fixtures and release output live in a fresh temporary directory.
set -eu
[ "$#" -le 1 ] || { printf 'usage: check-release.sh [new-release-output-dir]\n' >&2; exit 2; }
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$source_root"
for script in scripts/*.sh product/skills/cc-cli/scripts/*.sh; do sh -n "$script"; done
go test ./...
go vet ./...
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM
version=v$(cat VERSION)
release_output=${1:-$work/dist}
[ ! -e "$release_output" ] && [ ! -L "$release_output" ] || { printf 'FAIL: release output already exists\n' >&2; exit 1; }
mkdir -p "$release_output"
sh scripts/build-dist.sh "$version" "$release_output/workspace"
cli_version=$(cat CLI_VERSION)
sh scripts/build-cli.sh "$cli_version" "$release_output/cli"
release_output=$(CDPATH= cd -- "$release_output" && pwd -P)
(
  cd "$release_output/workspace"
  if command -v sha256sum >/dev/null 2>&1; then sha256sum -c SHA256SUMS; else shasum -a 256 -c SHA256SUMS; fi
)
(
  cd "$release_output/cli"
  if command -v sha256sum >/dev/null 2>&1; then sha256sum -c SHA256SUMS; else shasum -a 256 -c SHA256SUMS; fi
)
# Inspect all binary archives, including the Windows targets.
printf '%s\n' README.md LICENSE THIRD_PARTY_NOTICES.txt context-circuit-cli | LC_ALL=C sort > "$work/unix-inventory"
printf '%s\n' README.md LICENSE THIRD_PARTY_NOTICES.txt context-circuit-cli.exe | LC_ALL=C sort > "$work/windows-inventory"
for archive in "$release_output/cli"/*-darwin-*.tar.gz "$release_output/cli"/*-linux-*.tar.gz; do
  tar -tzf "$archive" | LC_ALL=C sort > "$work/inventory"
  cmp "$work/unix-inventory" "$work/inventory"
done
for archive in "$release_output/cli"/*-windows-*.zip; do
  unzip -Z1 "$archive" | LC_ALL=C sort > "$work/inventory"
  cmp "$work/windows-inventory" "$work/inventory"
done
mkdir "$work/native"
tar -xzf "$release_output/cli/context-circuit-cli-v$cli_version-$(go env GOOS)-$(go env GOARCH).tar.gz" -C "$work/native"
cc_binary="$work/native/context-circuit-cli"
[ "$("$cc_binary" version)" = "$cli_version" ]
"$cc_binary" --workspace "$work/workspace" init --name Release --purpose 'Release fixture' --member fixture --member-name Fixture >/dev/null
"$cc_binary" --workspace "$work/workspace" --json check
# The seed tar can also be extracted and initialized, with no source checkout.
mkdir "$work/seed"
tar -xzf "$release_output/workspace/context-circuit-$version.tar.gz" -C "$work/seed"
"$cc_binary" --workspace "$work/seed" init --name Seed --purpose Fixture --member fixture --member-name Fixture >/dev/null
"$cc_binary" --workspace "$work/seed" --json check
# Exercise the shipped installer with the real native CLI archive in isolation.
sh "$work/seed/.agents/skills/cc-cli/scripts/install.sh" --version "$cli_version" \
  --bin-dir "$work/installed" \
  --archive "$release_output/cli/context-circuit-cli-v$cli_version-$(go env GOOS)-$(go env GOARCH).tar.gz" \
  --checksums "$release_output/cli/SHA256SUMS"
"$work/installed/context-circuit-cli" --workspace "$work/seed" --json check
# Building into an existing output must fail without deleting its contents.
mkdir "$work/protected-dist"
printf 'preserve me\n' > "$work/protected-dist/sentinel"
if sh scripts/build-dist.sh "$version" "$work/protected-dist" > "$work/overwrite.log" 2>&1; then
  printf 'FAIL: build replaced existing output\n' >&2
  exit 1
fi
[ "$(cat "$work/protected-dist/sentinel")" = 'preserve me' ]
if sh scripts/build-cli.sh "$cli_version" "$work/protected-dist" > "$work/cli-overwrite.log" 2>&1; then
  printf 'FAIL: CLI build replaced existing output\n' >&2
  exit 1
fi
[ "$(cat "$work/protected-dist/sentinel")" = 'preserve me' ]
sh scripts/check-publication.sh "$cc_binary"
printf 'PASS: separate workspace and CLI packages, native CLI, six targets, checksums, and output protection\n'
