package cmd

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/cobra"
	"os"
)

var unshareCmd = &cobra.Command{
	Use:   "unshare",
	Short: "unshare specific object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := bucket.SetObjectACL(args[0], oss.ACLPrivate)
		if err != nil {
			fmt.Printf("failed to unshare: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("succeed to unshare\n")
	},
}
