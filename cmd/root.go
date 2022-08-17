package cmd

import (
	"envManager/secretsStorage"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var flagConfigFile string

// The name of the environment variable containing the loaded profile names
const envManagerLoadedProfilesName = "ENVMANAGER_LOADED"

var version = "unknown"

var homeDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "envManager",
	Short:   "Manage your environment variables",
	Long:    `A program to manage the environment variables of your shell, pulling the secrets from a secure secretsStorage like keepass.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	cobra.CheckErr(err)
	configPath := path.Join(homeDir, ".envManager.yml")
	rootCmd.PersistentFlags().StringVarP(
		&flagConfigFile,
		"config",
		"c",
		configPath,
		"Overrides the default config file to use.",
	)
	_ = rootCmd.MarkPersistentFlagFilename("config", "yml")
	_ = rootCmd.PersistentFlags().MarkDeprecated("config", "since the introduction of directory-aware loading.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config := secretsStorage.NewConfiguration()
	err := config.LoadFromFile(flagConfigFile)
	cobra.CheckErr(err)

	//region location aware config loading
	dir, err := os.Getwd()
	cobra.CheckErr(err)

	// use helper function to find all config files upward from here
	configFiles := discoverConfigFiles(dir, flagConfigFile)

	for i, configFile := range configFiles {
		if configFile == flagConfigFile {
			// do not merge the main config file as it was loaded with config.LoadFromFile()
			continue
		}
		err := config.MergeConfigFile(configFiles[i])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to merge configuration.\nAlready processed config files:")
			_, _ = fmt.Fprint(os.Stderr, formatList(configFiles[0:i], "\t- ", "\n", "\t<none>\n"))
			_, _ = fmt.Fprintf(os.Stderr, "Offending config file:\n\t%s\n", configFile)
			_, _ = fmt.Fprintln(os.Stderr, "Still to process:")
			_, _ = fmt.Fprint(os.Stderr, formatList(configFiles[i+1:], "\t- ", "\n", "\t<none>\n"))
			_, _ = fmt.Fprintf(os.Stderr, "Error message:\n\t%s\n", err.Error())
			os.Exit(1)
		}
	}
	//endregion

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

	for mappingPath, profiles := range config.DirectoryMapping {
		err := registry.AddDirectoryMapping(mappingPath, profiles)
		cobra.CheckErr(err)
	}
}
