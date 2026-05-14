# Agnos DevOps Candidate Assignment

This repository contains the implementation of a production-ready DevOps setup for a Go-based API and Background Worker.

## Architecture Overview

- **API Service**: A Go-based REST API that provides a `/health` endpoint and exposes Prometheus metrics at `/metrics`.
- **Worker Service**: A Go-based background worker that performs periodic tasks and exposes metrics at `:9090/metrics`.
- **Docker**: Both services use multi-stage builds for minimal image size and security.
- **Kubernetes**: 
  - **High Availability**: API is deployed with 3 replicas, Worker with 2 replicas.
  - **Auto-scaling**: HPA is configured for the API service based on CPU utilization.
  - **Reliability**: Readiness and liveness probes ensure only healthy pods serve traffic.
- **CI/CD**: GitHub Actions pipeline handles linting, testing, building, and security scanning.
- **Monitoring**: Prometheus metrics for request latency, error rates, and worker status.

## Project Structure

```
├── api/
│   ├── Dockerfile
│   ├── main.go
│   ├── main_test.go
│   └── go.mod
├── worker/
│   ├── Dockerfile
│   ├── main.go
│   ├── main_test.go
│   └── go.mod
├── k8s/
│   ├── api.yaml
│   ├── worker.yaml
│   ├── hpa.yaml
│   ├── configmap.yaml
│   └── prometheus.yaml
└── .github/workflows/
    └── main.yml
```

## Setup & Usage

### Local Development
1. Run API: `cd api && go run main.go`
2. Run Worker: `cd worker && go run main.go`

### Docker
Build images locally:
```bash
docker build -t agnos-api ./api
docker build -t agnos-worker ./worker
```

Alternatively, use Docker Compose to spin up the entire stack:
```bash
docker-compose up --build
```

### Kubernetes
1. **Enable Metrics Server** (Required for HPA):
   On Docker Desktop or Minikube, you may need to install the metrics server:
   ```bash
   kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
   # If using Docker Desktop, patch for insecure TLS:
   kubectl patch deployment metrics-server -n kube-system --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
   ```

2. **Apply manifests**:
   ```bash
   kubectl apply -f k8s/
   ```

3. **Access Prometheus**:
   ```bash
   kubectl port-forward svc/prometheus-service 9090:9090
   ```
   Open `http://localhost:9090` in your browser. (Alternatively, use the NodePort: `http://localhost:30090`)

## Failure Scenario Handling

### 1. API crashes during peak hours
- **Detection**: Liveness probe fails, or HighErrorRate alert triggers.
- **Handling**: Kubernetes automatically restarts crashed pods. HPA scales up replicas if CPU usage is high.

### 2. Worker fails and infinitely retries
- **Detection**: `WorkerStalled` alert triggers if the last success timestamp is too old.
- **Handling**: Logs will show the error in structured JSON format. Kubernetes will keep the pod running but the alert will notify engineers of the stall.

### 3. Bad deployment is released
- **Detection**: Readiness probe fails (pods never become "Ready") or `HighErrorRate` alert triggers immediately after deployment.
- **Handling**: Kubernetes `RollingUpdate` strategy (default) will stop the rollout if new pods fail probes. Use `kubectl rollout undo deployment/api-deployment` to revert.

### 4. Kubernetes node down
- **Detection**: `APIPodCrashLooping` or node-level alerts.
- **Handling**: Pods are automatically rescheduled by Kubernetes to other healthy nodes (HA configuration with multiple replicas ensures availability).

## Monitoring & Observability
- **Logs**: All logs are in structured JSON format.
- **Metrics**: 
    - `http_requests_total`: Tracks request count and status codes.
    - `http_request_duration_seconds`: Tracks latency.
    - `worker_last_success_timestamp_seconds`: Tracks worker health.
