# BUILD
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# RUNTIME
FROM alpine:latest

WORKDIR /app

# Create a non-root user and group
RUN addgroup -S appuser && adduser -S appuser -G appuser

# Copy the built application
COPY --from=builder /app/main .

# Change ownership of the application binary
RUN chown appuser:appuser /app/main

# Switch to the non-root user
USER appuser

ENV BASE_URL=https://mcp-api.depshub.com
ENV TRANSPORT=stdio
ENV PORT=8080

CMD ["./main"]
EXPOSE 8080
