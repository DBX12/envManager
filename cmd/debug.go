package cmd

import (
	"github.com/spf13/cobra"
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:              "debug",
	Short:            "Various commands for debugging your configuration",
	Long:             `With the debug command you can debug your configuration`,
	PersistentPreRun: InitConfig,
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
