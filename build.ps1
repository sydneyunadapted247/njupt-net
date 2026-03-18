$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot

New-Item -ItemType Directory -Force -Path "dist" | Out-Null

Write-Host "Building windows/amd64..."
$env:CGO_ENABLED = "0"
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -trimpath -ldflags "-s -w" -o dist/njupt-net-windows-amd64.exe ./cmd/njupt-net

Write-Host "Building linux/arm64..."
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -trimpath -ldflags "-s -w" -o dist/njupt-net-linux-arm64 ./cmd/njupt-net

Write-Host "Done. Artifacts are in ./dist"
