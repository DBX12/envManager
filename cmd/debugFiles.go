package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// debugFilesCmd represents the debug files command
var debugFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Shows which config files will be loaded in which order",
	// this empty PersistentPreRun function is only set to overwrite the inherited PersistedPreRun since loading the
	// config and merging files is not needed and will actually make this command less useful if there are collisions
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		cobra.CheckErr(err)

		_, _ = fmt.Fprintln(os.Stderr, "These files will be processed in this order (later files override earlier files):")
		_, _ = fmt.Fprint(
			os.Stderr,
			formatList(
				discoverConfigFiles(dir, flagConfigFile),
				"\t- ",
				"\n",
				"",
			),
		)
	},
}

func init() {
	debugCmd.AddCommand(debugFilesCmd)
}
