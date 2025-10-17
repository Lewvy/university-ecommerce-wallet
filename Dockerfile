FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM alpine:latest

WORKDIR /app

# Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/bin/api .

EXPOSE 8081

# The command to run your application
CMD ["./api"]
