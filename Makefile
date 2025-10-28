# Makefile for linuxwave project
# Variables for maintainability
GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# Binary names
SERVICE_BIN := linuxwave-service
PAM_BIN := linuxwave-pam
CLI_BIN := linuxwave-cli
ENROLL_BIN := linuxwave-enroll

# Build directories
BUILD_DIR := bin
CMD_DIR := cmd

# Installation paths
INSTALL_PREFIX := /usr/local
BIN_INSTALL_DIR := $(INSTALL_PREFIX)/bin
SYSTEMD_INSTALL_DIR := /etc/systemd/system

# All binaries
BINARIES := $(SERVICE_BIN) $(PAM_BIN) $(CLI_BIN) $(ENROLL_BIN)

.PHONY: all build test install clean lint help

# Default target
all: build

# Build all binaries
build:
	@echo "Building all binaries..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(SERVICE_BIN) ./$(CMD_DIR)/$(SERVICE_BIN)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(PAM_BIN) ./$(CMD_DIR)/$(PAM_BIN)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(CLI_BIN) ./$(CMD_DIR)/$(CLI_BIN)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(ENROLL_BIN) ./$(CMD_DIR)/$(ENROLL_BIN)
	@echo "Build complete: $(BINARIES)"

# Run all tests
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# Install binaries and systemd unit files
install: build
	@echo "Installing binaries to $(BIN_INSTALL_DIR)..."
	@mkdir -p $(BIN_INSTALL_DIR)
	@install -m 755 $(BUILD_DIR)/$(SERVICE_BIN) $(BIN_INSTALL_DIR)/
	@install -m 755 $(BUILD_DIR)/$(PAM_BIN) $(BIN_INSTALL_DIR)/
	@install -m 755 $(BUILD_DIR)/$(CLI_BIN) $(BIN_INSTALL_DIR)/
	@install -m 755 $(BUILD_DIR)/$(ENROLL_BIN) $(BIN_INSTALL_DIR)/
	@if [ -d systemd ] && [ -n "$$(ls -A systemd/*.service 2>/dev/null)" ]; then \
		echo "Installing systemd unit files to $(SYSTEMD_INSTALL_DIR)..."; \
		mkdir -p $(SYSTEMD_INSTALL_DIR); \
		install -m 644 systemd/*.service $(SYSTEMD_INSTALL_DIR)/; \
	fi
	@echo "Installation complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARIES)
	@$(GO) clean
	@echo "Clean complete"

# Run linting
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, using go vet..."; \
		$(GO) vet ./...; \
	fi

# Help target
help:
	@echo "linuxwave Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build    - Compile all binaries"
	@echo "  test     - Run all tests"
	@echo "  install  - Install binaries and systemd unit files"
	@echo "  clean    - Remove build artifacts"
	@echo "  lint     - Run linting (golangci-lint or go vet)"
	@echo "  help     - Show this help message"
