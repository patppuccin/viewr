package cmd

import (
	"github.com/patppuccin/viewr/src/out"
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
		out.Logger.Info("Installing Viewr as a system service")
	},
}

var serviceUninstallCmd = &cobra.Command{
	Use:           "uninstall",
	Short:         helpServiceUninstallCmd,
	Long:          out.Banner(helpServiceUninstallCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		out.Logger.Info("Uninstalling the Viewr service")
	},
}

var serviceStartCmd = &cobra.Command{
	Use:           "start",
	Short:         helpServiceStartCmd,
	Long:          out.Banner(helpServiceStartCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		out.Logger.Info("Starting the Viewr service")
	},
}

var serviceStopCmd = &cobra.Command{
	Use:           "stop",
	Short:         helpServiceStopCmd,
	Long:          out.Banner(helpServiceStopCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		out.Logger.Info("Stopping the Viewr service")
	},
}

var serviceRestartCmd = &cobra.Command{
	Use:           "restart",
	Short:         helpServiceRestartCmd,
	Long:          out.Banner(helpServiceRestartCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		out.Logger.Info("Restarting the Viewr service")
	},
}

var serviceStatusCmd = &cobra.Command{
	Use:           "status",
	Short:         helpServiceStatusCmd,
	Long:          out.Banner(helpServiceStatusCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		out.Logger.Info("Checking the current status of the Viewr service")
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
