ARGS := $(filter-out $@,$(MAKECMDGOALS))

ifneq ($(wildcard .env.development),)
	include .env.development
endif

AWS_REGION?=us-east-1
AWS_ACCOUNT_ID?=000000000000
AWS_REGISTRY_URL?=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
AWS_PROFILE?=default

IMAGE_VERSION=$(shell date +"%Y%m%dT%H%M%S")
IMAGE_URL?=${AWS_REGISTRY_URL}/go-rest-api

DB_SCHEMA?=public
DB_MIGRATION_FOLDER?=migrations

ifeq ($(DB_ENABLED_SSL),"false")
	DB_SSL=disable
else
	DB_SSL=require
endif

ifdef DB_HOST
	DB_URL=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}&search_path=${DB_SCHEMA}
endif

define run_migration_docker
  @if [ -z "${DB_HOST}" ]; then \
    echo "DB_URL is not set"; \
    exit 1; \
  fi

	docker run --rm -v $(shell pwd)/${DB_MIGRATION_FOLDER}:/migrations migrate/migrate -path /migrations/ -database "${DB_URL}" $(1)
endef

run:
	go run ./cmd/api/main.go

run_race:
	go run -race ./cmd/api/main.go

start_docker:
	docker compose -f docker-compose.yml down --remove-orphans
	docker compose -f docker-compose.yml up --build -d
	docker logs go-rest-api-api -f

start_development: check_build
	docker compose -f docker-compose.yml up redis postgres -d
	APP_ENV=development air -c .air.toml

start_production: check_build
	APP_ENV=production air -c .air.toml

docker_build:
	@if [ "$(word 2,$(ARGS))" != "dev" ] && [ "$(word 2,$(ARGS))" != "prod" ]; then \
		echo "invalid argument. Use: make docker_build <dev|prod> <local|aws>"; \
		exit 1; \
	fi

	@if [ "$(word 3,$(ARGS))" != "local" ] && [ "$(word 3,$(ARGS))" != "aws" ]; then \
		echo "invalid argument. Use: make docker_build <dev|prod> <local|aws>"; \
		exit 1; \
	fi

	@if [ "$(word 3,$(ARGS))" = "aws" ]; then \
		aws --profile ${AWS_PROFILE} ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_REGISTRY_URL}; \
		docker build --rm --no-cache --push --platform linux/amd64 --target prod -f ./Dockerfile -t ${IMAGE_URL}:$(word 2,$(ARGS))-${IMAGE_VERSION} .; \
	fi

	@if [ "$(word 3,$(ARGS))" = "local" ]; then \
		docker build --rm --no-cache --target prod -f ./Dockerfile -t ${IMAGE_URL}:$(word 2,$(ARGS))-${IMAGE_VERSION} .; \
	fi

check_build:
	go mod download
	go build -v ./...

create_migration:
	./create-migration-file.sh "$(name)"

migration_up:
	$(call run_migration_docker,up)

migration_down:
	$(call run_migration_docker,down 1)

migration_clean:
	$(call run_migration_docker,down -all)

generate_bin:
	rm -rf ./bin && mkdir -p ./bin

	@if [ "$(word 2,$(ARGS))" != "linux" ] && [ "$(word 2,$(ARGS))" != "local" ]; then \
		echo "invalid argument. Use: make generate_bin <linux|local>"; \
		exit 1; \
	fi

	@if [ "$(word 2,$(ARGS))" = "linux" ]; then \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./cmd/api/main.go; \
	fi

	@if [ "$(word 2,$(ARGS))" = "local" ]; then \
		CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/api ./cmd/api/main.go; \
	fi

update_modules:
	go get -u ./...
	go mod tidy
	make check_build

lint:
	@echo "🔍 Running golangci-lint..."
	golangci-lint run ./...

lint_install:
	@echo "📦 Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

security:
	@echo "🔒 Running security checks..."
	@echo "Running gosec..."
	gosec ./...
	@echo "Running govulncheck..."
	govulncheck -show verbose ./...

security_install:
	@echo "📦 Installing security tools..."
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

staticcheck:
	@echo "🔍 Running staticcheck..."
	staticcheck ./...

staticcheck_install:
	@echo "📦 Installing staticcheck..."
	go install honnef.co/go/tools/cmd/staticcheck@latest

format:
	@echo "🎨 Formatting code..."
	gofmt -s -w .
	goimports -w .

format_install:
	@echo "📦 Installing formatting tools..."
	go install golang.org/x/tools/cmd/goimports@latest

test:
	@echo "🧪 Running test without coverage..."
	APP_ENV=test go test -v ./...

test_race:
	@echo "🧪 Running tests with race detection..."
	APP_ENV=test go test -v --race ./...

test_coverage:
	@echo "📊 Running tests with coverage..."
	APP_ENV=test go test -v -coverprofile=coverage.out ./... -coverpkg=./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

install_tools: lint_install security_install staticcheck_install format_install
	@echo "✅ All development tools installed!"

quality: format lint staticcheck security
	@echo "✅ All quality checks completed!"

ci: check_build quality test_coverage
	@echo "🚀 CI pipeline completed successfully!"

help:
	@echo "📚 Available commands:"
	@echo ""
	@echo "🏗️  Build & Run:"
	@echo "  run                 - Run the application in development mode"
	@echo "  run_race            - Run the application in development mode with race detection"
	@echo "  check_build         - Verify build and dependencies"
	@echo "  generate_bin        - Generate binary <linux|local>"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  start_docker        - Start with Docker Compose"
	@echo "  docker_build        - Build Docker image <dev|prod> <local|aws>"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test               - Run all tests"
	@echo "  test_race          - Run tests with race detection"
	@echo "  test_coverage      - Run tests with coverage report"
	@echo ""
	@echo "🔍 Quality & Security:"
	@echo "  lint               - Run golangci-lint"
	@echo "  security           - Run security checks (gosec + govulncheck)"
	@echo "  staticcheck        - Run staticcheck analysis"
	@echo "  format             - Format code with gofmt and goimports"
	@echo "  quality            - Run all quality checks"
	@echo ""
	@echo "📦 Installation:"
	@echo "  install_tools      - Install all development tools"
	@echo "  lint_install       - Install golangci-lint"
	@echo "  security_install   - Install security tools"
	@echo "  staticcheck_install - Install staticcheck"
	@echo "  format_install     - Install formatting tools"
	@echo ""
	@echo "🚀 CI/CD:"
	@echo "  ci                 - Run complete CI pipeline"
	@echo ""
	@echo "🗄️  Database:"
	@echo "  create_migration   - Create new migration file"
	@echo "  migration_up       - Run database migrations"
	@echo "  migration_down     - Rollback database migrations"
	@echo "  migration_clean    - Rollback all database migrations"

.PHONY: run run_race start_docker start_development start_production docker_build check_build create_migration migration_up migration_down migration_clean generate_bin update_modules test test_race lint lint_install security security_install staticcheck staticcheck_install format format_install test_coverage install_tools quality ci help