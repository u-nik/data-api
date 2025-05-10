# Go-Variablen
GO := go
GOW := gow
APP_NAME := server
CMD_DIR := ./cmd/server
BUILD_DIR := ./bin
GOARCH := $(if $(GOARCH),$(GOARCH),$(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/'))

# Standardziele
.PHONY: all build run test clean watch generate

# Standardziel: Build
all: build

# Build des Projekts
build: generate
	@echo "Building $(APP_NAME) ($(GOARCH))..."
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME)-$(GOARCH) $(CMD_DIR)/

# Projekt ausführen
run:
	$(GO) run $(CMD_DIR)/

# Projekt ausführen mit Go Watch (GOW)
watch:
	$(GOW) run $(CMD_DIR)/

# Tests ausführen
test:
	@echo "Running tests for GOARCH=$(GOARCH)..."
	$(GO) test ./... -v 2>&1 | go-junit-report > test-results.xml

# Linter ausführen (z. B. mit golangci-lint)
lint:
	golangci-lint run ./...

# Clean: Entfernt generierte Dateien
clean:
	rm -rf $(BUILD_DIR)

# Generate: Generiert Code (z. B. mit go:generate)
generate:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
	$(GO) generate ./...
