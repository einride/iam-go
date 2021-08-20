SHELL := /usr/bin/env bash

.PHONY: all
all: \
	clang-format \
	buf-lint \
	buf-breaking \
	api-linter-lint \
	buf-generate \
	buf-generate-example

include tools/api-linter/rules.mk
include tools/buf/rules.mk
include tools/clang-format/rules.mk
include tools/protoc-gen-go-grpc/rules.mk

build/protoc-gen-go: ../go.mod
	$(info [$@] building binary...)
	@go build -o $@ google.golang.org/protobuf/cmd/protoc-gen-go

build/protoc-gen-go-aip: ../go.mod
	$(info [$@] building binary...)
	@go build -o $@ go.einride.tech/aip/cmd/protoc-gen-go-aip

.PHONY: build/protoc-gen-go-iam
build/protoc-gen-go-iam:
	$(info [$@] building binary...)
	@go build -o $@ ../cmd/protoc-gen-go-iam

build/descriptor.pb: $(shell find einride -type f -name '*.proto')
	$(info [$@] generating proto descriptor...)
	@mkdir -p $(dir $@)
	@$(buf) build -o $@

.PHONY: buf-lint
buf-lint: $(buf)
	$(info [$@] linting proto schemas...)
	@$(buf) lint

.PHONY: buf-breaking
buf-breaking: $(buf)
	$(info [$@] checking for breaking changes against origin/master)
	@$(buf) breaking --against buf.build/einride/iam --path einride

.PHONY: clang-format
clang-format: proto_files = $(shell find einride -type f -name '*.proto')
clang-format: $(clang_format)
	$(info [$@] formatting proto files...)
	@$(clang_format) \
		-i \
		--style='{BasedOnStyle: Google, ColumnLimit: 0, Language: Proto}' \
		$(proto_files)

.PHONY: api-linter-lint
api-linter-lint: $(api_linter) build/descriptor.pb
	$(info [$@] linting gRPC APIs...)
	@$(api_linter) \
		--set-exit-status \
		--config api-linter.yaml \
		--descriptor-set-in build/descriptor.pb \
		$(shell find einride/iam/example -type f -name '*.proto')

.PHONY: buf-generate
buf-generate: $(buf) build/protoc-gen-go
	$(info [$@] generating protobuf stubs...)
	@$(buf) generate \
		--template buf.gen.yaml \
		--path einride

example_plugins := \
	$(protoc_gen_go_grpc) \
	build/protoc-gen-go \
	build/protoc-gen-go-aip \
	build/protoc-gen-go-iam

.PHONY: buf-generate
buf-generate-example: $(buf) $(example_plugins)
	$(info [$@] generating protobuf stubs...)
	@$(buf) generate \
		--template buf.gen.example.yaml \
		--path einride
