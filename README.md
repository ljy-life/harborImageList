# harborImageList
Get the list of images in harbor, 获取 harbor 中的镜像列表

# Build Command

```shell
go build -o harborImageTag harbor.go
```

# Command Help
```shell
$ ./harborImageTag -h
Usage of ./harborImageTag:
  -file string
        运行结果输出文件 (default "harborImageList.txt")
  -passwd string
        harbor password, default 123456 (default "123456")
  -repositry string
        指定要选中的仓库，默认为 all，全部仓库 (default "all")
  -schema string
        http or https，default http (default "http")
  -url string
        harbor Address，deault harbor.k8s.local (default "harbor.k8s.local")
  -user string
        harbor admin, default admin (default "admin")
```

# Example
```shell
harborImageTag --schema http --url harbor.k8s.local --user admin --passwd 123456 --repositry all --file harborImageList.txt
```
