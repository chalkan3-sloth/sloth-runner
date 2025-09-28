.PHONY: proto

PROTO_PATH := proto/agent.proto
PROTOC := /opt/homebrew/bin/protoc
PROTOC_GEN_GO := /Users/chalkan3/go/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := /Users/chalkan3/go/bin/protoc-gen-go-grpc

proto:
	@echo "Generating protobuf Go code..."
	@$(PROTOC) --plugin=protoc-gen-go=$(PROTOC_GEN_GO) --go_out=. --plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) --go-grpc_out=. $(PROTO_PATH)
	@echo "Protobuf Go code generated."
