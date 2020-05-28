package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename specific object",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == args[1] {
			fmt.Printf("failed to rename: same name as before\n")
			os.Exit(1)
		}

		if _, err := bucket.CopyObject(args[0], args[1]); err != nil {
			fmt.Printf("failed to create new objcet: %v\n", err)
			os.Exit(1)
		}
		if err := bucket.DeleteObject(args[0]); err != nil {
			fmt.Printf("failed to delete origin object: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("succeed to rename\n")
	},
}
