# arguments
ARG PORT=3000
ARG GO_VERSION=1.25.1

# base image
FROM golang:${GO_VERSION}-alpine AS base
RUN apk add --no-cache tzdata ca-certificates

# dev image
FROM base AS dev

ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN apk add --no-cache git curl build-base
RUN go install github.com/air-verse/air@latest

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE ${PORT}

CMD [ "air", "-c", ".air.toml" ]

# builder image
FROM base AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o ./api ./cmd/api/main.go

# prod image
FROM gcr.io/distroless/static-debian12 AS prod

ENV TZ=UTC
ENV APP_ENV=production

WORKDIR /go/src

# Copy certificates and time zone
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo

# Copy migrations and api
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/api ./

# define non-root user
USER nonroot:nonroot

EXPOSE ${PORT}

CMD [ "/go/src/api" ]
