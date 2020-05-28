# Cloud - 命令行文件存储工具

### 这是什么?

此项目是基于`阿里云oss服务`的命令行文件存储程序，提供`上传`，`下载`，`查看所有文件`，`删除文件`，`分享文件`等功能。文件上传下载功能支持分片传输，断点续传，自动重试。使用上传下载功能可使用配置文件自行配置单个分片大小，线程数以及失败后等待时间。

### 如何编译?

**1. 修改此项目的`cmd/root.go`文件，填写`oss settings`**

**2. 编译**

```powershell
> cd .\cloud
> go build
> ls
Mode                LastWriteTime         Length Name
----                -------------         ------ ----
-a----        2020/2/29     11:58        8047104 cloud.exe
-a----        2020/2/29     11:02          11063 cloud.go
-a----        2020/2/29     11:57           3059 config.go
-a----        2020/2/28     20:04             84 config.json
-a----        2020/1/26     17:25           1050 progress.go
```

### 如何使用?

**查看帮助**

```powershell
> cloud.exe
Cloud provides easy interface to manage files. It Contains uploading,
downloading, removing, renaming, listing and sharing operations.
Based on aliyun-oss, you can set up your own net-disk rapidly.

Usage:
  cloud [command]

Available Commands:
  download    Download specific object
  help        Help about any command
  list        List all objects
  remove      Remove specific object
  rename      Rename specific object
  share       Share specific object
  unshare     Unshare specific object
  upload      Upload specific file

Flags:
  -h, --help   help for cloud

Use "cloud [command] --help" for more information about a command.

```

**默认配置文件**

```json
{
    "part_size_bytes": 2097152,
    "num_threads": 3,
    "wait_time_seconds": 5
}
```

|      配置名       |          说明          | 单位 |
| :---------------: | :--------------------: | :--: |
|  part_size_bytes  | 上传下载单个分片的大小 | 字节 |
|    num_threads    |  上传下载使用线程数量  |  个  |
| wait_time_seconds | 上传下载失败后等待时间 |  秒  |

**上传文件**

```powershell
> cloud.exe upload .\myfile.txt
config loaded:
  part_size_bytes:72
  num_threads:3
  wait_time_secondes:5s
uploading main.go, file size is 72 bytes
[====================] 72 / 72 100%
succeed to upload
```

**下载文件**

```powershell
> cloud.exe download myfile.txt
config loaded:
  part_size_bytes:72
  num_threads:3
  wait_time_secondes:5s
downloading myfile.txt, file size is 72 bytes
[====================] 72 / 72 100%
succeec to download
```

**查看所有文件**

```powershell
> cloud.exe list
myfile.txt 2020-04-23T15:23:23 72
1 objects in total
```

**删除文件**

```powershell
> cloud.exe remove myfile.txt
succeed to remove
```

**重命名文件**

```powershell
> cloud.exe rename myfile.txt newfile.txt
succeed to rename
```

**分享文件**

```powershell
> cloud.exe share myfile.txt
succeed to share, path: https://your-bucketName.your-endPoint.aliyuncs.com/myfile.txt
```

**取消分享**

```powershell
> cloud.exe unshare myfile.txt
succeed to unshare
```