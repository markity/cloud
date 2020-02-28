package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func initCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("init:参数错误, 输入help获取帮助信息\n")
		return
	}

	if err := prepareConfig(); err != nil {
		fmt.Printf("初始化配置文件失败:%v\n", err)
		return
	}
	fmt.Printf("初始化配置文件成功\n")
}

func uploadCmd(args []string) {
	if len(args) != 2 {
		fmt.Printf("upload:参数错误, 输入help获取帮助信息\n")
		return
	}
	filePath := args[1]

	// 打开文件, 这里不使用os.Stat, 保证文件不被中途篡改
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("获取文件信息失败:不存在的文件\n")
		} else if os.IsPermission(err) {
			fmt.Printf("获取文件信息失败:权限不足\n")
		} else {
			fmt.Printf("获取文件信息失败:%v\n", err)
		}
		return
	}
	defer file.Close()

	fileInfo, err := os.Stat(filePath)
	file.Stat()
	if err != nil {
		fmt.Printf("获取文件信息失败:%v\n", err)
		return
	}
	if fileInfo.IsDir() {
		fmt.Printf("不允许上传文件夹\n")
		return
	}

	// 获取文件名 文件大小
	fileName := fileInfo.Name()
	fileSize := fileInfo.Size()

	// 不允许空文件
	if !(fileSize > 0) {
		fmt.Printf("不允许上传空文件\n")
		return
	}

	// 检查对象名是否存在
	exists, err := bucket.IsObjectExist(fileName)
	if err != nil {
		fmt.Printf("查询对象失败:%v\n", err)
		return
	}
	if exists {
		fmt.Printf("同名对象已存在\n")
		return
	}

	// 加载配置文件
	if err := readConfig(); err != nil {
		fmt.Printf("加载配置文件失败:%v\n", err)
		return
	}
	fmt.Printf("成功加载配置:\n  分片大小:%v\n  线程数:%v\n  等待时长:%v\n", config.GetPartSize(), config.GetNumThreads(), config.GetWaitTime())
	fmt.Printf("当前操作:上传文件%v, 文件大小%v\n", fileName, fileSize)

	// 上传
	for {
		if err := bucket.UploadFile(fileName, filePath, config.GetPartSize(), oss.Progress(newProgressBar(fileSize)), oss.Routines(config.GetNumThreads()), oss.Checkpoint(true, fmt.Sprintf("%v.cp", fileName))); err != nil {
			fmt.Printf("\n上传文件失败:%v\n", err)
			fmt.Printf("%v后断点续传...\n", config.GetWaitTime())
			time.Sleep(config.GetWaitTime())
			continue
		}
		break
	}
	fmt.Printf("上传成功\n")
}

