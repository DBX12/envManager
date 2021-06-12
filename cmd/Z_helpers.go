package cmd

import (
	"envManager/helper"
	"envManager/secretsStorage"
	"github.com/spf13/cobra"
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

//InitConfig is a wrapper around the simple initConfig() method. With this adapter you can write
//PreRun: InitConfig, in your command object.
func InitConfig(_ *cobra.Command, _ []string) {
	initConfig()
}
