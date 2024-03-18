# Stage 1: Build the Go application
FROM golang:1.22.1-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the necessary files to the container
COPY go.mod go.sum ./

# Build the Go application
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /cmd/api/v1  ./cmd/api/v1/main.go

FROM alpine:latest as runner

WORKDIR /root/

COPY --from=builder /cmd/api/v1 .

CMD ["./v1"]