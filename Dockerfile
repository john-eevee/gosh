# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-X github.com/gosh/pkg/version.Version=$(git describe --tags --always)" \
    -o gosh ./cmd/gosh

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/gosh /usr/local/bin/gosh

# Set entrypoint
ENTRYPOINT ["gosh"]

# Default command
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="gosh" \
      org.opencontainers.image.description="HTTPie CLI alternative built with Go" \
      org.opencontainers.image.url="https://github.com/john-eevee/gosh" \
      org.opencontainers.image.source="https://github.com/john-eevee/gosh" \
      org.opencontainers.image.documentation="https://github.com/john-eevee/gosh/blob/main/README.md"
