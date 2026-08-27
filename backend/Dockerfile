# Stage 1 — Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Download dependencies first (layer cache efficiency)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /backend ./cmd/server

# Stage 2 — Minimal runtime image
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /backend .

EXPOSE 8080

ENTRYPOINT ["/app/backend"]
