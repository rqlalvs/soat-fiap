# Simple Makefile for a Go project

## Cross-team reproducible targets
.PHONY: all build run test clean watch docker-build docker-push k8s-apply k8s-delete tag

# Configurable image reference (override REGISTRY on CLI/CI)
REGISTRY ?= docker.io/youruser
IMAGE ?= soat-fiap
TAG ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo latest)
IMAGE_REF := $(REGISTRY)/$(IMAGE):$(TAG)

# Default target
all: build test

# Build the application (local)
build:
	@echo "Building..."
	@go build -o bin/app ./cmd/api

# Run the application (local)
run:
	@go run ./cmd/api

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bin/app || true

# Live Reload (optional)
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

# Docker image build/push with immutable tag and latest tag
docker-build:
	docker build -t $(IMAGE_REF) -t $(REGISTRY)/$(IMAGE):latest -f dockerfile .

docker-push:
	docker push $(IMAGE_REF)
	docker push $(REGISTRY)/$(IMAGE):latest

# Kubernetes apply/delete
k8s-apply:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml -f k8s/secret.yaml
	kubectl apply -f k8s/mysql-deployment.yaml
	kubectl apply -f k8s/app-deployment.yaml
	kubectl apply -f k8s/hpa.yaml

k8s-delete:
	kubectl delete -f k8s/hpa.yaml --ignore-not-found
	kubectl delete -f k8s/app-deployment.yaml --ignore-not-found
	kubectl delete -f k8s/mysql-deployment.yaml --ignore-not-found
	kubectl delete -f k8s/configmap.yaml -f k8s/secret.yaml --ignore-not-found
	kubectl delete -f k8s/namespace.yaml --ignore-not-found