func downloadCmd(args []string) {
	if len(args) != 2 {
		fmt.Printf("download:参数错误, 输入help获取帮助信息\n")
		return
	}
	objectName := args[1]

	// 检查对象是否存在
	exists, err := bucket.IsObjectExist(objectName)
	if err != nil {
		fmt.Printf("查询对象失败:%v\n", err)
		return
	}
	if !exists {
		fmt.Printf("对象不存在\n")
		return
	}

	// 获取对象元信息
	header, err := bucket.GetObjectMeta(objectName)
	if err != nil {
		if serviceErr, ok := err.(oss.ServiceError); ok && serviceErr.StatusCode == 404 {
			fmt.Printf("对象不存在\n")
			return
		}
		fmt.Printf("获取对象元信息失败:%v\n", err)
		return
	}

	// 获取对象大小
	fileSize, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	if err != nil {
		fmt.Printf("获取对象大小失败:%v\n", err)
		return
	}

	// 加载配置文件
	if err := readConfig(); err != nil {
		fmt.Printf("加载配置文件失败:%v\n", err)
		return
	}
	fmt.Printf("成功加载配置:\n  分片大小:%v\n  线程数:%v\n  等待时长:%v\n", config.GetPartSize(), config.GetNumThreads(), config.GetWaitTime())

	fmt.Printf("当前操作:下载文件%v, 文件大小%v\n", objectName, fileSize)

	// 下载
	for {
		if err := bucket.DownloadFile(objectName, objectName, config.GetPartSize(), oss.Progress(newProgressBar(fileSize)), oss.Routines(config.GetNumThreads()), oss.Checkpoint(true, fmt.Sprintf("%v.cp", objectName))); err != nil {
			fmt.Printf("\n下载文件失败:%v\n", err)
			fmt.Printf("%v后断点续传...\n", config.GetWaitTime())
			time.Sleep(config.GetWaitTime())
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

	// 获取所有对象信息
	objs := make([]oss.ObjectProperties, 0)
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			fmt.Printf("列举文件失败:%v\n", err)
			break
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

	// 循环遍历
	if len(objs) > 0 {
		fmt.Printf("==========================\n")
		for _, obj := range objs {
			fmt.Printf("%v %v %v\n", obj.Key, obj.LastModified.In(time.Local).Format("2006-01-02T15:04:05"), obj.Size)
		}
		fmt.Printf("==========================\n")
	} else {
		fmt.Printf("无任何对象\n")
	}
}

func removeCmd(args []string) {
	if len(args) != 2 {
		fmt.Printf("remove:参数错误, 输入help获取帮助信息\n")
		return
	}
	objectName := args[1]

	// 检查对象是否存在
	exists, err := bucket.IsObjectExist(objectName)
	if err != nil {
		fmt.Printf("查询对象失败:%v\n", err)
		return
	}
	if !exists {
		fmt.Printf("对象不存在\n")
		return
	}

	// 删除对象
	err = bucket.DeleteObject(objectName)
	if err != nil {
		fmt.Printf("删除文件失败:%v\n", err)
		return
	}
	fmt.Printf("删除成功\n")
}

func renameCmd(args []string) {
	if len(args) != 3 {
		fmt.Printf("rename:参数错误, 输入help获取帮助信息\n")
		return
	}

	if args[1] == args[2] {
		fmt.Println("重命名对象失败:新对象名不允许与原对象名相同")
		return
	}

	if _, err := bucket.CopyObject(args[1], args[2]); err != nil {
		fmt.Printf("创建新对象失败:%v\n", err)
		return
	}
	if err := bucket.DeleteObject(args[1]); err != nil {
		fmt.Printf("删除原对象失败:%v\n", err)
		return
	}

	fmt.Printf("重命名成功\n")
}

func comCmd(args []string) {
	if !(len(args) == 2 || len(args) == 3) {
		fmt.Printf("com:参数错误, 输入help获取帮助信息\n")
		return
	}

	// 检查文件类型
	srcInfo, err := os.Stat(args[1])
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("获取文件夹信息失败:不存在的文件夹\n")
		} else if os.IsPermission(err) {
			fmt.Printf("获取文件夹信息失败:权限不足\n")
		} else {
			fmt.Printf("获取文件夹信息失败:%v\n", err)
		}
		return
	}
	if !srcInfo.IsDir() {
		fmt.Printf("仅可压缩文件夹\n")
		return
	}

	// 目标压缩文件的文件名
	destName := srcInfo.Name() + ".tar"
	if len(args) == 3 {
		destName = args[2]
	}

	destFile, err := os.OpenFile(destName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("创建文件失败:权限不足\n")
		} else {
			fmt.Printf("创建文件失败:%v\n", err)
		}
		return
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			fmt.Printf("关闭目标压缩文件失败:%v\n", err)
		}
	}()

	tw := tar.NewWriter(destFile)
	defer func() {
		if err := tw.Close(); err != nil {
			fmt.Printf("关闭压缩写入器失败:%v\n", err)
		}
	}()

	err = filepath.Walk(args[1], func(path string, info os.FileInfo, err error) (e error) {
		if info.IsDir() {
			return nil
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		hdr.Name = path

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		// 最终关闭文件错误, 也视为失败
		defer func() {
			err := f.Close()
			if e == nil {
				e = err
			}
		}()

		_, err = io.Copy(tw, f)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("遍历文件夹失败:%v\n", err)
		return
	}
}

func helpCmd(args []string) {
	if len(args) != 1 {
		fmt.Printf("help:参数错误, 输入help获取帮助信息\n")
		return
	}

	message := "====================\ninit: 初始上传下载配置文件\nupload 文件路径: 上传文件\ndownload 文件名: 下载文件\nlist: 显示所有文件\nremove 文件名: 删除文件\nrename 文件名 新命名: 修改文件名\ncom 文件夹名 [目标文件名]: 打包一个文件夹为tar格式\nhelp: 查看帮助\n====================\n"

	fmt.Printf("%v", message)
}

// 分发命令
func handCommand(args []string) {
	if len(args) > 0 {
		mainCmd := args[0]
		switch mainCmd {
		case "init":
			initCmd(args)
		case "upload":
			uploadCmd(args)
		case "download":
			downloadCmd(args)
		case "list":
			listCmd(args)
		case "remove":
			removeCmd(args)
		case "rename":
			renameCmd(args)
		case "com":
			comCmd(args)
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
		fmt.Printf("初始化数据桶失败:%v\n", err)
		return
	}
	bucket = b

	args := os.Args[1:]
	handCommand(args)
}
