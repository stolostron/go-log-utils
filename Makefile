# Copyright (c) 2022 Red Hat, Inc.
# Copyright Contributors to the Open Cluster Management project

export PATH := $(PWD)/bin:$(PATH)
export GOBIN := $(PWD)/bin

############################################################
# format section
############################################################

GCI ?= $(GOBIN)/gci
$(GCI):
	go install github.com/daixiang0/gci@v0.2.9

GOFUMPT ?= $(GOBIN)/gofumpt
$(GOFUMPT):
	go install mvdan.cc/gofumpt@v0.2.0

.PHONY: fmt
fmt: $(GCI) $(GOFUMPT)
	find . -not \( -path "./.go" -prune \) -name "*.go" | xargs gofmt -s -w
	find . -not \( -path "./.go" -prune \) -name "*.go" | xargs gci -w -local "$(shell cat go.mod | head -1 | cut -d " " -f 2)"
	find . -not \( -path "./.go" -prune \) -name "*.go" | xargs gofumpt -l -w

############################################################
# lint section
############################################################

GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run
