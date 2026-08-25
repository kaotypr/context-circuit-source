#!/bin/sh
# Publish one template release into a checkout of context-circuit-template.
#
# This performs the deterministic, locally-runnable part of publication:
# assemble, diff-gate, guard, full-tree replace, version stamp, changelog,
# commit, and annotated tag. It does NOT push and does NOT create the GitHub
# Release: those external effects remain the caller's (the workflow's) gate,
# mirroring release-artifact.sh's "remaining_gate". See
#   sources/system-design/context-circuit/v0.5/10-template-publication.md
set -eu

usage() { printf 'usage: sh scripts/publish-template.sh <version> <template-checkout-dir>\n' >&2; exit 2; }
fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }

[ "$#" -eq 2 ] || usage
version=$1
template_dir=$2

# Reject anything that is not a plain version token.
case "$version" in ''|v*|*/*|*' '*|.*|-*|*..*) fail "invalid version: $version" ;; esac

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
source_root=$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null) || fail 'not a git source checkout'

binding="$source_root/release/binding.yaml"
[ -f "$binding" ] || fail "missing release binding: $binding"
DESTINATION_REF=$(sed -n 's/^destination_ref:[[:space:]]*//p' "$binding" | head -n1)
[ -n "$DESTINATION_REF" ] || fail "invalid release binding: $binding"

request="$source_root/release/requests/$version.md"
[ -f "$request" ] || fail "missing release request: release/requests/$version.md"
[ -d "$template_dir/.git" ] || fail "not a git checkout: $template_dir"

# Maintainer identity: never the CI bot. Fall back to source git config locally.
name=${MAINTAINER_NAME:-$(git -C "$source_root" config user.name 2>/dev/null || true)}
email=${MAINTAINER_EMAIL:-$(git -C "$source_root" config user.email 2>/dev/null || true)}
[ -n "$name" ] && [ -n "$email" ] || fail 'set MAINTAINER_NAME and MAINTAINER_EMAIL'

# Guard: a published tag is immutable and never re-published.
if git -C "$template_dir" rev-parse -q --verify "refs/tags/v$version" >/dev/null 2>&1; then
  fail "tag already published: v$version"
fi

# Assemble a source-versioned artifact into a temp area. The template version is
# stamped later; assembly stays source-versioned.
src_version=v$(sed -n 's/^runtime_version:[[:space:]]*//p' "$source_root/wrapper/manifest.yaml" | head -n1)
work=$(mktemp -d)
cleanup() { rm -rf "$work"; }
trap cleanup EXIT HUP INT TERM
sh "$source_root/scripts/release-artifact.sh" "$work/stage" "$work/out" "$src_version" >"$work/assemble.log" 2>&1 \
  || { cat "$work/assemble.log" >&2; fail 'assembly failed'; }
artifact="$work/out/context-circuit-$src_version"
[ -d "$artifact" ] || fail 'assembled artifact not found'

# Normalize a tree for content comparison: drop .git and the publish-only
# CHANGELOG.md, and neutralize the stamped template_version line.
normalize() {
  src=$1; dst=$2
  cp -R "$src" "$dst"
  rm -rf "$dst/.git" "$dst/CHANGELOG.md"
  mf="$dst/wrapper/manifest.yaml"
  if [ -f "$mf" ]; then
    sed 's/^template_version:.*/template_version: NORMALIZED/' "$mf" > "$mf.norm" && mv "$mf.norm" "$mf"
  fi
}

# Diff-gate: publication requires a real content change. Skipped on the first
# release, when the template repository has no commit yet.
if git -C "$template_dir" rev-parse -q --verify HEAD >/dev/null 2>&1; then
  normalize "$artifact" "$work/na"
  normalize "$template_dir" "$work/np"
  if diff -r "$work/na" "$work/np" >/dev/null 2>&1; then
    fail "nothing to publish for v$version (assembled tree matches the latest release)"
  fi
else
  printf 'first-release: template has no prior commit; diff-gate skipped\n'
fi

# Full-tree replace: the committed tree becomes exactly the assembled artifact,
# preserving only .git and the accumulated CHANGELOG.md.
find "$template_dir" -mindepth 1 -maxdepth 1 ! -name .git ! -name CHANGELOG.md -exec rm -rf {} +
(CDPATH= cd "$artifact" && tar -cf - .) | (CDPATH= cd "$template_dir" && tar -xf -)

# Stamp the published version into the shipped manifest.
mf="$template_dir/wrapper/manifest.yaml"
[ -f "$mf" ] || fail 'assembled artifact is missing wrapper/manifest.yaml'
sed "s/^template_version:.*/template_version: $version/" "$mf" > "$mf.stamp" && mv "$mf.stamp" "$mf"

# Prepend the release notes (request body, frontmatter stripped) to CHANGELOG.md.
notes=$(awk 'BEGIN{fm=0} /^---[[:space:]]*$/{fm++; next} fm>=2{print}' "$request")
changelog="$template_dir/CHANGELOG.md"
new_changelog="$work/CHANGELOG.md"
{
  printf '# Changelog\n\n'
  printf '## v%s — %s\n\n' "$version" "$(date +%Y-%m-%d)"
  printf '%s\n\n' "$notes"
  if [ -f "$changelog" ]; then
    tail -n +3 "$changelog"
  fi
} > "$new_changelog"
mv "$new_changelog" "$changelog"

# One commit, one annotated tag, authored by the maintainer. No attribution.
git -C "$template_dir" add -A
GIT_AUTHOR_NAME=$name GIT_AUTHOR_EMAIL=$email \
GIT_COMMITTER_NAME=$name GIT_COMMITTER_EMAIL=$email \
  git -C "$template_dir" commit -q -m "chore(release): v$version"
GIT_COMMITTER_NAME=$name GIT_COMMITTER_EMAIL=$email \
  git -C "$template_dir" tag -a "v$version" -m "v$version"

# Persist the archive next to dist/ for the caller to attach to the Release.
mkdir -p "$source_root/dist"
archive="$source_root/dist/context-circuit-v$version.tar.gz"
cp "$work/out/context-circuit-$src_version.tar.gz" "$archive"

printf 'version: v%s\n' "$version"
printf 'source_version: %s\n' "$src_version"
printf 'template_checkout: %s\n' "$template_dir"
printf 'commit: %s\n' "$(git -C "$template_dir" rev-parse HEAD)"
printf 'tag: v%s\n' "$version"
printf 'artifact_archive: %s\n' "$archive"
printf 'remaining_gate: push %s + tag v%s + GitHub Release\n' "$DESTINATION_REF" "$version"
