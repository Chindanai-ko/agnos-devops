# CI/CD Patterns for Agnos DevOps

## GitHub Actions Workflow
Standard pipeline for Go services.

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  lint-and-test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    - name: Lint
      run: go vet ./...
    - name: Test
      run: go test -v ./...

  build-and-push:
    needs: lint-and-test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Build Docker Image
      run: docker build -t agnos-api:${{ github.sha }} ./api
    - name: Security Scan (Optional)
      uses: aquasecurity/trivy-action@master
      with:
        image-ref: 'agnos-api:${{ github.sha }}'
        format: 'table'
        exit-code: '1'
        ignore-unfixed: true
        severity: 'CRITICAL,HIGH'

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    steps:
    - name: Deploy (Mocked)
      run: echo "Deploying to ${{ github.ref_name }} environment..."
```
