# Go-Variablen
GO := go
GOW := gow
APP_NAME := server
CMD_DIR := ./cmd
BUILD_DIR := ./bin
GOARCH := $(if $(GOARCH),$(GOARCH),$(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/'))

# Standardziele
.PHONY: all build run test clean watch generate migrate ui ui-build

# Standardziel: Build
all: build ui-build

# Build des Projekts
build: generate
	@echo "Building $(APP_NAME) ($(GOARCH))..."
	CGO_ENABLED=0 $(GO) build -o $(BUILD_DIR)/$(APP_NAME) -ldflags="-s -w" $(CMD_DIR)/

# Projekt ausführen
run: migrate
	@echo "Running $(APP_NAME) ($(GOARCH))..."
	$(GO) run $(CMD_DIR)/server/ run

migrate:
	@echo "Running migrations..."
	$(GO) run $(CMD_DIR)/migrate/

# Projekt ausführen mit Go Watch (GOW)
watch: migrate
	@echo "Running $(APP_NAME) ($(GOARCH)) with GOW..."
	$(GOW) run $(CMD_DIR)/server/ run

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

ui:
	@echo "Running UI in dev mode..."
	npm run dev --prefix ui

ui-build:
	@echo "Building UI..."
	npm install --prefix ui
	npm run build --prefix ui
