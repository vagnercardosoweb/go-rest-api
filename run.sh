#!/usr/bin/env sh

set -e

docker stop $(docker ps -a -q)
docker-compose -f docker-compose.yml up -d

PORT=8081 npx nodemon --exec go run main.go --signal SIGTERM
