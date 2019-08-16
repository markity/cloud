package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// oss基础设置
var endpoint = "oss-cn-chengdu.aliyuncs.com"
var accessKeyID = "LTAIsfVsnfB9GgSx"
var accessKeySecret = "yqSlknBGyIpe3iUr8jdbm2TqJqA8ni"
var bucketName = "cloud-netdisk"

// 默认配置
var baseConfig = `{
    "part_size_bytes": 2097152,
    "num_threads": 3,
    "wait_time_seconds": 5
}`

// 上传下载所需的可配置参数
var config *Config

type Config struct {
	PartSize        int64 `json:"part_size_bytes"`
	NumThreads      int   `json:"num_threads"`
	WaitTimeSeconds int   `json:"wait_time_seconds"`
}

func (c *Config) GetPartSize() int64 {
	return c.PartSize
}

func (c *Config) GetNumThreads() int {
	return c.NumThreads
}

func (c *Config) GetWaitTime() time.Duration {
	return time.Duration(c.WaitTimeSeconds) * time.Second
}

func getConfigPath() (string, error) {
	execable, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(execable), "config.json"), nil
}

func prepareConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("获取配置文件路径失败(%v)", err)
	}

	// 创建配置文件
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("创建配置文件失败(%v)", err)
	}

	// 写入默认配置
	_, err = file.Write([]byte(baseConfig))
	if err != nil {
		return err
	}

	return nil
}

// 读取可执行文件目录下的config.json配置
func readConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("获取配置文件路径失败(%v)", err)
	}

	// 打开配置文件
	file, err := os.OpenFile(configPath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("不存在配置文件(%v), 输入init初始化配置", configPath)
		}
		return err
	}

	// 读取配置文件
	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("读取配置文件失败(%v)", err)
	}

	// 解码配置文件
	c := &Config{}
	if err := json.Unmarshal(configBytes, c); err != nil {
		return fmt.Errorf("解码配置文件失败(%v)", err)
	}

	// 更新外部变量
	config = c

	return nil
}

// 获取数据桶, 通常只执行一次
func getBucket() (*oss.Bucket, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	b, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return b, nil
}
