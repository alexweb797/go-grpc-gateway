gateway-build: grpc-generate grpc-gateway-proxy grpc-openapi statik-it
	go build -ldflags "-X $(MODULE)/build.BuildNumber=$(BUILD_NUMBER) -X $(MODULE)/build.CommitSHA=$(COMMIT_SHA)" -o gateway ./cmd/gateway

statik-it:
	rm -f statik/statik.go
	statik -src=swagger/

grpc-generate:
	protoc -I . \
	   -I./external/grpc-gateway/third_party/googleapis \
       -I./external/grpc-gateway \
       --go_out .  \
       --go-grpc_out . \
       external/triton-proto/uaa.proto

grpc-gateway-proxy:
	protoc -I . \
		 -I./external/grpc-gateway/third_party/googleapis \
		 -I./external/grpc-gateway \
         --grpc-gateway_out . \
         --grpc-gateway_opt logtostderr=true \
         external/triton-proto/uaa.proto

grpc-openapi:
	protoc -I . \
	     -I./external/grpc-gateway/third_party/googleapis \
	     -I./external/grpc-gateway \
		 --openapiv2_out . \
		 --openapiv2_opt logtostderr=true \
		 --openapiv2_opt generate_unbound_methods=true \
		 external/triton-proto/uaa.proto
	mv external/triton-proto/uaa.swagger.json swagger/

install-lib:
	go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc