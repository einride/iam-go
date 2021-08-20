SHELL := /bin/bash

.PHONY: all
all: \
	proto \
	spanner-generate \
	go-lint \
	go-test \
	go-mod-tidy \
	go-install-iamctl \
	git-verify-nodiff

include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/semantic-release/rules.mk

.PHONY: proto
proto:
	$(info [$@] building protos...)
	@make -C proto

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

.PHONY: go-install-iamctl
go-install-iamctl:
	$(info [$@] installing iamctl...)
	@cd ./cmd/iamctl && go install .
