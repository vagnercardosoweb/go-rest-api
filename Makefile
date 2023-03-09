IMAGE_URL="vagnercardosoweb/go-rest-api"
IMAGE_VERSION=$$(date +"%Y%m%dT%H%M")

start:
	docker-compose -f docker-compose.yml up --build -d

build:
	docker build --rm --no-cache -f ./Dockerfile.production -t "${IMAGE_URL}:${IMAGE_VERSION}" .

generate-linux-bin:
	rm -rf ./bin && mkdir -p ./bin
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./cmd/api/main.go

sql-generate:
	sqlc -x -f ./sqlc.yaml generate

run:
	go run ./cmd/api/main.go

.PHONY: start