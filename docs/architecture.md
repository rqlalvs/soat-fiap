# Arquitetura do Sistema SOAT FIAP

## 🏗️ Diagrama da Arquitetura

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                CLIENTE                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   Mobile    │  │   Web App   │  │   Kiosk     │  │   Admin     │          │
│  │    App      │  │             │  │   Terminal  │  │   Panel     │          │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              API GATEWAY                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐  │
│  │                    Nginx Ingress Controller                              │  │
│  │  • Rate Limiting                                                        │  │
│  │  • SSL Termination                                                      │  │
│  │  • Load Balancing                                                       │  │
│  │  • Request Routing                                                      │  │
│  └─────────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           AUTENTICAÇÃO                                         │
│  ┌─────────────────────────────────────────────────────────────────────────┐  │
│  │                    Google Cloud Function                                │  │
│  │  • CPF Validation                                                      │  │
│  │  • JWT Token Generation                                               │  │
│  │  • Client Authentication                                               │  │
│  │  • Auto-scaling                                                        │  │
│  └─────────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        APLICAÇÃO PRINCIPAL                                     │
│  ┌─────────────────────────────────────────────────────────────────────────┐  │
│  │                    Kubernetes Cluster (GKE)                            │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │   Pedidos   │  │  Produtos   │  │  Clientes   │  │  Pagamentos │  │  │
│  │  │   Service   │  │   Service   │  │   Service   │  │   Service    │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │   Health    │  │   Metrics   │  │   Logs      │  │   Config    │  │  │
│  │  │   Check     │  │   Service    │  │   Service   │  │   Service    │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  └─────────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           BANCO DE DADOS                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐  │
│  │                    Google Cloud SQL (MySQL 8.0)                        │  │
│  │  • High Availability                                                  │  │
│  │  • Automated Backups                                                  │  │
│  │  │  • Automated Failover                                               │  │
│  │  │  • Read Replicas                                                    │  │
│  │  │  • Encryption at Rest                                               │  │
│  └─────────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           MONITORAMENTO                                        │
│  ┌─────────────────────────────────────────────────────────────────────────┐  │
│  │                    Google Cloud Monitoring                            │  │
│  │  • Application Metrics                                                │  │
│  │  • Infrastructure Metrics                                             │  │
│  │  • Log Aggregation                                                    │  │
│  │  • Alerting                                                           │  │
│  └─────────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## 🔄 Fluxo de Dados

### **1. Autenticação**
```
Cliente → API Gateway → Cloud Function → Cloud SQL → JWT Token → Cliente
```

### **2. Pedido**
```
Cliente → API Gateway → Aplicação K8s → Cloud SQL → Resposta → Cliente
```

### **3. Monitoramento**
```
Aplicação → Cloud Monitoring → Alertas → Equipe
```

## 🏗️ Componentes Detalhados

### **API Gateway (Nginx Ingress)**
- **Função**: Ponto de entrada único
- **Tecnologia**: Nginx Ingress Controller
- **Localização**: GKE Cluster
- **Recursos**:
  - SSL Termination
  - Rate Limiting
  - Load Balancing
  - Request Routing
  - Health Checks

### **Cloud Function (Autenticação)**
- **Função**: Autenticação por CPF
- **Tecnologia**: Google Cloud Functions
- **Runtime**: Go 1.22
- **Recursos**:
  - Auto-scaling
  - Pay-per-use
  - Serverless
  - JWT Generation

### **Aplicação Principal (Kubernetes)**
- **Função**: API REST para gerenciamento
- **Tecnologia**: Go + GKE
- **Recursos**:
  - Horizontal Pod Autoscaler
  - Rolling Updates
  - Health Checks
  - Resource Limits

### **Banco de Dados (Cloud SQL)**
- **Função**: Armazenamento persistente
- **Tecnologia**: MySQL 8.0
- **Recursos**:
  - High Availability
  - Automated Backups
  - Read Replicas
  - Encryption at Rest

## 🔒 Segurança

### **Autenticação**
- **JWT Tokens**: Stateless authentication
- **CPF Validation**: Brazilian document validation
- **Rate Limiting**: Protection against abuse

### **Autorização**
- **RBAC**: Role-based access control
- **Service Accounts**: Kubernetes service accounts
- **Network Policies**: Network segmentation

### **Criptografia**
- **TLS**: Transport layer security
- **Encryption at Rest**: Database encryption
- **Secrets Management**: GCP Secret Manager

## 📊 Monitoramento

### **Métricas**
- **Application**: Response time, error rate
- **Infrastructure**: CPU, memory, disk
- **Database**: Connection pool, query performance

### **Logs**
- **Application Logs**: Structured logging
- **Access Logs**: Request/response logs
- **Error Logs**: Exception tracking

### **Alertas**
- **High Error Rate**: > 5% errors
- **High Response Time**: > 2 seconds
- **Resource Usage**: > 80% CPU/Memory

## 🚀 Escalabilidade

### **Horizontal Scaling**
- **Kubernetes HPA**: Pod autoscaling
- **Cloud Function**: Auto-scaling
- **Database**: Read replicas

### **Vertical Scaling**
- **Resource Limits**: CPU/Memory limits
- **Database**: Instance size adjustment
- **Monitoring**: Performance optimization

## 🔧 Manutenção

### **Deployments**
- **Blue-Green**: Zero-downtime deployments
- **Rolling Updates**: Gradual pod updates
- **Canary Releases**: Gradual traffic shifting

### **Backups**
- **Database**: Automated daily backups
- **Configuration**: Git-based versioning
- **Secrets**: Encrypted storage

## 📈 Performance

### **Caching**
- **Application**: In-memory caching
- **Database**: Query result caching
- **CDN**: Static content delivery

### **Optimization**
- **Database**: Index optimization
- **Application**: Code optimization
- **Infrastructure**: Resource optimization

## 🔍 Troubleshooting

### **Common Issues**
- **High Latency**: Check database connections
- **Memory Issues**: Check resource limits
- **Authentication**: Check JWT tokens

### **Debugging**
- **Logs**: Application and system logs
- **Metrics**: Performance metrics
- **Traces**: Request tracing
