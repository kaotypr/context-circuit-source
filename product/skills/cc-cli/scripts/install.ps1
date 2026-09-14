param(
  [Parameter(Mandatory=$true)][string]$Version,
  [string]$BinDir = (Join-Path $env:LOCALAPPDATA 'ContextCircuit\bin'),
  [string]$Archive,
  [string]$Checksums
)
$ErrorActionPreference = 'Stop'
if ($Version -notmatch '^2\.[A-Za-z0-9.+-]+$' -or $Version.Contains('..')) { throw 'Supply an exact compatible v2 CLI version without a v prefix.' }
if ($env:OS -ne 'Windows_NT') { throw 'Use install.sh on macOS/Linux.' }
$architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
$arch = switch ($architecture) { 'x64' { 'amd64' }; 'arm64' { 'arm64' }; default { throw "Unsupported architecture: $_" } }
$package = "context-circuit-cli-v$Version-windows-$arch.zip"
[void](New-Item -ItemType Directory -Force -Path $BinDir)
$BinDir = (Resolve-Path -LiteralPath $BinDir).Path
$versions = Join-Path $BinDir '.context-circuit-versions'
if ((Test-Path -LiteralPath $versions) -and ((Get-Item -LiteralPath $versions).Attributes -band [IO.FileAttributes]::ReparsePoint)) { throw 'Version directory must not be a link.' }
[void](New-Item -ItemType Directory -Force -Path $versions)
$lockPath = Join-Path $versions 'install.lock'
$lock = [IO.File]::Open($lockPath, [IO.FileMode]::OpenOrCreate, [IO.FileAccess]::ReadWrite, [IO.FileShare]::None)
$work = Join-Path $versions ('.install-' + [Guid]::NewGuid().ToString('N'))
try {
  [void](New-Item -ItemType Directory -Path $work)
  $commandPath = Join-Path $BinDir 'context-circuit.cmd'
  $commandHash = Join-Path $versions 'launcher.sha256'
  if (Test-Path -LiteralPath (Join-Path $BinDir 'context-circuit.exe')) { throw 'An existing executable would shadow the launcher; use another -BinDir.' }
  if (Test-Path -LiteralPath $commandPath) {
    if (!(Test-Path -LiteralPath $commandHash) -or ((Get-FileHash -Algorithm SHA256 -LiteralPath $commandPath).Hash -ne (Get-Content -Raw -LiteralPath $commandHash).Trim())) { throw 'Existing command is unmanaged or modified; preserved.' }
  }
  $download = Join-Path $work $package
  $sums = Join-Path $work 'SHA256SUMS'
  if ($Archive -or $Checksums) {
    if (!$Archive -or !$Checksums) { throw 'Offline installation needs both -Archive and -Checksums.' }
    Copy-Item -LiteralPath $Archive -Destination $download
    Copy-Item -LiteralPath $Checksums -Destination $sums
  } else {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    $base = "https://github.com/kaotypr/context-circuit-source/releases/download/cli-v$Version"
    Invoke-WebRequest -UseBasicParsing -Uri "$base/$package" -OutFile $download
    Invoke-WebRequest -UseBasicParsing -Uri "$base/SHA256SUMS" -OutFile $sums
  }
  $candidateHashes = @(Get-Content -LiteralPath $sums | ForEach-Object {
    if ($_ -match '^([a-fA-F0-9]{64})\s+\*?(?:\./)?(.+)$' -and $Matches[2] -eq $package) { $Matches[1] }
  })
  if ($candidateHashes.Count -ne 1 -or (Get-FileHash -Algorithm SHA256 -LiteralPath $download).Hash -ne $candidateHashes[0]) { throw 'Missing, ambiguous, or mismatched checksum.' }
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  $zip = [IO.Compression.ZipFile]::OpenRead($download)
  try {
    $names = @($zip.Entries | ForEach-Object { $_.FullName } | Sort-Object)
    $expected = @('README.md', 'THIRD_PARTY_NOTICES.txt', 'context-circuit.exe' | Sort-Object)
    if (($names -join "`n") -cne ($expected -join "`n")) { throw 'Unexpected archive contents.' }
  } finally { $zip.Dispose() }
  $staged = Join-Path $work 'package'
  [IO.Compression.ZipFile]::ExtractToDirectory($download, $staged)
  $binary = Join-Path $staged 'context-circuit.exe'
  $reported = & $binary version
  if ($LASTEXITCODE -ne 0 -or $reported -ne $Version) { throw 'Executable version does not match the release.' }
  $target = Join-Path $versions "$Version-windows-$arch"
  if (Test-Path -LiteralPath $target) {
    if ((Get-Item -LiteralPath $target).Attributes -band [IO.FileAttributes]::ReparsePoint) { throw 'Existing version path is a link.' }
    if ((Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $target 'context-circuit.exe')).Hash -ne (Get-FileHash -Algorithm SHA256 -LiteralPath $binary).Hash) { throw 'Existing version differs; preserved.' }
  } else { Move-Item -LiteralPath $staged -Destination $target }
  # A small command shim allows updates while an older executable is in use.
  $launcher = Join-Path $work 'context-circuit.cmd'
  $content = '@"%~dp0.context-circuit-versions\' + "$Version-windows-$arch" + '\context-circuit.exe" %*' + "`r`n"
  [IO.File]::WriteAllText($launcher, $content, [Text.Encoding]::ASCII)
  if (Test-Path -LiteralPath $commandPath) { [IO.File]::Replace($launcher, $commandPath, $null) }
  else { [IO.File]::Move($launcher, $commandPath) }
  [IO.File]::WriteAllText($commandHash, (Get-FileHash -Algorithm SHA256 -LiteralPath $commandPath).Hash)
  Write-Output "version: $Version`nplatform: windows/$arch`ncommand: $commandPath`nPATH directory: $BinDir"
} finally {
  if (Test-Path -LiteralPath $work) { Remove-Item -Recurse -Force -LiteralPath $work }
  $lock.Dispose()
}
