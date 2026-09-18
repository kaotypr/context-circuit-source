param(
  [Parameter(Mandatory=$true)][string]$Version,
  [string]$BinDir = (Join-Path $env:LOCALAPPDATA 'ContextCircuit\bin'),
  [string]$Archive,
  [string]$Checksums,
  [string]$Token,
  [string]$GitlabUrl
)
$ErrorActionPreference = 'Stop'
foreach ($fallback in $env:CONTEXT_CIRCUIT_TOKEN, $env:GH_TOKEN, $env:GITHUB_TOKEN) {
  if (!$Token) { $Token = $fallback }
}
if ($Token -and $Token -notmatch '^[A-Za-z0-9_-]+$') { throw 'Token contains unexpected characters.' }
# An organization that mirrors CLI releases into its own GitLab project names it
# here. Unset, the installer reads the product's own GitHub releases. The API
# base and project path are derived from that URL, so this shipped script carries
# no organization's own address, and the credential reaches only one registry.
if (!$GitlabUrl) { $GitlabUrl = $env:CONTEXT_CIRCUIT_GITLAB_URL }
if ($GitlabUrl) {
  if ($GitlabUrl -notmatch '^https://') { throw 'Mirror URL must begin with https://.' }
  $rest = $GitlabUrl.TrimEnd('/')
  if ($rest.EndsWith('.git')) { $rest = $rest.Substring(0, $rest.Length - 4) }
  $rest = $rest.Substring(8)
  $slash = $rest.IndexOf('/')
  if ($slash -lt 1 -or $slash -ge $rest.Length - 1) { throw 'Mirror URL needs a host and a project path.' }
  $gitlabHost = $rest.Substring(0, $slash)
  $gitlabProject = $rest.Substring($slash + 1)
  if ("$gitlabHost$gitlabProject" -notmatch '^[A-Za-z0-9._:/-]+$') { throw 'Mirror URL contains unexpected characters.' }
  $gitlabProject = $gitlabProject.Replace('/', '%2F')
}
if ($Version -notmatch '^2\.[A-Za-z0-9.+-]+$' -or $Version.Contains('..')) { throw 'Supply an exact compatible v2 CLI version without a v prefix.' }
if ($env:OS -ne 'Windows_NT') { throw 'Use install.sh on macOS/Linux.' }
$architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
$arch = switch ($architecture) { 'x64' { 'amd64' }; 'arm64' { 'arm64' }; default { throw "Unsupported architecture: $_" } }
# The release is tagged `cli-v<version>`; the assets under it keep the
# executable's own name, so a downloaded archive still says what it holds.
$releaseTag = "cli-v$Version"
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
  $commandPath = Join-Path $BinDir 'context-circuit-cli.cmd'
  $commandHash = Join-Path $versions 'launcher.sha256'
  if (Test-Path -LiteralPath (Join-Path $BinDir 'context-circuit-cli.exe')) { throw 'An existing executable would shadow the launcher; use another -BinDir.' }
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
    if ($GitlabUrl) {
      # A generic package is addressed by version and file name, so a mirror
      # needs no release lookup and no asset ids.
      $base = "https://$gitlabHost/api/v4/projects/$gitlabProject/packages/generic/context-circuit-cli/$Version"
      $mirrorHeaders = @{}
      if ($Token) { $mirrorHeaders = @{ 'PRIVATE-TOKEN' = $Token } }
      foreach ($wanted in @(@{ Name = $package; File = $download }, @{ Name = 'SHA256SUMS'; File = $sums })) {
        try {
          Invoke-WebRequest -UseBasicParsing -Uri "$base/$($wanted.Name)" -Headers $mirrorHeaders -OutFile $wanted.File
        } catch { throw "Cannot download $($wanted.Name) from the mirror; confirm the version is published there and the token grants access." }
      }
    } elseif ($Token) {
      # A private repository serves release assets only through the API, by
      # asset id; the public download path answers 404. Invoke-WebRequest drops
      # the Authorization header on the redirect to signed storage, which is
      # what that storage requires.
      $host_ = if ($env:CONTEXT_CIRCUIT_API) { $env:CONTEXT_CIRCUIT_API } else { 'https://api.github.com' }
      $api = "$host_/repos/kaotypr/context-circuit-source"
      $auth = @{ Authorization = "Bearer $Token" }
      try {
        $release = Invoke-RestMethod -UseBasicParsing -Uri "$api/releases/tags/$releaseTag" -Headers ($auth + @{ Accept = 'application/vnd.github+json'; 'X-GitHub-Api-Version' = '2022-11-28' })
      } catch { throw "Cannot read release $releaseTag; confirm it exists and the token grants access." }
      foreach ($wanted in @(@{ Name = $package; File = $download }, @{ Name = 'SHA256SUMS'; File = $sums })) {
        $asset = @($release.assets | Where-Object { $_.name -eq $wanted.Name })[0]
        if (!$asset) { throw "Release $releaseTag publishes no asset named $($wanted.Name)." }
        Invoke-WebRequest -UseBasicParsing -Uri $asset.url -Headers ($auth + @{ Accept = 'application/octet-stream' }) -OutFile $wanted.File
      }
    } else {
      $base = "https://github.com/kaotypr/context-circuit-source/releases/download/$releaseTag"
      Invoke-WebRequest -UseBasicParsing -Uri "$base/$package" -OutFile $download
      Invoke-WebRequest -UseBasicParsing -Uri "$base/SHA256SUMS" -OutFile $sums
    }
  }
  $candidateHashes = @(Get-Content -LiteralPath $sums | ForEach-Object {
    if ($_ -match '^([a-fA-F0-9]{64})\s+\*?(?:\./)?(.+)$' -and $Matches[2] -eq $package) { $Matches[1] }
  })
  if ($candidateHashes.Count -ne 1 -or (Get-FileHash -Algorithm SHA256 -LiteralPath $download).Hash -ne $candidateHashes[0]) { throw 'Missing, ambiguous, or mismatched checksum.' }
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  $zip = [IO.Compression.ZipFile]::OpenRead($download)
  try {
    $names = @($zip.Entries | ForEach-Object { $_.FullName } | Sort-Object)
    $expected = @('README.md', 'LICENSE', 'THIRD_PARTY_NOTICES.txt', 'context-circuit-cli.exe' | Sort-Object)
    if (($names -join "`n") -cne ($expected -join "`n")) { throw 'Unexpected archive contents.' }
  } finally { $zip.Dispose() }
  $staged = Join-Path $work 'package'
  [IO.Compression.ZipFile]::ExtractToDirectory($download, $staged)
  $binary = Join-Path $staged 'context-circuit-cli.exe'
  $reported = & $binary version
  if ($LASTEXITCODE -ne 0 -or $reported -ne $Version) { throw 'Executable version does not match the release.' }
  $target = Join-Path $versions "$Version-windows-$arch"
  if (Test-Path -LiteralPath $target) {
    if ((Get-Item -LiteralPath $target).Attributes -band [IO.FileAttributes]::ReparsePoint) { throw 'Existing version path is a link.' }
    if ((Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $target 'context-circuit-cli.exe')).Hash -ne (Get-FileHash -Algorithm SHA256 -LiteralPath $binary).Hash) { throw 'Existing version differs; preserved.' }
  } else { Move-Item -LiteralPath $staged -Destination $target }
  # A small command shim allows updates while an older executable is in use.
  $launcher = Join-Path $work 'context-circuit-cli.cmd'
  $content = '@"%~dp0.context-circuit-versions\' + "$Version-windows-$arch" + '\context-circuit-cli.exe" %*' + "`r`n"
  [IO.File]::WriteAllText($launcher, $content, [Text.Encoding]::ASCII)
  # Replace needs a real backup path: PowerShell binds $null to a string
  # parameter as an empty string, which is not a legal path. The backup stays in
  # the working directory and is discarded with it.
  $backup = Join-Path $work 'context-circuit-cli.cmd.backup'
  if (Test-Path -LiteralPath $commandPath) { [IO.File]::Replace($launcher, $commandPath, $backup) }
  else { [IO.File]::Move($launcher, $commandPath) }
  [IO.File]::WriteAllText($commandHash, (Get-FileHash -Algorithm SHA256 -LiteralPath $commandPath).Hash)
  # Report the version-specific command as well, so a caller serving several
  # workspaces invokes the version each one pins instead of the shared launcher.
  Write-Output "version: $Version`nplatform: windows/$arch`ncommand: $commandPath`nversioned command: $(Join-Path $target 'context-circuit-cli.exe')`nversion store: $versions`nPATH directory: $BinDir"
} finally {
  if (Test-Path -LiteralPath $work) { Remove-Item -Recurse -Force -LiteralPath $work }
  $lock.Dispose()
}
