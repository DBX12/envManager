package cmd

import (
	"envManager/secretsStorage"
	"fmt"
	"github.com/spf13/cobra"
)

// debugProfileCmd represents the profile command
var debugProfileCmd = &cobra.Command{
	Use:   "profile [name]",
	Short: "Shows configuration of a profile",
	Long:  `The debug profile command shows the configuration of an entry specified by its name`,
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		registry := secretsStorage.GetRegistry()
		profileName := args[0]
		profile, err := registry.GetProfile(profileName)
		cobra.CheckErr(err)
		profileDependencies, err := profile.GetDependencies([]string{profileName})
		cobra.CheckErr(err)
		fmt.Printf(
			"Profile: %s\nStorage adapter: %s\nStorage adapter exists: %t\nPath in adapter: %s\n",
			profileName,
			profile.Storage,
			registry.HasStorage(profile.Storage),
			profile.Path,
		)
		fmt.Printf("Profile depends on other profiles: %t\n", len(profileDependencies) > 0)
		if len(profileDependencies) > 0 {
			for _, dependency := range profileDependencies {
				fmt.Printf(" - %s\n", dependency)
			}
		}
		fmt.Printf("Provides static environment variables: %t\n", len(profile.ConstEnv) > 0)
		if len(profile.ConstEnv) > 0 {
			debugProfilePrintEnv(profile.ConstEnv)
		}

		fmt.Printf("Provides dynamic environment variables: %t\n", len(profile.Env) > 0)
		if len(profile.Env) > 0 {
			debugProfilePrintEnv(profile.Env)
		}
	},
}

func debugProfilePrintEnv(in map[string]string) {
	for key, value := range in {
		fmt.Printf(" %s : %s\n", key, value)
	}
}

func init() {
	debugCmd.AddCommand(debugProfileCmd)
}
