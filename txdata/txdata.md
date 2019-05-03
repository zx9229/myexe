说明:  
`protoc.exe`源自`protoc-3.6.1-win32.zip`。  
  文件名: protoc-3.6.1-win32.zip  
文件大小: 1007473 字节  
链接地址: https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-win32.zip  
    MD5: EE9C100E84A6A6A64636306185993538  

在执行命令前请确认以下项:  
确保`GOPATH`的第一个目录是`%USERPROFILE%\go`
已经执行过 go get -u -v google.golang.org/grpc  
已经执行过 go get -u -v github.com/golang/protobuf/protoc-gen-go  

曾经的命令:  
protoc  --proto_path=./  txdata.proto  --go_out=plugins=grpc:./  
protoc  --proto_path=./  txdata.proto  --cpp_out=../businessClient/protobuf/  

现在的命令:  
protoc  --proto_path=./  --proto_path=../businessClient/protobuf/protobuf-3.6.1/src/  txdata.proto  --go_out=plugins=grpc:./  
protoc  --proto_path=./  --proto_path=../businessClient/protobuf/protobuf-3.6.1/src/  txdata.proto  --cpp_out=../businessClient3/protobuf/  
