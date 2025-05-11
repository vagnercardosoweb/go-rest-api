# Go REST API

Uma API REST robusta e escalável desenvolvida em Golang, seguindo boas práticas de arquitetura e design de software.

## Tecnologias Utilizadas

- **Go 1.24+**: Linguagem de programação principal
- **Gin**: Framework web para construção de APIs REST
- **PostgreSQL**: Banco de dados relacional
- **Redis**: Armazenamento em cache e gerenciamento de filas
- **JWT**: Autenticação baseada em tokens
- **Docker & Docker Compose**: Containerização e orquestração
- **Golang Migrate**: Gerenciamento de migrações de banco de dados
- **AWS SDK**: Integração com serviços da AWS (SES, etc.)
- **Bcrypt**: Hash de senhas seguro
- **Air**: Hot-reload para desenvolvimento
- **Testcontainers**: Testes de integração com containers

## Pré-requisitos

- Go 1.24 ou superior
- Docker e Docker Compose
- Make (opcional, para facilitar o uso dos comandos)

## Configuração do Ambiente

1. Clone o repositório:

   ```bash
   git clone https://github.com/vagnercardosoweb/go-rest-api.git
   cd go-rest-api
   ```

2. Copie o arquivo de exemplo de variáveis de ambiente:

   ```bash
   cp .env.example .env.development
   ```

3. Ajuste as variáveis de ambiente conforme necessário no arquivo `.env.development`

## Executando o Projeto

### Usando Docker (recomendado)

```bash
make start_docker
```

Este comando irá iniciar todos os serviços necessários (API, PostgreSQL e Redis) em containers Docker.

### Localmente com Hot-Reload

```bash
make start_development
```

Este comando utiliza o Air para fornecer hot-reload durante o desenvolvimento.

## Banco de Dados

### Migrações

Criar uma nova migração:

```bash
make create_migration name="nome_da_migracao"
```

Executar migrações:

```bash
make migration_up
```

Reverter todas as migrações:

```bash
make migration_down
```

## Testes

Executar todos os testes:

```bash
make test
```

Executar testes com detecção de race conditions:

```bash
make test_race
```

## Build

### Build para ambiente local

```bash
make generate_local_bin
```

### Build para Linux (produção)

```bash
make generate_linux_bin
```

### Build para Docker

```bash
make docker_build_local
```

### Build e publicação para AWS ECR

```bash
make docker_build_aws
```

## Estrutura do Projeto

- **cmd/api**: Ponto de entrada da aplicação
- **internal**: Código específico da aplicação
  - **events**: Sistema de eventos da aplicação
  - **handlers**: Manipuladores HTTP
  - **repositories**: Camada de acesso a dados
  - **schedules**: Tarefas agendadas
  - **services**: Lógica de negócio
  - **types**: Definições de tipos e estruturas
- **migrations**: Arquivos de migração do banco de dados
- **pkg**: Pacotes reutilizáveis
  - **api**: Configuração e utilitários da API REST
  - **aws**: Integrações com serviços AWS
  - **env**: Gerenciamento de variáveis de ambiente
  - **errors**: Tratamento padronizado de erros
  - **logger**: Sistema de logging
  - **password**: Utilitários para hash de senhas
  - **postgres**: Cliente e utilitários para PostgreSQL
  - **redis**: Cliente e utilitários para Redis
  - **token**: Implementação de JWT
  - **utils**: Funções utilitárias diversas
- **resources**: Recursos estáticos
- **tests**: Testes de integração e utilitários para testes

## Funcionalidades Principais

- Autenticação via JWT
- Sistema de eventos para desacoplamento de operações
- Tarefas agendadas
- Integração com AWS (SES, etc.)
- Logging estruturado
- Tratamento padronizado de erros
- Monitoramento e profiling

## Variáveis de Ambiente

O arquivo `.env.example` contém todas as variáveis de ambiente necessárias para executar a aplicação. Algumas das principais são:

- `PORT`: Porta em que a API será executada
- `APP_ENV`: Ambiente da aplicação (local, production, staging, test)
- `JWT_SECRET_KEY`: Chave secreta para assinatura de tokens JWT
- `DB_*`: Configurações do PostgreSQL
- `REDIS_*`: Configurações do Redis
- `AWS_*`: Configurações da AWS

## Licença

Este projeto está sob a licença MIT.

## Autor

Desenvolvido por [Vagner Cardoso](https://github.com/vagnercardosoweb).
