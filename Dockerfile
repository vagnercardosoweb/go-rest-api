# arguments
ARG PORT=3000
ARG GO_VERSION=1.24.3

# base image
FROM golang:${GO_VERSION}-bullseye AS base

RUN apt-get update -y

WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download all && go mod verify
COPY . .

# dev image
FROM base AS dev

RUN go install github.com/air-verse/air@latest
EXPOSE ${PORT}

CMD [ "air", "-c", ".air.toml" ]

# builder image
FROM base AS builder

WORKDIR /go/src
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./api /go/src/cmd/api/main.go

# prod image
FROM scratch AS prod

WORKDIR /go/src
COPY --from=builder /go/src/api ./

EXPOSE ${PORT}

CMD [ "/go/src/api" ]
