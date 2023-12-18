生成模板文件
```bash
goctl api go --dir=./ --api applet.api

goctl rpc protoc ./user.proto --go_out=. --go-grpc_out=. --zrpc_out=./
```

创建sql
```shell
goctl model mysql datasource --dir ./internal/model --table user --cache true --url "root:123456@tcp(127.0.0.1:3307)/beyond_user"
```