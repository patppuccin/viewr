package models

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
