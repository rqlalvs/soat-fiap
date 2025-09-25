# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main.exe cmd/api/main.go

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t soat-fiap:latest .

# Build and push to local registry (for development)
docker-build-local:
	@echo "Building and pushing to local registry..."
	docker build -t localhost:5000/soat-fiap:latest .
	docker push localhost:5000/soat-fiap:latest

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

.PHONY: all build docker-build docker-build-local run test clean watch
