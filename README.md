# golang-basic-app

Minimal Go template for future reference/showcase. Includes local pipeline commands (lint/test/build), Docker image, and a quick k8s (k3s via k3d) deploy.

## Overview

Basic GoLang application template:

1. "Hello, World!" web server on port 8080.
2. Web server exposing scraping point on 9090 port.
3. Graceful termination on SIGTERM signal.

## Local Environment

Basic commands:

- Run linter

```sh
golangci-lint run
```

- Run test

```sh
go test -v -coverprofile=out.cover ./...
```

- Build and run the application

```sh
go build -o bin/myapp cmd/myapp/main.go && ./bin/myapp
```

To ensure consistent results between your local environment and the repository’s CI pipeline, use the same tool versions.
For example:

```
- name: Run golangci-lint
  uses: golangci/golangci-lint-action@v8
  with:
    version: v2.4.0
```

The action [golangci-lint-action@v8](https://github.com/golangci/golangci-lint-action/tree/v8) defaults to [golangci-lint v2.1.0](https://github.com/golangci/golangci-lint/tree/v2.1.0), but the `version` field overrides it to use [v2.4.0](https://github.com/golangci/golangci-lint/tree/v2.4.0).

### Make Targets

If you prefer `Make`, add a `Makefile` with these:

```make
.PHONY: lint test build run clean
lint:  ; golangci-lint run
test:  ; go test -v -coverprofile=out.cover  ./...
build: ; mkdir -p bin && go build -o bin/myapp ./cmd/myapp/main.go
run:   ; go run ./cmd/myapp
clean: ; rm -rf bin
```

## Docker

Sometimes it's easier to use Docker containers instead of reconfiguring your local environment for a specific repository:

- Run linter

```docker
docker run --rm \
    -v "$(pwd)":/myapp \
    -w /myapp \
    golangci/golangci-lint:v2.4.0-alpine \
    golangci-lint run ./...
```

- Run test

```docker
docker run --rm \
    -v "$(pwd)":/myapp \
    -w /myapp \
    golang:1.25.0-alpine \
    go test -v -coverprofile=out.cover ./...
```

- Build local image

```docker
docker build -t myapp:dev .
```

- Run container

```docker
docker run --rm -p 8080:8080 -p 9090:9090 myapp:dev
```

## Local k8s Cluster

### Install

The local Kubernetes cluster is powered by [k3d](https://k3d.io/stable/#what-is-k3d), which allows us to create multi-node [k3s](https://github.com/k3s-io/k3s) clusters in Docker.

To install k3d, follow the official [installation instructions](https://k3d.io/stable/#installation).

For easier local development and deployment, use [tilt](https://docs.tilt.dev/).

To install Tilt, follow the official [installation guide](https://docs.tilt.dev/install.html).

To install kubectl, follow the official [installation guide](https://kubernetes.io/docs/tasks/tools/).

### Setup

```
# Create k3d cluster.
k3d cluster create --config ./k3d/cluster.yaml

# Start Tilt to build, deploy, and watch the app.
# Tilt will auto-redeploy on code changes.
tilt up
```

### Testing Command

```
# Bind the application's service port to local port.
kubectl port-forward svc/myapp -n myapp 8888:80

# Test the HTTP server.
curl localhost:8888/hello

# Bind the application's metric service port to local port.
kubectl port-forward svc/myapp-metrics -n myapp 9090:9090

# Test the metrics scraping point.
curl localhost:9090/metrics
```

### Clean Up

1. Stop Tilt.

Press `Ctrl`+`C` on the tilt console tab to stop it.

2. Stop cluster.

```
# Stop the cluster.
# The cluster name is ctrl due to ./k3d/cluster.yaml.
k3d cluster stop ctrl
```

3. Delete the cluster.

```
# Delete the cluster.
k3d cluster delete ctrl
```

### Additional Information

Take a look [here](./k3d/README.md).

## Metrics

More about the provided metrics [here](./internal/metrics/README.md).

More about how to access the metrics in the local environment [here](./k3d/README.md)

## Tests

By default, all tests are executed, but you can skip the tests related to the `handlers`.

To skip them, set the environment variable `SKIP_HANDLER_TEST` to `true` when running the tests:

```sh
SKIP_HANDLER_TEST=true go test -v -coverprofile=out.cover ./...
```

## License

See [LICENSE](./LICENSE).