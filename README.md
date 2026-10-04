# Mini Load Balancer

A small HTTP reverse proxy written in Go. It distributes requests across a configured list of backend servers using round robin.

## Requirements

- Go 1.26 or later

This first version uses only the Go standard library, so there are no third-party dependencies to install.

## Run

```sh
go run ./cmd/loadbalancer \
  -listen :8080 \
  -backends http://localhost:9001,http://localhost:9002
```

Start two HTTP services on ports 9001 and 9002 before sending a request to `http://localhost:8080`. The load balancer forwards each request to the next backend in the list.

Options can also be supplied with `LB_LISTEN_ADDR` and `LB_BACKENDS` environment variables. Command-line flags take precedence.

## Layout

```text
cmd/loadbalancer/       executable entry point
internal/config/        flag and environment configuration
internal/balancer/      backend selection and reverse proxy
```

`internal` keeps these packages private to this module. The executable is separate from the balancing package so the latter can be tested and reused inside the project without coupling it to process startup.

## Current scope

This starter implements round-robin selection and HTTP reverse proxying. It does not yet include active health checks, graceful shutdown, TLS termination, or metrics. Those are natural next steps after the initial request path is clear.
