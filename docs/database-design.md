# 🗄️ Design de Banco de Dados - SOAT FIAP

## 📋 **Visão Geral**

### **Banco Escolhido: MySQL 8.0 (Cloud SQL)**
- **ACID Compliance**: Garantia de consistência transacional
- **Relacional**: Estrutura normalizada para integridade referencial
- **Performance**: Índices otimizados para queries complexas
- **Escalabilidade**: Suporte a read replicas e sharding
- **Ecosystem**: Ferramentas maduras e comunidade ativa
- **Cloud Native**: Cloud SQL oferece gerenciamento automático

### **Modelagem de Dados**
Seguimos os padrões de **Normalização** e **Domain-Driven Design (DDD)** para garantir:
- **Integridade Referencial**: Foreign keys e constraints
- **Consistência**: Transações ACID
- **Performance**: Índices estratégicos
- **Escalabilidade**: Estrutura preparada para crescimento

---

## 🏗️ **Arquitetura do Banco**

### **Diagrama ER**

![Diagrama ER](diagrama%20bd.png)

*Diagrama visual das entidades e relacionamentos do sistema SOAT FIAP*


### **Cenários de Uso do Sistema**

#### **1. Fluxo de Pedido Completo**
```
Cliente → Autenticação (CPF) → Seleção de Produtos → Montagem do Pedido → 
Pagamento (QR Code) → Preparação → Entrega → Finalização
```

#### **2. Gestão de Produtos**
```
Admin → Categorias → Produtos → Preços → Disponibilidade → Menu Atualizado
```

#### **3. Relatórios e Analytics**
```
Pedidos → Itens → Produtos → Categorias → Relatórios de Vendas → 
Análise de Performance → Otimização do Menu
```

### **Relacionamentos do Diagrama**

#### **Relacionamentos Principais**
- **CLIENTES → PEDIDOS**: Um cliente pode fazer vários pedidos (1:N)
- **CATEGORIAS → PRODUTOS**: Uma categoria pode ter vários produtos (1:N)
- **PEDIDOS → PEDIDO_ITENS**: Um pedido pode ter vários itens (1:N)
- **PRODUTOS → PEDIDO_ITENS**: Um produto pode estar em vários itens (1:N)
- **PEDIDOS → PAGAMENTOS**: Um pedido pode ter vários pagamentos (1:N)

#### **Cardinalidades**
- **CLIENTES**: 1 cliente → N pedidos
- **CATEGORIAS**: 1 categoria → N produtos
- **PEDIDOS**: 1 pedido → N itens
- **PRODUTOS**: 1 produto → N itens (em diferentes pedidos)
- **PAGAMENTOS**: 1 pedido → N pagamentos (parcelado)

---

## 📊 **Estrutura das Tabelas**

### **1. CLIENTES**
```sql
CREATE TABLE clientes (
    id VARCHAR(36) PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(11) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    telefone VARCHAR(20),
    endereco TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_clientes_cpf (cpf),
    INDEX idx_clientes_email (email),
    INDEX idx_clientes_created_at (created_at)
);
```
- **UUID**: Identificador único global
- **CPF UNIQUE**: Garantia de unicidade
- **Soft Delete**: Preservação de dados históricos
- **Timestamps**: Auditoria de criação e atualização

### **2. CATEGORIAS**
```sql
CREATE TABLE categorias (
    id VARCHAR(36) PRIMARY KEY,
    nome VARCHAR(100) NOT NULL UNIQUE,
    descricao TEXT,
    ativa BOOLEAN DEFAULT TRUE,
    ordem INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_categorias_nome (nome),
    INDEX idx_categorias_ativa (ativa),
    INDEX idx_categorias_ordem (ordem)
);
```
- **Categorias Fixas**: Lanche, Acompanhamento, Bebida, Sobremesa
- **Ordem**: Controle de exibição no menu
- **Soft Delete**: Preservação de dados históricos

### **3. PRODUTOS**
```sql
CREATE TABLE produtos (
    id VARCHAR(36) PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    categoria_id VARCHAR(36) NOT NULL,
    preco DECIMAL(10,2) NOT NULL,
    descricao TEXT,
    imagem_url VARCHAR(500),
    ativo BOOLEAN DEFAULT TRUE,
    disponivel BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (categoria_id) REFERENCES categorias(id),
    INDEX idx_produtos_categoria (categoria_id),
    INDEX idx_produtos_ativo (ativo),
    INDEX idx_produtos_disponivel (disponivel),
    INDEX idx_produtos_preco (preco)
);
```
- **Foreign Key**: Integridade referencial com categorias
- **DECIMAL**: Precisão monetária
- **Status Duplo**: Ativo (cadastrado) vs Disponível (estoque)
- **Imagem URL**: Referência para assets

