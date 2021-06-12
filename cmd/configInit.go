package cmd

import (
	"envManager/secretsStorage"
	"fmt"
	"github.com/spf13/cobra"
)

var flagForce bool

// configInitCmd represents the init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes your config file",
	Long: "This command creates an empty config file for you. It will not overwrite an existing config file if not " +
		"forced to do so",
	Run: func(cmd *cobra.Command, args []string) {
		emptyConfig := secretsStorage.NewConfiguration()
		configPath, err := rootCmd.PersistentFlags().GetString("config")
		cobra.CheckErr(err)
		err = emptyConfig.WriteToFile(configPath, flagForce)
		cobra.CheckErr(err)
		fmt.Printf("Configuration initialized in %s\n", configPath)
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().BoolVarP(&flagForce, "force", "f", false, "Force initialization (will overwrite existing config)")
}
