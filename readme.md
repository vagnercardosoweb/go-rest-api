# Go REST API

Uma API REST robusta e escalável desenvolvida em Golang, seguindo boas práticas de arquitetura limpa e design patterns modernos.

## 🚀 Tecnologias Utilizadas

- **Go 1.26**: Linguagem de programação principal
- **Gin**: Framework web para construção de APIs REST
- **PostgreSQL**: Banco de dados relacional
- **Redis**: Cache e gerenciamento de filas
- **JWT**: Autenticação baseada em tokens
- **Docker & Docker Compose**: Containerização e orquestração
- **Golang Migrate**: Gerenciamento de migrações de banco de dados
- **AWS SDK v2**: Integração com serviços AWS (SES, S3, SNS, SQS)
- **Bcrypt**: Hash de senhas seguro
- **Air**: Hot-reload para desenvolvimento
- **Testcontainers**: Testes de integração com containers
- **Slack Integration**: Sistema de alertas e notificações

## 📋 Pré-requisitos

- Go 1.26 ou superior
- Docker e Docker Compose
- Make (opcional, para facilitar o uso dos comandos)

## ⚙️ Configuração do Ambiente

1. **Clone o repositório:**

   ```bash
   git clone https://github.com/vagnercardosoweb/go-rest-api.git
   cd go-rest-api
   ```

2. **Configure as variáveis de ambiente:**

   ```bash
   # Crie o arquivo de ambiente baseado no de exemplo
   cp .env.example .env.development
   # Edite o arquivo conforme necessário
   ```

3. **Instale as ferramentas de desenvolvimento (opcional):**
   ```bash
   make install_tools
   ```

## 🏃‍♂️ Executando o Projeto

### Usando Docker (Recomendado) com Hot-Reload

```bash
make start_docker
```

### Desenvolvimento Local com Hot-Reload

```bash
make start_development
```

### Execução Direta

```bash
make run
```

## 🗄️ Gerenciamento do Banco de Dados

### Migrações

```bash
# Criar nova migração
make create_migration name="nome_da_migracao"

# Executar todas as migrações
make migration_up

# Reverter a última migração
make migration_down

# Reverter todas as migrações
make migration_clean
```

## 🧪 Testes

```bash
# Executar todos os testes
make test

# Testes com detecção de race conditions
make test_race

# Testes com relatório de cobertura
make test_coverage
```

## 🔍 Qualidade de Código

```bash
# Executar todas as verificações de qualidade
make quality

# Verificações individuais
make lint          # Linting com golangci-lint
make security      # Verificações de segurança (gosec + govulncheck)
make staticcheck   # Análise estática
make format        # Formatação de código
```

## 🏗️ Build e Deploy

### Build Local

```bash
make generate_bin local
```

### Build para Linux (Produção)

```bash
make generate_bin linux
```

### Build Docker

```bash
# Build local
make docker_build dev local

# Build e push para AWS ECR
make docker_build prod aws
```

## 📁 Arquitetura do Projeto

```
go-rest-api/
├── cmd/api/                    # Ponto de entrada da aplicação
├── internal/                   # Código específico da aplicação
│   ├── events/                 # Sistema de eventos
│   ├── handlers/               # Handlers HTTP por domínio
│   │   └── user/
│   ├── repositories/           # Camada de acesso a dados
│   │   └── user/
│   ├── schedules/              # Tarefas agendadas
│   ├── services/               # Lógica de negócio
│   │   └── user/
│   └── types/                  # Tipos e estruturas específicas
├── pkg/                        # Pacotes reutilizáveis
│   ├── api/                    # Framework REST customizado
│   │   ├── context/            # Contexto da API
│   │   ├── handlers/           # Handlers genéricos
│   │   ├── middlewares/        # Middlewares
│   │   ├── request/            # Utilitários de request
│   │   └── response/           # Utilitários de response
│   ├── aws/                    # Clientes AWS (SES, S3, SNS, SQS)
│   ├── env/                    # Gerenciamento de variáveis de ambiente
│   ├── errors/                 # Sistema de tratamento de erros
│   ├── events/                 # Sistema de eventos
│   ├── logger/                 # Sistema de logging estruturado
│   ├── mailer/                 # Sistema de envio de emails
│   ├── monitoring/             # Profiling e monitoramento
│   ├── password/               # Utilitários para hash de senhas
│   ├── postgres/               # Cliente PostgreSQL
│   ├── redis/                  # Cliente Redis
│   ├── slack/                  # Integração com Slack
│   ├── token/                  # Implementação JWT
│   └── utils/                  # Funções utilitárias
├── migrations/                 # Migrações do banco de dados
├── resources/                  # Recursos estáticos
│   ├── aws_ses_templates/      # Templates de email
│   └── kubernetes/             # Manifests Kubernetes
└── tests/                      # Utilitários para testes
```

