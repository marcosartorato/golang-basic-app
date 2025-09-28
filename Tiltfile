load('ext://helm_resource', 'helm_resource', 'helm_repo')
k8s_yaml(['k3d/prometheus/manifests/ns.yaml'])

# --- Prometheus --------------------------------------------------------------
helm_repo('prometheus-community', 'https://prometheus-community.github.io/helm-charts')
helm_resource('prometheus', 'prometheus-community/prometheus', namespace='monitoring')

# --- Grafana -----------------------------------------------------------------
helm_repo('grafana', 'https://grafana.github.io/helm-charts')
helm_resource('grafana-ui', 'grafana/grafana', namespace='monitoring')
k8s_resource(workload='grafana-ui', port_forwards=['3000:3000'])

# --- Cluster & namespace -----------------------------------------------------
# Tilt will apply resources into the "myapp" namespace declared in your YAML.
# It’s fine if the namespace doesn’t exist yet; Tilt applies it first.
k8s_yaml([
  'k3d/k8s/ns.yaml',
  'k3d/k8s/depl.yaml',
  'k3d/k8s/svc.yaml',
  'k3d/k8s/svc-metrics.yaml',
  'k3d/k8s/busybox.yaml',
])

# --- Image Build -------------------------------------------------------------
# We keep the image name exactly as in your Deployment: k3d-tilt-registry:5000/myapp:latest
# This plays nicely with the k3d registry defined in k3d/cluster.yaml.

# Keep the image ref simple; Tilt will rewrite to the default registry above
docker_build('myapp', context='.')

# --- Port Forwards -----------------------------------------------------------
# Forward the app Service and metrics Service to your host so you can curl locally too.
# (You can also just use the BusyBox pod inside the cluster.)
k8s_resource('myapp', port_forwards=['8080:80'])

# BusyBox is just a helper pod; no ports
k8s_resource('busybox-curl')

