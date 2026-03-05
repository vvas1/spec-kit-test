# Multi-stage build for issue-tracker backend (003-containerize-backend)
# Build stage: compile Go binary
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Run stage: minimal image, no ports published (reachable only from same Compose network)
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /server .

ENV MONGODB_URI=mongodb://db:27017
EXPOSE 8080

ENTRYPOINT ["./server"]
