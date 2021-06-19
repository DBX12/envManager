package cmd

import (
	"envManager/secretsStorage"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/errgo.v2/fmt/errors"
)

// debugEntryCmd represents the entry command
var debugEntryCmd = &cobra.Command{
	Use:   "entry [storage] [path]",
	Short: "Shows attributes of an entry",
	Long:  "The debug entry command shows the attribute names of an entry specified by storage name and path",
	Args: func(cmd *cobra.Command, args []string) error {
		initConfig()
		if len(args) != 2 {
			return errors.New("Storage and path are required.")
		}
		if !secretsStorage.GetRegistry().HasStorage(args[0]) {
			return errors.Newf("The storage %s is not configured.", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		storageName := args[0]
		path := args[1]
		storagePtr, err := secretsStorage.GetRegistry().GetStorage(storageName)
		cobra.CheckErr(err)
		entry, err := (*storagePtr).GetEntry(path)
		cobra.CheckErr(err)
		fmt.Printf("Attributes of %s in %s:\n", path, storageName)
		attributeNames := entry.GetAttributeNames()
		for i := 0; i < len(attributeNames); i++ {
			fmt.Printf("- %s\n", attributeNames[i])
		}
		if (*storagePtr).IsCaseSensitive() {
			fmt.Println("This storage provider is case-sensitive!")
		}
	},
}

func init() {
	debugCmd.AddCommand(debugEntryCmd)
}
