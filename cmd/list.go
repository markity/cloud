package cmd

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"time"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all objects",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		objs := make([]oss.ObjectProperties, 0)
		marker := ""
		for {
			lsRes, err := bucket.ListObjects(oss.Marker(marker))
			if err != nil {
				fmt.Printf("failed to list objects: %v\n", err)
				os.Exit(1)
			}

			for _, obj := range lsRes.Objects {
				objs = append(objs, obj)
			}

			if lsRes.IsTruncated {
				marker = lsRes.NextMarker
			} else {
				break
			}
		}

		objAcls := make([]string, 0, len(objs))
		for _, obj := range objs {
			acl, err := bucket.GetObjectACL(obj.Key)
			if err != nil {
				fmt.Printf("failed to get object ACL: %v\n", err)
				os.Exit(1)
			}
			objAcls = append(objAcls, acl.ACL)
		}

		if len(objs) > 0 {
			fmt.Printf("==========================\n")
			for i := 0; i < len(objs); i++ {
				obj := objs[i]
				objAcl := objAcls[i]
				if objAcl == string(oss.ACLPublicRead) || objAcl == string(oss.ACLPublicReadWrite) {
					fmt.Printf("%v %v %v %v\n",
						obj.Key, obj.LastModified.In(time.Local).Format("2006-01-02T15:04:05"),
						obj.Size, fmt.Sprintf("https://%v.%v/%v", bucketName, endpoint, url.PathEscape(obj.Key)))
				} else {
					fmt.Printf("%v %v %v\n",
						obj.Key, obj.LastModified.In(time.Local).Format("2006-01-02T15:04:05"), obj.Size)
				}
			}
			fmt.Printf("==========================\n")
		} else {
			fmt.Printf("No objects\n")
		}
	},
}
