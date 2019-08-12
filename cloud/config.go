package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var endpoint = "oss-cn-chengdu.aliyuncs.com"
var accessKeyID = "LTAIsfVsnfB9GgSx"
var accessKeySecret = "yqSlknBGyIpe3iUr8jdbm2TqJqA8ni"
var bucketName = "cloud-netdisk"

// 获取数据桶, 通常只执行一次
func getBucket() (*oss.Bucket, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
