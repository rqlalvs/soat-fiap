#!/bin/bash

echo "🚀 Setting up SOAT FIAP API deployment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl is not installed. Please install kubectl and try again."
    exit 1
fi

# Start local registry if not running
echo "📦 Setting up local registry..."
if ! docker ps | grep -q "registry:2"; then
    echo "Starting local registry..."
    docker run -d --name registry -p 5000:5000 registry:2
else
    echo "Local registry is already running."
fi

# Build and push Docker image
echo "🔨 Building Docker image..."
docker build -t localhost:5000/soat-fiap:latest ..

echo "📤 Pushing image to local registry..."
docker push localhost:5000/soat-fiap:latest

# Update deployment to use local registry
echo "📝 Updating deployment configuration..."
sed -i 's|image: soat-fiap:latest|image: localhost:5000/soat-fiap:latest|g' app-deployment.yaml

# Create namespace if it doesn't exist
echo "🏗️  Creating namespace..."
kubectl create namespace soat-fiap --dry-run=client -o yaml | kubectl apply -f -

# Apply Kubernetes configurations
echo "🚀 Deploying to Kubernetes..."
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f mysql-pvc.yaml
kubectl apply -f mysql-deployment.yaml
kubectl apply -f mysql-service.yaml
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml

echo "⏳ Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=soat-fiap -n soat-fiap --timeout=300s

echo "✅ Deployment complete!"
echo "🌐 Access your application at: http://localhost:30080"
echo "📊 Check status with: kubectl get pods -n soat-fiap" 