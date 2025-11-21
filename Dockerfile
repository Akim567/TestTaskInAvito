# ---------- Stage 1: build ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/app

# ---------- Stage 2: runtime ----------
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /app/server

# Порт из конфига по умолчанию 8080
EXPOSE 8080

ENV HTTP_PORT=8080

CMD ["/app/server"]
