package cmd

import (
	"envManager/secretsStorage"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envManager",
	Short: "Manage your environment variables",
	Long:  `A program to manage the environment variables of your shell, pulling the secrets from a secure secretsStorage like keepass.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Printf("%#v",args)
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only runLoad
	// when this action is called directly.
	//rootCmd.Flags().StringArray("load",nil, "List of profiles to load")
	rootCmd.Flags().Bool("verbose", false, "Enable verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//TODO replace hardcoded dev path with one from home dir
	config, err := secretsStorage.LoadConfigurationFromFile("/home/dbx12/GoLandProjects/envManager/data/envManager.yml")
	cobra.CheckErr(err)

	registry := secretsStorage.GetRegistry()
	for name, storageConfig := range config.Storages {
		adapter, err := secretsStorage.CreateStorageAdapter(name, storageConfig)
		cobra.CheckErr(err)
		err = registry.AddStorage(name, adapter)
		cobra.CheckErr(err)
	}

	for name, profile := range config.Profiles {
		err := registry.AddProfile(name, profile)
		cobra.CheckErr(err)
	}
}