### Proto generation

SRC_DIR=.
GO_DST_DIR=../internal/helloworld

setup: ./helloworld.proto
	mkdir -p $(GO_DST_DIR)

go_stub: ./helloworld.proto setup
	protoc -I$(SRC_DIR) \
	       -I$(GOPATH)/src \
	       --go_out=plugins=grpc:$(GO_DST_DIR) \
	       helloworld.proto