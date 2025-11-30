# Stage 1: Build the Go application
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o omalib .

EXPOSE 8080

# Stage 2: Create a minimal runtime image
FROM debian:bookworm-slim

WORKDIR /app

# 1. ca-certificates: потрібні для HTTPS (наприклад, для AWS S3)
# 2. libpq5: потрібна для роботи драйвера github.com/lib/pq (PostgreSQL)
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libpq5 \
 && rm -rf /var/lib/apt/lists/*
 
COPY --from=builder /app/omalib .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/configs ./configs

RUN chmod +x /app/omalib 

EXPOSE 8080

CMD [ "./omalib" ]

