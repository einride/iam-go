.PHONY: all
all: \
	buf-lint \
	buf-generate-authorization \
	buf-generate-example \
	go-mod-tidy

include tools/buf/rules.mk
include tools/protoc-gen-go-grpc/rules.mk

.PHONY: build/protoc-gen-go-authorization-policy
build/protoc-gen-go-authorization-policy:
	$(info [$@] rebuilding plugin...)
	@go build -o $@ .
	@touch $@

build/protoc-gen-go: go.mod
	$(info [$@] rebuilding plugin...)
	@go build -o $@ google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: buf-lint
buf-lint: $(buf)
	$(info [$@] linting proto files...)
	@$(buf) lint

proto_plugins := \
	build/protoc-gen-go \
	build/protoc-gen-go-authorization-policy \
	$(protoc_gen_go_grpc)

.PHONY: buf-generate-authorization
buf-generate-authorization: $(buf) build/protoc-gen-go
	$(info [$@] generating authorization proto stubs...)
	@rm -rf proto/gen/einride/authorization/v1
	@$(buf) generate --path proto/src/einride/authorization/v1 --template buf.gen.authorization.yaml

.PHONY: buf-generate-example
buf-generate-example: $(buf) build/protoc-gen-go build/protoc-gen-go-authorization-policy $(protoc_gen_go_grpc)
	$(info [$@] generating proto stubs...)
	@rm -rf proto/gen/einride/authorization/example/v1
	@$(buf) generate --path proto/src/einride/authorization/example/v1 --template buf.gen.example.yaml

.PHONY: go-mod-tidy
go-mod-tidy:
	$(info [$@] tidying Go module files...)
	@go mod tidy
