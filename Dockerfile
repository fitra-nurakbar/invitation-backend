# ─── Stage 1: Builder ────────────────────────────────────────
FROM golang:1.25-alpine AS builder

# Install dependencies system
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod dulu — supaya layer cache efisien
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o invitation-app .

# ─── Stage 2: Runner ─────────────────────────────────────────
FROM alpine:latest

# Install ca-certificates untuk HTTPS request (Xendit, Google OAuth)
RUN apk add --no-cache ca-certificates tzdata

# Set timezone Indonesia
ENV TZ=Asia/Jakarta

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/invitation-app .

# Copy folder migrations
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Jalankan app
CMD ["./invitation-app"]