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
	"github.com/patppuccin/viewr/src/models"
	"github.com/spf13/pflag"
	"go.yaml.in/yaml/v3"
)

var (
	GlobalConfig    models.AppConfig
	GlobalConfigSrc string
	GlobalConfigErr error
)

func Load(configFilePath string, flags *pflag.FlagSet) {

	// Default configuration — always valid
	GlobalConfig = models.AppConfig{
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

	// Try loading YAML file
	if data, err := os.ReadFile(absPath); err == nil {
		decoder := yaml.NewDecoder(bytes.NewReader(data))
		decoder.KnownFields(true)

		if err := decoder.Decode(&GlobalConfig); err != nil {
			GlobalConfigErr = helpers.SafeErr("error parsing YAML config", err)
		}

		GlobalConfigSrc = absPath
	} else if !os.IsNotExist(err) {
		GlobalConfigErr = helpers.SafeErr("error reading config file", err)
	}

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
