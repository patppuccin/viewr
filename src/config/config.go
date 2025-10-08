package config

import (
	"bytes"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/patppuccin/viewr/src/helpers"
	"github.com/patppuccin/viewr/src/include"
	"github.com/patppuccin/viewr/src/models"
	"github.com/spf13/pflag"
	"go.yaml.in/yaml/v3"
)

var (
	GlobalConfig    *models.AppConfig
	GlobalConfigSrc string
	GlobalConfigErr error
)

func Load(configFilePath string, flags *pflag.FlagSet) {

	// Default configuration — always valid
	GlobalConfig = &models.AppConfig{
		Server: models.ServerConfig{
			LogLevel: "info",
			Port:     5567,
			Address:  "127.0.0.1",
		},
		Paths: []models.PathConfig{},
	}
	GlobalConfigSrc = "defaults"
	GlobalConfigErr = nil

	// Resolve config path
	if configFilePath == "" {
		if rootPath, err := helpers.GetRootPath(); err == nil {
			configFilePath = filepath.Join(rootPath, "viewr-config.yaml")
		} else {
			GlobalConfigErr = helpers.SafeErr("unable to resolve the application root path", err)
		}
	}

	absPath, err := filepath.Abs(configFilePath)
	if err != nil {
		GlobalConfigErr = helpers.SafeErr("invalid config path - "+configFilePath, err)
	}

	cfg, cfgSrc, err := readYAMLConfig(absPath)
	if err != nil {
		GlobalConfigErr = err
		return
	}
	GlobalConfig = &cfg
	GlobalConfigSrc = cfgSrc

	// Track & handle per-field Overrides
	configOverrides := map[string]string{}
	setConfigOverride := func(field, src string) { configOverrides[field] = src }

	// Parse overrides from ENV Variables
	if val := strings.ToLower(os.Getenv("VIEWR_LOG_LEVEL")); val != "" {
		if helpers.IsValidLogLevel(val) {
			GlobalConfig.Server.LogLevel = val
			setConfigOverride("log-level", "env:VIEWR_LOG_LEVEL")
		}
	}

	if val := os.Getenv("VIEWR_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil && helpers.IsValidPort(port) {
			GlobalConfig.Server.Port = port
			setConfigOverride("port", "env:VIEWR_PORT")
		}
	}

	if val := os.Getenv("VIEWR_ADDRESS"); val != "" {
		if helpers.IsValidAddress(val) {
			GlobalConfig.Server.Address = val
			setConfigOverride("address", "env:VIEWR_ADDRESS")
		}
	}

	// Parse overrides from CLI Flags
	if flags != nil {
		if flags.Changed("log-level") {
			if val, _ := flags.GetString("log-level"); helpers.IsValidLogLevel(val) {
				GlobalConfig.Server.LogLevel = val
				setConfigOverride("log-level", "flag:log-level")
			}
		}

		if flags.Changed("port") {
			if port, _ := flags.GetInt("port"); helpers.IsValidPort(port) {
				GlobalConfig.Server.Port = port
				setConfigOverride("port", "flag:port")
			}
		}

		if flags.Changed("address") {
			if addr, _ := flags.GetString("address"); helpers.IsValidAddress(addr) {
				GlobalConfig.Server.Address = addr
				setConfigOverride("address", "flag:address")
			}
		}
	}

	// Summarize overrides deterministically
	if len(configOverrides) > 0 {
		var overrides []string
		fields := slices.Collect(maps.Keys(configOverrides))
		slices.Sort(fields)
		for _, field := range fields {
			overrides = append(overrides, field+"="+configOverrides[field])
		}
		GlobalConfigSrc += " (overrides: " + strings.Join(overrides, ", ") + ")"
	}
}

func Validate(configFilePath string) (string, error) {
	if configFilePath == "" {
		if rootPath, err := helpers.GetRootPath(); err == nil {
			configFilePath = filepath.Join(rootPath, "viewr-config.yaml")
		} else {
			return "", helpers.SafeErr("unable to resolve the application root path", err)
		}
	}

	absPath, err := filepath.Abs(configFilePath)
	if err != nil {
		return "", helpers.SafeErr("invalid config path - "+configFilePath, err)
	}

	_, cfgSrc, err := readYAMLConfig(absPath)
	if err != nil {
		return "", err
	}

	return cfgSrc, nil
}

func ExportTemplate(destPath string, overwrite bool) (string, error) {

	// Handle default destination path
	if destPath == "" {
		if rootPath, err := helpers.GetRootPath(); err == nil {
			destPath = filepath.Join(rootPath, "viewr-config-template.yaml")
		} else {
			return "", helpers.SafeErr("unable to resolve the application root path", err)
		}
	}

	// Handle overwrite protection
	if helpers.DoesYAMLFileExist(destPath) && !overwrite {
		return "", helpers.SafeErr("file already exists at "+destPath+" (overwrite disabled)", nil)
	}

	// Ensure parent directories exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return "", helpers.SafeErr("failed to create parent directories to "+destPath, err)
	}

	// Write default config
	if err := os.WriteFile(destPath, include.DefaultConfig, 0644); err != nil {
		return "", helpers.SafeErr("error writing to dest template file at "+destPath, err)
	}

	// Return destination path & no error
	return destPath, nil

}

// Local helpers

func readYAMLConfig(configFilePath string) (models.AppConfig, string, error) {
	var cfg models.AppConfig

	// Resolve to absolute path
	absPath, err := filepath.Abs(configFilePath)
	if err != nil {
		return cfg, "", helpers.SafeErr("invalid config path: "+configFilePath, err)
	}

	// Verify file existence before reading
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Graceful fallback: config file missing, return empty config
			return cfg, "", nil
		}
		return cfg, "", helpers.SafeErr("failed to access config file", err)
	}
	if info.IsDir() {
		return cfg, "", helpers.SafeErr("config path points to a directory", nil)
	}

	// Read file contents
	data, err := os.ReadFile(absPath)
	if err != nil {
		return cfg, "", helpers.SafeErr("failed to read config file", err)
	}

	// Strict YAML decoding
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, "", helpers.SafeErr("failed to parse YAML config", err)
	}

	return cfg, absPath, nil
}
