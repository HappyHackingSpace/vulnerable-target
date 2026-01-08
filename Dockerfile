# Build stage for development tools
FROM golang:1.24.0-bullseye AS dev-tools

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    git=1:2.30.* \
    make=4.3-4.1 \
    python3-pip \
    python3-setuptools \
    && rm -rf /var/lib/apt/lists/* && \
    pip3 install --no-cache-dir pre-commit

# Install Go tools with versions compatible with Go 1.24.0
RUN go install mvdan.cc/gofumpt@v0.5.0 && \
    go install github.com/kisielk/errcheck@v1.6.3 && \
    go install github.com/go-delve/delve/cmd/dlv@v1.22.0 && \
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest && \
    go install github.com/securego/gosec/v2/cmd/gosec@v2.19.0

# Builder stage
FROM golang:1.24.0-bullseye AS builder

WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go/bin/vulnerable-target ./cmd/vt

# Final stage for application
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    docker-cli \
    git \
    bash \
    curl \
    jq \
    yamllint \
    && rm -rf /var/cache/apk/*

# Copy Go tools from dev-tools stage
COPY --from=dev-tools /go/bin/* /usr/local/bin/

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /go/bin/vulnerable-target .

# Set the entrypoint to bash by default
ENTRYPOINT ["/bin/bash"]
