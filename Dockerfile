FROM golang:1.24.4-alpine
WORKDIR /app

# Copy executable
COPY server .

# Document port
EXPOSE 8080

# Entrypoint for Docker
ENTRYPOINT ["./server"]
