.PHONY: .protoc

# go install build tools
PKG+="google.golang.org/protobuf/cmd/protoc-gen-go@latest"
PKG+="google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
PKG+="github.com/favadi/protoc-go-inject-tag@latest"

%.install: ## go install build tools
	@/bin/bash ./make.sh install $*

protoc: $(PKG:=.install) ## build protobuf
	@/bin/bash ./make.sh protoc