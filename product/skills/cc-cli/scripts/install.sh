#!/bin/sh
# Install an exact CLI release without requiring a language runtime or root.
set -eu
fail() { printf 'CLI installation failed: %s\n' "$1" >&2; exit 1; }
version= bin_dir=${HOME:?}/.local/bin archive= checksums=
while [ "$#" -gt 0 ]; do
  [ "$#" -ge 2 ] || fail 'options require a value'
  case "$1" in
    --version) version=$2 ;; --bin-dir) bin_dir=$2 ;;
    --archive) archive=$2 ;; --checksums) checksums=$2 ;;
    *) fail "unknown option: $1" ;;
  esac
  shift 2
done
case "$version" in 2.*) ;; *) fail 'supply an exact compatible v2 CLI version without a v prefix' ;; esac
case "$version" in *[!A-Za-z0-9.+-]*|*..*) fail 'invalid version' ;; esac
case "$(uname -s)" in Darwin) platform=darwin ;; Linux) platform=linux ;; *) fail 'use install.ps1 on native Windows' ;; esac
case "$(uname -m)" in arm64|aarch64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; *) fail 'unsupported architecture' ;; esac
package="context-circuit-cli-v$version-$platform-$arch.tar.gz"
mkdir -p "$bin_dir"
bin_dir=$(CDPATH= cd -- "$bin_dir" && pwd -P)
versions="$bin_dir/.context-circuit-versions"
[ ! -L "$versions" ] || fail 'version directory must not be a symlink'
mkdir -p "$versions"
lock="$versions/install.lock"
mkdir "$lock" 2>/dev/null || fail "another install is running (or inspect a stale lock): $lock"
work=$(mktemp -d "$versions/.install-XXXXXX")
cleanup() { rm -rf "$work"; rmdir "$lock"; }
trap cleanup EXIT HUP INT TERM
command_path="$bin_dir/context-circuit"
if [ -e "$command_path" ] || [ -L "$command_path" ]; then
  [ -L "$command_path" ] || fail 'existing command is not managed by this installer; use another --bin-dir'
  case "$(readlink "$command_path")" in "$versions"/*/context-circuit) ;; *) fail 'existing command points outside managed versions' ;; esac
fi
hash() {
  if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d ' ' -f1
  elif command -v shasum >/dev/null 2>&1; then shasum -a 256 "$1" | cut -d ' ' -f1
  else fail 'sha256sum or shasum is required'; fi
}
if [ -n "$archive" ] || [ -n "$checksums" ]; then
  [ -n "$archive" ] && [ -n "$checksums" ] || fail 'offline installation needs both archive and checksums'
  cp "$archive" "$work/$package"
  cp "$checksums" "$work/SHA256SUMS"
else
  command -v curl >/dev/null 2>&1 || fail 'curl is required to download the CLI'
  base="https://github.com/kaotypr/context-circuit-source/releases/download/cli-v$version"
  curl --fail --location --silent --show-error --proto '=https' --proto-redir '=https' "$base/$package" -o "$work/$package"
  curl --fail --location --silent --show-error --proto '=https' --proto-redir '=https' "$base/SHA256SUMS" -o "$work/SHA256SUMS"
fi
expected=$(awk -v name="$package" '{file=$2; sub(/^\*/, "", file); sub(/^\.\//, "", file); if(file==name) print $1}' "$work/SHA256SUMS")
[ "${#expected}" -eq 64 ] || fail 'missing or ambiguous checksum entry'
[ "$(hash "$work/$package")" = "$expected" ] || fail 'checksum mismatch'
printf '%s\n' README.md THIRD_PARTY_NOTICES.txt context-circuit | LC_ALL=C sort > "$work/expected"
tar -tzf "$work/$package" | LC_ALL=C sort > "$work/actual"
cmp "$work/expected" "$work/actual" >/dev/null || fail 'unexpected archive contents'
mkdir "$work/package"
tar -xzf "$work/$package" -C "$work/package"
for file in context-circuit README.md THIRD_PARTY_NOTICES.txt; do
  [ -f "$work/package/$file" ] && [ ! -L "$work/package/$file" ] || fail 'archive contains non-regular files'
done
chmod 755 "$work/package/context-circuit"
[ "$("$work/package/context-circuit" version)" = "$version" ] || fail 'executable version does not match the release'
target="$versions/$version-$platform-$arch"
if [ -e "$target" ] || [ -L "$target" ]; then
  [ -d "$target" ] && [ ! -L "$target" ] || fail 'existing version path is not a directory'
  cmp "$target/context-circuit" "$work/package/context-circuit" >/dev/null || fail 'existing version differs; preserved'
else
  mv "$work/package" "$target"
fi
ln -s "$target/context-circuit" "$work/context-circuit"
mv -f "$work/context-circuit" "$command_path"
printf 'version: %s\nplatform: %s/%s\ncommand: %s\nPATH directory: %s\n' "$version" "$platform" "$arch" "$command_path" "$bin_dir"
