FROM golang:1.26.2-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/condominio-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /bin/condominio-api /app/condominio-api
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/condominio-api"]