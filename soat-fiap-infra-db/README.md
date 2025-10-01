# SOAT FIAP - Infraestrutura de Banco de Dados

## 🎯 Função
Terraform para provisionar e gerenciar a infraestrutura de banco de dados MySQL no Google Cloud Platform.

## 🏗️ Arquitetura

### **Componentes**
- **Cloud SQL**: MySQL 8.0 gerenciado
- **Database**: `soat` com tabelas estruturadas
- **Users**: Usuários com permissões específicas
- **Backups**: Backup automático configurado
- **High Availability**: Configuração para alta disponibilidade

### **Recursos Provisionados**
```
Cloud SQL Instance
├── MySQL 8.0
├── High Availability
├── Automated Backups
├── Database: soat
├── Users: appuser, testuser
└── Network: Authorized networks
```

## 🚀 Deploy

### **Pré-requisitos**
- Google Cloud SDK
- Terraform >= 1.0
- MySQL client (opcional)

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

### **3. Configurar Banco**
```bash
# Executar script de setup
chmod +x scripts/setup-database.sh
./scripts/setup-database.sh
```

## 🔧 Configuração

### **Variáveis do Terraform**
```hcl
# terraform.tfvars
project_id = "sylvan-willow-473218-f8"
region     = "us-central1"
instance_name = "soat-fiap-mysql"
database_version = "MYSQL_8_0"
tier = "db-f1-micro"
```

### **Recursos Criados**
- **Cloud SQL Instance**: `soat-fiap-mysql`
- **Database**: `soat`
- **Users**: `appuser`, `testuser`
- **Network**: Authorized networks configuradas
- **Backups**: Automated daily backups

## 📊 Estrutura do Banco

### **Tabelas**
```sql
-- Clientes
CREATE TABLE clientes (
    id VARCHAR(36) PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(11) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Produtos
CREATE TABLE produtos (
    id VARCHAR(36) PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    categoria VARCHAR(100) NOT NULL,
    preco DECIMAL(10,2) NOT NULL,
    descricao TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pedidos
CREATE TABLE pedidos (
    id VARCHAR(36) PRIMARY KEY,
    cliente_id VARCHAR(36),
    status VARCHAR(50) DEFAULT 'RECEBIDO',
    total DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cliente_id) REFERENCES clientes(id)
);

-- Itens do Pedido
CREATE TABLE pedido_itens (
    id VARCHAR(36) PRIMARY KEY,
    pedido_id VARCHAR(36),
    produto_id VARCHAR(36),
    quantidade INT NOT NULL,
    preco DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id),
    FOREIGN KEY (produto_id) REFERENCES produtos(id)
);
```

### **Índices**
```sql
-- Índices para performance
CREATE INDEX idx_clientes_cpf ON clientes(cpf);
CREATE INDEX idx_produtos_categoria ON produtos(categoria);
CREATE INDEX idx_pedidos_status ON pedidos(status);
CREATE INDEX idx_pedidos_cliente ON pedidos(cliente_id);
```

## 🔒 Segurança

### **Usuários e Permissões**
- **appuser**: Usuário da aplicação (SELECT, INSERT, UPDATE, DELETE)
- **testuser**: Usuário para testes (SELECT, INSERT, UPDATE, DELETE)
- **root**: Administrador (todas as permissões)

### **Network Security**
- **Authorized Networks**: IPs autorizados
- **SSL**: Conexões criptografadas
- **Firewall**: Regras de firewall configuradas

### **Backup e Recovery**
- **Automated Backups**: Diários às 02:00 UTC
- **Retention**: 7 dias
- **Point-in-time Recovery**: Habilitado
- **High Availability**: Configurado

## 📊 Monitoramento

### **Métricas**
- **CPU**: Utilização de CPU
- **Memory**: Utilização de memória
- **Connections**: Conexões ativas
- **Queries**: Performance de queries

### **Logs**
- **Error Logs**: Logs de erro
- **Slow Query Log**: Queries lentas
- **General Log**: Log geral

## 🚀 CI/CD

### **GitHub Actions**
- **Trigger**: Push para main branch
- **Terraform**: Plan e Apply
- **Database**: Setup automático
- **Verification**: Health checks

### **Pipeline**
1. **Checkout**: Código do repositório
2. **Setup**: Terraform e MySQL client
3. **Plan**: Terraform plan
4. **Apply**: Terraform apply
5. **Setup DB**: Executar scripts SQL
6. **Verify**: Health checks

## 🔧 Desenvolvimento

### **Estrutura**
```
soat-fiap-infra-db/
├── terraform/
│   ├── main.tf              # Configuração principal
│   ├── variables.tf         # Variáveis
│   ├── outputs.tf          # Outputs
│   └── terraform.tfvars    # Valores das variáveis
├── sql/
│   ├── 01-create-database.sql    # Criação do banco
│   ├── 02-insert-sample-data.sql # Dados de exemplo
│   └── 03-create-indexes.sql     # Criação de índices
├── scripts/
│   └── setup-database.sh        # Script de setup
└── .github/
    └── workflows/
        └── deploy.yml             # CI/CD pipeline
```

### **Comandos Úteis**
```bash
# Conectar ao banco
gcloud sql connect soat-fiap-mysql --user=appuser

# Backup manual
gcloud sql backups create --instance=soat-fiap-mysql

# Restaurar backup
gcloud sql backups restore BACKUP_ID --instance=soat-fiap-mysql
```

## 📈 Performance

### **Otimizações**
- **Índices**: Criados para queries frequentes
- **Connection Pooling**: Configurado na aplicação
- **Query Optimization**: Queries otimizadas
- **Resource Limits**: CPU e memória configurados

### **Scaling**
- **Vertical**: Aumentar tier da instância
- **Horizontal**: Read replicas (futuro)
- **Connection Limits**: Configurados

## 🔍 Troubleshooting

### **Common Issues**
- **Connection Timeout**: Verificar network
- **Permission Denied**: Verificar usuários
- **Slow Queries**: Verificar índices

### **Debug**
```bash
# Status da instância
gcloud sql instances describe soat-fiap-mysql

# Logs da instância
gcloud sql operations list --instance=soat-fiap-mysql

# Conectar e testar
gcloud sql connect soat-fiap-mysql --user=appuser
```

## 🧹 Limpeza

### **Destruir Infraestrutura**
```bash
# Destruir com Terraform
terraform destroy

# Verificar recursos
gcloud sql instances list
```

## 📚 Documentação

### **Design de Banco de Dados**
- [Modelagem de Dados](docs/database-design.md) - Padrões de normalização e DDD
- [Estrutura das Tabelas](docs/database-design.md#estrutura-das-tabelas) - Schema completo
- [Índices e Performance](docs/database-design.md#índices-estratégicos) - Otimizações
- [Segurança e Compliance](docs/database-design.md#segurança-e-compliance) - LGPD e auditoria


