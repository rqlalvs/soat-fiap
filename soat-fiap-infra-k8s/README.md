# SOAT FIAP - Infraestrutura Kubernetes

## 🎯 Função
Terraform para provisionar e gerenciar a infraestrutura Kubernetes no Google Cloud Platform.

## 🏗️ Arquitetura

### **Componentes**
- **GKE Cluster**: Google Kubernetes Engine
- **Nginx Ingress**: API Gateway
- **HPA**: Horizontal Pod Autoscaler
- **Secrets**: Gerenciamento de secrets
- **ConfigMaps**: Configurações da aplicação

### **Recursos Provisionados**
```
GKE Cluster
├── Node Pool (Auto-scaling)
├── Nginx Ingress Controller
├── Application Namespace
├── Secrets (Database, JWT)
├── ConfigMaps (App Config)
└── HPA (Horizontal Pod Autoscaler)
```

## 🚀 Deploy

### **Pré-requisitos**
- Google Cloud SDK
- Terraform >= 1.0
- kubectl configurado

### **1. Configurar GCP**
```bash
# Autenticação
gcloud auth login
gcloud auth application-default login

# Configurar projeto
gcloud config set project sylvan-willow-473218-f8
```

### **2. Deploy com Terraform**
```bash
# Inicializar Terraform
terraform init

# Planejar infraestrutura
terraform plan

# Aplicar mudanças
terraform apply
```

### **3. Configurar kubectl**
```bash
# Obter credenciais do cluster
gcloud container clusters get-credentials soat-fiap --region=us-central1

# Verificar cluster
kubectl get nodes
```


### **Recursos Criados**
- **GKE Cluster**: `soat-fiap`
- **Node Pool**: Auto-scaling (1-10 nodes)
- **Nginx Ingress**: Controller instalado
- **Namespace**: `soat-fiap`
- **Secrets**: Database e JWT
- **ConfigMaps**: Configurações da app

## 📊 Monitoramento

### **Métricas**
- **Cluster**: CPU, Memory, Disk
- **Nodes**: Resource utilization
- **Pods**: Application metrics
- **Ingress**: Request metrics

### **Logs**
- **Cluster**: GKE logs
- **Application**: Pod logs
- **Ingress**: Nginx logs

## 🔒 Segurança

### **RBAC**
- **Service Accounts**: Configurados
- **Roles**: Permissões mínimas
- **Network Policies**: Isolamento de rede

### **Secrets**
- **Database**: Credenciais do Cloud SQL
- **JWT**: Chave secreta
- **TLS**: Certificados SSL

## 🚀 CI/CD

### **GitHub Actions**
- **Trigger**: Push para main branch
- **Terraform**: Plan e Apply
- **Kubectl**: Deploy da aplicação
- **Verification**: Health checks

### **Pipeline**
1. **Checkout**: Código do repositório
2. **Setup**: Terraform e kubectl
3. **Plan**: Terraform plan
4. **Apply**: Terraform apply
5. **Deploy**: kubectl apply
6. **Verify**: Health checks

## 🔧 Desenvolvimento

### **Estrutura**
```
soat-fiap-infra-k8s/
├── terraform/
│   ├── main.tf              # Configuração principal
│   ├── variables.tf         # Variáveis
│   ├── outputs.tf          # Outputs
│   └── terraform.tfvars    # Valores das variáveis
├── k8s/
│   ├── namespace.yaml      # Namespace
│   ├── configmap.yaml      # Configurações
│   ├── secret.yaml         # Secrets
│   ├── app-deployment.yaml # Deploy da aplicação
│   ├── hpa.yaml            # Horizontal Pod Autoscaler
│   └── nginx-ingress.yaml  # Ingress controller
└── .github/
    └── workflows/
        └── deploy.yml       # CI/CD pipeline
```

### **Comandos Úteis**
```bash
# Status do cluster
kubectl get nodes
kubectl get pods -n soat-fiap

# Logs da aplicação
kubectl logs -l app=soat-fiap-app -n soat-fiap

# Port forward para teste
kubectl port-forward svc/soat-fiap-app 8080:80 -n soat-fiap
```

## 📈 Escalabilidade

### **Horizontal Pod Autoscaler**
- **CPU**: 70% threshold
- **Memory**: 80% threshold
- **Min Pods**: 2
- **Max Pods**: 10

### **Node Auto-scaling**
- **Min Nodes**: 1
- **Max Nodes**: 10
- **Machine Type**: e2-medium

## 🔍 Troubleshooting

### **Common Issues**
- **Cluster Creation**: Verificar quotas
- **Node Scaling**: Verificar recursos
- **Ingress**: Verificar DNS

### **Debug**
```bash
# Status do cluster
kubectl get nodes
kubectl describe nodes

# Status dos pods
kubectl get pods -n soat-fiap
kubectl describe pods -n soat-fiap

# Logs do ingress
kubectl logs -n ingress-nginx
```

## 🧹 Limpeza

### **Destruir Infraestrutura**
```bash
# Destruir com Terraform
terraform destroy

# Verificar recursos
gcloud container clusters list
```
