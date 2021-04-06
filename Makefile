.PHONY: all
all: \
	buf-lint \
	buf-generate

include tools/buf/rules.mk
include tools/protobuf-go/rules.mk

.PHONY: buf-lint
buf-lint: $(buf)
	$(info [$@] linting proto files...)
	@$(buf) lint


.PHONY: buf-generate
buf-generate: $(buf) $(protoc_gen_go)
	$(info [$@] generating proto stubs...)
	@$(buf) generate
