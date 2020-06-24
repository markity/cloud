package cmd

import (
	"cloud/util"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

// oss settings
const endpoint = "Your-Endpoint" // like oss-cn-chengdu.aliyuncs.com
const accessKeyID = "Your-AccessKeyID"
const accessKeySecret = "Your-AccessKeySecret"
const bucketName = "Your-BuckName"

var bucket *oss.Bucket
var cfg *util.Config

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
		fmt.Printf("failed to get executable file path: %v\n", err)
		os.Exit(1)
	}
	cfgPath := filepath.Join(filepath.Dir(executable), util.CfgName)
	_, err = os.Stat(cfgPath)
	if err != nil {
		// unknown error
		if !os.IsNotExist(err) {
			fmt.Printf("failed to stat profile: %v\n", err)
			os.Exit(1)
		}
		// err is "existed", write file
		f, err := os.Create(cfgPath)
		if err != nil {
			fmt.Printf("failed to create profile: %v\n", err)
			os.Exit(1)
		}
		_, err = f.Write(util.CfgBase)
		if err != nil {
			fmt.Printf("failed to write profile: %v\n", err)
			os.Exit(1)
		}
		err = f.Close()
		if err != nil {
			fmt.Printf("failed to close profile: %v\n", err)
			os.Exit(1)
		}
	}
	cfgBytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		fmt.Printf("failed to read profile: %v\n", err)
		os.Exit(1)
	}
	cfg = &util.Config{}
	err = json.Unmarshal(cfgBytes, cfg)
	if err != nil {
		fmt.Printf("failed to unmarshal profile: %v\n", err)
		os.Exit(1)
	}
	if cfg.NumThreads < 1 || cfg.PartSize < 1 || cfg.WaitTimeSeconds < 0 {
		fmt.Printf("invalid config, please check config file\n")
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
