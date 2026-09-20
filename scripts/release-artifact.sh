#!/bin/sh
# Source-only release assembly. Go is required here, never in a user's workspace.
set -eu
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
[ "$#" -eq 3 ] || fail 'usage: release-artifact.sh <staging-dir> <new-output-dir> <version>'
staging_dir=$1
output_dir=$2
version=$3
case "$version" in ''|*[!A-Za-z0-9.+-]*|.*|-*) fail "invalid version: $version" ;; esac
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
[ ! -e "$output_dir" ] && [ ! -L "$output_dir" ] || fail "output already exists: $output_dir"
mkdir -p "$staging_dir"
staging_dir=$(CDPATH= cd -- "$staging_dir" && pwd)
[ ! -e "$staging_dir/native" ] || fail 'staging area already used'
mkdir -p "$output_dir"
output_dir=$(CDPATH= cd -- "$output_dir" && pwd)
cd "$source_root"

# A native helper exports its embedded blank seed. Cross-built binaries are
# never run as evidence of native support.
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(cat CLI_VERSION)" -o "$staging_dir/native" ./cmd/context-circuit
artifact_name="context-circuit-$version"
artifact_dir="$output_dir/$artifact_name"
"$staging_dir/native" template export --path "$artifact_dir" >/dev/null
printf '%s\n' "${version#v}" > "$artifact_dir/.context-circuit/VERSION"

# Record which registry this template's workspaces install their CLI from. The
# seed carries none, which the installer reads as its own default; a release
# published anywhere else writes the registry it was published to, because a
# developer's machine has no CI variables to learn it from. This is done here
# rather than after assembly so the published tree and the release archive are
# the same bytes.
if [ -n "${CLI_REGISTRY_SOURCE:-}" ]; then
  registry_repository=${CLI_REGISTRY_REPOSITORY:-}
  [ -n "$registry_repository" ] || fail 'set CLI_REGISTRY_REPOSITORY beside CLI_REGISTRY_SOURCE'
  case "$CLI_REGISTRY_SOURCE" in
    github) ;;
    gitlab) [ -n "${CLI_REGISTRY_API:-}" ] || fail 'a gitlab CLI_REGISTRY_SOURCE needs CLI_REGISTRY_API' ;;
    *) fail "unknown CLI_REGISTRY_SOURCE: $CLI_REGISTRY_SOURCE" ;;
  esac
  {
    printf 'cli_registry:\n'
    printf '  source: %s\n' "$CLI_REGISTRY_SOURCE"
    printf '  repository: %s\n' "$registry_repository"
    if [ -n "${CLI_REGISTRY_API:-}" ]; then printf '  api: %s\n' "$CLI_REGISTRY_API"; fi
  } >> "$artifact_dir/workspace.yaml"
fi
awk 'NF && $1 !~ /^#/ {print $2}' scripts/release-manifest.txt | LC_ALL=C sort > "$staging_dir/expected"
(cd "$artifact_dir" && find . -type f -print | sed 's|^./||' | LC_ALL=C sort) > "$staging_dir/actual"
cmp "$staging_dir/expected" "$staging_dir/actual" || fail 'seed contains missing or unexpected files'
(cd "$artifact_dir" && tar -cf - .) | gzip -n > "$output_dir/$artifact_name.tar.gz"

(cd "$output_dir" && if command -v sha256sum >/dev/null 2>&1; then sha256sum ./*.tar.gz; else shasum -a 256 ./*.tar.gz; fi) > "$output_dir/SHA256SUMS"
printf 'version: %s\n' "$version"
printf 'artifact_dir: %s\n' "$artifact_dir"
printf 'artifact_archive: %s/%s.tar.gz\n' "$output_dir" "$artifact_name"
printf 'checksums: %s/SHA256SUMS\n' "$output_dir"
