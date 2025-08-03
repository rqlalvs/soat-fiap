# SOAT FIAP - API Lanchonete 

API REST completa para gerenciamento de pedidos de uma lanchonete com sistema de pagamentos integrado, desenvolvida com arquitetura hexagonal em Go.

## 📋 Requisitos

- Git Bash
- Docker
- Go


## 🚀 Tecnologias

- Go 1.23
- MySQL 8.0
- Docker
- Gorilla Mux (Routing)
- Swagger/OpenAPI (Documentação)
- UUID (Identificadores únicos)
- Testify (Testes unitários)

## 📁 Estrutura do Projeto

```
.
├── cmd/                    # Ponto de entrada da aplicação
│   └── api/               # API principal
├── configs/               # Configurações
├── docs/                  # Documentação
├── internal/              # Código interno da aplicação
│   ├── adapters/         # Adaptadores (primários e secundários)
│   │   ├── primary/     # Adaptadores primários (handlers HTTP)
│   │   │   └── handlers/ # ClienteHandler, ProdutoHandler, PedidoHandler, PagamentoHandler
│   │   └── secondary/   # Adaptadores secundários (repositórios)
│   │       └── repositories/ # ClienteRepository, ProdutoRepository, PedidoRepository, PagamentoRepository
│   ├── core/            # Núcleo da aplicação
│   │   ├── domain/     # Entidades: Cliente, Produto, Pedido, Pagamento
│   │   ├── ports/      # Interfaces (portas)
│   │   └── services/   # Serviços de aplicação + testes
│   └── routes/          # Configuração de rotas
└── pkg/                 # Pacotes compartilhados
    └── database/       # Configuração de banco de dados
```

## 🔧 Configuração e Execução

1. Clone o repositório utilizando Git Bash
```bash
git clone https://github.com/rqlalvs/soat-fiap
cd soat-fiap
```

2. Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:
```env
PORT=8080
APP_ENV=local
BLUEPRINT_DB_URL=./test.db
# Server
SERVER_PORT=8080

# Database
MYSQL_ROOT_PASSWORD=root
MYSQL_DATABASE=soat_fiap
MYSQL_USER=soat
MYSQL_PASSWORD=fiap123

# API Environment
DB_HOST=mysql
DB_PORT=3306
DB_USER=soat
DB_PASSWORD=fiap123
DB_NAME=soat_fiap
LOG_LEVEL=info
SWAGGER_ENABLE=true
```

3. Execute o Docker Desktop, e execute projeto usando Docker Compose:
```bash
# Construir e iniciar os containers
docker compose up --build

# OU para executar em background
docker compose up --build -d
```

4. Verifique se os containers estão rodando:
```bash
docker ps
```

5. A API estará disponível em:
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/api/v1/health

6. Para parar os containers:
```bash
docker compose down
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


## 📡 Endpoints da API

### **Checkout e Pedidos**

#### Checkout Integrado com Pagamento
```http
POST /api/v1/checkout
Content-Type: application/json

{
  "cliente_id": "cliente-123",
  "itens": [
    {
      "produto_id": "produto-456",
      "quantidade": 2,
      "observacao": "Sem cebola"
    }
  ]
}
```

**Resposta:**
```json
{
  "pedido_id": "pedido-789",
  "status": "RECEBIDO",
  "valor": 25.50,
  "pagamento_id": "pagamento-abc",
  "qr_code": "00020126580014br.gov.bcb.pix...",
  "external_id": "MP_12345678",
  "expires_at": "2025-08-03T01:00:00Z",
  "message": "Pedido criado com sucesso. Realize o pagamento via PIX para prosseguir."
}
```

#### Pedidos
- `GET /api/v1/pedidos` - **Listar pedidos com ordenação especial** (Pronto > Em Preparação > Recebido, mais antigos primeiro, sem finalizados)
- `GET /api/v1/pedidos?status=EM_PREPARACAO` - Filtrar por status
- `GET /api/v1/pedidos?cliente_id=123` - Filtrar por cliente
- `GET /api/v1/pedidos/{id}` - Buscar pedido por ID
- `PATCH /api/v1/pedidos/{id}/status` - **Atualizar status com validações** (cozinha só inicia após pagamento aprovado)

### **Sistema de Pagamentos (Novo)**

#### Consultar Status de Pagamento
```http
GET /api/v1/pedidos/{pedido_id}/pagamento
```

#### Webhook de Pagamento (Mercado Pago)
```http
POST /api/v1/webhook/pagamento
Content-Type: application/json

{
  "external_id": "MP_12345678",
  "status": "APROVADO"
}
```

#### Simulação de Pagamento (Para Testes)
```http
POST /api/v1/mock/pagamento/{external_id}/simular?status=APROVADO
```

#### Outros Endpoints de Pagamento
- `POST /api/v1/pagamentos` - Criar pagamento manual
- `GET /api/v1/pagamentos/{id}` - Buscar pagamento por ID
- `PATCH /api/v1/pagamentos/{id}/cancelar` - Cancelar pagamento

### **Clientes**
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

### Pedidos
- `POST /api/v1/checkout` - Realizar checkout do pedido
- `POST /api/v1/pedidos` - Criar pedido
- `GET /api/v1/pedidos` - Listar pedidos
- `GET /api/v1/pedidos/{id}` - Buscar pedido por ID
- `PATCH /api/v1/pedidos/{id}/status` - Atualizar status do pedido

### Health Check
- `GET /api/v1/health` - Verificar status da API

A documentação Swagger está disponível em:
```
http://localhost:8080/swagger/index.html
```

Feito por Raquel Alves e Alexandre Pereira :D
