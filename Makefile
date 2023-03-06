IMAGE_VERSION=$$(date +"%Y%m%dT%H%M")

start:
	docker-compose -f docker-compose.yml up --build

build-image:
	docker build --rm --no-cache -f ./Dockerfile.production -t "vagnercardosoweb/my-personal-finances:${IMAGE_VERSION}" .

build-bin:
	rm -rf ./bin && mkdir -p ./bin
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./cmd/api/main.go

run-bin:
	go run ./bin/api

run:
	go run ./cmd/api/main.go

.PHONY: start