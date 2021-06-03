package cmd

import (
	"envManager/secretsStorage"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// debugStorageCmd represents the storage command
var debugStorageCmd = &cobra.Command{
	Use:   "storage [name]",
	Short: "Shows the configuration of the storages",
	Long: `Shows which storage adapters are configured. If called without a name,
a list of all storages is shown.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// show configured storages
			fmt.Println("Configured storages:")
			storages := secretsStorage.GetRegistry().GetAllStorages()
			for name, _ := range storages {
				fmt.Printf("- %s\n", name)
			}
			return
		}

		// show details for one storage
		storageName := args[0]
		storagePtr, err := secretsStorage.GetRegistry().GetStorage(storageName)
		fmt.Printf(
			"Storage %s\nIs configured: %t\n",
			storageName,
			err == nil,
		)
		if err != nil {
			return
		}
		fmt.Println("Running storage dependent checks:")
		_, checks := (*storagePtr).Validate()
		fmt.Println(strings.Join(checks, "\n"))
	},
}

func init() {
	debugCmd.AddCommand(debugStorageCmd)
}
