# Cloud - 命令行文件存储工具

### 这是什么?

此项目是基于`阿里云oss服务`的命令行文件存储程序，提供`上传`，`下载`，`查看所有文件`，`删除文件`，`分享文件`等功能。文件上传下载功能支持分片传输，断点续传，自动重试。使用上传下载功能可使用配置文件自行配置单个分片大小，线程数以及失败后等待时间。

### 如何编译?

**1. 修改此项目的`cloud/config.go`文件，填写`oss基础配置`**

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
> cloud.exe help
====================
init: 初始上传下载配置文件
upload 文件路径: 上传文件
download 文件名: 下载文件
list: 显示所有文件
remove 文件名: 删除文件
rename 文件名 新命名: 修改文件名
com 文件夹名 [目标文件名]: 打包一个文件夹为tar格式
share 文件名: 分享文件
unshare 文件名: 取消分享
help: 查看帮助
====================
```

**初始化配置文件**

```powershell
> cloud.exe init
初始化配置文件成功
```

此时`cloud.exe`所在路径将会创建一个名为`config.json`的配置文件，默认配置参数如下

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

注意：此配置文件仅仅使用于`上传`与`下载`，`查看所有文件`与`删除文件`不需初始化配置文件

**上传文件**

```powershell
> cloud.exe upload .\myfile.txt
成功加载配置:
  分片大小:2097152
  线程数:3
  等待时长:5s
当前操作:上传文件myfile.txt, 文件大小205
[====================] 205 / 205 100%
上传成功
```

**下载文件**

```powershell
> cloud.exe download myfile.txt
成功加载配置:
  分片大小:2097152
  线程数:3
  等待时长:5s
当前操作:下载文件myfile.txt, 文件大小205
[====================] 205 / 205 100%
下载成功
```

**查看所有文件**

```powershell
> cloud.exe list
==========================
myfile.txt 2019-11-10T12:29:39 205
==========================
```


第一列是文件名，第二列为上传时间，第三列为文件大小（单位字节）

**删除文件**

```powershell
> cloud.exe remove myfile.txt
删除成功
```

**重命名文件**

```powershell
> cloud.exe rename myfile.txt newfile.txt
重命名成功
```

**压缩文件夹**

```powershell
> cloud.exe com .\imagefloder images.tar
```

将`imagefloder`文件夹压缩为`images.tar`文件

**分享文件**

```powershell
> cloud.exe share myfile.txt
分享文件成功, 文件地址:https://your-bucketName.your-endPoint.aliyuncs.com/myfile.txt
```

**取消分享**

```powershell
> cloud.exe unshare myfile.txt
取消分享成功
```