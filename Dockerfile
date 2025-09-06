# Multi-stage Dockerfile for OpenX
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go modules first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION:-dev} -X main.commit=${COMMIT:-unknown} -X main.date=${BUILD_DATE:-unknown}" \
    -o openx ./cmd/openx

# Final stage
FROM scratch

# Import certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/openx /usr/local/bin/openx

# Create a non-root user
COPY --from=builder /etc/passwd /etc/passwd
USER nobody

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/openx"]

# Default command
CMD ["--help"]

# Metadata
LABEL org.opencontainers.image.title="OpenX"
LABEL org.opencontainers.image.description="Cross-platform developer environment control tool"
LABEL org.opencontainers.image.url="https://github.com/muthuishere/openx"
LABEL org.opencontainers.image.source="https://github.com/muthuishere/openx"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.vendor="Muthu Kumar"
