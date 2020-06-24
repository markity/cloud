package cmd

import (
	"cloud/util"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var shared bool

func init() {
	uploadCmd.Flags().BoolVarP(&shared, "shared", "s", false, "true if share this object")
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload specific file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Printf("failed to stat file: %v\n", err)
			os.Exit(1)
		}
		if fileInfo.IsDir() {
			fmt.Printf("folder is not allowed to upload\n")
			os.Exit(1)
		}

		fileName := fileInfo.Name()
		fileSize := fileInfo.Size()

		if !(fileSize > 0) {
			fmt.Printf("empty file is not allowed to upload\n")
			os.Exit(1)
		}

		exists, err := bucket.IsObjectExist(fileName)
		if err != nil {
			fmt.Printf("failed to check existence of object: %v\n", err)
			os.Exit(1)
		}
		if exists {
			fmt.Printf("objcet with same name exists")
			os.Exit(1)
		}

		fmt.Printf("config loaded:\n  part_size_bytes:%v\n  num_threads:%v\n  wait_time_secondes:%v\n",
			cfg.GetPartSize(), cfg.GetNumThreads(), cfg.GetWaitTime())
		if shared {
			fmt.Printf("  acl: PublicRead\n")
		} else {
			fmt.Printf("  acl: Private")
		}
		fmt.Printf("uploading %v, file size is %v bytes\n", fileName, fileSize)

		for {
			if err := bucket.UploadFile(fileName, filePath, cfg.GetPartSize(),
				oss.Progress(util.NewProgressBar(fileSize)), oss.Routines(cfg.GetNumThreads()),
				oss.Checkpoint(true, fmt.Sprintf("%v.cp", fileName))); err != nil {
				fmt.Printf("\nfailed to upload: %v\n", err)
				fmt.Printf("retrying %v latter...\n", cfg.GetWaitTime())
				time.Sleep(cfg.GetWaitTime())
				continue
			}
			break
		}
		fmt.Printf("succeed to upload\n")

		if err := bucket.SetObjectACL(fileName, oss.ACLPublicRead); err != nil {
			fmt.Printf("falied to share object: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("succeed to share object")
	},
}
