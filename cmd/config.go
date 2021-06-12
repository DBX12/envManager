package cmd

import (
	"github.com/spf13/cobra"
)

//flagForceConfig represents the status of the boolean --force flag
var flagForceConfig bool

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage your configuration",
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().BoolVarP(
		&flagForceConfig,
		"force",
		"f",
		false,
		"Force operation (will overwrite existing data)",
	)
}
