package cmd

import (
	"github.com/spf13/cobra"
)

// configAddCmd represents the add command
var configAddCmd = &cobra.Command{
	Use:              "add",
	Short:            "Add storages and profiles to your configuration",
	PersistentPreRun: InitConfig,
}

func init() {
	configCmd.AddCommand(configAddCmd)
}
