#!/bin/sh
# Install an exact CLI release without requiring a language runtime or root.
set -eu
fail() { printf 'CLI installation failed: %s\n' "$1" >&2; exit 1; }
version= bin_dir=${HOME:?}/.local/bin archive= checksums=
token=${CONTEXT_CIRCUIT_TOKEN:-${GH_TOKEN:-${GITHUB_TOKEN:-}}}
# An organization that mirrors CLI releases into its own GitLab project names it
# here. Unset, the installer reads the product's own GitHub releases.
gitlab_url=${CONTEXT_CIRCUIT_GITLAB_URL:-}
usage() {
  cat <<'USAGE'
Usage: install.sh --version <2.x.y> [options]

  --version     exact CLI version to install, without a v prefix (required)
  --bin-dir     user-writable command directory (default: $HOME/.local/bin)
  --token       registry credential; also read from CONTEXT_CIRCUIT_TOKEN,
                GH_TOKEN, or GITHUB_TOKEN
  --gitlab-url  GitLab project mirroring the CLI releases; also read from
                CONTEXT_CIRCUIT_GITLAB_URL
  --archive     install this local package instead of downloading one
  --checksums   the SHA256SUMS it is verified against, required with --archive
USAGE
}
while [ "$#" -gt 0 ]; do
  # Validate the name before asking for its value, so an option nobody defined
  # is reported as unknown rather than as one missing an argument, and a bare
  # `--version` names itself rather than reporting that options in general take
  # values. Asking how to run this is answered, not refused.
  case "$1" in
    -h|--help) usage; exit 0 ;;
    --version|--bin-dir|--archive|--checksums|--token|--gitlab-url) ;;
    *) fail "unknown option: $1" ;;
  esac
  [ "$#" -ge 2 ] || fail "$1 requires a value"
  case "$1" in
    --version) version=$2 ;; --bin-dir) bin_dir=$2 ;;
    --archive) archive=$2 ;; --checksums) checksums=$2 ;;
    --token) token=$2 ;; --gitlab-url) gitlab_url=$2 ;;
  esac
  shift 2
done
# The token is written to a curl configuration line, which has no escape for a
# quote or backslash. Reject anything outside the character set GitHub and GitLab
# issue rather than build a malformed request from it.
case "$token" in *[!A-Za-z0-9_-]*) fail 'token contains unexpected characters' ;; esac
# Derive the API base and the project path from the mirror URL, so this shipped
# script carries no organization's own address. The credential is presented to
# whichever registry is selected below and never to both.
if [ -n "$gitlab_url" ]; then
  case "$gitlab_url" in https://*) ;; *) fail 'mirror URL must begin with https://' ;; esac
  rest=${gitlab_url%/}
  rest=${rest%.git}
  rest=${rest#https://}
  gitlab_host=${rest%%/*}
  gitlab_project=${rest#*/}
  [ -n "$gitlab_host" ] && [ -n "$gitlab_project" ] && [ "$gitlab_project" != "$rest" ] ||
    fail 'mirror URL needs a host and a project path'
  case "$gitlab_host$gitlab_project" in *[!A-Za-z0-9._:/-]*) fail 'mirror URL contains unexpected characters' ;; esac
  gitlab_project=$(printf '%s' "$gitlab_project" | sed 's|/|%2F|g')
fi
case "$version" in 2.*) ;; *) fail 'supply an exact compatible v2 CLI version without a v prefix' ;; esac
case "$version" in *[!A-Za-z0-9.+-]*|*..*) fail 'invalid version' ;; esac
case "$(uname -s)" in Darwin) platform=darwin ;; Linux) platform=linux ;; *) fail 'use install.ps1 on native Windows' ;; esac
case "$(uname -m)" in arm64|aarch64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; *) fail 'unsupported architecture' ;; esac
# The release is tagged `cli-v<version>`; the assets under it keep the
# executable's own name, so a downloaded archive still says what it holds.
release_tag="cli-v$version"
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
  fetch_private() { printf 'header = "PRIVATE-TOKEN: %s"\n' "$token" | fetch --config - "$@"; }
  if [ -n "$gitlab_url" ]; then
    # A generic package is addressed by version and file name, so a mirror needs
    # no release lookup and no asset ids.
    base="https://$gitlab_host/api/v4/projects/$gitlab_project/packages/generic/context-circuit-cli/$version"
    fetch_mirror() { if [ -n "$token" ]; then fetch_private "$@"; else fetch "$@"; fi; }
    for name in "$package" SHA256SUMS; do
      fetch_mirror "$base/$name" -o "$work/$name" ||
        fail "cannot download $name from the mirror; confirm the version is published there and the token grants access"
    done
  elif [ -n "$token" ]; then
    # A private repository serves release assets only through the API, by asset
    # id; the public download path answers 404. curl does not carry the
    # Authorization header across the redirect to signed storage.
    api="${CONTEXT_CIRCUIT_API:-https://api.github.com}/repos/kaotypr/context-circuit-source"
    fetch_auth --header 'Accept: application/vnd.github+json' \
      --header 'X-GitHub-Api-Version: 2022-11-28' \
      "$api/releases/tags/$release_tag" -o "$work/release.json" ||
      fail "cannot read release $release_tag; confirm it exists and the token grants access"
    # Read each brace-delimited object as one record, so the asset URL and the
    # name beside it are matched together whether the API pretty-prints the
    # payload or returns it compact. An asset's own URL precedes the nested
    # uploader object that ends the record.
    asset_url() {
      awk -v want="$1" '
        BEGIN { RS = "{" }
        {
          record = $0
          gsub(/[ \t\r\n]+/, "", record)
          if (index(record, "\"name\":\"" want "\"") == 0) next
          if (match(record, /"url":"[^"]*\/releases\/assets\/[0-9]+"/) == 0) next
          print substr(record, RSTART + 7, RLENGTH - 8)
          exit
        }
      ' "$work/release.json"
    }
    for name in "$package" SHA256SUMS; do
      url=$(asset_url "$name")
      [ -n "$url" ] || fail "release $release_tag publishes no asset named $name"
      fetch_auth --header 'Accept: application/octet-stream' "$url" -o "$work/$name"
    done
  else
    base="https://github.com/kaotypr/context-circuit-source/releases/download/$release_tag"
    fetch "$base/$package" -o "$work/$package"
    fetch "$base/SHA256SUMS" -o "$work/SHA256SUMS"
  fi
fi
expected=$(awk -v name="$package" '{file=$2; sub(/^\*/, "", file); sub(/^\.\//, "", file); if(file==name) print $1}' "$work/SHA256SUMS")
[ "${#expected}" -eq 64 ] || fail 'missing or ambiguous checksum entry'
[ "$(hash "$work/$package")" = "$expected" ] || fail 'checksum mismatch'
printf '%s\n' README.md LICENSE THIRD_PARTY_NOTICES.txt context-circuit-cli | LC_ALL=C sort > "$work/expected"
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
