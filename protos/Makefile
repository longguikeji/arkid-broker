.PHONY: all

INCLUDES := $(wildcard include/*.proto)

all: apistore-rpc apistore-proxy

apistore-rpc: apistore.proto $(INCLUDES)
	@env PATH=$$PATH:~/go/bin protoc \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I/usr/local/include \
		-I$$GOPATH/src \
		-I.\
		--go_out=plugins=grpc:apistore \
		$<

apistore-proxy: apistore.proto $(INCLUDES)
	@env PATH=$$PATH:~/go/bin protoc \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I/usr/local/include \
		-I$$GOPATH/src \
		-I.\
		--grpc-gateway_out=logtostderr=true:apistore \
		$<