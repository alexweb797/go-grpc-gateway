module triton-grpc-gateway

go 1.14

replace github.com/grpc-ecosystem/grpc-gateway/v2 v2.1.0 => ./external/grpc-gateway

require (
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.1.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506
	google.golang.org/grpc v1.35.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.25.0
)
