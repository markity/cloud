package main

import (
	"fmt"

	minio "github.com/minio/minio-go"
)

// minio客户端基本信息
var endpoint = "127.0.0.1:9000"
var accessKeyID = "xuri"
var secretAccessKey = "jx2004119"
var secure = false

// 数据桶信息
var bucketName = "music"
var location = "us-east-1"

// 只执行一次, 用来创建数据桶并返回客户端对象
func getClient() (*minio.Client, error) {
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		return nil, err
	}

	// 检查数据桶是否存在
	exists, err := client.BucketExists(bucketName)
	if err != nil {
		return client, fmt.Errorf("查询数据桶错误(%v)", err)
	}
	// 不存在则创建数据桶
	if !exists {
		err := client.MakeBucket(bucketName, location)
		if err != nil {
			return client, fmt.Errorf("创建数据桶错误(%v)", err)
		}
	}

	return client, nil
}
