# Go-Variablen
GO := go
GOW := gow
APP_NAME := server
CMD_DIR := ./cmd/server
BUILD_DIR := ./bin
GOARCH ?= $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

# Standardziele
.PHONY: all build run test clean watch generate

# Standardziel: Build
all: build

# Build des Projekts
build: generate test lint
	@echo "Building $(APP_NAME) ($(GOARCH))..."
	GOARCH=$(GOARCH) $(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/

# Projekt ausführen
run:
	GOARCH=$(GOARCH) $(GO) run $(CMD_DIR)/

# Projekt ausführen mit Go Watch (GOW)
watch:
	GOARCH=$(GOARCH) $(GOW) run $(CMD_DIR)/

# Tests ausführen
test:
	@echo "Running tests for GOARCH=$(GOARCH)..."
    GOARCH=$(GOARCH) $(GO) test ./... -v 2>&1 | go-junit-report > test-results.xml
    @echo "Test results written to test-results.xml"

# Linter ausführen (z. B. mit golangci-lint)
lint:
	golangci-lint run ./...

# Clean: Entfernt generierte Dateien
clean:
	rm -rf $(BUILD_DIR)

# Generate: Generiert Code (z. B. mit go:generate)
generate:
	GOARCH=$(GOARCH) $(GO) mod tidy
	GOARCH=$(GOARCH) $(GO) mod vendor
	GOARCH=$(GOARCH) $(GO) mod verify
	GOARCH=$(GOARCH) $(GO) generate ./...
