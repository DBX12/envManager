package cmd

import (
	"envManager/environment"
	"envManager/secretsStorage"
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:               "load",
	Short:             "Load profiles",
	Long:              `Load one or more profiles to this shell's environment`,
	Run:               runLoad,
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: CompleteProfiles,
}

func runLoad(_ *cobra.Command, args []string) {
	registry := secretsStorage.GetRegistry()
	env := environment.NewEnvironment()

	for i := 0; i < len(args); i++ {
		profile, err := registry.GetProfile(args[i])
		cobra.CheckErr(err)
		err = profile.AddToEnvironment(&env)
		cobra.CheckErr(err)
	}
	print(env.WriteStatements())
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
