# SOAT FIAP - API Lanchonete

API REST para gerenciamento de pedidos de uma lanchonete, desenvolvida com arquitetura hexagonal em Go.

## 🚀 Tecnologias

- Go 1.21
- MySQL 8.0
- Docker

## 📁 Estrutura do Projeto

```
.
├── cmd/                    # Ponto de entrada da aplicação
├── configs/               # Configurações
├── internal/              # Código interno da aplicação
│   ├── adapters/         # Adaptadores (primários e secundários)
│   │   ├── primary/     # Adaptadores primários (handlers HTTP)
│   │   └── secondary/   # Adaptadores secundários (repositórios)
│   ├── core/            # Núcleo da aplicação
│   │   ├── domain/     # Entidades e regras de negócio
│   │   ├── ports/      # Interfaces (portas)
│   │   └── services/   # Serviços de aplicação
│   └── routes/          # Configuração de rotas
└── pkg/                 # Pacotes compartilhados
    └── database/       # Configuração de banco de dados
```

## 🔧 Configuração

1. Clone o repositório
```bash
git clone https://github.com/seu-usuario/soat-fiap.git
cd soat-fiap
```

2. Crie um arquivo `.env` baseado no `.env.example`:
```bash
cp .env.example .env
```

3. Configure as variáveis de ambiente no arquivo `.env`:
```env
# Server
SERVER_PORT=8080

# Database
MYSQL_ROOT_PASSWORD=sua_senha_root
MYSQL_DATABASE=soat_fiap
MYSQL_USER=seu_usuario
MYSQL_PASSWORD=sua_senha

# API Environment
DB_HOST=mysql
DB_PORT=3306
DB_USER=seu_usuario
DB_PASSWORD=sua_senha
DB_NAME=soat_fiap
LOG_LEVEL=info
SWAGGER_ENABLE=true
```

## 🚀 Executando o Projeto

### Com Docker

```bash
# Construir e iniciar os containers
docker-compose up --build

# Executar em background
docker-compose up -d

# Parar os containers
docker-compose down
```

### Localmente

```bash
# Instalar dependências
go mod download

# Executar a aplicação
go run cmd/api/main.go
```

## 📡 Endpoints da API

### Clientes
- `POST /api/v1/clientes` - Criar cliente
- `GET /api/v1/clientes` - Listar clientes
- `GET /api/v1/clientes/cpf/{cpf}` - Buscar cliente por CPF
- `GET /api/v1/clientes/{id}` - Buscar cliente por ID
- `PUT /api/v1/clientes/{id}` - Atualizar cliente
- `DELETE /api/v1/clientes/{id}` - Deletar cliente

### Produtos
- `POST /api/v1/produtos` - Criar produto
- `GET /api/v1/produtos` - Listar produtos
- `GET /api/v1/produtos?categoria=LANCHE` - Listar produtos por categoria
- `GET /api/v1/produtos/{id}` - Buscar produto por ID
- `PUT /api/v1/produtos/{id}` - Atualizar produto
- `DELETE /api/v1/produtos/{id}` - Deletar produto

### Categorias de Produtos
- `LANCHE`
- `ACOMPANHAMENTO`
- `BEBIDA`
- `SOBREMESA`

## 🧪 Testes

```bash
# Executar todos os testes
make test

# Executar testes com coverage
make test-coverage
```

## 📚 Makefile

```bash
make all           # Build com testes
make build         # Build da aplicação
make run           # Executar a aplicação
make watch         # Live reload
make test          # Executar testes
make clean         # Limpar binários
```

## 🔍 Health Check

A API possui um endpoint de health check:
```
GET /api/v1/health
```

## 📖 Swagger

A documentação Swagger está disponível em:
```
http://localhost:8080/swagger/index.html
```