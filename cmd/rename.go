package cmd

import "github.com/spf13/cobra"

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename specific object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
