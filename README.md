# 启动

在 cmd/app 目录下执行 go run .

## wire 依赖注入

在 cmd/app 目录下执行 wire

## swagger 文档生成

根目录执行

swag init -g ./cmd/app/main.go -o ./docs


## 未解决问题

wire生成的InitApp函数如果在根目录执行run会提示InitApp未定义。先粘贴到main里面后面再看看。
```
go run cmd/app/main.go 
```