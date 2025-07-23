# Build stage
FROM golang:1.21-alpine AS builder

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs and build
RUN swag init -g cmd/api/main.go -o docs
RUN CGO_ENABLED=0 GOOS=linux go build -o apiserver cmd/api/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary and docs
COPY --from=builder /app/apiserver .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/static ./static

# Create non-root user
RUN adduser -D -s /bin/sh apiuser
USER apiuser

EXPOSE 3000

CMD ["./apiserver"]