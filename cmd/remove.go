package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove specific object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		objectName := args[0]

		exists, err := bucket.IsObjectExist(objectName)
		if err != nil {
			fmt.Printf("failed to query object: %v\n", err)
			os.Exit(1)
		}
		if !exists {
			fmt.Printf("object not exists\n")
			os.Exit(1)
		}

		err = bucket.DeleteObject(objectName)
		if err != nil {
			fmt.Printf("failed to remove: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("succeed to remove\n")
	},
}
