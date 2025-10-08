package server

import (
	"context"
	"os"

	"github.com/kardianos/service"
	"github.com/patppuccin/viewr/src/config"
	"github.com/patppuccin/viewr/src/constants"
	"github.com/patppuccin/viewr/src/out"
)

// Service declarations & definition
var svcConfig = &service.Config{
	Name:        constants.AppAbbrName,
	DisplayName: constants.AppFullName,
	Description: constants.AppDescription,
	Option: service.KeyValue{
		"LogOutput": true, // Enables logging to default service logs
	},
}

type Program struct {
	ctx          context.Context
	cancel       context.CancelFunc
	logLevel     string
	address      string
	port         int
	logToConsole bool
}

func (p *Program) Start(s service.Service) error {
	go func() {
		if err := Run(p.ctx, p.logLevel, p.address, p.port, p.logToConsole); err != nil {
			os.Stdout.WriteString("[service start error] " + err.Error() + "\n")
		}
	}()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	p.cancel()
	return nil
}

func NewProgram(port int, address, logLevel string, logToConsole bool) *Program {
	ctx, cancel := context.WithCancel(context.Background())
	return &Program{
		ctx:          ctx,
		cancel:       cancel,
		logLevel:     logLevel,
		address:      address,
		port:         port,
		logToConsole: logToConsole,
	}
}

func GetService(port int, address, logLevel string, logToConsole bool) (service.Service, error) {
	return service.New(NewProgram(port, address, logLevel, logToConsole), svcConfig)
}

func RunServerService() {
	cfg := config.GlobalConfig
	if cfg == nil {
		out.Logger.Error("Failed to fetch service: configuration is not loaded")
		os.Exit(1)
	}

	svc, err := GetService(cfg.Server.Port, cfg.Server.Address, cfg.Server.LogLevel, false)
	if err != nil {
		os.Stdout.WriteString("[service fetch error] " + err.Error() + "\n")
		os.Exit(1)
	}

	if err := svc.Run(); err != nil {
		os.Stdout.WriteString("[service run error] " + err.Error() + "\n")
		os.Exit(1)
	}
}
