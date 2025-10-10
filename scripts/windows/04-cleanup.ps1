# ==== Setup script for Windows ============================

# ---- Global flags and params -----------------------------
[CmdletBinding()]
param ()

$ErrorActionPreference = 'Continue'

# ---- Helpers and utilities -------------------------------
function Write-ConsoleLog {
    param(
        [ValidateSet("INF", "DBG", "WRN", "ERR")][string]$Level = "INF",
        [Parameter(Mandatory)][string]$Message
    )

    $Action = "[CLN]"

    $LevelColor = @{ INF = 'Green'; DBG = 'Blue'; WRN = 'Yellow'; ERR = 'Red' }
    Write-Host $Action -ForegroundColor DarkGray -NoNewline
    Write-Host " $($Level.ToUpperInvariant()) " -ForegroundColor $LevelColor[$Level] -NoNewline
    Write-Host $Message
}

# ---- Entrypoint of Execution -----------------------------

Write-ConsoleLog -Level INF "Cleaning up..."
exit 0