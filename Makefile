IMAGE_URL="vagnercardosoweb/go-rest-api"
POSTGRESQL_LOCAL_URL="postgres://root:root@host.docker.internal:5432/development?sslmode=disable&search_path=public"
IMAGE_VERSION=$$(date +"%Y%m%dT%H%M")

define run_migration_docker
	docker run --rm -v $(shell pwd)/migrations:/migrations migrate/migrate -path /migrations/ -database ${POSTGRESQL_LOCAL_URL} $(1)
endef

start:
	docker-compose -f docker-compose.yml up --build -d

build:
	docker build --rm --no-cache -f ./Dockerfile.production -t "${IMAGE_URL}:${IMAGE_VERSION}" .

migration_up:
	$(call run_migration_docker,up)

migration_down:
	$(call run_migration_docker,down -all)

generate_linux_bin:
	rm -rf ./bin && mkdir -p ./bin
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./cmd/api/main.go

sql_generate:
	~/go/bin/sqlc -x -f ./sqlc.yaml generate

run:
	go run ./cmd/api/main.go

.PHONY: start
