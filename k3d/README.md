# Local Environment

This folder contains everything you need to run a local k3s (via k3d) cluster,
a local image registry, deploy the app, and test it from inside the cluster.

## Prerequisites

- Docker
- k3d (`brew install k3d` or see https://k3d.io)

## Commands

The following commands can be run from the repository root directory.
This setup uses the ctrl cluster (as defined in `./k3d/cluster.yaml`).

```sh
# Create the cluster (with registry)
k3d cluster create --config ./k3d/cluster.yaml

# Start Tilt to build, deploy, and watch the controller.
# Tilt automatically rebuilds and redeploys on source code changes.
tilt up
```

The cluster is up! The cluster can be explored using `k9s` or `kubectl`.

To clean the local environment:

- Press `Ctrl`+`C` to stop tilt.
- Run `k3d cluster stop ctrl` to stop the cluster.
- Run `k3d cluster delete ctrl` to delete the cluster.
