# Go-Variablen
GO := go
GOW := gow
APP_NAME := server
CMD_DIR := ./cmd/server
BUILD_DIR := ./bin

# Standardziele
.PHONY: all build run test clean watch

# Standardziel: Build
all: build

# Build des Projekts
build:
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/

# Projekt ausführen
run:
	$(GO) run $(CMD_DIR)/

# Projekt ausführen mit Go Watch (GOW)
watch:
	$(GOW) run $(CMD_DIR)/

# Tests ausführen
test:
	$(GO) test ./... -v

# Linter ausführen (z. B. mit golangci-lint)
lint:
	golangci-lint run ./...

# Clean: Entfernt generierte Dateien
clean:
	rm -rf $(BUILD_DIR)
