---
name: agnos-devops-expert
description: Guidance for completing the Agnos DevOps assignment. Use when setting up Docker, Kubernetes, CI/CD, and Monitoring according to the assignment's specific requirements (HA, multi-stage builds, environment variables for DEV/UAT/PROD, etc.).
---

# Agnos DevOps Expert

This skill provides specialized workflows and patterns for the Agnos Candidate Assignment.

## Workflows

### 1. Dockerization
Implement multi-stage builds for Go services with environment variable support.
- **Reference**: [docker-patterns.md](references/docker-patterns.md)

### 2. Kubernetes Orchestration
Design high-availability manifests with resource limits, probes, and HPA.
- **Reference**: [k8s-patterns.md](references/k8s-patterns.md)

### 3. CI/CD Pipeline
Set up GitHub Actions for linting, testing, building, and security scanning.
- **Reference**: [ci-cd-patterns.md](references/ci-cd-patterns.md)

### 4. Observability & Monitoring
Implement structured JSON logging and Prometheus metrics/alerts.
- **Reference**: [observability-patterns.md](references/observability-patterns.md)

## Core Principles
- **Reliability**: Use HA patterns and readiness/liveness probes.
- **Security**: Implement multi-stage builds and security scans.
- **Observability**: Ensure all logs are structured and key metrics are tracked.
- **Environment Parity**: Use the same Dockerfile across DEV, UAT, and PROD, varying only via env vars.
