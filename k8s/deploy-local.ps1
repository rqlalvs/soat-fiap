# PowerShell script for SOAT FIAP API deployment

Write-Host "🚀 Setting up SOAT FIAP API deployment..." -ForegroundColor Green

# Check if Docker is running
try {
    docker info | Out-Null
} catch {
    Write-Host "❌ Docker is not running. Please start Docker and try again." -ForegroundColor Red
    exit 1
}

# Check if kubectl is available
try {
    kubectl version --client | Out-Null
} catch {
    Write-Host "❌ kubectl is not installed. Please install kubectl and try again." -ForegroundColor Red
    exit 1
}

# Start local registry if not running
Write-Host "📦 Setting up local registry..." -ForegroundColor Yellow
$registryRunning = docker ps | Select-String "registry:2"
if (-not $registryRunning) {
    Write-Host "Starting local registry..." -ForegroundColor Yellow
    docker run -d --name registry -p 5000:5000 registry:2
} else {
    Write-Host "Local registry is already running." -ForegroundColor Green
}

# Build and push Docker image
Write-Host "🔨 Building Docker image..." -ForegroundColor Yellow
docker build -t localhost:5000/soat-fiap:latest ..

Write-Host "📤 Pushing image to local registry..." -ForegroundColor Yellow
docker push localhost:5000/soat-fiap:latest

# Update deployment to use local registry
Write-Host "📝 Updating deployment configuration..." -ForegroundColor Yellow
$content = Get-Content "app-deployment.yaml"
$content = $content -replace "image: soat-fiap:latest", "image: localhost:5000/soat-fiap:latest"
$content | Set-Content "app-deployment.yaml"

# Create namespace if it doesn't exist
Write-Host "🏗️  Creating namespace..." -ForegroundColor Yellow
kubectl create namespace soat-fiap --dry-run=client -o yaml | kubectl apply -f -

# Apply Kubernetes configurations
Write-Host "🚀 Deploying to Kubernetes..." -ForegroundColor Yellow
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f mysql-pvc.yaml
kubectl apply -f mysql-deployment.yaml
kubectl apply -f mysql-service.yaml
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml

Write-Host "⏳ Waiting for pods to be ready..." -ForegroundColor Yellow
kubectl wait --for=condition=ready pod -l app=soat-fiap -n soat-fiap --timeout=300s

Write-Host "✅ Deployment complete!" -ForegroundColor Green
Write-Host "🌐 Access your application at: http://localhost:30080" -ForegroundColor Cyan
Write-Host "📊 Check status with: kubectl get pods -n soat-fiap" -ForegroundColor Cyan 