SHELL := /bin/bash
PROJECT := mail-testserver

PKG := ./...
GO := go

GOFUMPT_VERSION := v0.9.2
GOVULNCHECK_VERSION := v1.1.4


.PHONY: help init tools fmt format lint test vuln vet tidy ci precommit

help:
	@echo "Targets:"
	@echo "  init       - install tools + setup hooks"
	@echo "  fmt        - gofmt/goimports (and gofumpt if installed)"
	@echo "  lint       - golangci-lint"
	@echo "  test       - go test"
	@echo "  vuln       - govulncheck"
	@echo "  ci         - fmt + lint + test + vuln"

init: tools precommit

precommit:
	@git config core.hooksPath .githooks
	@chmod +x .githooks/pre-commit
	@echo "Enabled git hooks at .githooks/"

tools:
	@echo "Installing tools..."
	@echo install goimports
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo install gofumpt
	@$(GO) install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	@echo install govulncheck
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	@echo install golangci-lint



fmt: format
format:
	@echo "Formatting..."
	@gofmt -w $(find . -name '*.go')
	@goimports -w $(find . -name '*.go')



test: 
	@echo "Running tests..."
	@$(GO) test -v ./internal/httpapi
