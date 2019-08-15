package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func uploadCmd(args []string) {
	var err error

	if len(args) != 2 {
		fmt.Printf("upload:参数错误, 输入help获取帮助信息\n")
		return
	}
	filePath := args[1]

	// Get file information
	fileInfo, err := func() (os.FileInfo, error) {
		var err error

		fInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("不存在的文件")
			}
			return nil, err
		}

		return fInfo, nil
	}()
	if err != nil {
		fmt.Printf("获取文件信息失败:%v\n", err)
		return
	}

	// Get file name and size
	fileName := fileInfo.Name()
	fileSize := fileInfo.Size()

	// Empty file is not allowed
	if !(fileSize > 0) {
		fmt.Printf("不允许上传空文件\n")
		return
	}

	// Check if the object exists
	exists, err := bucket.IsObjectExist(fileName)
	if err != nil {
		fmt.Printf("查询对象失败:%v\n", err)
		return
	}
	if exists {
		fmt.Printf("同名对象已存在\n")
		return
	}

	fmt.Printf("当前操作:上传文件%v, 文件大小%v\n", fileName, fileSize)

	// Do upload file
	for {
		if err := bucket.UploadFile(fileName, filePath, partSize, oss.Progress(newProgressBar(fileSize)), oss.Routines(numThreads), oss.Checkpoint(true, fmt.Sprintf("%v.cp", fileName))); err != nil {
			fmt.Printf("\n上传文件失败:%v\n", err)
			fmt.Printf("%v后断点续传...\n", waitTime)
			time.Sleep(waitTime)
			continue
		}
		break
	}
	fmt.Printf("上传成功\n")
}

func downloadCmd(args []string) {
	var err error

	if len(args) != 2 {
		fmt.Printf("download:参数错误, 输入help获取帮助信息\n")
		return
	}
	objectName := args[1]

	// Cheak if the object exists
	exists, err := bucket.IsObjectExist(objectName)
	if err != nil {
		fmt.Printf("查询对象失败:%v\n", err)
		return
	}
	if !exists {
		fmt.Printf("对象不存在\n")
		return
	}

	// Check the local file exists
	_, err = os.Stat(objectName)
	if err == nil {
		fmt.Printf("当前目录同名文件已存在\n")
		return
	} else {
		if !os.IsNotExist(err) {
			fmt.Printf("获取本地文件信息失败:%v\n", err)
			return
		}
	}

	// Get object meta
	header, err := bucket.GetObjectMeta(objectName)
	if err != nil {
		fmt.Printf("获取对象元信息失败:%v\n", err)
	}

	// Get object size
	fileSize, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	if err != nil {
		fmt.Printf("获取对象大小失败:%v\n", err)
		return
	}

	fmt.Printf("当前操作:下载文件%v, 文件大小%v\n", objectName, fileSize)

	// Do download file
	for {
		if err := bucket.DownloadFile(objectName, objectName, partSize, oss.Progress(newProgressBar(fileSize)), oss.Routines(numThreads), oss.Checkpoint(true, fmt.Sprintf("%v.cp", objectName))); err != nil {
			fmt.Printf("\n下载文件失败:%v\n", err)
			fmt.Printf("%v后断点续传...\n", waitTime)
			time.Sleep(waitTime)
			continue
		}
		break
	}
	fmt.Printf("下载成功\n")
}

func listCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("list:参数错误, 输入help获取帮助信息\n")
		return
	}

	fmt.Printf("==========================\n")
	marker := ""
	for {
		var err error

		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			fmt.Printf("列举文件失败:%v\n", err)
			break
		}

		for _, object := range lsRes.Objects {
			fmt.Printf("%v %v %v\n", object.Key, object.LastModified.In(time.Local).Format("2006-01-02T15:04:05"), object.Size)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	fmt.Printf("==========================\n")
}

func removeCmd(args []string) {
	var err error

	if len(args) != 2 {
		fmt.Printf("remove:参数错误, 输入help获取帮助信息\n")
		return
	}
	objectName := args[1]

	// Check if the object exists
	exists, err := bucket.IsObjectExist(objectName)
	if err != nil {
		fmt.Printf("查询对象失败:%v\n", err)
		return
	}
	if !exists {
		fmt.Printf("对象不存在\n")
		return
	}

	// Do delete object
	err = bucket.DeleteObject(objectName)
	if err != nil {
		fmt.Printf("删除文件失败:%v\n", err)
		return
	}

	fmt.Printf("删除成功\n")
}

func helpCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("help:参数错误, 输入help获取帮助信息\n")
		return
	}

	message := "====================\nupload 文件路径: 上传文件\ndownload 文件名: 下载文件\nlist: 显示所有文件\nremove 文件名: 删除文件\nhelp: 查看帮助\n====================\n"

	fmt.Printf("%v", message)
}

// Distribute commands
func handCommand(args []string) {
	if len(args) > 0 {
		mainCmd := args[0]
		switch mainCmd {
		case "upload":
			uploadCmd(args)
		case "download":
			downloadCmd(args)
		case "list":
			listCmd(args)
		case "remove":
			removeCmd(args)
		case "help":
			helpCmd(args)
		default:
			fmt.Printf("未知的命令, 输入help获取帮助信息\n")
		}
	} else {
		fmt.Printf("命令参数有误, 输入help获取帮助信息\n")
	}
}

var bucket *oss.Bucket

func main() {
	b, err := getBucket()
	if err != nil {
		fmt.Printf("初始化数据桶错误:%v\n", err)
		return
	}
	bucket = b

	args := os.Args[1:]
	handCommand(args)
}
