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

# A native helper exports its embedded blank seed. Cross-built executables are
# never run as evidence of native support.
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${version#v}" -o "$staging_dir/native" ./cmd/context-circuit
artifact_name="context-circuit-$version"
artifact_dir="$output_dir/$artifact_name"
"$staging_dir/native" template export --path "$artifact_dir" >/dev/null
awk 'NF && $1 !~ /^#/ {print $2}' scripts/release-manifest.txt | LC_ALL=C sort > "$staging_dir/expected"
(cd "$artifact_dir" && find . -type f -print | sed 's|^./||' | LC_ALL=C sort) > "$staging_dir/actual"
cmp "$staging_dir/expected" "$staging_dir/actual" || fail 'seed contains missing or unexpected files'
(cd "$artifact_dir" && tar -cf - .) | gzip -n > "$output_dir/$artifact_name.tar.gz"

# Include the licenses for dependencies and the Go runtime in each binary asset.
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
  binary=context-circuit
  [ "$target_os" != windows ] || binary=context-circuit.exe
  CGO_ENABLED=0 GOOS=$target_os GOARCH=$target_arch go build -trimpath \
    -ldflags "-s -w -X main.version=${version#v}" -o "$package_dir/$binary" ./cmd/context-circuit
  cp product/README.md "$package_dir/README.md"
  cp "$notices" "$package_dir/THIRD_PARTY_NOTICES.txt"
  package_name="$artifact_name-$target_os-$target_arch"
  if [ "$target_os" = windows ]; then
    (cd "$package_dir" && zip -q "$output_dir/$package_name.zip" "$binary" README.md THIRD_PARTY_NOTICES.txt)
  else
    (cd "$package_dir" && tar -cf - "$binary" README.md THIRD_PARTY_NOTICES.txt) | gzip -n > "$output_dir/$package_name.tar.gz"
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
printf 'version: %s\n' "$version"
printf 'artifact_dir: %s\n' "$artifact_dir"
printf 'artifact_archive: %s/%s.tar.gz\n' "$output_dir" "$artifact_name"
printf 'checksums: %s/SHA256SUMS\n' "$output_dir"
