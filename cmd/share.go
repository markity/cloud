package cmd

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/cobra"
	"net/url"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "share specific object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := bucket.SetObjectACL(args[0], oss.ACLPublicRead)
		if err != nil {
			fmt.Printf("failed to share:%v\n", err)
			return
		}

		fmt.Printf("succeed to share, path: %v\n",
			fmt.Sprintf("https://%v.%v/%v", bucketName, endpoint, url.PathEscape(args[1])))
	},
}
