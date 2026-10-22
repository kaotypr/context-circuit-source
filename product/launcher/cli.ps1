# Run the Context Circuit CLI this workspace pins, from the installer's version
# store. It never installs, never reads PATH, and runs only a binary the
# installer verified. Call it as:
#   powershell -NoProfile -ExecutionPolicy Bypass -File <workspace>/.context-circuit/cli.ps1 ARGS...
$ErrorActionPreference = 'Stop'
function Refuse([string]$why) { [Console]::Error.WriteLine("refused: $why; nothing was run"); exit 1 }
$here = (Resolve-Path -LiteralPath $PSScriptRoot).ProviderPath
$root = Split-Path -Parent $here
$pinFile = Join-Path $here 'CLI_VERSION'
if (-not (Test-Path -LiteralPath $pinFile -PathType Leaf)) { Refuse "$pinFile is missing" }
# An empty file reads as $null, which the string conversion turns into ''.
$pin = "$(Get-Content -Raw -LiteralPath $pinFile)".Trim()
# The installer's own rule.
if ($pin -notmatch '^2\.[A-Za-z0-9.+-]+$' -or $pin.Contains('..')) { Refuse "invalid pin '$pin'" }
$bin = if ($env:CONTEXT_CIRCUIT_BIN_DIR) { $env:CONTEXT_CIRCUIT_BIN_DIR } else { Join-Path $env:LOCALAPPDATA 'ContextCircuit\bin' }
# The installer's architecture rule.
$arch = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()) {
  'x64' { 'amd64' } 'arm64' { 'arm64' } default { Refuse "unsupported architecture: $_" }
}
$cli = Join-Path $bin ".context-circuit-versions\$pin-windows-$arch\context-circuit-cli.exe"
$item = Get-Item -LiteralPath $cli -ErrorAction SilentlyContinue
if (-not $item -or $item.PSIsContainer -or ($item.Attributes -band [IO.FileAttributes]::ReparsePoint)) {
  [Console]::Error.WriteLine("refused: this workspace pins CLI $pin, which is not installed at $cli; nothing was run.")
  [Console]::Error.WriteLine("Install exactly $pin with $root\.agents\skills\cc-cli\SKILL.md, using the cli_registry in $root\workspace.yaml.")
  [Console]::Error.WriteLine('If it was installed with -BinDir, set CONTEXT_CIRCUIT_BIN_DIR to that directory.')
  exit 1
}
# The pin belongs to this workspace. Another workspace has its own launcher.
# GetFullPath expands 8.3 short names (C:\Users\RUNNER~1) of an existing path,
# on .NET Framework and .NET alike, and turns / into \; -ne ignores case.
function Key([string]$path) { [IO.Path]::GetFullPath($path).TrimEnd('\') }
$rootKey = Key $root
for ($i = 0; $i -lt $args.Count; $i++) {
  $value = $null
  if ($args[$i] -in '-workspace', '--workspace') { if ($i + 1 -lt $args.Count) { $value = $args[$i + 1] } }
  elseif ($args[$i] -match '^--?workspace=(.*)$') { $value = $Matches[1] }
  if ($null -ne $value) {
    $resolved = Resolve-Path -LiteralPath $value -ErrorAction SilentlyContinue
    if (-not $resolved -or (Key $resolved.ProviderPath) -ne $rootKey) {
      Refuse "--workspace $value is another workspace; run its own .context-circuit\cli.ps1"
    }
  }
}
& $cli --workspace $root @args
exit $LASTEXITCODE
