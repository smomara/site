# Build stage
FROM golang:1.22.4-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

ENV GOOGLE_CLOUD_RUN_SERVICE_URL=https://site-m5svpyta6q-uc.a.run.app

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./main"]