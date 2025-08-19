# k3d + myapp

This folder contains everything you need to run a local k3s (via k3d) cluster,
a local image registry, deploy the app, and test it from inside the cluster.

## Prereqs

- Docker
- k3d (`brew install k3d` or see https://k3d.io)

## 1) Create the cluster (with registry)

```sh
k3d cluster create --config ./k3d/cluster.yaml
```