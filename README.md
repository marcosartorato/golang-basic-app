# golang-basic-app

Basic GoLang application template:

1. "Hello, World!" web server on 8080 port.
2. Web server exposing scraping point on 9090 port.
3. Graceful termination on SIGTERM signal.

## Instructions

Run locally:

```
go run ./cmd/myapp
```

Build Docker image:

```
docker build -t myapp .
```

Run with Docker:

```
docker run -p 8080:8080 -p 9090:9090 myapp
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