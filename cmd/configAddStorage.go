package cmd

import (
	"envManager/helper"
	"envManager/secretsStorage"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/errgo.v2/fmt/errors"
	"regexp"
	"strings"
)

var storageType string
var storageName string

// configAddStorageCmd represents the storage command
var configAddStorageCmd = &cobra.Command{
	Use:   "storage [type] [name]",
	Short: "Add a storage adapter to your configuration",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("type and name are required")
		}
		storageType = args[0]
		storageName = args[1]
		supportedTypes := secretsStorage.GetStorageAdapterTypes()
		if !helper.SliceStringContains(storageType, supportedTypes) {
			return errors.Newf("type %s is not known. Supported types: %s", storageType, strings.Join(supportedTypes, "\n"))
		}
		validName, err := regexp.MatchString("^[a-zA-Z0-9]+$", storageName)
		if err != nil {
			return err
		}
		if !validName {
			return errors.Newf("the name must only consist out of letters and numbers")
		}
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			// completing type
			return secretsStorage.GetStorageAdapterTypes(), cobra.ShellCompDirectiveNoFileComp
		}
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := rootCmd.PersistentFlags().GetString("config")
		cobra.CheckErr(err)

		config := secretsStorage.NewConfiguration()
		cobra.CheckErr(
			config.LoadFromFile(configPath),
		)

		_, storageExists := config.Storages[storageName]
		if storageExists && !flagForceConfig {
			cobra.CheckErr(
				errors.Newf(
					"the storage %s already exists in the config file %s. Set --force to overwrite it.",
					storageName,
					configPath,
				),
			)
		}
		defaultConfig, err := secretsStorage.GetStorageAdapterDefaultConfig(storageType)
		cobra.CheckErr(err)
		config.Storages[storageName] = secretsStorage.Storage{
			StorageType: storageType,
			Config:      defaultConfig,
		}

		cobra.CheckErr(
			config.WriteToFile(configPath, true),
		)
		fmt.Printf("Storage %s has been added to the configuration at %s\n", storageName, configPath)
	},
}

func init() {
	configAddCmd.AddCommand(configAddStorageCmd)
}
