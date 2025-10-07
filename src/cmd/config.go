package cmd

import (
	"github.com/patppuccin/viewr/src/out"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:           "config",
	Short:         helpConfigCmd,
	Long:          out.Banner(helpConfigCmd),
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		if flagConfigInit != "" {
			// TODO: Handle the initialization of the configuration file with overwrites
			out.Logger.Info("Initializing the configuration file: " + flagConfigInit)
			return
		}

		if flagConfigValidate {
			// TODO: Validate the configuration file
			out.Logger.Info("Validating the configuration file")
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&flagConfigInit, "init", "i", "", "initialize the configuration file")
	configCmd.Flags().BoolVarP(&flagConfigValidate, "validate", "v", false, "validate the configuration file")
	configCmd.Flags().BoolVarP(&flagConfigOverwrite, "overwrite", "o", false, "overwrite the configuration file")
}
