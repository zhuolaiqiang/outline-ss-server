# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23 AS builder
WORKDIR /app

# Pre-cache dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy sources
COPY . .

# Inject version via ldflags
ARG VERSION=dev

# Build static Linux binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /app/bin/outline-ss-server \
    ./cmd/outline-ss-server

# Runtime stage (distroless, non-root)
FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=builder /app/bin/outline-ss-server /usr/local/bin/outline-ss-server
# Default config inside image
COPY cmd/outline-ss-server/config_example.yml /etc/outline/config.yml

USER nonroot
EXPOSE 9000/tcp 9000/udp 9001/tcp 9001/udp 9090/tcp
ENTRYPOINT ["/usr/local/bin/outline-ss-server"]
# Default runtime arguments (override by providing args to docker run)
CMD ["-config", "/etc/outline/config.yml", "-metrics", ":9090"]
