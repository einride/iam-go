.PHONY: all
all: \
	buf-lint \
	buf-generate-iam \
	buf-generate-example \
	go-test \
	go-mod-tidy

include tools/buf/rules.mk
include tools/protoc-gen-go-grpc/rules.mk

.PHONY: build/protoc-gen-go-iam
build/protoc-gen-go-iam:
	$(info [$@] rebuilding plugin...)
	@go build -o $@ ./cmd/protoc-gen-go-iam
	@touch $@

build/protoc-gen-go: go.mod
	$(info [$@] rebuilding plugin...)
	@go build -o $@ google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: buf-lint
buf-lint: $(buf)
	$(info [$@] linting proto files...)
	@$(buf) lint

.PHONY: buf-generate-iam
buf-generate-iam: $(buf) build/protoc-gen-go
	$(info [$@] generating iam proto stubs...)
	@rm -rf proto/gen/einride/iam/v1
	@$(buf) generate --path proto/src/einride/iam/v1 --template buf.gen.iam.yaml

.PHONY: buf-generate-example
buf-generate-example: $(buf) build/protoc-gen-go build/protoc-gen-go-iam $(protoc_gen_go_grpc)
	$(info [$@] generating proto stubs...)
	@rm -rf proto/gen/einride/iam/example/v1
	@$(buf) generate --path proto/src/einride/iam/example/v1 --template buf.gen.example.yaml

.PHONY: spanner-aip-go-generate
spanner-aip-go-generate:
	$(info [$@] generating Spanner database APIs...)
	@go run go.einride.tech/spanner-aip/cmd/spanner-aip-go generate

.PHONY: go-mod-tidy
go-mod-tidy:
	$(info [$@] tidying Go module files...)
	@go mod tidy

.PHONY: go-test
go-test:
	$(info [$@] running Go test suites...)
	go test -count=1 -race ./...
