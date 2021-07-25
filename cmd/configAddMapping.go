package cmd

import (
	"envManager/environment"
	"envManager/helper"
	"envManager/secretsStorage"
	"fmt"
	"github.com/josa42/go-prompt/prompt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var flagAddMappingSelect bool

// configAddMappingCmd represents the mapping command
var configAddMappingCmd = &cobra.Command{
	Use:   "mapping",
	Short: "Add a directory mapping to your config",
	Long:  `A directory mapping links one or more profiles to the current directory.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := getConfigPath()
		config := secretsStorage.NewConfiguration()
		env := environment.NewEnvironment()
		env.Load()
		cobra.CheckErr(config.LoadFromFile(configPath))

		var profilesToMap []string

		if flagAddMappingSelect {
			options := prompt.Options{}
			for _, name := range secretsStorage.GetRegistry().GetProfileNames() {
				options = append(options, [2]string{name, name})
			}
			var canceled bool
			profilesToMap, canceled = prompt.MultiSelect("Map these profiles:", options)
			if canceled {
				fmt.Printf("Selection was cancelled, not mapping anything")
				return
			}
		} else {
			loadedProfiles := strings.Split(env.GetCurrent(envManagerLoadedProfilesName, ""), ",")
			profilesToMap = helper.SliceStringRemove("", loadedProfiles)
			if len(profilesToMap) == 0 {
				cobra.CheckErr("You have no profiles loaded at the moment. Either load some or call again with --select to select from your configured profiles.")
				return
			}
		}

		workingDir, err := os.Getwd()
		cobra.CheckErr(err)
		config.DirectoryMapping[workingDir] = profilesToMap

		cobra.CheckErr(
			config.WriteToFile(configPath, true),
		)
		fmt.Printf("Mapped the profiles (%s) to your current working directory.", strings.Join(profilesToMap, ", "))
	},
}

func init() {
	configAddCmd.AddCommand(configAddMappingCmd)
	configAddMappingCmd.Flags().BoolVarP(&flagAddMappingSelect, "select", "s", false, "Select profiles interactively")
}
