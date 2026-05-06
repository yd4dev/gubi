# Stage 1: Build the binary
FROM golang:1.26-bookworm AS builder

# Set the working directory
WORKDIR /app

# Copy dependency files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the app.
# Note: We keep CGO enabled for SQLite support.
RUN go build -o gubi-bot ./cmd/main.go

# Stage 2: Run the binary
FROM debian:bookworm-slim

# Install CA certificates (required for HTTPS/Discord API)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/gubi-bot .

# Create a folder for the database
RUN mkdir data

# Run the bot
CMD ["./gubi-bot"]
