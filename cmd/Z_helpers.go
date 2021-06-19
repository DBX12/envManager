package cmd

import (
	"envManager/helper"
	"envManager/secretsStorage"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/errgo.v2/fmt/errors"
	"regexp"
)

//CompleteProfiles provides completion for a command which expects at least one
//profile.
func CompleteProfiles(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	initConfig()
	// complete profiles
	completions := helper.Completion(
		// all profiles of this storage adapter
		secretsStorage.GetRegistry().GetProfileNames(),
		// all already specified profiles
		args,
		toComplete,
	)
	return completions, cobra.ShellCompDirectiveNoFileComp
}

//CompleteStorages provides completion for a command which expects at least one
//storage.
func CompleteStorages(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	initConfig()
	// complete profiles
	completions := helper.Completion(
		// all profiles of this storage adapter
		secretsStorage.GetRegistry().GetStorageNames(),
		// all already specified profiles
		args,
		toComplete,
	)
	return completions, cobra.ShellCompDirectiveNoFileComp
}

//InitConfig is a wrapper around the simple initConfig() method. With this adapter you can write
//PreRun: InitConfig, in your command object.
func InitConfig(_ *cobra.Command, _ []string) {
	initConfig()
}

//getConfigPath gets the config path from the persistent flags of the rootCmd
func getConfigPath() string {
	configPath, err := rootCmd.PersistentFlags().GetString("config")
	cobra.CheckErr(err)
	return configPath
}

//promptForEnvVariables shows a prompt for environment variables (the keys) and their values or the mapping (the values)
func promptForEnvVariables(prompt string, labelKey string, labelValue string) (map[string]string, error) {
	out := map[string]string{}
	fmt.Println(prompt)

	keyPrompt := promptui.Prompt{
		Label: labelKey,
		Validate: func(input string) error {
			match, err := regexp.MatchString("^[a-zA-Z0-9_]+$", input)
			if err != nil {
				return err
			}
			if !match {
				return errors.New("Only letters, numbers and underscore is allowed")
			}
			return nil
		},
	}
	valuePrompt := promptui.Prompt{
		Label: labelValue,
	}
	for {
		key, keyErr := keyPrompt.Run()
		if keyErr != nil {
			return nil, keyErr
		}
		value, valueErr := valuePrompt.Run()
		if valueErr != nil {
			return nil, valueErr
		}

		out[key] = value
		if !promptYesNo("Add another one?") {
			break
		}
	}
	return out, nil
}

//promptYesNo shows a prompt for a yes / no question. The (Y|N) is added to the prompt automatically.
func promptYesNo(prompt string) bool {
	yesNoPrompt := promptui.Prompt{
		Label: fmt.Sprintf("%s (Y|N)", prompt),
		Validate: func(input string) error {
			match, err := regexp.MatchString("^[YyNn]$", input)
			if err != nil {
				return err
			}
			if !match {
				return errors.New("")
			}
			return nil
		},
	}
	response, err := yesNoPrompt.Run()
	if err != nil {
		return false
	}
	return response == "Y" || response == "y"
}
