protoc_gen_go_grpc_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
protoc_gen_go_grpc_version := 1.1.0
protoc_gen_go_grpc_dir := $(protoc_gen_go_grpc_cwd)/$(protoc_gen_go_grpc_version)
protoc_gen_go_grpc := $(protoc_gen_go_grpc_dir)/protoc-gen-go-grpc
export PATH := $(dir $(protoc_gen_go_grpc)):$(PATH)

ifeq ($(shell uname),Linux)
protoc_gen_go_grpc_archive_url := https://github.com/grpc/grpc-go/releases/download/cmd%2Fprotoc-gen-go-grpc%2Fv$(protoc_gen_go_grpc_version)/protoc-gen-go-grpc.v$(protoc_gen_go_grpc_version).linux.amd64.tar.gz
else ifeq ($(shell uname),Darwin)
protoc_gen_go_grpc_archive_url := https://github.com/grpc/grpc-go/releases/download/cmd%2Fprotoc-gen-go-grpc%2Fv$(protoc_gen_go_grpc_version)/protoc-gen-go-grpc.v$(protoc_gen_go_grpc_version).darwin.amd64.tar.gz
else
$(error unsupported OS: $(shell uname))
endif

$(protoc_gen_go_grpc):
	$(info [protoc-gen-go-grpc] fetching version $(protoc_gen_go_grpc_version)...)
	@mkdir -p $(protoc_gen_go_grpc_dir)
	@curl -sSL $(protoc_gen_go_grpc_archive_url) -o - | tar -xz --directory $(protoc_gen_go_grpc_dir)
	@chmod +x $@
	@touch $@
