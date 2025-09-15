# golang-basic-app

Minimal Go template for future reference/showcase. Includes local pipeline commands (lint/test/build), Docker image, and a quick k8s (k3s via k3d) deploy.

## Overview

Basic GoLang application template:

1. "Hello, World!" web server on 8080 port.
2. Web server exposing scraping point on 9090 port.
3. Graceful termination on SIGTERM signal.

## Local Pipeline (lint/test/build/run)

```
# Run linter
golangci-lint run

# Run test
go test -v -coverprofile=out.cover ./...

# Build and run binary
go build -o bin/myapp cmd/myapp/main.go
./bin/myapp
```

## Make Targets

If you prefer `Make`, add a `Makefile` with these:

```make
.PHONY: lint test build run clean
lint:  ; golangci-lint run
test:  ; go test -v -cover ./...
build: ; mkdir -p bin && go build -o bin/myapp ./cmd/myapp/main.go
run:   ; go run ./cmd/myapp
clean: ; rm -rf bin
```

## Docker

```
# Run linter
docker run --rm \
    -v "$(pwd)":/myapp \
    -w /myapp \
    golangci/golangci-lint:v2.4.0-alpine \
    golangci-lint run ./...

# Run test
docker run --rm \
    -v "$(pwd)":/myapp \
    -w /myapp \
    golang:1.25.0-alpine \
    go test -v -coverprofile=out.cover ./...

# Build local image
docker build -t myapp:dev .

# Run container
docker run -rm -p 8080:8080 -p 9090:9090 myapp:dev
```

## Local k8s Cluster

### Install

The local Kubernetes cluster is powered by [k3d](https://k3d.io/stable/#what-is-k3d), which allows us to create multi-node [k3s](https://github.com/k3s-io/k3s) clusters in Docker.

To install k3d, follow the official [installation instructions](https://k3d.io/stable/#installation).

For easier local development and deployment, use [tilt](https://docs.tilt.dev/).

To install Tilt, follow the official [installation guide](https://docs.tilt.dev/install.html).

### Setup

```
# Clean local folders from old artifacts
sudo rm -rf /tmp/k3d

# Create new folders
mkdir -p /tmp/k3d/myapp/server/kubelet /tmp/k3d/myapp/server/containerd
mkdir -p /tmp/k3d/myapp/agent/kubelet /tmp/k3d/myapp/agent/containerd

# Create k3d cluster
k3d cluster create --config ./k3d/cluster.yaml

# Run Tilt
tilt up
```

## License

See [LICENSE](./LICENSE).