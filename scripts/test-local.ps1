param(
    [string]$ExePath = "./dist/njupt-net-windows-amd64.exe",
    [string]$CredentialsPath = "./credentials.json",
    [string]$IP = "",
    [switch]$IncludeWriteOps
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot\..

$passed = 0
$failed = 0

function Write-Section {
    param([string]$Title)
    Write-Host ""
    Write-Host "=== $Title ===" -ForegroundColor Cyan
}

function Invoke-Step {
    param(
        [string]$Name,
        [ScriptBlock]$Action
    )

    Write-Host "[RUN] $Name" -ForegroundColor Yellow
    try {
        & $Action
        Write-Host "[PASS] $Name" -ForegroundColor Green
        $script:passed++
    }
    catch {
        Write-Host "[FAIL] $Name" -ForegroundColor Red
        Write-Host "       $($_.Exception.Message)" -ForegroundColor Red
        $script:failed++
        throw
    }
}

function Require-File {
    param([string]$Path)
    if (-not (Test-Path $Path)) {
        throw "Required file not found: $Path"
    }
}

function Detect-WifiIPv4 {
    $candidate = (
        Get-NetIPAddress -AddressFamily IPv4 |
        Where-Object {
            $_.IPAddress -notlike "169.*" -and
            $_.IPAddress -ne "127.0.0.1" -and
            $_.InterfaceAlias -match "Wi-Fi|WLAN|无线|wlan"
        } |
        Select-Object -First 1 -ExpandProperty IPAddress
    )
    return $candidate
}

function Exec-Cli {
    param([string[]]$Args)
    & $script:Exe @Args
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed with exit code $LASTEXITCODE: $script:Exe $($Args -join ' ')"
    }
}

Write-Section "Precheck"
Require-File $ExePath
Require-File $CredentialsPath

$script:Exe = (Resolve-Path $ExePath).Path
$cred = Get-Content $CredentialsPath -Raw | ConvertFrom-Json

$WUser = $cred.accounts.W.username
$WPwd = $cred.accounts.W.password
$BUser = $cred.accounts.B.username
$BPwd = $cred.accounts.B.password
$CmccUser = $cred.cmcc.account
$CmccPwd = $cred.cmcc.password

if ([string]::IsNullOrWhiteSpace($IP)) {
    $IP = Detect-WifiIPv4
}
if ([string]::IsNullOrWhiteSpace($IP)) {
    throw "Cannot auto-detect Wi-Fi IPv4. Please pass -IP explicitly."
}

Write-Host "Executable: $script:Exe"
Write-Host "Using IP: $IP"
Write-Host "IncludeWriteOps: $IncludeWriteOps"

Write-Section "Smoke"
Invoke-Step "root help" { Exec-Cli @("--help") }
Invoke-Step "service help" { Exec-Cli @("service", "--help") }
Invoke-Step "portal help" { Exec-Cli @("portal", "--help") }
Invoke-Step "raw help" { Exec-Cli @("raw", "--help") }

Write-Section "Core Read/Auth"
Invoke-Step "service login (W)" {
    Exec-Cli @("service", "login", "-u", $WUser, "-p", $WPwd)
}
Invoke-Step "portal login (W, mobile)" {
    Exec-Cli @("portal", "login", "-u", $WUser, "-p", $WPwd, "-i", $IP, "--isp", "mobile")
}

Write-Section "Raw Probe"
Invoke-Step "raw get login page" {
    Exec-Cli @("raw", "get", "/Self/login/?302=LI")
}
Invoke-Step "raw get public help page" {
    Exec-Cli @("raw", "get", "/Self/unlogin/help")
}
Invoke-Step "raw portal logout probe" {
    Exec-Cli @("raw", "get", "https://10.10.244.11:802/eportal/portal/logout?callback=dr1003&login_method=1&wlan_user_ip=$IP")
}

if ($IncludeWriteOps) {
    Write-Section "Write Ops (Side Effects)"

    Invoke-Step "service bind (cmcc to current account)" {
        Exec-Cli @("service", "bind", "--fld3", $CmccUser, "--fld4", $CmccPwd)
    }

    Invoke-Step "service migrate (W -> B)" {
        Exec-Cli @(
            "service", "migrate",
            "--from-user", $WUser,
            "--from-password", $WPwd,
            "--to-user", $BUser,
            "--to-password", $BPwd,
            "--fld3", $CmccUser,
            "--fld4", $CmccPwd
        )
    }
}

Write-Section "Summary"
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor Red

if ($failed -gt 0) {
    exit 1
}

exit 0
