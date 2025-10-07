package include

import "embed"

//go:embed assets/*
var Assets embed.FS

//go:embed config/viewr-config.yaml
var DefaultConfig []byte
