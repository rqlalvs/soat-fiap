# SOAT FIAP - Cloud Function de Autenticação

## 🎯 Função
Cloud Function para autenticação de clientes por CPF no sistema de lanchonete.

## 🏗️ Arquitetura

### **Tecnologia**
- **Runtime**: Go 1.22
- **Platform**: Google Cloud Functions (Gen2)
- **Trigger**: HTTP
- **Authentication**: JWT

### **Fluxo de Autenticação**
```
Cliente → Cloud Function → Cloud SQL → JWT Token → Cliente
```

## 🚀 Deploy

### **Automático (CI/CD)**
```bash
git push origin main  # Trigger automático via GitHub Actions
```

### **Manual**
```bash
gcloud functions deploy authByCPF \
  --gen2 \
  --runtime=go122 \
  --region=us-central1 \
  --source=. \
  --entry-point=AuthByCPF \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars="DB_HOST=,DB_PORT=,DB_USER=,DB_PASSWORD=,DB_NAME=,JWT_SECRET="
```

## 🔧 Configuração

### **Variáveis de Ambiente**
- `DB_HOST`: Host do Cloud SQL
- `DB_PORT`: Porta do banco
- `DB_USER`: Usuário do banco
- `DB_PASSWORD`: Senha do banco
- `DB_NAME`: Nome do banco
- `JWT_SECRET`: Chave secreta para JWT

### **Dependências**
```go
github.com/go-sql-driver/mysql
github.com/golang-jwt/jwt/v5
```

## 📡 API

### **Endpoint**
```
POST https://us-central1-sylvan-willow-473218-f8.cloudfunctions.net/authByCPF
```

### **Request**
```json
{
  "cpf": "12345678901"
}
```

### **Response (Sucesso)**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "cliente": {
    "id": "1",
    "nome": "João Silva",
    "cpf": "12345678901",
    "email": "joao@email.com"
  }
}
```

### **Response (Erro)**
```json
{
  "error": "Cliente não encontrado"
}
```

## 📊 Monitoramento

### **Logs**
- **Cloud Logging**: Logs estruturados
- **Error Tracking**: Exceções capturadas
- **Performance**: Tempo de resposta

### **Métricas**
- **Invocations**: Número de chamadas
- **Errors**: Taxa de erro
- **Duration**: Tempo de execução

## 🔧 Desenvolvimento

### **Estrutura**
```
soat-fiap-auth-function/
├── main.go              # Function entry point
├── go.mod               # Dependencies
├── go.sum               # Dependency checksums
└── .github/
    └── workflows/
        └── deploy.yml   # CI/CD pipeline
```

### **Local Testing**
```bash
# Instalar dependências
go mod tidy

# Executar localmente
go run main.go
```

## 🚀 CI/CD

### **GitHub Actions**
- **Trigger**: Push para main branch
- **Build**: Go build
- **Deploy**: gcloud functions deploy
- **Test**: Automated testing

### **Pipeline**
1. **Checkout**: Código do repositório
2. **Setup Go**: Configuração do ambiente
3. **Build**: Compilação do código
4. **Deploy**: Deploy para Cloud Functions
5. **Test**: Teste da função

## 📚 Documentação

### **Swagger/OpenAPI**
- **Endpoint**: `/docs`
- **Format**: OpenAPI 3.0
- **Authentication**: JWT Bearer

### **Exemplos**
- **cURL**: Exemplos de uso
- **Postman**: Collection disponível
- **SDK**: Go client library

