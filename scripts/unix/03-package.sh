#!/usr/bin/env bash
# ==== Setup Script for Unix/Linux ==========================

# ---- Global flags and params ------------------------------

# ---- Helpers and utilities --------------------------------
log() {
  local level="${1:-INF}"
  shift
  local message="$*"

  local action="\033[38;5;240m[PKG]\033[0m"  # Gray
  local color reset="\033[0m"

  case "$level" in
    DBG) color="\033[34m" ;; # Blue
    INF) color="\033[32m" ;; # Green
    WRN) color="\033[33m" ;; # Yellow
    ERR) color="\033[31m" ;; # Red
    *)   color="\033[0m"  ;; # Default
  esac

  printf "%b %b%s%b %s\n" "$action" "$color" "$level" "$reset" "$message"
}

# ---- Entrypoint of execution ------------------------------
log INF "Packaging the application for release..."
exit 0