## ✨ Funcionalidades Principais

- **🔐 Autenticação JWT**: Sistema completo de autenticação baseado em tokens
- **📧 Sistema de Email**: Integração com AWS SES e templates
- **📊 Sistema de Eventos**: Arquitetura orientada a eventos para desacoplamento
- **⏰ Tarefas Agendadas**: Scheduler para execução de jobs em background
- **🔔 Alertas Slack**: Notificações automáticas de eventos importantes
- **📝 Logging Estruturado**: Sistema de logs com metadados e redação de dados sensíveis
- **🛡️ Tratamento de Erros**: Sistema padronizado de tratamento e propagação de erros
- **📈 Monitoramento**: Profiling e métricas de performance
- **🧪 Testes Integrados**: Suporte completo a testes com containers
- **🐳 Containerização**: Suporte completo ao Docker e Kubernetes

## 🌍 Variáveis de Ambiente

### Aplicação

- `APP_ENV`: Ambiente da aplicação (`development`, `production`, `test`)
- `PORT`: Porta da aplicação (padrão: `3000`)
- `IS_LOCAL`: Execução local (padrão: `false`)

### Autenticação

- `JWT_SECRET_KEY`: Chave secreta para assinatura JWT
- `JWT_EXPIRES_IN_SECONDS`: Tempo de expiração do token (padrão: `86400`)

### Banco de Dados

- `DB_HOST`: Host do PostgreSQL
- `DB_PORT`: Porta do PostgreSQL (padrão: `5432`)
- `DB_NAME`: Nome do banco de dados
- `DB_USERNAME`: Usuário do banco
- `DB_PASSWORD`: Senha do banco
- `DB_ENABLED_SSL`: Habilitar SSL (padrão: `false`)
- `DB_AUTO_MIGRATE`: Executar migrações automaticamente (padrão: `false`)

### Redis

- `REDIS_HOST`: Host do Redis
- `REDIS_PORT`: Porta do Redis (padrão: `6379`)
- `REDIS_PASSWORD`: Senha do Redis
- `REDIS_DATABASE`: Database do Redis (padrão: `0`)

### AWS

- `AWS_REGION`: Região AWS (padrão: `us-east-1`)
- `AWS_ACCESS_KEY_ID`: Chave de acesso AWS
- `AWS_SECRET_ACCESS_KEY`: Chave secreta AWS
- `AWS_SES_SOURCE`: Email remetente para SES

### Slack

- `SLACK_ENABLED`: Habilitar integração Slack (padrão: `true`)
- `SLACK_TOKEN`: Token do bot Slack
- `SLACK_CHANNEL`: Canal para alertas (padrão: `alerts`)
- `SLACK_USERNAME`: Nome do usuário bot (padrão: `go-rest-api`)

### Sistema

- `SCHEDULER_ENABLED`: Habilitar scheduler (padrão: `false`)
- `SCHEDULER_SLEEP`: Intervalo do scheduler em segundos (padrão: `60`)
- `LOGGER_ENABLED`: Habilitar logging (padrão: `true`)
- `PROFILER_ENABLED`: Habilitar profiler (padrão: `false`)

## 🚀 Pipeline CI/CD

```bash
# Executar pipeline completo
make ci
```

O pipeline inclui:

1. Verificação de build e dependências
2. Formatação de código
3. Linting
4. Análise estática
5. Verificações de segurança
6. Testes com cobertura

## 📚 Comandos Disponíveis

Execute `make help` para ver todos os comandos disponíveis organizados por categoria:

- **🏗️ Build & Run**: `run`, `check_build`, `generate_bin`
- **🐳 Docker**: `start_docker`, `docker_build`
- **🧪 Testing**: `test`, `test_race`, `test_coverage`
- **🔍 Quality & Security**: `lint`, `security`, `staticcheck`, `format`, `quality`
- **📦 Installation**: `install_tools`, `lint_install`, `security_install`
- **🚀 CI/CD**: `ci`
- **🗄️ Database**: `create_migration`, `migration_up`, `migration_down`

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feat/AmazingFeature`)
3. Execute os testes e verificações de qualidade (`make ci`)
4. Commit suas mudanças (`git commit -m 'feat: add some AmazingFeature'`)
5. Push para a branch (`git push origin feat/AmazingFeature`)
6. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👨‍💻 Autor

Desenvolvido por [Vagner Cardoso](https://github.com/vagnercardosoweb)

---

⭐ Se este projeto foi útil para você, considere dar uma estrela no repositório!
