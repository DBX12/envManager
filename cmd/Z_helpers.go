package cmd

import (
	"envManager/environment"
	"envManager/helper"
	"envManager/secretsStorage"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/errgo.v2/fmt/errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//CompleteProfiles provides completion for a command which expects at least one
//profile.
func CompleteProfiles(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	initConfig()

	var possibleValues, excludedValues []string
	switch cmd.Use {
	// Special behavior for envManager unload: Suggest loaded profiles
	case "unload":
		env := environment.NewEnvironment()
		env.Load()
		possibleValues = strings.Split(env.GetCurrent(envManagerLoadedProfilesName, ""), ",")
		excludedValues = args
	default:
		possibleValues = secretsStorage.GetRegistry().GetProfileNames()
		excludedValues = args
	}

	// complete profiles
	completions := helper.Completion(
		possibleValues,
		excludedValues,
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

// formatList formats an input slice by adding the prefix and suffix to every item. If the slice
// is empty, the emptyPlaceholder will be returned without added prefix or suffix
func formatList(items []string, prefix string, suffix string, emptyPlaceholder string) string {
	if len(items) == 0 {
		return emptyPlaceholder
	}
	output := make([]string, len(items))
	for i, item := range items {
		output[i] = prefix + item + suffix
	}
	return strings.Join(output, "")
}

//discoverConfigFiles traverses from startDir to the file system root and list
//all envManager config files. If it encounters the mainConfigFile, it will not
//add the file to the list of discovered files again (since it is already added
//as first file). The file with the highest precedence comes last, which means
//that the main config file has the lowest precedence.
func discoverConfigFiles(startDir string, mainConfigFile string) []string {
	var configFiles []string
	dir := startDir
	// traverse up to the root directory
	for dir != "/" {
		cfgFile := filepath.Join(dir, ".envManager.yml")

		// do not add the main config file to configPaths, will be done at the end
		if cfgFile != mainConfigFile {
			_, statErr := os.Stat(cfgFile)
			if !os.IsNotExist(statErr) {
				// note the path if it exists
				configFiles = append(configFiles, cfgFile)
			}
		}
		dir = filepath.Dir(dir)
	}

	// add the main config file as final file, so it becomes the first (least precedence) after reversing the slice
	configFiles = append(configFiles, mainConfigFile)
	return helper.SliceStringReverse(configFiles)
}
