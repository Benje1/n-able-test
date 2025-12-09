# Service Health Monitoring Program

This is a basic service health monitoring tool written in Go.  
It is designed to set itself up using only a `services.yaml` configuration file.

In production you would normally mount your own `services.yaml` rather than bundling a default one inside the container. For testing purposes, a default file is included in the repo.

A GitHub Actions pipeline is included. It runs all tests, and if they succeed (and the branch is `main`), it should build a Docker image automatically.

A minimal example of how this could be deployed to AWS (via a simple curl-based Lambda) is also included, but it is not a full AWS infrastructure implementation.

I did use the yaml as described in the task outline

## Assumptions
- Services will not change frequently.
- The program does not require persistent state or a database.
- Restarting the container with an updated `services.yaml` is sufficient for configuration changes.

---

## Usage

```
To run as a non dockerd server:
go run main.go

Install Go dependencies
go mod download

To run the tests:
go test .\ServiceMonitor\

Load a container image from a .tar file:
docker load -i benjamin-test.tar

To run as a docketed server:
docker run -p 8080:8080 nable-benjamin:latest

To run docker with a mounterd services.yaml:
docker run -p 8080:8080 \
  -v $(pwd)/services.yaml:/app/services.yaml \
  nable-benjamin:latest

To create new docker container:
docker build -f docker -t nable-benjamin .

To hit the endoint:
curl http://localhost:8080/health/aggregate
```