```bash
protoc --proto_path=$GOPATH/src --proto_path=. --go_out=paths=source_relative:. status.proto

```