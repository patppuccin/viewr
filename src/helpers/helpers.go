package helpers

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strconv"
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
		execErr = errors.Join(errors.New("failed to resolve current executable path"), err)
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
		return errors.Join(errors.New(userMsg), internalErr)
	}
	return errors.New(userMsg)
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

func CheckTCPBind(addr string, port int) error {
	var errs []string

	// Validate address
	if !IsValidAddress(addr) {
		errs = append(errs, "invalid address: "+addr+" (must be valid IP or hostname)")
	}

	// Validate port
	if !IsValidPort(port) {
		errs = append(errs, "invalid port: "+strconv.Itoa(port)+" (must be between 1–65535)")
	}

	// Combine validation errors if any
	if len(errs) > 0 {
		return SafeErr("invalid bind parameters: "+strings.Join(errs, "; "), nil)
	}

	// Try binding to address
	hostPort := net.JoinHostPort(addr, strconv.Itoa(port))
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		return SafeErr("failed to bind to "+hostPort, err)
	}
	defer func() {
		_ = listener.Close()
	}()

	return nil
}

func DoesYAMLFileExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}
