package cmd

import (
	"envManager/environment"
	"envManager/helper"
	"envManager/secretsStorage"
	"github.com/spf13/cobra"
	"strings"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:               "load",
	Short:             "Load profiles",
	Long:              `Load one or more profiles to this shell's environment`,
	Run:               runLoad,
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: CompleteProfiles,
	PreRun:            InitConfig,
}

func runLoad(_ *cobra.Command, args []string) {
	registry := secretsStorage.GetRegistry()
	env := environment.NewEnvironment()
	env.Load()
	var profilesToLoad []string
	for _, name := range args {
		// get the profile from the registry
		profile, err := registry.GetProfile(name)
		cobra.CheckErr(err)

		if helper.SliceStringContains(name, profilesToLoad) {
			// this profile is already loaded, thus all its dependencies are
			// selected for loading too
			continue
		}
		// select the profile for loading
		profilesToLoad = append(profilesToLoad, name)
		// get the dependencies of this profile
		dependencies, err := profile.GetDependencies(profilesToLoad)
		cobra.CheckErr(err)
		// select the dependencies for loading too
		profilesToLoad = append(profilesToLoad, dependencies...)
	}

	// load every profile selected for loading
	for _, name := range profilesToLoad {
		profile, err := registry.GetProfile(name)
		cobra.CheckErr(err)
		err = profile.AddToEnvironment(&env)
		cobra.CheckErr(err)
	}
	loadedProfiles := strings.Split(env.GetCurrent(envManagerLoadedProfilesName, ""), ",")
	newEnvManagerLoadedValue := helper.SliceStringUnique(append(loadedProfiles, profilesToLoad...))
	newEnvManagerLoadedValue = helper.SliceStringRemove("", newEnvManagerLoadedValue)
	_ = env.Set(envManagerLoadedProfilesName, strings.Join(newEnvManagerLoadedValue, ","))
	print(env.WriteStatements())
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
