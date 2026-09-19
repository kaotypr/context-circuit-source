#!/bin/sh
# Explicit publication preparation: assemble a release, replace a clean template
# checkout, commit, and tag. Push and hosted releases remain separate caller steps.
set -eu
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
[ "$#" -eq 2 ] || fail 'usage: publish-template.sh <version-without-v> <template-checkout-dir>'
version=$1
template_dir=$2
case "$version" in ''|v*|*[!A-Za-z0-9.+-]*|.*|-*) fail "invalid version: $version" ;; esac
source_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
request="$source_root/release/requests/$version.md"
[ -f "$request" ] || fail "missing release request: release/requests/$version.md"
[ -d "$template_dir/.git" ] || fail "expected a standalone template checkout: $template_dir"
template_dir=$(CDPATH= cd -- "$template_dir" && pwd -P)
[ "$template_dir" != "$source_root" ] || fail 'cannot publish over the source checkout'
[ "$(git -C "$template_dir" rev-parse --show-toplevel)" = "$template_dir" ] || fail 'target must be its checkout root'
[ -z "$(git -C "$template_dir" status --porcelain --untracked-files=all --ignored)" ] || fail 'template checkout has changed, untracked, or ignored files; preserve them before publication'
original_revision=$(git -C "$template_dir" rev-parse -q --verify HEAD || true)
binding="$source_root/release/binding.yaml"
destination_ref=$(sed -n 's/^destination_ref:[[:space:]]*//p' "$binding" | head -n1)
[ -n "$destination_ref" ] || fail 'missing destination ref'
[ "$(git -C "$template_dir" symbolic-ref --short HEAD)" = "$destination_ref" ] || fail "select template branch $destination_ref"
if git -C "$template_dir" rev-parse -q --verify "refs/tags/v$version" >/dev/null 2>&1; then fail "tag exists: v$version"; fi
# Files the destination repository owns rather than a workspace. They describe the
# published repository to whoever lands on it, and GitHub renders them as tabs
# beside the README. They are deliberately absent from the release manifest,
# because a workspace is somebody else's repository: a CONTRIBUTING.md about
# Context Circuit is noise in their checkout, and a LICENSE at their root would
# make GitHub label their own project with this one's license.
repo_owned_dir="$source_root/release/template-repo"
[ -d "$repo_owned_dir" ] || fail 'missing release/template-repo'
repo_owned=''
for path in "$repo_owned_dir"/*; do
  [ -f "$path" ] || fail "expected regular files in release/template-repo: $path"
  owned=${path##*/}
  case "$owned" in *[!A-Za-z0-9._-]*) fail "unexpected name in release/template-repo: $owned" ;; esac
  # The published README is the product guide, assembled from the manifest.
  case "$owned" in README.md) fail 'release/template-repo must not define README.md' ;; esac
  repo_owned="$repo_owned $owned"
done
[ -n "$repo_owned" ] || fail 'release/template-repo is empty'
# The license is not duplicated there. The repository shows the same 0BSD text a
# workspace receives at .context-circuit/LICENSE, so the terms on display and the
# terms travelling with the files cannot disagree.
[ -f "$source_root/product/LICENSE" ] || fail 'missing product/LICENSE'

name=${MAINTAINER_NAME:-$(git -C "$source_root" config user.name 2>/dev/null || true)}
email=${MAINTAINER_EMAIL:-$(git -C "$source_root" config user.email 2>/dev/null || true)}
[ -n "$name" ] && [ -n "$email" ] || fail 'set MAINTAINER_NAME and MAINTAINER_EMAIL'

# Version-specific output is new and retained for upload; never overwrite dist.
output_dir="$source_root/dist/v$version"
[ ! -e "$output_dir" ] && [ ! -L "$output_dir" ] || fail "release output exists: $output_dir"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT HUP INT TERM
sh "$source_root/scripts/release-artifact.sh" "$work/stage" "$output_dir" "v$version"
artifact="$output_dir/context-circuit-v$version"

# Compare only the committed tree, excluding release notes and version stamps.
normalize() {
  directory=$1
  rm -f "$directory/CHANGELOG.md" "$directory/LICENSE"
  for owned in $repo_owned; do rm -f "$directory/$owned"; done
  if [ -f "$directory/.context-circuit/VERSION" ]; then printf 'NORMALIZED\n' > "$directory/.context-circuit/VERSION"; fi
}
if git -C "$template_dir" rev-parse -q --verify HEAD >/dev/null 2>&1; then
  cp -R "$artifact" "$work/new"
  mkdir "$work/old"
  git -C "$template_dir" archive HEAD | tar -xf - -C "$work/old"
  normalize "$work/new"
  normalize "$work/old"
  if diff -r "$work/new" "$work/old" >/dev/null 2>&1; then fail 'nothing to publish except a version change'; fi
fi
# Assembly can take time. Recheck immediately before replacing the target so a
# new edit or commit made during the build is preserved as well.
[ -z "$(git -C "$template_dir" status --porcelain --untracked-files=all --ignored)" ] || fail 'template checkout changed during assembly; prepared assets were retained'
[ "$(git -C "$template_dir" rev-parse -q --verify HEAD || true)" = "$original_revision" ] || fail 'template HEAD changed during assembly'
[ "$(git -C "$template_dir" symbolic-ref --short HEAD)" = "$destination_ref" ] || fail 'template branch changed during assembly'
find "$template_dir" -mindepth 1 -maxdepth 1 ! -name .git ! -name CHANGELOG.md -exec rm -rf {} +
(cd "$artifact" && tar -cf - .) | tar -xf - -C "$template_dir"
# Restore the destination's own landing-page files over the assembled artifact.
for owned in $repo_owned; do cp "$repo_owned_dir/$owned" "$template_dir/$owned"; done
cp "$source_root/product/LICENSE" "$template_dir/LICENSE"
{
  printf '# Changelog\n\n## v%s — %s\n\n' "$version" "$(date +%Y-%m-%d)"
  awk 'BEGIN{fm=0} /^---[[:space:]]*$/{fm++; next} fm>=2{print}' "$request"
  printf '\n'
  if [ -f "$template_dir/CHANGELOG.md" ]; then tail -n +3 "$template_dir/CHANGELOG.md"; fi
} > "$work/CHANGELOG.md"
mv "$work/CHANGELOG.md" "$template_dir/CHANGELOG.md"
# Disable hooks for this release commit so hosts cannot inject agent attribution.
mkdir "$work/no-hooks"
git -C "$template_dir" add -A
GIT_AUTHOR_NAME=$name GIT_AUTHOR_EMAIL=$email GIT_COMMITTER_NAME=$name GIT_COMMITTER_EMAIL=$email \
  git -C "$template_dir" -c core.hooksPath="$work/no-hooks" commit -q -m "chore(release): publish v$version"
[ "$(git -C "$template_dir" log -1 --format='%B')" = "chore(release): publish v$version" ] || fail 'unexpected commit attribution or message'
GIT_COMMITTER_NAME=$name GIT_COMMITTER_EMAIL=$email git -C "$template_dir" tag -a "v$version" -m "v$version"
printf 'artifact_archive: %s/context-circuit-v%s.tar.gz\n' "$output_dir" "$version"
printf 'release_assets_dir: %s\n' "$output_dir"
printf 'commit: %s\n' "$(git -C "$template_dir" rev-parse HEAD)"
printf 'remaining_gate: push %s, tag v%s, and hosted releases\n' "$destination_ref" "$version"
