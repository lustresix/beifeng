```bash
goctl api go --dir=./ --api applet.api

goctl rpc protoc ./user.proto --go_out=. --go-grpc_out=. --zrpc_out=./
```