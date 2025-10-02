.PHONY: proto

PROTO_PATH := proto/agent.proto
PROTOC := /opt/homebrew/bin/protoc
PROTOC_GEN_GO := /Users/chalkan3/go/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := /Users/chalkan3/go/bin/protoc-gen-go-grpc

proto:
	@echo "Generating protobuf Go code..."
	@PATH=$(dir $(PROTOC_GEN_GO)):$$PATH $(PROTOC) --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. $(PROTO_PATH)
	@echo "Protobuf Go code generated."
