package cmd

import (
	"os"
	"strings"

	"github.com/kardianos/service"
	"github.com/patppuccin/viewr/src/config"
	"github.com/patppuccin/viewr/src/helpers"
	"github.com/patppuccin/viewr/src/out"
	"github.com/patppuccin/viewr/src/server"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:           "service",
	Short:         helpServiceCmd,
	Long:          out.Banner(helpServiceCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var serviceInstallCmd = &cobra.Command{
	Use:           "install",
	Short:         helpServiceInstallCmd,
	Long:          out.Banner(helpServiceInstallCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := serviceFetch()
		if err != nil {
			out.Logger.Error(err.Error())
			os.Exit(1)
		}

		if err := svc.Install(); err != nil {
			out.Logger.Error(helpers.SafeErr("failed to install service", err).Error())
			os.Exit(1)
		} else {
			out.Logger.Info("Service installed successfully")
		}
	},
}

var serviceUninstallCmd = &cobra.Command{
	Use:           "uninstall",
	Short:         helpServiceUninstallCmd,
	Long:          out.Banner(helpServiceUninstallCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := serviceFetch()
		if err != nil {
			out.Logger.Error(err.Error())
			os.Exit(1)
		}

		status, _ := svc.Status()
		switch status {
		case service.StatusRunning:
			if err := svc.Stop(); err != nil {
				out.Logger.Error(helpers.SafeErr("failed to stop service", err).Error())
				os.Exit(1)
			}

		case service.StatusStopped:
		default:
			out.Logger.Warn("Service is neither running nor stopped, hence cannot be uninstalled")
			os.Exit(1)
		}

		out.Logger.Info("Uninstalling service...")
		if err := svc.Uninstall(); err != nil {
			out.Logger.Error(helpers.SafeErr("failed to uninstall service", err).Error())
			os.Exit(1)
		} else {
			out.Logger.Info("Service uninstalled successfully")
		}
	},
}

var serviceStartCmd = &cobra.Command{
	Use:           "start",
	Short:         helpServiceStartCmd,
	Long:          out.Banner(helpServiceStartCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := serviceFetch()
		if err != nil {
			out.Logger.Error(err.Error())
			os.Exit(1)
		}

		status, _ := svc.Status()
		switch status {
		case service.StatusRunning:
			out.Logger.Info("Service is already running")
			os.Exit(0)

		case service.StatusStopped:
			if err := svc.Start(); err != nil {
				out.Logger.Error(helpers.SafeErr("failed to start service", err).Error())
				os.Exit(1)
			}

		default:
			out.Logger.Warn("Service is neither running nor stopped, hence cannot be started")
			os.Exit(1)
		}

		out.Logger.Info("Service started successfully")
	},
}

var serviceStopCmd = &cobra.Command{
	Use:           "stop",
	Short:         helpServiceStopCmd,
	Long:          out.Banner(helpServiceStopCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := serviceFetch()
		if err != nil {
			out.Logger.Error(err.Error())
			os.Exit(1)
		}

		status, _ := svc.Status()
		switch status {
		case service.StatusRunning:
			if err := svc.Stop(); err != nil {
				out.Logger.Error(helpers.SafeErr("failed to stop service", err).Error())
				os.Exit(1)
			}

		case service.StatusStopped:
			out.Logger.Info("Service is already stopped")
			os.Exit(0)

		default:
			out.Logger.Warn("Service is neither running nor stopped, hence cannot be stopped")
			os.Exit(1)
		}

		out.Logger.Info("Service stopped successfully")
	},
}

var serviceRestartCmd = &cobra.Command{
	Use:           "restart",
	Short:         helpServiceRestartCmd,
	Long:          out.Banner(helpServiceRestartCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := serviceFetch()
		if err != nil {
			out.Logger.Error(err.Error())
			os.Exit(1)
		}

		status, _ := svc.Status()
		switch status {
		case service.StatusRunning:
			if err := svc.Restart(); err != nil {
				out.Logger.Error(helpers.SafeErr("failed to restart service", err).Error())
				os.Exit(1)
			}

		case service.StatusStopped:
			if err := svc.Start(); err != nil {
				out.Logger.Error(helpers.SafeErr("failed to start service", err).Error())
				os.Exit(1)
			}

		default:
			out.Logger.Warn("Service is neither running nor stopped, hence cannot be restarted")
			os.Exit(1)
		}

		out.Logger.Info("Service restarted successfully")
	},
}

var serviceStatusCmd = &cobra.Command{
	Use:           "status",
	Short:         helpServiceStatusCmd,
	Long:          out.Banner(helpServiceStatusCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := serviceFetch()
		if err != nil {
			out.Logger.Error(err.Error())
			os.Exit(1)
		}

		status, err := svc.Status()
		if err != nil {
			if strings.Contains(err.Error(), "does not exist") ||
				strings.Contains(err.Error(), "not installed") ||
				strings.Contains(err.Error(), "not found") {
				out.Logger.Warn("Service Status: NOT INSTALLED")
				os.Exit(3)
			}
			out.Logger.Error("Failed to determine service status")
			os.Exit(1)
		}

		switch status {
		case service.StatusRunning:
			out.Logger.Info("Service Status: RUNNING")
			os.Exit(0)

		case service.StatusStopped:
			out.Logger.Warn("Service Status: STOPPED")
			os.Exit(2)

		case service.StatusUnknown:
			out.Logger.Warn("Service Status: UNKNOWN")
			os.Exit(3)

		default:
			out.Logger.Warn("Service Status: UNEXPECTED")
			os.Exit(4)
		}
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.AddCommand(serviceInstallCmd)
	serviceCmd.AddCommand(serviceUninstallCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceRestartCmd)
	serviceCmd.AddCommand(serviceStatusCmd)
}

// Service-related helpers

var serviceFetch = func() (service.Service, error) {
	cfg := config.GlobalConfig
	if cfg == nil {
		out.Logger.Error("Failed to fetch service: configuration is not loaded")
		os.Exit(1)
	}
	svc, err := server.GetService(cfg.Server.Port, cfg.Server.Address, cfg.Server.LogLevel, false)
	if err != nil {
		out.Logger.Error(err.Error())
		os.Exit(1)
	}

	return svc, nil
}
