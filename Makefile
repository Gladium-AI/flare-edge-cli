SHELL := /bin/sh

GO ?= go
BINARY ?= flare-edge-cli
PACKAGE ?= ./cmd/flare-edge-cli
BINDIR ?= bin
BUILD_OUTPUT := $(BINDIR)/$(BINARY)

# Prefer standard user-local bin directories and fall back to ~/.local/bin.
INSTALL_DIR ?= $(shell \
	if [ -n "$$XDG_BIN_HOME" ]; then \
		printf '%s\n' "$$XDG_BIN_HOME"; \
	elif [ -d "$$HOME/.local/bin" ]; then \
		printf '%s\n' "$$HOME/.local/bin"; \
	elif [ -d "$$HOME/bin" ]; then \
		printf '%s\n' "$$HOME/bin"; \
	else \
		printf '%s\n' "$$HOME/.local/bin"; \
	fi)

.PHONY: build test install

build:
	@mkdir -p "$(BINDIR)"
	$(GO) build -o "$(BUILD_OUTPUT)" $(PACKAGE)

test:
	$(GO) test ./...

install: build
	@mkdir -p "$(INSTALL_DIR)"
	install -m 0755 "$(BUILD_OUTPUT)" "$(INSTALL_DIR)/$(BINARY)"
	@printf 'Installed %s to %s\n' "$(BINARY)" "$(INSTALL_DIR)"