### **4. PEDIDOS**
```sql
CREATE TABLE pedidos (
    id VARCHAR(36) PRIMARY KEY,
    cliente_id VARCHAR(36),
    status ENUM('RECEBIDO', 'EM_PREPARACAO', 'PRONTO', 'FINALIZADO', 'CANCELADO') DEFAULT 'RECEBIDO',
    total DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    taxa_servico DECIMAL(10,2) DEFAULT 0.00,
    desconto DECIMAL(10,2) DEFAULT 0.00,
    observacoes TEXT,
    tempo_estimado INT, -- em minutos
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (cliente_id) REFERENCES clientes(id),
    INDEX idx_pedidos_cliente (cliente_id),
    INDEX idx_pedidos_status (status),
    INDEX idx_pedidos_created_at (created_at),
    INDEX idx_pedidos_total (total)
);
```
- **ENUM**: Validação de status no banco
- **Campos Monetários**: Separação de responsabilidades
- **Tempo Estimado**: Para acompanhamento do cliente
- **Cliente Opcional**: Pedidos sem identificação

### **5. PEDIDO_ITENS**
```sql
CREATE TABLE pedido_itens (
    id VARCHAR(36) PRIMARY KEY,
    pedido_id VARCHAR(36) NOT NULL,
    produto_id VARCHAR(36) NOT NULL,
    quantidade INT NOT NULL,
    preco_unitario DECIMAL(10,2) NOT NULL,
    preco_total DECIMAL(10,2) NOT NULL,
    observacoes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id) ON DELETE CASCADE,
    FOREIGN KEY (produto_id) REFERENCES produtos(id),
    INDEX idx_pedido_itens_pedido (pedido_id),
    INDEX idx_pedido_itens_produto (produto_id),
    INDEX idx_pedido_itens_quantidade (quantidade)
);
```
- **CASCADE DELETE**: Limpeza automática de itens
- **Preço Capturado**: Preservação do preço no momento do pedido
- **Observações**: Personalizações do item

### **6. PAGAMENTOS**
```sql
CREATE TABLE pagamentos (
    id VARCHAR(36) PRIMARY KEY,
    pedido_id VARCHAR(36) NOT NULL,
    metodo ENUM('QR_CODE', 'CARTAO', 'DINHEIRO', 'PIX') NOT NULL,
    status ENUM('PENDENTE', 'APROVADO', 'REJEITADO', 'CANCELADO') DEFAULT 'PENDENTE',
    valor DECIMAL(10,2) NOT NULL,
    transacao_id VARCHAR(255),
    qr_code_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id),
    INDEX idx_pagamentos_pedido (pedido_id),
    INDEX idx_pagamentos_status (status),
    INDEX idx_pagamentos_metodo (metodo),
    INDEX idx_pagamentos_transacao (transacao_id)
);
```
- **Separação de Responsabilidades**: Pagamento independente do pedido
- **Múltiplos Métodos**: Flexibilidade de pagamento
- **Transação ID**: Rastreamento de pagamentos externos

---

## 🔍 **Índices Estratégicos**

### **Índices Estratégicos por Tabela**
```sql
-- Índices para CLIENTES
CREATE INDEX idx_clientes_cpf ON clientes(cpf);
CREATE INDEX idx_clientes_email ON clientes(email);
CREATE INDEX idx_clientes_created_at ON clientes(created_at);

-- Índices para CATEGORIAS
CREATE INDEX idx_categorias_nome ON categorias(nome);
CREATE INDEX idx_categorias_ativa ON categorias(ativa);
CREATE INDEX idx_categorias_ordem ON categorias(ordem);

-- Índices para PRODUTOS
CREATE INDEX idx_produtos_categoria ON produtos(categoria_id);
CREATE INDEX idx_produtos_ativo ON produtos(ativo);
CREATE INDEX idx_produtos_disponivel ON produtos(disponivel);
CREATE INDEX idx_produtos_preco ON produtos(preco);

-- Índices para PEDIDOS
CREATE INDEX idx_pedidos_cliente ON pedidos(cliente_id);
CREATE INDEX idx_pedidos_status ON pedidos(status);
CREATE INDEX idx_pedidos_created_at ON pedidos(created_at);
CREATE INDEX idx_pedidos_total ON pedidos(total);

-- Índices para PEDIDO_ITENS
CREATE INDEX idx_pedido_itens_pedido ON pedido_itens(pedido_id);
CREATE INDEX idx_pedido_itens_produto ON pedido_itens(produto_id);
CREATE INDEX idx_pedido_itens_quantidade ON pedido_itens(quantidade);

-- Índices para PAGAMENTOS
CREATE INDEX idx_pagamentos_pedido ON pagamentos(pedido_id);
CREATE INDEX idx_pagamentos_status ON pagamentos(status);
CREATE INDEX idx_pagamentos_metodo ON pagamentos(metodo);
CREATE INDEX idx_pagamentos_transacao ON pagamentos(transacao_id);
```

