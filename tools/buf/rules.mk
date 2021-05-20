buf_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
buf_version := 0.41.0
buf := $(buf_cwd)/$(buf_version)/bin/buf
export PATH := $(dir $(buf)):$(PATH)

buf_arch = $(shell uname -s)-$(shell uname -m)
# enforce x86 arch if mac m1 until tool has official support
ifeq ($(buf_arch),Darwin-arm64)
buf_arch = Darwin-x86_64
endif
buf_bin_url := https://github.com/bufbuild/buf/releases/download/v$(buf_version)/buf-$(buf_arch)

$(buf): $(buf_cwd)/rules.mk
	$(info [buf] feching version $(buf_version)...)
	@mkdir -p $(dir $@)
	@curl -sSL $(buf_bin_url) -o $@
	@chmod +x $@
	@touch $@
