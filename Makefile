# Copyright (c) 2022 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

export PATH := $(PWD)/bin:$(PATH)
export GOBIN := $(PWD)/bin

## CLI versions (with links to the latest releases)
# https://github.com/golangci/golangci-lint/releases/latest
GOLANGCI_VERSION := v1.64.8
# https://github.com/mvdan/gofumpt/releases/latest
GOFUMPT_VERSION := v0.7.0
# https://github.com/daixiang0/gci/releases/latest
GCI_VERSION := v0.13.5
# https://github.com/securego/gosec/releases/latest
GOSEC_VERSION := v2.22.2

# go-get-tool will 'go install' any package $1 and install it to LOCAL_BIN.
define go-get-tool
@set -e ;\
echo "Checking installation of $(1)" ;\
GOBIN=$(GOBIN) go install $(1)
endef

############################################################
# format section
############################################################

.PHONY: fmt-dependencies
fmt-dependencies:
	$(call go-get-tool,github.com/daixiang0/gci@$(GCI_VERSION))
	$(call go-get-tool,mvdan.cc/gofumpt@$(GOFUMPT_VERSION))

.PHONY: fmt
fmt: fmt-dependencies
	find . -not \( -path "./.go" -prune \) -name "*.go" | xargs gofmt -s -w
	find . -not \( -path "./.go" -prune \) -name "*.go" | xargs gci write -s standard -s default -s "prefix($(shell cat go.mod | head -1 | cut -d " " -f 2))"
	find . -not \( -path "./.go" -prune \) -name "*.go" | xargs gofumpt -l -w

############################################################
# lint section
############################################################

.PHONY: lint-dependencies
lint-dependencies:
	$(call go-get-tool,github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION))

GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
.PHONY: lint
lint: lint-dependencies
	$(GOLANGCI_LINT) run

############################################################
# test section
############################################################
GOSEC = $(GOBIN)/gosec

.PHONY: gosec
gosec:
	$(call go-get-tool,github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION))

.PHONY: gosec-scan
gosec-scan: gosec
	$(GOSEC) -fmt sonarqube -out gosec.json -no-fail -exclude-dir=.go ./...
