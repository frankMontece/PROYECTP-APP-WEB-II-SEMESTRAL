# ============================================================================
# Build multi-stage: builder compila el binario, runner solo lo ejecuta.
# Resultado: imagen pequeña sin el toolchain de Go.
# ============================================================================

# ---- Etapa 1: builder ----
FROM golang:1.22-alpine AS builder
WORKDIR /src

# Cachear dependencias primero
COPY go.mod go.sum ./
RUN go mod download

# Compilar
COPY . .
# CGO_ENABLED=0 produce binario estático (los drivers SQLite que usamos
# son Go puro, no necesitan CGO). GOOS=linux para el runner.
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/condominio-api ./cmd/api

# ---- Etapa 2: runner (imagen final mínima) ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
# Usuario no-root por seguridad
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /bin/condominio-api /app/condominio-api
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/condominio-api"]