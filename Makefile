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
GO_TEST_FLAGS := -v -race

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

## Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GO_TEST_FLAGS) ./...
	@echo "Tests complete"

## Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	@echo "Dependencies installed"

## Show help
help:
	@echo "Available targets:"
	@echo "  build   - compile the application"
	@echo "  start     - run the application"
	@echo "  test    - run tests"
	@echo "  clean   - remove build artifacts"
	@echo "  deps    - install dependencies"
	@echo "  help    - show this help message"