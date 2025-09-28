# --- Prometheus --------------------------------------------------------------
# Ensure the repo exists locally
local('helm repo add prometheus-community https://prometheus-community.github.io/helm-charts || true')
local('helm repo update')

# Render the chart to YAML that Tilt will apply
local('helm template prometheus prometheus-community/prometheus --namespace monitoring -f ./k3d/prometheus/prom-values.yaml --create-namespace > k3d/prometheus/manifests/prometheus.yaml')

# Apply with Tilt
k8s_yaml([
  'k3d/prometheus/manifests/ns.yaml',
  'k3d/prometheus/manifests/prometheus.yaml',
])

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

