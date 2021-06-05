package cmd

import (
	"envManager/helper"
	"envManager/secretsStorage"
	"github.com/spf13/cobra"
)

//CompleteProfiles provides completion for a command which expects at least one
//profile.
func CompleteProfiles(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
