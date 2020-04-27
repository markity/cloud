package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// oss settings
const endpoint = "oss-cn-chengdu.aliyuncs.com" // like oss-cn-chengdu.aliyuncs.com
const accessKeyID = "LTAI4GHri1t1mHtWHK5RRyf9"
const accessKeySecret = "e0PNEk94aVKMlHRrDfzb2rUEZIy69O"
const bucketName = "cloud-netdisk"

var bucket *oss.Bucket

// config settings
type config struct {
	PartSize        int64 `json:"part_size_bytes"`
	NumThreads      int   `json:"num_threads"`
	WaitTimeSeconds int   `json:"wait_time_seconds"`
}

var cfg *config
var cfgName = "cloud-cfg.json"
var cfgBase = []byte(`{
    "part_size_bytes": 2097152,
    "num_threads": 3,
    "wait_time_seconds": 5
}`)

func (c *config) GetPartSize() int64 {
	return c.PartSize
}
func (c *config) GetNumThreads() int {
	return c.NumThreads
}
func (c *config) GetWaitTime() time.Duration {
	return time.Duration(c.WaitTimeSeconds) * time.Second
}

// progress bar
func newProgressBar(total int64) *progressBar {
	return &progressBar{total: total, current: 0, percent: 0}
}

type progressBar struct {
	total   int64
	current int64
	percent int
}

func (pb *progressBar) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferDataEvent:
		if pb.total == event.TotalBytes {
			pb.current = event.ConsumedBytes
			percent := int(float64(pb.current) * 100 / float64(pb.total))
			if percent != pb.percent {
				pb.percent = percent
				pb.draw()
				if percent == 100 {
					fmt.Println()
				}
			}
		}
	default:
	}
}
func (pb *progressBar) draw() {
	num := pb.percent / 5
	fmt.Printf("\r[%v%v] %v / %v %v%%", multiString("=", num), multiString(" ", 20-num), pb.current, pb.total, pb.percent)
}
func multiString(s string, num int) string {
	var buffer bytes.Buffer
	for i := 0; i < num; i++ {
		buffer.WriteString(s)
	}

	return buffer.String()
}

var rootCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Cloud is a net-disk program based on aliyun-oss",
	Long: `Cloud provides easy interface to manage files. It Contains uploading, 
downloading, removing, renaming, listing and sharing operations.
Based on aliyun-oss, you can set up your own net-disk rapidly.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("failed to execute command: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(shareCmd)
	rootCmd.AddCommand(unshareCmd)
}

func initConfig() {
	executable, err := os.Executable()
	if err != nil {
		fmt.Printf("failed to get executable file path: %v", err)
		os.Exit(1)
	}
	cfgPath := filepath.Join(filepath.Dir(executable), cfgName)
	_, err = os.Stat(cfgPath)
	if err != nil {
		// unknown error
		if !os.IsNotExist(err) {
			fmt.Printf("failed to init profile: %v\n", err)
			os.Exit(1)
		}
		// err not exists, write file
		f, err := os.Create(cfgPath)
		if err != nil {
			fmt.Printf("failed to init profile: %v\n", err)
			os.Exit(1)
		}
		_, err = f.Write(cfgBase)
		if err != nil {
			fmt.Printf("failed to init profile: %v\n", err)
			os.Exit(1)
		}
		err = f.Close()
		if err != nil {
			fmt.Printf("failed to init profile: %v\n", err)
			os.Exit(1)
		}
	}
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		fmt.Printf("failed to load profile: %v\n", err)
		os.Exit(1)
	}
	cfg = &config{}
	err = json.Unmarshal(cfgBytes, cfg)
	if err != nil {
		fmt.Printf("failed to load profile: %v\n", err)
		os.Exit(1)
	}
	if cfg.NumThreads < 1 || cfg.PartSize < 1 || cfg.WaitTimeSeconds < 0 {
		fmt.Printf("invalid config, please check the %v\n", cfgPath)
		os.Exit(1)
	}
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		fmt.Printf("failed to init aliyun-oss client: %v\n", err)
		os.Exit(1)
	}
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("failed to init aliyun-oss bucket: %v\n", err)
		os.Exit(1)
	}
}
