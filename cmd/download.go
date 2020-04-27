package cmd

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download specific object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		objectName := args[0]

		exists, err := bucket.IsObjectExist(objectName)
		if err != nil {
			fmt.Printf("failed to query: %v\n", err)
			os.Exit(1)
		}
		if !exists {
			fmt.Printf("the object does't exist\n")
			os.Exit(1)
		}

		hdr, err := bucket.GetObjectMeta(objectName)
		if err != nil {
			fmt.Printf("failed to query: %v\n", err)
			os.Exit(1)
		}

		fileSize, err := strconv.ParseInt(hdr.Get("Content-Length"), 10, 64)
		if err != nil {
			fmt.Printf("failed to query: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("config loaded:\n  part_size_bytes:%v\n  num_threads:%v\n  wait_time_secondes:%v\n",
			cfg.GetPartSize(), cfg.GetNumThreads(), cfg.GetWaitTime())

		fmt.Printf("downloading %v, file size is %v bytes\n", objectName, fileSize)

		for {
			if err := bucket.DownloadFile(objectName, objectName, cfg.GetPartSize(),
				oss.Progress(newProgressBar(fileSize)), oss.Routines(cfg.GetNumThreads()),
				oss.Checkpoint(true, fmt.Sprintf("%v.cp", objectName))); err != nil {
				fmt.Printf("\nfailed to download: %v\n", err)
				fmt.Printf("retring %v later...\n", cfg.GetWaitTime())
				time.Sleep(cfg.GetWaitTime())
				continue
			}
			break
		}
		fmt.Printf("succeec to download\n")
	},
}
