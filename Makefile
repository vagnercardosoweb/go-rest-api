AWS_REGION?=us-east-1
AWS_ACCOUNT_ID?=000000000000
AWS_REGISTRY_URL?=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
AWS_PROFILE?=golang

IMAGE_VERSION=$(shell date +"%Y%m%dT%H%M%S")
IMAGE_URL?=${AWS_REGISTRY_URL}/golang-api

DB_PORT?=5432
DB_USERNAME?=root
DB_PASSWORD?=root
DB_HOST?=host.docker.internal
DB_NAME?=development
DB_ENABLED_SSL?=false
DB_SCHEMA?=public

ifeq ($(DB_ENABLED_SSL),false)
	DB_SSL?=disable
else
	DB_SSL?=require
endif

DB_URL?=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}&search_path=${DB_SCHEMA}

define run_migration_docker
	docker run --rm -v $(shell pwd)/migrations:/migrations migrate/migrate -path /migrations/ -database "${DB_URL}" $(1)
endef

run:
	go run -race ./cmd/api/main.go

start_docker:
	docker-compose -f docker-compose.yml down --remove-orphans
	docker-compose -f docker-compose.yml up --build -d
	docker logs go-rest-api.api -f

start_local: check_build
	APP_ENV=local ~/go/bin/air -c .air.toml

start_production: check_build
	APP_ENV=production ~/go/bin/air -c .air.toml

start_staging: check_build
	APP_ENV=staging ~/go/bin/air -c .air.toml

docker_build_local:
	docker build --rm --no-cache -f ./Dockerfile.production -t ${IMAGE_URL}.${IMAGE_VERSION} .

docker_build_aws:
	docker build --rm --no-cache --platform linux/amd64 -f ./Dockerfile.production -t ${IMAGE_URL}.${IMAGE_VERSION} .
	aws --profile ${AWS_PROFILE} ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_REGISTRY_URL}
	docker push ${IMAGE_URL}.${IMAGE_VERSION}

check_build:
	go mod tidy
	go build -v ./...

migration_up:
	$(call run_migration_docker,up)

migration_down:
	$(call run_migration_docker,down -all)

generate_linux_bin:
	rm -rf ./bin && mkdir -p ./bin
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./cmd/api/main.go

generate_local_bin:
	rm -rf ./bin && mkdir -p ./bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/api ./cmd/api/main.go

update_modules:
	go get -u ./...
	go mod tidy
	make check_build

test: check_build
	go test -v ./...

.PHONY: run start_docker start_local start_production start_staging docker_build_local docker_build_aws check_build migration_up migration_down generate_linux_bin generate_local_bin update_modules test
