FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache postgresql-client bash

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/config /app/config

COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]