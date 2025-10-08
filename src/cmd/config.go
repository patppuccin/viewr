package cmd

import (
	"github.com/patppuccin/viewr/src/config"
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
		if flagConfigInit {
			destPath, err := config.ExportTemplate("", flagConfigOverwrite)
			if err != nil {
				out.Logger.Error("Failed initialization - " + err.Error())
				return
			}
			out.Logger.Info("Configuration file initialized at: " + destPath)
			return
		}

		if flagConfigValidate {
			cfgSrc, err := config.Validate("")
			if err != nil {
				out.Logger.Error("Failed to validate the configuration file: " + err.Error())
				return
			}
			out.Logger.Info("Configuration Source: " + cfgSrc)
			out.Logger.Info("Configuration is valid")
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&flagConfigInit, "init", "i", false, "initialize the configuration file")
	configCmd.Flags().BoolVarP(&flagConfigValidate, "validate", "v", false, "validate the configuration file")
	configCmd.Flags().BoolVarP(&flagConfigOverwrite, "overwrite", "o", false, "overwrite the configuration file")
}
