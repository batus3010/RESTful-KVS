# ─── Stage 1: Build the Go binary ─────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

# Install git (needed if you import any modules via git)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (so dependencies are cached unless they change)
COPY go.mod ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build your server binary.
# Adjust the path if your entrypoint is under cmd/RESTful or you want to build server.go instead.
RUN go build -o kvs-server ./cmd/RESTful

# ─── Stage 2: Create a minimal runtime image ───────────────────────────────────
FROM alpine:3.18

# (Optional) add CA certs if your app makes HTTPS calls
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/kvs-server .

# Expose whichever port your server listens on
EXPOSE 5000

# Run the server
ENTRYPOINT ["./kvs-server"]
