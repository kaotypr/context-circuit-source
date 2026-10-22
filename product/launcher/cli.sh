#!/bin/sh
# Run the Context Circuit CLI this workspace pins, from the installer's version
# store. It never installs, never reads PATH, and runs only a binary the
# installer verified. Call it as: sh <workspace>/.context-circuit/cli ARGS...
set -eu
refuse() { printf 'refused: %s; nothing was run\n' "$1" >&2; exit 1; }
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
root=$(dirname -- "$here")
[ -f "$here/CLI_VERSION" ] || refuse "$here/CLI_VERSION is missing"
pin=$(tr -d '[:space:]' < "$here/CLI_VERSION")
# The installer's own rule: no path separators, no "..".
case "$pin" in 2.*) ;; *) refuse "this workspace pins '$pin', not a v2 CLI version" ;; esac
case "$pin" in *[!A-Za-z0-9.+-]*|*..*) refuse "invalid pin '$pin'" ;; esac
# The installers' default directories, and the override they share.
case "$(uname -s)" in
  Darwin) os=darwin exe= bin=${CONTEXT_CIRCUIT_BIN_DIR:-$HOME/.local/bin} ;;
  Linux) os=linux exe= bin=${CONTEXT_CIRCUIT_BIN_DIR:-$HOME/.local/bin} ;;
  MINGW*|MSYS*|CYGWIN*) os=windows exe=.exe
    bin=$(cygpath -u "${CONTEXT_CIRCUIT_BIN_DIR:-${LOCALAPPDATA:?}\\ContextCircuit\\bin}") ;;
  *) refuse "unsupported OS: $(uname -s)" ;;
esac
case "$(uname -m)" in
  arm64|aarch64) arch=arm64 ;;
  x86_64|amd64) arch=amd64 ;;
  *) refuse "unsupported architecture: $(uname -m)" ;;
esac
store=$bin/.context-circuit-versions
cli=$store/$pin-$os-$arch/context-circuit-cli$exe
if [ "$os" = windows ] && [ "$arch" = amd64 ] && [ ! -f "$cli" ]; then
  # An x64 Git Bash on Windows on Arm reports x86_64, while install.ps1 picks
  # the operating system's architecture.
  alt=$store/$pin-windows-arm64/context-circuit-cli.exe
  if [ -f "$alt" ] && [ ! -L "$alt" ]; then cli=$alt; fi
fi
if [ ! -f "$cli" ] || [ -L "$cli" ]; then
  printf 'refused: this workspace pins CLI %s, which is not installed at %s; nothing was run.\n' "$pin" "$cli" >&2
  printf 'Install exactly %s with %s/.agents/skills/cc-cli/SKILL.md, using the cli_registry in %s/workspace.yaml.\n' "$pin" "$root" "$root" >&2
  printf 'If it was installed with --bin-dir, set CONTEXT_CIRCUIT_BIN_DIR to that directory.\n' >&2
  exit 1
fi
# The pin belongs to this workspace. Another workspace has its own launcher.
# On Windows, pwd keeps an 8.3 short name (RUNNER~1) as it was typed, and case
# does not tell directories apart.
key() {
  if [ "$os" = windows ]; then cygpath -m -l "$1" | tr '[:upper:]' '[:lower:]'; else printf '%s\n' "$1"; fi
}
same() {
  dir=$(CDPATH= cd -- "$1" 2>/dev/null && pwd -P) || refuse "--workspace $1 is not a directory"
  [ "$(key "$dir")" = "$(key "$root")" ] || refuse "--workspace $1 is another workspace; run $1/.context-circuit/cli for it"
}
want=
for arg in "$@"; do
  if [ -n "$want" ]; then same "$arg"; want=; continue; fi
  case "$arg" in
    -workspace|--workspace) want=1 ;;
    -workspace=*|--workspace=*) same "${arg#*=}" ;;
  esac
done
exec "$cli" --workspace "$root" "$@"
