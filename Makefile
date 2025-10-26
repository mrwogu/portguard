.PHONY: build test clean install run docker-build docker-run help

BINARY_NAME=portguard
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=dist
INSTALL_DIR=/opt/portguard
CONFIG_DIR=/etc/portguard
LDFLAGS=-ldflags "-s -w -X main.appVersion=$(VERSION)"

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@echo "All builds complete in $(BUILD_DIR)/"

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	$(GOCMD) tool cover -html=coverage.txt

fmt: ## Format code
	$(GOFMT) ./...

vet: ## Run go vet
	$(GOVET) ./...

lint: ## Run golangci-lint
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

mod-download: ## Download dependencies
	$(GOMOD) download

mod-tidy: ## Tidy dependencies
	$(GOMOD) tidy

mod-verify: ## Verify dependencies
	$(GOMOD) verify

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.txt
	@echo "Clean complete"

run: build ## Build and run
	$(BUILD_DIR)/$(BINARY_NAME) --config config.yaml.example

install: build ## Install to system (requires root)
	@echo "Installing $(BINARY_NAME)..."
	@mkdir -p $(INSTALL_DIR)
	@mkdir -p $(CONFIG_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@if [ ! -f $(CONFIG_DIR)/config.yaml ]; then \
		cp config.yaml.example $(CONFIG_DIR)/config.yaml; \
		echo "Config installed to $(CONFIG_DIR)/config.yaml"; \
	else \
		echo "Config already exists at $(CONFIG_DIR)/config.yaml"; \
	fi
	@cp portguard.service /etc/systemd/system/
	@systemctl daemon-reload
	@echo "Installation complete. Start with: systemctl start portguard"

uninstall: ## Uninstall from system (requires root)
	@echo "Uninstalling $(BINARY_NAME)..."
	@systemctl stop portguard 2>/dev/null || true
	@systemctl disable portguard 2>/dev/null || true
	@rm -f /etc/systemd/system/portguard.service
	@rm -rf $(INSTALL_DIR)
	@systemctl daemon-reload
	@echo "Uninstall complete. Config remains at $(CONFIG_DIR)/"

version: ## Show version
	@echo "$(VERSION)"

.DEFAULT_GOAL := help
