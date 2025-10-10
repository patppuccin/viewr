# ==== Setup script for Windows ============================

# ---- Global flags and params -----------------------------
[CmdletBinding()]
param ()

$ErrorActionPreference = 'Continue'

$RootDir = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)

# ---- Helpers and utilities -------------------------------
function Write-ConsoleLog {
    param(
        [ValidateSet("INF", "DBG", "WRN", "ERR")][string]$Level = "INF",
        [Parameter(Mandatory)][string]$Message
    )

    $Action = "[BLD]"

    $LevelColor = @{ INF = 'Green'; DBG = 'Blue'; WRN = 'Yellow'; ERR = 'Red' }
    Write-Host $Action -ForegroundColor DarkGray -NoNewline
    Write-Host " $($Level.ToUpperInvariant()) " -ForegroundColor $LevelColor[$Level] -NoNewline
    Write-Host $Message
}

# ---- Entrypoint of Execution -----------------------------

Write-ConsoleLog -Level INF "Building the project..."

Write-ConsoleLog -Level INF "Bundling the assets..."
$AssetsDir = Join-Path $RootDir "src/include/assets"
$StaticDir = Join-Path $RootDir "src/web/static"
New-Item -ItemType Directory -Force -Path $AssetsDir | Out-Null
Copy-Item -Path (Join-Path $StaticDir '*') -Destination $AssetsDir -Recurse -Force


Write-ConsoleLog -Level INF "Building templates with a-h/templ"
Write-Host ""
& tools/templ.exe generate

exit 0