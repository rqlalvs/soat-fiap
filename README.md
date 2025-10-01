# SOAT FIAP - Sistema de Lanchonete

## 🏗️ Arquitetura

### **Visão Geral**
Sistema de autoatendimento para lanchonete com arquitetura em microserviços, utilizando Google Cloud Platform (GCP) e Kubernetes.

### **Componentes da Arquitetura**

#### **1. API Gateway (Nginx Ingress)**
- **Função**: Recebe requisições externas e roteia para os serviços
- **Tecnologia**: Nginx Ingress Controller no GKE
- **Justificativa**: Solução nativa do Kubernetes, alta performance e flexibilidade

#### **2. Cloud Function (Autenticação)**
- **Função**: Autentica clientes por CPF
- **Tecnologia**: Google Cloud Functions (Go 1.22)
- **Justificativa**: Serverless, escalabilidade automática, custo-benefício

#### **3. Aplicação Principal (Kubernetes)**
- **Função**: API REST para gerenciamento de pedidos
- **Tecnologia**: Go + GKE
- **Justificativa**: Controle total, alta disponibilidade, escalabilidade

#### **4. Banco de Dados (Cloud SQL)**
- **Função**: Armazenamento persistente
- **Tecnologia**: MySQL 8.0 (Cloud SQL)
- **Justificativa**: Banco gerenciado, backup automático, alta disponibilidade

### **Fluxo de Dados**

```
Cliente → API Gateway → Cloud Function (Auth) → Aplicação K8s → Cloud SQL
```

### **Segurança**
- **JWT** para autenticação
- **Secrets** gerenciados pelo GCP
- **Network policies** no Kubernetes
- **RBAC** configurado

## 🚀 Deploy

### **Pré-requisitos**
- Google Cloud SDK
- kubectl configurado
- Terraform (para infraestrutura)

### **1. Deploy da Infraestrutura**
```bash
# Banco de dados
cd soat-fiap-infra-db
terraform init
terraform plan
terraform apply

# Kubernetes
cd soat-fiap-infra-k8s
terraform init
terraform plan
terraform apply
```

### **2. Deploy da Aplicação**
```bash
# Aplicação principal
cd soat-fiap
git push origin main  # Trigger automático

# Cloud Function
cd soat-fiap-auth-function
git push origin main  # Trigger automático
```

### **3. Verificação**
```bash
# Health check
curl https://soat-fiap-gateway-dtxpqilj.uc.gateway.dev/api/v1/health

# Autenticação
curl -X POST https://us-central1-sylvan-willow-473218-f8.cloudfunctions.net/authByCPF \
  -H "Content-Type: application/json" \
  -d '{"cpf":"12345678901"}'
```

## 📊 Monitoramento

### **Logs**
- **Cloud Function**: Cloud Logging
- **Kubernetes**: kubectl logs
- **Banco**: Cloud SQL logs

### **Métricas**
- **CPU/Memory**: GKE metrics
- **Requests**: Nginx metrics
- **Database**: Cloud SQL metrics

## 🔧 Desenvolvimento

### **Estrutura do Projeto**
```
soat-fiap/                 # Aplicação principal
├── cmd/api/              # Entry point
├── internal/             # Código interno
│   ├── adapters/         # Adapters (handlers, repositories)
│   ├── core/             # Domain logic
│   └── services/         # Business services
├── pkg/                  # Packages compartilhados
└── k8s/                  # Kubernetes manifests

soat-fiap-auth-function/  # Cloud Function
├── main.go              # Function entry point
└── go.mod               # Dependencies

soat-fiap-infra-k8s/     # Kubernetes infrastructure
├── terraform/           # Terraform configuration
└── k8s/                 # K8s manifests

soat-fiap-infra-db/      # Database infrastructure
├── terraform/           # Terraform configuration
└── sql/                 # Database scripts
```

### **CI/CD**
- **GitHub Actions** para cada repositório
- **Deploy automático** na main branch
- **Pull Request** obrigatório
- **Secrets** gerenciados pelo GitHub

## 🎯 Funcionalidades

### **Cliente**
- ✅ Autenticação por CPF
- ✅ Montagem de combo (Lanche + Acompanhamento + Bebida)
- ✅ Pagamento via QR Code (Mercado Pago)
- ✅ Acompanhamento do pedido

### **Administrativo**
- ✅ Gerenciamento de clientes
- ✅ Gerenciamento de produtos
- ✅ Acompanhamento de pedidos
- ✅ Relatórios e métricas

## 📈 Escalabilidade

### **Horizontal**
- **Kubernetes HPA** para pods
- **Cloud Function** auto-scaling
- **Cloud SQL** read replicas

### **Vertical**
- **Resource limits** configurados
- **Monitoring** para ajustes

## 🔒 Segurança

### **Autenticação**
- **JWT** tokens
- **CPF** validation
- **Rate limiting**

### **Autorização**
- **RBAC** no Kubernetes
- **Service accounts** configurados
- **Network policies**

## 📚 Documentação Adicional

- [Arquitetura Detalhada](docs/architecture.md)
- [API Documentation](docs/api.md)
- [Deployment Guide](docs/deployment.md)
- [Troubleshooting](docs/troubleshooting.md)

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

