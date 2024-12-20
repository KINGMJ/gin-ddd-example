# 启动

在 cmd/app 目录下执行 go run .


## wire 依赖注入

在 cmd/app 目录下执行 wire

## swagger 文档生成

根目录执行

swag init -g ./cmd/app/main.go -o ./docs

## 未解决问题

wire 生成的 InitApp 函数如果在根目录执行 run 会提示 InitApp 未定义。先粘贴到 main 里面后面再看看。

```
go run cmd/app/main.go
```

example 里面是各种中间件的示例代码

```
go run example/rabbitmq_example/send/send.go 
```


