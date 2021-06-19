package cmd

import (
	"envManager/helper"
	"envManager/secretsStorage"
	"fmt"
	"github.com/josa42/go-prompt/prompt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/errgo.v2/fmt/errors"
	"regexp"
)

var flagAddProfileDependencies bool
var flagAddProfileConstEnv bool
var flagAddProfileEnv bool

// configAddProfileCmd represents the profile command
var configAddProfileCmd = &cobra.Command{
	Use:   "profile [storageAdapter] [profileName]",
	Short: "Add a profile to your config",
	Args: func(cmd *cobra.Command, args []string) error {
		initConfig()
		if len(args) < 2 {
			return errors.Newf("storageAdapter and profileName are required")
		}
		storageAdapter := args[0]
		profileName := args[1]

		if !secretsStorage.GetRegistry().HasStorage(storageAdapter) {
			return errors.Newf("the storage adapter %s is not known. Please configure it first.", storageAdapter)
		}
		validName, err := regexp.MatchString("^[a-zA-Z0-9]+$", profileName)
		if err != nil {
			return err
		}
		if !validName {
			return errors.Newf("the name must only consist out of letters and numbers")
		}
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		initConfig()
		if len(args) == 0 {
			return secretsStorage.GetRegistry().GetStorageNames(), cobra.ShellCompDirectiveNoFileComp
		}
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		storageAdapter := args[0]
		profileName := args[1]

		configPath := getConfigPath()

		config := secretsStorage.NewConfiguration()
		cobra.CheckErr(
			config.LoadFromFile(configPath),
		)

		_, profileExists := config.Profiles[profileName]
		if profileExists && !flagForceConfig {
			cobra.CheckErr(
				errors.Newf(
					"the profile %s already exists in the config file %s. Set --force to overwrite it.",
					profileName,
					configPath,
				),
			)
		}

		needsPath := !flagAddProfileConstEnv || flagAddProfileEnv

		profile := secretsStorage.Profile{
			Env:       map[string]string{},
			ConstEnv:  map[string]string{},
			DependsOn: []string{},
			Storage:   storageAdapter,
		}

		if needsPath {
			var err error
			profile.Path, err = helper.GetInput().PromptString("Enter the path to the entry")
			cobra.CheckErr(err)
		}

		if flagAddProfileConstEnv {
			var err error
			profile.ConstEnv, err = promptForConstEnv()
			cobra.CheckErr(err)
		} else {
			fmt.Println("Not defining constant environment variables. Use --constEnv to configure them interactively.")
		}

		if flagAddProfileEnv {
			cobra.CheckErr(
				promptAndAddEnvMapping(&profile),
			)
		} else {
			fmt.Println("Not defining environment variable mapping. Use --env to configure it interactively.")
		}

		if flagAddProfileDependencies {

			profile.DependsOn = promptForDependencies(profileName)
		}

		config.Profiles[profileName] = profile
		cobra.CheckErr(
			config.WriteToFile(configPath, true),
		)
		fmt.Printf("Profile %s has been added to the configuration at %s\n", profileName, configPath)
	},
}

func promptForConstEnv() (map[string]string, error) {
	out := map[string]string{}
	fmt.Println("Define constant environment variables")

	keyPrompt := promptui.Prompt{
		Label: "Variable name",
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
		Label: "Variable value",
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

func promptAndAddEnvMapping(profile *secretsStorage.Profile) error {
	fmt.Println("Define environment mapping")

	keyPrompt := promptui.Prompt{
		Label: "Variable name",
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

	for {
		options, optionsErr := makeOptions(profile.Storage, profile.Path)
		if optionsErr != nil {
			return optionsErr
		}
		key, keyErr := keyPrompt.Run()
		if keyErr != nil {
			return keyErr
		}
		value, canceled := prompt.Select("Attribute", options)
		if !canceled {
			profile.Env[key] = value
		} else {
			fmt.Println("Last attribute selection was cancelled, not adding this one.")
		}
		if !promptYesNo("Add another one?") {
			break
		}
	}
	return nil
}

func makeOptions(storageName string, entryPath string) (prompt.Options, error) {
	storagePtr, err := secretsStorage.GetRegistry().GetStorage(storageName)
	if err != nil {
		return nil, err
	}
	entryPtr, err := (*storagePtr).GetEntry(entryPath)
	if err != nil {
		return nil, err
	}
	var options prompt.Options
	for _, attrName := range (*entryPtr).GetAttributeNames() {
		options = append(options, [2]string{attrName, attrName})
	}
	return options, nil
}

func promptForDependencies(exceptProfileName string) []string {
	options := prompt.Options{}
	for _, name := range secretsStorage.GetRegistry().GetProfileNames() {
		if name == exceptProfileName {
			continue
		}
		options = append(options, [2]string{name, name})
	}
	selection, _ := prompt.MultiSelect("This profile depends on:", options)
	return selection
}

func init() {
	configAddCmd.AddCommand(configAddProfileCmd)
	configAddProfileCmd.Flags().BoolVarP(&flagAddProfileDependencies, "dependencies", "d", false, "Select dependencies interactively")
	configAddProfileCmd.Flags().BoolVarP(&flagAddProfileConstEnv, "constEnv", "o", false, "Specify constant environment variables interactively")
	configAddProfileCmd.Flags().BoolVarP(&flagAddProfileEnv, "env", "e", false, "Specify environment variable mapping interactively")
}
