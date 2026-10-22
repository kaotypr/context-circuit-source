#!/bin/sh
# Install an exact CLI release without requiring a language runtime or root.
set -eu
fail() { printf 'CLI installation failed: %s\n' "$1" >&2; exit 1; }
version= bin_dir=${CONTEXT_CIRCUIT_BIN_DIR:-${HOME:?}/.local/bin} archive= checksums=
token=${CONTEXT_CIRCUIT_TOKEN:-${GH_TOKEN:-${GITHUB_TOKEN:-}}}
source=${CONTEXT_CIRCUIT_SOURCE:-github}
repo=${CONTEXT_CIRCUIT_REPO:-}
api=${CONTEXT_CIRCUIT_API:-}
usage() {
  cat <<'USAGE'
Usage: install.sh --version <2.x.y> [options]

  --version     exact CLI version to install, without a v prefix (required)
  --bin-dir     user-writable command directory (default: CONTEXT_CIRCUIT_BIN_DIR,
                else $HOME/.local/bin)
  --source      github (default) or gitlab; also read from CONTEXT_CIRCUIT_SOURCE
  --repo        owner/repo (github) or group/project (gitlab); also read from
                CONTEXT_CIRCUIT_REPO. Defaults to kaotypr/context-circuit-source
                for github; required for gitlab.
  --api         API base URL; also read from CONTEXT_CIRCUIT_API. Defaults to
                https://api.github.com for github; required for gitlab (e.g.
                https://gitlab.example.com/api/v4).
  --token       registry credential; also read from CONTEXT_CIRCUIT_TOKEN,
                GH_TOKEN, or GITHUB_TOKEN. Required for gitlab.
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
    --version|--bin-dir|--archive|--checksums|--token|--source|--repo|--api) ;;
    *) fail "unknown option: $1" ;;
  esac
  [ "$#" -ge 2 ] || fail "$1 requires a value"
  case "$1" in
    --version) version=$2 ;; --bin-dir) bin_dir=$2 ;;
    --archive) archive=$2 ;; --checksums) checksums=$2 ;;
    --token) token=$2 ;;
    --source) source=$2 ;; --repo) repo=$2 ;; --api) api=$2 ;;
  esac
  shift 2
done
# The token is written to a curl configuration line, which has no escape for a
# quote or backslash. Reject anything outside the character set GitHub and
# GitLab issue rather than build a malformed request from it.
case "$token" in *[!A-Za-z0-9_-]*) fail 'token contains unexpected characters' ;; esac
case "$source" in
  github) repo=${repo:-kaotypr/context-circuit-source}; api=${api:-https://api.github.com} ;;
  gitlab)
    [ -n "$repo" ] || fail 'gitlab installs require --repo (the group/project path)'
    [ -n "$api" ] || fail 'gitlab installs require --api (the instance API base, e.g. https://gitlab.example.com/api/v4)'
    [ -n "$token" ] || fail 'gitlab installs require --token (a personal or project access token); this registry has no public download path'
    ;;
  *) fail 'unknown --source: expected github or gitlab' ;;
esac
case "$version" in 2.*) ;; *) fail 'supply an exact compatible v2 CLI version without a v prefix' ;; esac
case "$version" in *[!A-Za-z0-9.+-]*|*..*) fail 'invalid version' ;; esac
case "$(uname -s)" in Darwin) platform=darwin ;; Linux) platform=linux ;; *) fail 'use install.ps1 on native Windows' ;; esac
case "$(uname -m)" in arm64|aarch64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; *) fail 'unsupported architecture' ;; esac
# A variant version X.Y.Z-<label>.N installs under its own command name, so it
# never repoints the public command. The binary in the version store keeps the
# public name. The same rule is in build-cli.sh and install.ps1.
# Read by a Go test; keep each on one line.
reserved_labels='alpha beta rc dev test pre preview snapshot'
label_pattern='^[a-z]{1,20}$'
name=context-circuit-cli
base=${version%%-*} rest=${version#*-}
if [ "$rest" != "$version" ] && printf '%s\n' "$base" | grep -Eq '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$'; then
  label=${rest%%.*} counter=${rest#*.}
  if [ "$label.$counter" = "$rest" ] && printf '%s\n' "$label" | grep -Eq "$label_pattern" &&
    printf '%s\n' "$counter" | grep -Eq '^[1-9][0-9]*$'; then
    case " $reserved_labels " in *" $label "*) ;; *) name=$label-context-circuit-cli ;; esac
  fi
