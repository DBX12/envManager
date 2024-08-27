package cmd

import (
	"envManager/environment"
	"envManager/helper"
	"envManager/secretsStorage"
	"github.com/spf13/cobra"
	"os"
	"slices"
	"strings"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load profiles",
	Long: `Load one or more profiles to this shell's environment.
If called without profiles, the directory mapping for the current working directory will be loaded.`,
	Run:               runLoad,
	ValidArgsFunction: CompleteProfiles,
	PreRun:            InitConfig,
}

func runLoad(_ *cobra.Command, args []string) {
	registry := secretsStorage.GetRegistry()
	env := environment.NewEnvironment()
	env.Load()

	if len(args) == 0 {
		workingDir, err := os.Getwd()
		cobra.CheckErr(err)
		if registry.HasDirectoryMapping(workingDir) {
			args, err = registry.GetDirectoryMapping(workingDir)
			cobra.CheckErr(err)
		} else {
			cobra.CheckErr("No profiles specified and no mapping for this path found")
		}
	}

	var profilesToLoad []string
	for _, name := range args {
		// get the profile from the registry
		profile, err := registry.GetProfile(name)
		cobra.CheckErr(err)

		if slices.Contains(profilesToLoad, name) {
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
