package cmd

import (
	"envManager/environment"
	"fmt"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

var currentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"loaded"},
	Short:   "Shows currently loaded profiles",
	Run: func(cmd *cobra.Command, args []string) {
		env := environment.NewEnvironment()
		env.Load()
		envValue := env.GetCurrent(envManagerLoadedProfilesName, "")
		if len(envValue) == 0 {
			fmt.Print("No profiles loaded")
			return
		}
		loadedProfiles := strings.Split(envValue, ",")
		sort.Strings(loadedProfiles)
		fmt.Print(
			strings.Join(loadedProfiles, "\n"),
		)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
