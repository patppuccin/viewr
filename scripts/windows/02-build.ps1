# ==== Setup script for Windows ============================

# ---- Global flags and params -----------------------------
[CmdletBinding()]
param ()

$ErrorActionPreference = 'Continue'
$RootDir = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
$env:TAILWIND_DISABLE_WARNING = "true"

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

# Step 01: Build TailwindCSS
Write-ConsoleLog INF "Building TailwindCSS..."
$TWExec = Join-Path $RootDir "tools\tailwind.exe"
$TWEntrypointPath = Join-Path $RootDir ".\src\web\app.css"
$TWBuildPath = Join-Path $RootDir ".\src\include\assets\styles\global.min.css"
if (Test-Path $TWExec -PathType Leaf -ErrorAction SilentlyContinue) {
    Write-Host ""
    & $TWExec -i $TWEntrypointPath -o $TWBuildPath --minify
    if ($LASTEXITCODE -ne 0) {
        Write-ConsoleLog -Level ERR "TailwindCSS build failed"
        exit $LASTEXITCODE
    }
    Write-Host ""
}
else {
    Write-ConsoleLog -Level ERR "TailwindCSS binary not found at $Tailwind"
    exit 1
}

# Step 02: Build templates with a-h/templ
Write-ConsoleLog -Level INF "Building templates with a-h/templ"
Write-Host ""
& tools/templ.exe generate

exit 0