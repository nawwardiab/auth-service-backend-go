# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk update && apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source
COPY . .

# ADD THIS LINE - Tidy modules after copying source
RUN go mod tidy

# Then build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o server \
    ./internal/cmd/main.go

# Stage 2: Runtime (keep as-is)
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]