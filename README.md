# Sidecar Library

The Sidecar library is a robust and efficient Go package that enables effortless deployment of a server for forwarding HTTP requests to an incremented port within the same host.

## Prerequisites

-   [Go](https://go.dev/doc/install)

## Installation

To install the Sidecar library, execute the following command:

```bash
go get github.com/k8sidecar/go-lib@v0.0.4
```

Ensure you have Go installed on your machine and your `GOPATH` is set.

## Usage

he Sidecar library provides two custom function types for handling HTTP requests and responses: `TriFunction` and `QuaFunction`.

-   `TriFunction` takes in an HTTP request, an HTTP response writer, and a FilterChain.
-   `QuaFunction` takes in an HTTP request, an HTTP response writer, a Cloud Event, and a FilterChain.

To use the Sidecar library, define your functions based on the `TriFunction` or `QuaFunction` type. Then, instantiate a `SidecarFilter` struct and assign your function to the `TriFunction` or `QuaFunction` field. Finally, call the Listen method on your SidecarFilter instance.

Here is a high-level example:

```go
filter := &sidecar.SidecarFilter{
    TriFunction: yourTriFunction,
}
filter.Listen()
```

For a detailed [usage example](https://github/).
