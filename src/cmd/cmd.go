package cmd

import (
	"os"

	"github.com/patppuccin/viewr/src/config"
	"github.com/patppuccin/viewr/src/constants"
	"github.com/patppuccin/viewr/src/out"
	"github.com/spf13/cobra"
)

const (
	helpRootCmd             = "Manage the Viewr application"
	helpRunCmd              = "Run the Viewr application on the console"
	helpConfigCmd           = "Manage the Viewr configuration"
	helpServiceCmd          = "Manage the Viewr service"
	helpServiceInstallCmd   = "Install Viewr as a system service"
	helpServiceUninstallCmd = "Uninstall the Viewr service"
	helpServiceStartCmd     = "Start the Viewr service"
	helpServiceStopCmd      = "Stop the Viewr service"
	helpServiceRestartCmd   = "Restart the Viewr service"
	helpServiceStatusCmd    = "Check the current status of the Viewr service"
)

var (
	flagConfigInit      bool
	flagConfigValidate  bool
	flagConfigOverwrite bool
	flagRunLogLevel     string
	flagRunPort         int
	flagRunAddress      string
)

var rootCmd = &cobra.Command{
	Use:           constants.AppAbbrName,
	Short:         helpRootCmd,
	Long:          out.Banner(helpRootCmd),
	Version:       constants.AppVersion,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {

	// Root Command Configurations
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		out.Logger.Error("Flag parse error: " + err.Error())
		os.Exit(2) // Exit code 2 indicates a command-line flag error
		return nil // Unreachable, but the compiler requires it
	})

	// Root Command Flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to configuration file")

	// Initialize and Load the configuration
	cobra.OnInitialize(func() {
		configPath, _ := rootCmd.PersistentFlags().GetString("config")
		config.Load(configPath, runCmd.Flags())
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Stdout.WriteString("Error: " + err.Error() + "\n")
		os.Exit(1)
	}
}
