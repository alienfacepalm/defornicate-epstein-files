.PHONY: build test clean install run help lint fmt vet

# Build variables
BINARY_NAME=epstein-files-defornicator
MAIN_PATH=main.go
BUILD_DIR=bin
RELEASE_DIR=releases

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-release: ## Build release binaries for all platforms
	@echo "Building release binaries..."
	@if not exist $(RELEASE_DIR) mkdir $(RELEASE_DIR)
	@echo "Building Windows amd64..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(RELEASE_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Building Windows 386..."
	@GOOS=windows GOARCH=386 $(GOBUILD) -o $(RELEASE_DIR)/$(BINARY_NAME)-windows-386.exe $(MAIN_PATH)
	@echo "Building Linux amd64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(RELEASE_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "Building Linux 386..."
	@GOOS=linux GOARCH=386 $(GOBUILD) -o $(RELEASE_DIR)/$(BINARY_NAME)-linux-386 $(MAIN_PATH)
	@echo "Building Linux arm64..."
	@GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(RELEASE_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "Building macOS amd64..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(RELEASE_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "Building macOS arm64..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(RELEASE_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "Release binaries built in $(RELEASE_DIR)/"

build-windows: ## Build Windows binary
	@echo "Building Windows binary..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)

build-linux: ## Build Linux binary
	@echo "Building Linux binary..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

build-darwin: ## Build macOS binary
	@echo "Building macOS binary..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin $(MAIN_PATH)

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR) 2>nul || true
	@if exist $(RELEASE_DIR) rmdir /s /q $(RELEASE_DIR) 2>nul || true
	@if exist coverage.out del /q coverage.out 2>nul || true
	@if exist coverage.html del /q coverage.html 2>nul || true
	@echo "Clean complete"

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@echo "Formatting complete"

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "Vet complete"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	$(GOMOD) tidy
	@echo "Tidy complete"

install: build ## Install the binary to GOPATH/bin
	@echo "Installing..."
	$(GOCMD) install $(MAIN_PATH)
	@echo "Installation complete"

run: ## Run the application
	@echo "Running application..."
	$(GOCMD) run $(MAIN_PATH)

.DEFAULT_GOAL := help

