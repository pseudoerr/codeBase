FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth-api ./cmd/auth-service

# prod
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/auth-api /app/
COPY --from=builder /app/.env /app/
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8081

ENTRYPOINT ["/app/auth-api"]