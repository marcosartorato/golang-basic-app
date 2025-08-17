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