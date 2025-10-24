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

# Delete resources created by "tilt up".
tilt down

# Delete the cluster.
k3d cluster delete myapp
```

The cluster is up! The cluster can be explored using `k9s` or `kubectl`.

To clean the local environment:

- Press `Ctrl`+`C` to stop tilt.
- Run `k3d cluster stop ctrl` to stop the cluster.
- Run `k3d cluster delete ctrl` to delete the cluster.

## Metrics and Logs

The local cluster uses Prometheus, Loki, and Grafana as the metrics/logs stack. They are deployed via Helm and Tilt. 

The Grafana user interface is exposed and can be accessed [here](http://localhost:3000/). Use Grafana to explore information related to logs and metrics. Grafana’s dashboard is usually sufficient for exploring logs and metrics. If you need direct access to Prometheus or Loki, you can port-forward their Kubernetes services to a local port.

To get Grafana username and password, run:

```
k get secret/grafana-ui -n monitoring -o jsonpath='{.data.admin-user}' | base64 -d
k get secret/grafana-ui -n monitoring -o jsonpath='{.data.admin-password}' | base64 -d
```

Loki and Prometheus are already configured as Grafana datasources by the related [value](monitoring/grafana/values.yaml) file. To setup additional data sources, goes on the [related section](http://localhost:3000/connections/datasources) in the Grafana user interface.