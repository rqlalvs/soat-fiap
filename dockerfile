FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata gcc musl-dev

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o soat-fiap ./cmd/api

FROM alpine:3.18

RUN apk update && apk add --no-cache ca-certificates tzdata sqlite

RUN adduser -D -g '' appuser

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/soat-fiap /app/soat-fiap

RUN mkdir -p /app/data && \
    chown -R appuser:appuser /app

USER appuser

WORKDIR /app

EXPOSE 8080

ENV SERVER_PORT=8080
ENV DATABASE_PATH=/app/data/db.sqlite
ENV LOG_LEVEL=info
ENV SWAGGER_ENABLE=true

ENTRYPOINT ["/app/soat-fiap"]