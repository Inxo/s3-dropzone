# Makefile for Go project

# Environment variables
ENV_FILE := .env

# Load environment variables from .env file
#include $(ENV_FILE)
export

# Variables
BINARY_NAME := dropZone
GO_FILES := $(wildcard *.go)

# Commands
.PHONY: build run clean

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o build/$(BINARY_NAME) -ldflags='-s -w' ./$(GO_FILES)

buildarm:
	@echo "Building $(BINARY_NAME)..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o build/$(BINARY_NAME) -ldflags='-s -w' ./$(GO_FILES)


package:
	@echo "Package $(BINARY_NAME)..."
	@fyne package -os darwin -icon Icon.png -name $(BINARY_NAME)App ./$(GO_FILES)

build-win:
	@echo "Building windows $(BINARY_NAME)..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	go build -a -ldflags '-extldflags "-static" -s -w' -o build/$(BINARY_NAME).exe  ./$(GO_FILES)

test:
	@echo "Testing $(BINARY_NAME)..."
	@go test ./src/$(GO_FILES)

lint:
	@golangci-lint run -v

run:
	@echo "Running $(BINARY_NAME)..."
	@./build/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# Target to run the application
start: build run

# Target to build and run the application
restart: clean start

# Help target
help:
	@echo "Available targets:"
	@echo "  build    : Build the Go application"
	@echo "  run      : Run the Go application"
	@echo "  clean    : Clean up generated files"
	@echo "  docs     : Generated documentation files"
	@echo "  start    : Build and run the Go application"
	@echo "  restart  : Clean, build, and run the Go application"
