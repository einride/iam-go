SHELL := /bin/bash

.PHONY: all
all: \
	buf-lint \
	buf-generate \
	buf-generate-example \
	spanner-generate \
	go-lint \
	go-test \
	go-mod-tidy

include tools/buf/rules.mk
include tools/golangci-lint/rules.mk
include tools/protoc-gen-go-grpc/rules.mk
include tools/semantic-release/rules.mk

build/protoc-gen-go: go.mod
	$(info [$@] rebuilding plugin...)
	@go build -o $@ google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: build/protoc-gen-go-iam
build/protoc-gen-go-iam:
	$(info [$@] rebuilding plugin...)
	@go build -o $@ ./cmd/protoc-gen-go-iam

.PHONY: buf-lint
buf-lint: $(buf)
	$(info [$@] linting proto files...)
	@$(buf) lint

.PHONY: buf-generate
buf-generate: $(buf) build/protoc-gen-go
	$(info [$@] generating proto stubs...)
	@rm -rf proto/gen/einride/iam/v1
	@$(buf) generate --path proto/src/einride/iam/v1 --template buf.gen.yaml

protoc_plugins := \
 build/protoc-gen-go \
 build/protoc-gen-go-iam \
 $(protoc_gen_go_grpc)

.PHONY: buf-generate-example
buf-generate-example: $(buf) $(protoc_plugins)
	$(info [$@] generating proto stubs...)
	@rm -rf proto/gen/einride/iam/example/v1
	@$(buf) generate --path proto/src/einride/iam/example/v1 --template buf.gen.example.yaml

.PHONY: spanner-generate
spanner-generate:
	$(info [$@] generating Spanner database APIs...)
	@go run -mod=mod go.einride.tech/spanner-aip/cmd/spanner-aip-go generate

.PHONY: go-mod-tidy
go-mod-tidy:
	$(info [$@] tidying Go module files...)
	@go mod tidy
	@cd cmd/iamctl && go mod tidy

.PHONY: go-test
go-test:
	$(info [$@] running Go test suites...)
	go test -count=1 -race ./...
