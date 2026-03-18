param(
	[ValidateSet("local", "all")]
	[string]$Mode = "local"
)

$ErrorActionPreference = "Stop"

Set-Location (Join-Path $PSScriptRoot "..")

New-Item -ItemType Directory -Force -Path "dist" | Out-Null

function Build-Target {
	param(
		[Parameter(Mandatory = $true)][string]$GoOS,
		[Parameter(Mandatory = $true)][string]$GoArch,
		[string]$Extension = ""
	)

	$out = "dist/njupt-net-$GoOS-$GoArch$Extension"
	Write-Host "Building $GoOS/$GoArch -> $out"

	$env:CGO_ENABLED = "0"
	$env:GOOS = $GoOS
	$env:GOARCH = $GoArch
	go build -trimpath -ldflags "-s -w" -o $out ./cmd/njupt-net
}

try {
	if ($Mode -eq "local") {
		Build-Target -GoOS "windows" -GoArch "amd64" -Extension ".exe"
		Build-Target -GoOS "linux" -GoArch "arm64"
	}
	else {
		Build-Target -GoOS "windows" -GoArch "amd64" -Extension ".exe"
		Build-Target -GoOS "linux" -GoArch "amd64"
		Build-Target -GoOS "linux" -GoArch "arm64"
		Build-Target -GoOS "darwin" -GoArch "arm64"
	}

	Write-Host "Done. Artifacts are in ./dist"
}
finally {
	Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
	Remove-Item Env:GOOS -ErrorAction SilentlyContinue
	Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
}
