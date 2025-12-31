FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod (and go.sum if present)
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o paxos-web .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/paxos-web .

EXPOSE 8080
CMD ["./paxos-web"]
