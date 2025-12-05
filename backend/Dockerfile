# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install basic dependencies (in case needed for some go packages)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/

# Install certificates for HTTPS calls (Gemini/YouTube APIs)
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
# Optional: Copy .env if you want to bake it in, but better to supply via docker-compose
# COPY .env . 

EXPOSE 3000

CMD ["./main"]
