FROM  golang:1.16-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o backend ./cmd/sensor-stream

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/backend /app/backend

WORKDIR /app/data

# Run the web service on container startup.
CMD ["/app/backend"]