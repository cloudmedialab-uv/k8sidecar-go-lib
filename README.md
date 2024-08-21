# Sidecar Library

The Sidecar library is a robust and efficient Go package that enables the deployment of sidecar proxies that receive a HTTP request, perform some action, forwards the request to an incremented port within the same host, and then can perform some action before returning the response.

The library can be used as is or can be modified to suit the needs.

## Prerequisites

-   [Go](https://go.dev/doc/install)

## Installation

To install the Sidecar library, execute the following command:

```bash
go get github.com/cloudmedialab-uv/k8sidecar-go-lib
```

Ensure you have Go installed on your machine and your `GOPATH` is set.

### Common instalation fails

If `go get` encounters issues, check the page [go.dev/doc/faq#git_https](https://go.dev/doc/faq#git_https) for potential errors in the configuration of Go git modules.

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

## Example: ratelimiter sidecar

The folder [example/ratelimiter](example/ratelimiter) contains an example of the use of this library to build a sidecar that performs rate limit. The number of requests allowed per second can be passed in the environment variable `RATE` (by default 100).

If the number of request is bigger than the configured rate then a  HTTP 429 Too Many Requests client error response is returned.

We provide a Dockerfile to build the image and we also provide an image in dockerhub: `cloudmedialab/sidecar-ratelimiter:1.0.0`.

See [k8sidecar](https://github.com/cloudmedialab-uv/k8sidecar) for an example of usage to define a Filter using this sidecar.
