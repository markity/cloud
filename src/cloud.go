package main

// @TODO 上传前命名检验

import (
	"fmt"
	"os"
	"time"

	progress "github.com/markity/minio-progress"

	minio "github.com/minio/minio-go"
)

func uploadCmd(args []string) {
	if len(args) != 2 {
		fmt.Printf("upload:参数错误,输入help获得帮助信息\n")
		return
	}
	filePath := args[1]

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("不存在的文件\n")
		}
		fmt.Printf("打开文件失败:%v\n", err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("获取文件信息失败:%v\n", err)
		return
	}

	fmt.Printf("当前操作:上传文件%v,文件大小%v\n", fileInfo.Name(), fileInfo.Size())

	progressBar := progress.NewUploadProgress(fileInfo.Size())
	_, err = client.PutObject(bucketName, fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progressBar})
	if err != nil {
		fmt.Printf("上传文件失败:%v\n", err)
		return
	}
	fmt.Printf("上传完毕\n")
}

func downloadCmd(args []string) {
	if len(args) != 2 {
		fmt.Printf("download:参数错误,输入help获得帮助信息\n")
		return
	}
	fileName := args[1]

	obj, err := client.GetObject(bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Printf("获取对象文件失败:%v\n", err)
		return
	}

	objInfo, err := obj.Stat()
	if err != nil {
		fmt.Printf("获取对象信息失败:%v\n", err)
		return
	}

	if _, err := os.Stat(objInfo.Key); err != nil && !os.IsNotExist(err) {
		fmt.Printf("该目录下存在同名文件,请删除或改名后重试\n")
		return
	}

	fmt.Printf("当前操作:下载文件%v,文件大小%v\n", objInfo.Key, objInfo.Size)

	file, err := os.OpenFile(objInfo.Key, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("创建文件失败:%v\n", err)
		return
	}

	if _, err = progress.CopyWithProgress(file, obj); err != nil {
		fmt.Printf("下载文件失败:%v\n", err)
		return
	}

	fmt.Printf("下载完毕\n")
}

func showCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("help:参数错误,输入help查看帮助信息\n")
		return
	}

	fmt.Printf("当前操作:列出所有文件\n")

	doneCh := make(chan struct{})
	defer close(doneCh)

	fmt.Printf("==========================\n")
	for objInfo := range client.ListObjectsV2(bucketName, "", true, doneCh) {
		if objInfo.Err != nil {
			fmt.Printf("获取对象信息失败:%v\n", objInfo.Err)
			break
		}
		fmt.Printf("%v %v %v\n", objInfo.Key, objInfo.LastModified.In(time.Local).Format("2006-01-02T15:04:05"), objInfo.Size)
	}
	fmt.Printf("==========================\n")
}

func removeCmd(args []string) {
	var err error

	if len(args) != 2 {
		fmt.Printf("remove:参数错误,输入help获得帮助信息\n")
		return
	}
	fileName := args[1]

	objInfo, err := client.StatObject(bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		fmt.Printf("获取文件信息失败:%v\n", err)
		return
	}

	fmt.Printf("当前操作:删除文件%v,文件大小%v\n", objInfo.Key, objInfo.Size)

	err = client.RemoveObject(bucketName, fileName)
	if err != nil {
		fmt.Printf("删除文件失败:%v\n", err)
		return
	}

	fmt.Printf("删除完毕\n")
}

func helpCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("help:参数错误,输入help查看帮助信息\n")
		return
	}

	message := "====================\nupload 文件路径: 上传文件\ndownload 文件名: 下载文件\nshow: 显示所有文件\nremove 文件名: 删除文件\nhelp: 查看帮助\n====================\n"

	fmt.Printf("%v", message)
}

// 分发任务
func handCommand(args []string) {
	if len(args) > 0 {
		mainCmd := args[0]
		switch mainCmd {
		case "upload":
			uploadCmd(args)
		case "download":
			downloadCmd(args)
		case "show":
			showCmd(args)
		case "remove":
			removeCmd(args)
		case "help":
			helpCmd(args)
		default:
			fmt.Printf("未知的命令,请输入help查看帮助信息\n")
		}
	} else {
		fmt.Printf("命令参数有误,请输入help查看帮助信息\n")
	}
}

var client *minio.Client

func main() {
	c, err := getClient()
	if err != nil {
		fmt.Printf("初始化客户端错误:%v\n", err)
		return
	}
	client = c

	args := os.Args[1:]
	handCommand(args)
}