fi
# The release is tagged `cli-v<version>`; the assets under it carry the
# command's name, so a downloaded archive still says what it holds.
release_tag="cli-v$version"
package="$name-v$version-$platform-$arch.tar.gz"
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
command_path="$bin_dir/$name"
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
  fetch_bearer() { printf 'header = "Authorization: Bearer %s"\n' "$token" | fetch --config - "$@"; }
  fetch_private_token() { printf 'header = "PRIVATE-TOKEN: %s"\n' "$token" | fetch --config - "$@"; }
  case "$source" in
  gitlab)
    # GitLab's project path can itself contain slashes, which the API expects
    # percent-encoded rather than literal in the path segment.
    encoded_repo=$(printf '%s' "$repo" | sed 's,/,%2F,g')
    fetch_private_token "$api/projects/$encoded_repo/releases/$release_tag" -o "$work/release.json" ||
      fail "cannot read release $release_tag; confirm it exists and the token grants access"
    # Read each brace-delimited object as one record, so an asset link's name
    # and its own url are matched together whether the API pretty-prints the
    # payload or returns it compact.
    #
    # The url must be a generic-package one, which is where publication puts
    # these assets. A name alone is not enough to identify a record: the
    # release's author and its commit each carry a "name" and a "url" of their
    # own, and a release whose author happened to match would otherwise send
    # this installer to that author's profile.
    asset_url() {
      awk -v want="$1" '
        BEGIN { RS = "{" }
        {
          record = $0
          gsub(/[ \t\r\n]+/, "", record)
          if (index(record, "\"name\":\"" want "\"") == 0) next
          if (match(record, /"url":"[^"]*\/packages\/generic\/[^"]*"/) == 0) next
          value = substr(record, RSTART + 7, RLENGTH - 8)
          gsub(/\\\//, "/", value)
          print value
          exit
        }
      ' "$work/release.json"
    }
    for asset in "$package" SHA256SUMS; do
      url=$(asset_url "$asset")
      [ -n "$url" ] || fail "release $release_tag publishes no asset named $asset"
      fetch_private_token "$url" -o "$work/$asset"
    done
    ;;
  github|*)
    if [ -n "$token" ]; then
      # A private repository serves release assets only through the API, by asset
      # id; the public download path answers 404. curl does not carry the
      # Authorization header across the redirect to signed storage.
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
      fetch_bearer --header 'Accept: application/vnd.github+json' \
        --header 'X-GitHub-Api-Version: 2022-11-28' \
        "$api/repos/$repo/releases/tags/$release_tag" -o "$work/release.json" ||
        fail "cannot read release $release_tag; confirm it exists and the token grants access"
      for asset in "$package" SHA256SUMS; do
        url=$(asset_url "$asset")
        [ -n "$url" ] || fail "release $release_tag publishes no asset named $asset"
        fetch_bearer --header 'Accept: application/octet-stream' "$url" -o "$work/$asset"
      done
    else
      base="https://github.com/$repo/releases/download/$release_tag"
      fetch "$base/$package" -o "$work/$package"
      fetch "$base/SHA256SUMS" -o "$work/SHA256SUMS"
    fi
    ;;
  esac
fi
expected=$(awk -v name="$package" '{file=$2; sub(/^\*/, "", file); sub(/^\.\//, "", file); if(file==name) print $1}' "$work/SHA256SUMS")
[ "${#expected}" -eq 64 ] || fail 'missing or ambiguous checksum entry'
[ "$(hash "$work/$package")" = "$expected" ] || fail 'checksum mismatch'
printf '%s\n' README.md LICENSE THIRD_PARTY_NOTICES.txt context-circuit-cli | LC_ALL=C sort > "$work/expected"
tar -tzf "$work/$package" | LC_ALL=C sort > "$work/actual"
cmp "$work/expected" "$work/actual" >/dev/null || fail 'unexpected archive contents'
mkdir "$work/package"
tar -xzf "$work/$package" -C "$work/package"
for file in context-circuit-cli README.md LICENSE THIRD_PARTY_NOTICES.txt; do
  [ -f "$work/package/$file" ] && [ ! -L "$work/package/$file" ] || fail 'archive contains non-regular files'
done
chmod 755 "$work/package/context-circuit-cli"
[ "$("$work/package/context-circuit-cli" version)" = "$version" ] || fail 'CLI version does not match the release'
target="$versions/$version-$platform-$arch"
if [ -e "$target" ] || [ -L "$target" ]; then
  [ -d "$target" ] && [ ! -L "$target" ] || fail 'existing version path is not a directory'
  cmp "$target/context-circuit-cli" "$work/package/context-circuit-cli" >/dev/null || fail 'existing version differs; preserved'
else
  mv "$work/package" "$target"
fi
ln -s "$target/context-circuit-cli" "$work/$name"
mv -f "$work/$name" "$command_path"
# Report the version-specific command as well. Workspaces pin their own CLI
# version, so a caller serving several workspaces invokes that path directly
# instead of relying on which version the shared command currently selects.
printf 'version: %s\nplatform: %s/%s\ncommand: %s\nversioned command: %s\nversion store: %s\nPATH directory: %s\n' \
  "$version" "$platform" "$arch" "$command_path" "$target/context-circuit-cli" "$versions" "$bin_dir"
