# Metrics Server

This service exposes Prometheus-compatible metrics on `:8080/metrics`.  
The instrumentation follows the **RED pattern** (Rate, Errors, Duration), plus resource/runtime metrics useful for SRE and debugging.

---

## Application metrics

### Requests

- **`http_requests_total{method,route,status}`** 

  Counter. Total number of HTTP requests handled, labeled by:

  - `method`: HTTP verb (e.g., `GET`, `POST`)  
  - `route`: the path (e.g., `/hello`, `/api/message`)  
  - `status`: HTTP response code (e.g., `200`, `500`)  

  **Usage**: traffic rate, error rate, per-endpoint request distribution.

### Request duration

- **`http_request_duration_seconds_bucket{method,route,status,le}`**  
  Histogram. Distribution of request latencies in seconds.  
- **`http_request_duration_seconds_count/sum`**  
  Aggregated counts and totals.

  **Usage**: percentile latency (p50, p95, p99) per route, alerting on slow endpoints.

### Inflight requests

- **`http_inflight_requests{route}`**  
  Gauge. Number of requests currently being processed.

  **Usage**: detect backlog buildup, long-lived requests, saturation.

### Request size

- **`http_request_size_bytes_bucket{method,route}`**  
  Histogram. Approximate request body sizes (from `Content-Length`).  

  **Usage**: spot unusually large uploads or payload trends.

### Response size

- **`http_response_size_bytes_bucket{method,route,status}`**  
  Histogram. Distribution of response payload sizes.  

  **Usage**: track API response growth, detect unexpectedly large responses.

### Panics

- **`http_panics_total{route}`**  
  Counter. Number of recovered panics while serving requests.

  **Usage**: alert if application code crashes inside handlers.

---

## Runtime metrics (from collectors)

These are **standard Prometheus Go client collectors**, automatically refreshed at scrape time.

### Go runtime (`go_*`)

Examples:

- `go_goroutines` – Number of goroutines  
- `go_memstats_alloc_bytes` – Heap allocations  
- `go_gc_duration_seconds` – GC pause times  

**Usage**: monitor memory pressure, goroutine leaks, GC overhead.

### Process (`process_*`)

Examples:

- `process_cpu_seconds_total` – Total user+system CPU time  
- `process_resident_memory_bytes` – Resident memory (RSS)  
- `process_open_fds` – Number of open file descriptors  
- `process_start_time_seconds` – Unix timestamp when process started  

**Usage**: monitor CPU/memory usage, detect FD leaks, track uptime.

---

## Example PromQL queries

- **Error rate (5xx %)**  
  ```promql
  sum(rate(http_requests_total{status=~"5.."}[5m]))
    /
  sum(rate(http_requests_total[5m]))
  ```

---

[Go Back](../../README.md)