### **Índices Compostos para Queries Complexas**
```sql
-- Pedidos por cliente e status
CREATE INDEX idx_pedidos_cliente_status ON pedidos(cliente_id, status);

-- Produtos por categoria e status
CREATE INDEX idx_produtos_categoria_ativo ON produtos(categoria_id, ativo);

-- Itens por pedido e produto
CREATE INDEX idx_pedido_itens_pedido_produto ON pedido_itens(pedido_id, produto_id);

-- Pedidos por data e status (para relatórios)
CREATE INDEX idx_pedidos_data_status ON pedidos(created_at, status);
```

---

## 🔒 **Constraints e Validações**

### **Constraints de Integridade**
```sql
-- Validação de CPF (11 dígitos numéricos)
ALTER TABLE clientes ADD CONSTRAINT chk_cpf_format 
CHECK (LENGTH(cpf) = 11 AND cpf REGEXP '^[0-9]+$');

-- Validação de Email
ALTER TABLE clientes ADD CONSTRAINT chk_email_format 
CHECK (email REGEXP '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Validação de Preços Positivos
ALTER TABLE produtos ADD CONSTRAINT chk_preco_positivo 
CHECK (preco > 0);

-- Validação de Quantidades Positivas
ALTER TABLE pedido_itens ADD CONSTRAINT chk_quantidade_positiva 
CHECK (quantidade > 0);

-- Validação de Valores Monetários
ALTER TABLE pedidos ADD CONSTRAINT chk_total_positivo 
CHECK (total >= 0);

ALTER TABLE pagamentos ADD CONSTRAINT chk_valor_positivo 
CHECK (valor > 0);

-- Validação de Tempo Estimado
ALTER TABLE pedidos ADD CONSTRAINT chk_tempo_estimado 
CHECK (tempo_estimado >= 0);
```

### **Triggers para Auditoria**
```sql
-- Trigger para atualizar updated_at
DELIMITER $$
CREATE TRIGGER tr_clientes_updated_at 
BEFORE UPDATE ON clientes 
FOR EACH ROW 
BEGIN 
    SET NEW.updated_at = CURRENT_TIMESTAMP; 
END$$
DELIMITER ;
```

---

## 📈 **Otimizações de Performance**

### **1. Particionamento**
```sql
-- Particionamento por data (futuro)
ALTER TABLE pedidos PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);
```

### **2. Views Materializadas**
```sql
-- View para relatórios de vendas
CREATE VIEW vw_vendas_diarias AS
SELECT 
    DATE(created_at) as data,
    COUNT(*) as total_pedidos,
    SUM(total) as receita_total,
    AVG(total) as ticket_medio
FROM pedidos 
WHERE status = 'FINALIZADO'
GROUP BY DATE(created_at);
```

### **3. Stored Procedures**
```sql
-- Procedure para calcular total do pedido
DELIMITER $$
CREATE PROCEDURE sp_calcular_total_pedido(IN p_pedido_id VARCHAR(36))
BEGIN
    UPDATE pedidos 
    SET total = (
        SELECT SUM(preco_total) 
        FROM pedido_itens 
        WHERE pedido_id = p_pedido_id
    )
    WHERE id = p_pedido_id;
END$$
DELIMITER ;
```

---

## 🔄 **Estratégias de Backup e Recovery**

### **Backup Automático**
- **Frequência**: Diário às 02:00 UTC
- **Retenção**: 7 dias
- **Point-in-time Recovery**: Habilitado
- **Cross-region**: Backup em região diferente

### **Recovery Testing**
```sql
-- Script de teste de recovery
-- 1. Criar backup de teste
-- 2. Simular falha
-- 3. Restaurar backup
-- 4. Validar integridade
```

---

## 📊 **Monitoramento e Métricas**

### **Queries de Monitoramento**
```sql
-- Performance de queries
SELECT 
    query_time,
    lock_time,
    rows_sent,
    rows_examined,
    sql_text
FROM mysql.slow_log 
WHERE start_time > DATE_SUB(NOW(), INTERVAL 1 HOUR)
ORDER BY query_time DESC;

-- Uso de índices
SELECT 
    table_name,
    index_name,
    cardinality,
    sub_part,
    packed
FROM information_schema.statistics 
WHERE table_schema = 'soat'
ORDER BY table_name, seq_in_index;
```

### **Métricas de Negócio**
```sql
-- Relatório de vendas
SELECT 
    DATE(created_at) as data,
    COUNT(*) as pedidos,
    SUM(total) as receita,
    AVG(total) as ticket_medio
FROM pedidos 
WHERE status = 'FINALIZADO'
GROUP BY DATE(created_at)
ORDER BY data DESC;

-- Produtos mais vendidos
SELECT 
    p.nome,
    SUM(pi.quantidade) as total_vendido,
    SUM(pi.preco_total) as receita
FROM produtos p
JOIN pedido_itens pi ON p.id = pi.produto_id
JOIN pedidos pe ON pi.pedido_id = pe.id
WHERE pe.status = 'FINALIZADO'
GROUP BY p.id, p.nome
ORDER BY total_vendido DESC;
```

