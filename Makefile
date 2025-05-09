# Makefile for Go application

# Project settings
APP_NAME := sso
APP_PATH := cmd/sso/main.go
BUILD_DIR := bin
BINARY := $(BUILD_DIR)/$(APP_NAME)
CONFIG_FILE ?= config/local.yaml  # Default config file

# Go settings
GO := go
GO_BUILD_FLAGS := -v
GO_RUN_FLAGS :=
GO_TEST_FLAGS :=

.PHONY: all build run clean test

all: build

## Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GO_BUILD_FLAGS) -o $(BINARY) $(APP_PATH)
	@echo "Build complete: $(BINARY)"

## Run the application
start:
	@echo "Starting $(APP_NAME)..."
	$(GO) run $(GO_RUN_FLAGS) $(APP_PATH) --config $(CONFIG_FILE)

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## Run generate
generate:
	@echo "Running generate..."
	$(GO) generate ./...
	@echo "Generate complete"

## Run tests
unit-test:
	@echo "Running tests..."
	$(GO) test $(GO_TEST_FLAGS) ./internal/...
	@echo "Tests complete"

## Run test with cover
unit-test-cover:
	@echo "Running tests..."
	$(GO) test ./internal/... -coverprofile=coverage.txt
	$(GO) tool cover -html coverage.txt -o index.html
	@echo "Tests complete"

## Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	@echo "Dependencies installed"

## Show help
help:
	@echo "Available targets:"
	@echo "  build              - compile the application"
	@echo "  start              - run the application"
	@echo "  unit-test   	    - run tests"
	@echo "  unit-test-cover    - run tests with cover"
	@echo "  clean              - remove build artifacts"
	@echo "  generate           - generate"
	@echo "  deps               - install dependencies"
	@echo "  help               - show this help message"