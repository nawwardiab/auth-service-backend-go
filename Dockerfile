# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install build dependencies (git may be needed for some Go modules)
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=0 creates a statically-linked binary (no C dependencies)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o server \
    ./internal/cmd/main.go

# Stage 2: Create minimal runtime image
FROM scratch

# Copy ca-certificates from builder (for HTTPS)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data (optional, but good to have)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder
COPY --from=builder /build/server /server

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/server"]