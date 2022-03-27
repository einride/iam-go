# Code generated by go.einride.tech/sage. DO NOT EDIT.
# To learn more, see .sage/sagefile.go and https://github.com/einride/sage.

.DEFAULT_GOAL := all

sagefile := .sage/bin/sagefile

$(sagefile): .sage/go.mod .sage/*.go
	@cd .sage && go mod tidy && go run .

.PHONY: sage
sage:
	@git clean -fxq $(sagefile)
	@$(MAKE) $(sagefile)

.PHONY: update-sage
update-sage:
	@cd .sage && go get -d go.einride.tech/sage@latest && go mod tidy && go run .

.PHONY: clean-sage
clean-sage:
	@git clean -fdx .sage/tools .sage/bin .sage/build

.PHONY: all
all: $(sagefile)
	@$(sagefile) All

.PHONY: convco-check
convco-check: $(sagefile)
	@$(sagefile) ConvcoCheck

.PHONY: format-markdown
format-markdown: $(sagefile)
	@$(sagefile) FormatMarkdown

.PHONY: format-yaml
format-yaml: $(sagefile)
	@$(sagefile) FormatYAML

.PHONY: git-verify-no-diff
git-verify-no-diff: $(sagefile)
	@$(sagefile) GitVerifyNoDiff

.PHONY: go-iamctl
go-iamctl: $(sagefile)
	@$(sagefile) GoIamctl

.PHONY: go-lint
go-lint: $(sagefile)
	@$(sagefile) GoLint

.PHONY: go-mod-tidy
go-mod-tidy: $(sagefile)
	@$(sagefile) GoModTidy

.PHONY: go-review
go-review: $(sagefile)
	@$(sagefile) GoReview

.PHONY: go-test
go-test: $(sagefile)
	@$(sagefile) GoTest

.PHONY: proto
proto:
	$(MAKE) -C proto -f Makefile
