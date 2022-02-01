package cmd

import (
	"envManager/environment"
	"envManager/helper"
	"envManager/secretsStorage"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// unloadCmd represents the unload command
var unloadCmd = &cobra.Command{
	Use:               "unload",
	Short:             "Unload profiles",
	Long:              `Unload one or more profiles from this shell's environment`,
	Run:               runUnload,
	ValidArgsFunction: CompleteProfiles,
	PreRun:            InitConfig,
}

func runUnload(cmd *cobra.Command, args []string) {
	registry := secretsStorage.GetRegistry()
	env := environment.NewEnvironment()
	env.Load()
	loadedProfiles := strings.Split(env.GetCurrent(envManagerLoadedProfilesName, ""), ",")
	unloadAllFlagSet, _ := cmd.Flags().GetBool("all")
	if unloadAllFlagSet {
		// pretend all loaded profiles were listed as args
		args = loadedProfiles
	} else if len(args) == 0 {
		// no --all flag and no profile name specified
		fmt.Println("You must specify at least one profile to unload")
	}
	for i := 0; i < len(args); i++ {
		profile, err := registry.GetProfile(args[i])
		cobra.CheckErr(err)
		err = profile.RemoveFromEnvironment(&env)
		cobra.CheckErr(err)
		loadedProfiles = helper.SliceStringRemove(args[i], loadedProfiles)
	}
	loadedProfiles = helper.SliceStringRemove("", loadedProfiles)
	_ = env.Set(envManagerLoadedProfilesName, strings.Join(loadedProfiles, ","))
	print(env.WriteStatements())
}

func init() {
	rootCmd.AddCommand(unloadCmd)
	unloadCmd.Flags().Bool("all", false, "Select all currently loaded profiles for unloading")
}
