FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata gcc musl-dev

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o soat-fiap ./cmd/api

FROM alpine:3.18

RUN apk update && apk add --no-cache ca-certificates tzdata

RUN adduser -D -g '' appuser

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/soat-fiap /app/soat-fiap

USER appuser

WORKDIR /app

EXPOSE ${SERVER_PORT}

ENV SERVER_PORT=${SERVER_PORT}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV LOG_LEVEL=${LOG_LEVEL}
ENV SWAGGER_ENABLE=${SWAGGER_ENABLE}

ENTRYPOINT ["/app/soat-fiap"]
