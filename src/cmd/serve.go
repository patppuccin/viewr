package cmd

import (
	"strconv"
	"strings"

	"github.com/patppuccin/viewr/src/config"
	"github.com/patppuccin/viewr/src/constants"
	"github.com/patppuccin/viewr/src/out"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:           "run",
	Short:         helpRunCmd,
	Long:          out.Banner(helpRunCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		logLevel := config.GlobalConfig.Server.LogLevel
		port := config.GlobalConfig.Server.Port
		address := config.GlobalConfig.Server.Address
		out.Logger.Info("Running Viewr with log level: " + logLevel + " on port: " + strconv.Itoa(port) + " and address: " + address)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&flagRunLogLevel, "log-level", "l", "", "log levels: "+strings.Join(constants.LogLevels, ", "))
	runCmd.Flags().IntVarP(&flagRunPort, "port", "p", 0, "server port")
	runCmd.Flags().StringVarP(&flagRunAddress, "address", "a", "", "server address")
}
