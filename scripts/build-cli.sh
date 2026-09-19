#!/bin/sh
# CLI-only packages. CLI versions and publication are independent of the template.
# The default output under dist/ is clean-rebuilt each run; an explicit output
# directory must be new and is never cleaned or replaced.
set -eu
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
[ "$#" -le 2 ] || fail 'usage: build-cli.sh [version-without-v] [new-output-dir]'
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
version=${1:-$(cat "$source_root/CLI_VERSION")}
case "$version" in ''|v*|*[!A-Za-z0-9.+-]*|.*|-*) fail "invalid CLI version: $version" ;; esac
if [ "$#" -ge 2 ]; then
  output_dir=$2
  [ ! -e "$output_dir" ] && [ ! -L "$output_dir" ] || fail "output already exists: $output_dir"
else
  # Own subdirectory per build so a CLI rebuild never removes workspace assets.
  output_dir=$source_root/dist/cli-$version
  case "$output_dir" in ''|/|"$source_root"|"$source_root"/) fail "refusing to clean: $output_dir" ;; esac
  rm -rf "$output_dir"
fi
mkdir -p "$output_dir"
output_dir=$(CDPATH= cd -- "$output_dir" && pwd)
staging_dir=$(mktemp -d)
trap 'rm -rf "$staging_dir"' EXIT HUP INT TERM
cd "$source_root"
artifact_name="context-circuit-cli-v$version"
# Include the licenses for dependencies and the Go runtime in each binary asset.
# The CLI's own license ships beside them: Apache-2.0 asks a redistributor
# to pass the License on with the work, and an archive is read wherever it was
# downloaded, where nothing else says what its terms are.
notices="$staging_dir/THIRD_PARTY_NOTICES.txt"
printf 'Context Circuit third-party notices\n\n' > "$notices"
for module in github.com/goccy/go-yaml github.com/gofrs/flock golang.org/x/sys; do
  module_dir=$(go list -m -f '{{.Dir}}' "$module")
  [ -f "$module_dir/LICENSE" ] || fail "missing license for $module"
  printf '\n%s\n\n' "$module" >> "$notices"
  cat "$module_dir/LICENSE" >> "$notices"
done
printf '\nGo runtime and standard library\n\n' >> "$notices"
go_root=$(go env GOROOT)
go_license="$go_root/LICENSE"
# Homebrew places the license beside libexec; official distributions use GOROOT.
[ -f "$go_license" ] || go_license="$go_root/../LICENSE"
[ -f "$go_license" ] || fail 'Go distribution license was not found'
cat "$go_license" >> "$notices"

for target in darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64; do
  target_os=${target%/*}
  target_arch=${target#*/}
  package_dir="$staging_dir/$target_os-$target_arch"
  mkdir "$package_dir"
  binary=context-circuit-cli
  [ "$target_os" != windows ] || binary=context-circuit-cli.exe
  CGO_ENABLED=0 GOOS=$target_os GOARCH=$target_arch go build -trimpath \
    -ldflags "-s -w -X main.version=${version#v}" -o "$package_dir/$binary" ./cmd/context-circuit
  cp CLI.md "$package_dir/README.md"
  cp LICENSE "$package_dir/LICENSE"
  cp "$notices" "$package_dir/THIRD_PARTY_NOTICES.txt"
  package_name="$artifact_name-$target_os-$target_arch"
  if [ "$target_os" = windows ]; then
    (cd "$package_dir" && zip -q "$output_dir/$package_name.zip" "$binary" README.md LICENSE THIRD_PARTY_NOTICES.txt)
  else
    (cd "$package_dir" && tar -cf - "$binary" README.md LICENSE THIRD_PARTY_NOTICES.txt) | gzip -n > "$output_dir/$package_name.tar.gz"
  fi
  printf 'binary_target: %s\n' "$target"
done
(
  cd "$output_dir"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum ./*.tar.gz ./*.zip
  else
    shasum -a 256 ./*.tar.gz ./*.zip
  fi
) > "$output_dir/SHA256SUMS"
printf 'cli_version: %s\n' "$version"
printf 'cli_assets_dir: %s\n' "$output_dir"
