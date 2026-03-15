FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy module file first for better layer caching
COPY go.mod ./
RUN go mod download

# Copy the full module source (handlers, services, repository, models, main.go, etc.)
COPY . .

# Build the API binary
RUN go build -o api .

EXPOSE 8080

ENTRYPOINT ["./api"]