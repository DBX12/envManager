package cmd

import (
	"envManager/environment"
	"envManager/helper"
	"envManager/secretsStorage"
	"github.com/spf13/cobra"
	"strings"
)

// unloadCmd represents the unload command
var unloadCmd = &cobra.Command{
	Use:               "unload",
	Short:             "Unload profiles",
	Long:              `Unload one or more profiles from this shell's environment`,
	Run:               runUnload,
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: CompleteProfiles,
	PreRun:            InitConfig,
}

func runUnload(_ *cobra.Command, args []string) {
	registry := secretsStorage.GetRegistry()
	env := environment.NewEnvironment()
	env.Load()
	loadedProfiles := strings.Split(env.GetCurrent(envManagerLoadedProfilesName, ""), ",")
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
}
