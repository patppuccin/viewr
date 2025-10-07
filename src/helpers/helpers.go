package helpers

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/patppuccin/viewr/src/constants"
)

var (
	execPath string
	execDir  string
	execErr  error
)

func init() {
	p, err := os.Executable()
	if err != nil {
		execErr = fmt.Errorf("failed to resolve current executable path: %w", err)
		return
	}
	execPath = p
	execDir = filepath.Dir(p)
}

func IsDevMode() bool { return strings.Contains(strings.ToLower(execPath), "go-build") }

func GetRootPath() (string, error) {

	// Handle exec path fetch error
	if execErr != nil {
		return "", execErr
	}

	// Handle go run temp path scenarios (contains "go-build")
	if IsDevMode() {
		return os.Getwd()
	}

	// Return resolved root path (parent directory of executable)
	return execDir, nil
}

func SafeErr(userMsg string, internalErr error) error {
	if IsDevMode() && internalErr != nil {
		return fmt.Errorf("%s: %w", userMsg, internalErr)
	}
	return fmt.Errorf("%s", userMsg)
}

func IsValidLogLevel(level string) bool {
	return slices.Contains(constants.LogLevels, level)
}

func IsValidPort(port int) bool {
	// Check for valid port range + common web ports
	return (port >= 1024 && port <= 65535) || port == 80 || port == 443
}

func IsValidAddress(addr string) bool {
	// Empty check
	if addr == "" {
		return false
	}

	// Basic trim + sanity check
	addr = strings.TrimSpace(addr)
	if strings.ContainsAny(addr, " \t\n") {
		return false
	}

	// Try parsing as IP
	if ip := net.ParseIP(addr); ip != nil {
		return true
	}

	// Fallback: simple hostname check (no spaces, dots allowed)
	if len(addr) > 255 {
		return false
	}
	return true
}
