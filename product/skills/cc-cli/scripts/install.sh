#!/bin/sh
# Install an exact CLI release without requiring a language runtime or root.
set -eu
fail() { printf 'CLI installation failed: %s\n' "$1" >&2; exit 1; }
version= bin_dir=${HOME:?}/.local/bin archive= checksums=
token=${CONTEXT_CIRCUIT_TOKEN:-${GH_TOKEN:-${GITHUB_TOKEN:-}}}
while [ "$#" -gt 0 ]; do
  [ "$#" -ge 2 ] || fail 'options require a value'
  case "$1" in
    --version) version=$2 ;; --bin-dir) bin_dir=$2 ;;
    --archive) archive=$2 ;; --checksums) checksums=$2 ;;
    --token) token=$2 ;;
    *) fail "unknown option: $1" ;;
  esac
  shift 2
done
# The token is written to a curl configuration line, which has no escape for a
# quote or backslash. Reject anything outside the character set GitHub issues
# rather than build a malformed request from it.
case "$token" in *[!A-Za-z0-9_-]*) fail 'token contains unexpected characters' ;; esac
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
command_path="$bin_dir/context-circuit-cli"
if [ -e "$command_path" ] || [ -L "$command_path" ]; then
  [ -L "$command_path" ] || fail 'existing command is not managed by this installer; use another --bin-dir'
  case "$(readlink "$command_path")" in "$versions"/*/context-circuit-cli) ;; *) fail 'existing command points outside managed versions' ;; esac
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
  fetch() { curl --fail --location --silent --show-error --proto '=https' --proto-redir '=https' "$@"; }
  # Pass the credential on stdin so it stays out of the process list.
  fetch_auth() { printf 'header = "Authorization: Bearer %s"\n' "$token" | fetch --config - "$@"; }
  if [ -n "$token" ]; then
    # A private repository serves release assets only through the API, by asset
    # id; the public download path answers 404. curl does not carry the
    # Authorization header across the redirect to signed storage.
    api="${CONTEXT_CIRCUIT_API:-https://api.github.com}/repos/kaotypr/context-circuit-source"
    fetch_auth --header 'Accept: application/vnd.github+json' \
      --header 'X-GitHub-Api-Version: 2022-11-28' \
      "$api/releases/tags/cli-v$version" -o "$work/release.json" ||
      fail "cannot read release cli-v$version; confirm it exists and the token grants access"
    # Each asset object begins a line once the JSON is split on braces, so the
    # asset URL read beside a matching name belongs to that asset.
    asset_url() {
      sed 's/" *: */":/g' "$work/release.json" | tr '{' '\n' |
        grep -F "\"name\":\"$1\"" |
        sed -n 's|.*"url":"\([^"]*/releases/assets/[0-9][0-9]*\)".*|\1|p' | head -n 1
    }
    for name in "$package" SHA256SUMS; do
      url=$(asset_url "$name")
      [ -n "$url" ] || fail "release cli-v$version publishes no asset named $name"
      fetch_auth --header 'Accept: application/octet-stream' "$url" -o "$work/$name"
    done
  else
    base="https://github.com/kaotypr/context-circuit-source/releases/download/cli-v$version"
    fetch "$base/$package" -o "$work/$package"
    fetch "$base/SHA256SUMS" -o "$work/SHA256SUMS"
  fi
fi
expected=$(awk -v name="$package" '{file=$2; sub(/^\*/, "", file); sub(/^\.\//, "", file); if(file==name) print $1}' "$work/SHA256SUMS")
[ "${#expected}" -eq 64 ] || fail 'missing or ambiguous checksum entry'
[ "$(hash "$work/$package")" = "$expected" ] || fail 'checksum mismatch'
printf '%s\n' README.md THIRD_PARTY_NOTICES.txt context-circuit-cli | LC_ALL=C sort > "$work/expected"
tar -tzf "$work/$package" | LC_ALL=C sort > "$work/actual"
cmp "$work/expected" "$work/actual" >/dev/null || fail 'unexpected archive contents'
mkdir "$work/package"
tar -xzf "$work/$package" -C "$work/package"
for file in context-circuit-cli README.md THIRD_PARTY_NOTICES.txt; do
  [ -f "$work/package/$file" ] && [ ! -L "$work/package/$file" ] || fail 'archive contains non-regular files'
done
chmod 755 "$work/package/context-circuit-cli"
[ "$("$work/package/context-circuit-cli" version)" = "$version" ] || fail 'executable version does not match the release'
target="$versions/$version-$platform-$arch"
if [ -e "$target" ] || [ -L "$target" ]; then
  [ -d "$target" ] && [ ! -L "$target" ] || fail 'existing version path is not a directory'
  cmp "$target/context-circuit-cli" "$work/package/context-circuit-cli" >/dev/null || fail 'existing version differs; preserved'
else
  mv "$work/package" "$target"
fi
ln -s "$target/context-circuit-cli" "$work/context-circuit-cli"
mv -f "$work/context-circuit-cli" "$command_path"
# Report the version-specific command as well. Workspaces pin their own CLI
# version, so a caller serving several workspaces invokes that path directly
# instead of relying on which version the shared command currently selects.
printf 'version: %s\nplatform: %s/%s\ncommand: %s\nversioned command: %s\nversion store: %s\nPATH directory: %s\n' \
  "$version" "$platform" "$arch" "$command_path" "$target/context-circuit-cli" "$versions" "$bin_dir"
