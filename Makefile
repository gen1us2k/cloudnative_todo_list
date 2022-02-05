.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## install binaries
	brew tap bufbuild/buf
	brew install buf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

build_grpc:
	protoc -I . \
	    --go_out ./grpc --go_opt paths=source_relative \
	    --go-grpc_out ./grpc --go-grpc_opt paths=source_relative \
	    grpc/todolist.proto
