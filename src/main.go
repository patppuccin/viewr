package main

import (
	"github.com/kardianos/service"
	"github.com/patppuccin/viewr/src/cmd"
	"github.com/patppuccin/viewr/src/config"
)

func main() {
	if service.Interactive() {
		cmd.Execute()
	} else {
		config.Load("", nil)
		// TODO: Set up handling for service
	}
}
