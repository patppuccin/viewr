package models

import "github.com/rs/zerolog"

type AppConfig struct {
	Server ServerConfig `yaml:"server"`
	Paths  []PathConfig `yaml:"paths"`
}

type ServerConfig struct {
	LogLevel string `yaml:"logLevel"`
	Port     int    `yaml:"port"`
	Address  string `yaml:"address"`
}

type PathConfig struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	Disable bool   `yaml:"disable"`
}

type AppContext struct {
	Config *AppConfig
	// DBConn *DBConn
	Logger *zerolog.Logger
}
