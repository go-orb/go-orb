# ![go-orb Logo](docs/logo-header.png) [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/go-orb/go-orb?tab=doc) [![Go Report Card](https://goreportcard.com/badge/github.com/go-orb/go-orb)](https://goreportcard.com/report/github.com/go-orb/go-orb) [![Discord](https://dcbadge.vercel.app/api/server/sggGS389qb?style=flat-square&theme=default-inverted)](https://discord.gg/4n6E4NYjnR) ![Matrix](https://img.shields.io/matrix/go-orb%3Ajochum.dev?server_fqdn=matrix.jochum.dev&fetchMode=guest&logo=Matrix)

Go Orb is a framework for distributed systems development, it can be seen as the successor of [go-micro.dev/v4](https://github.com/go-micro/go-micro).

The core of go-orb has been completely refactored, to support the removal of reflect and introduction of wire.

## In active Development

While its possible to try out go-orb currently, it is in active development and not ready for production use.
Please have a look at our [roadmap](https://github.com/orgs/go-orb/projects/1) for more details.

## Overview

Go Orb provides the core requirements for distributed systems development including RPC and Event driven communication.
The Go Orb philosophy is sane defaults with a pluggable architecture. We provide defaults to get you started quickly
but everything can be easily swapped out.

## Features

Go Orb abstracts away the details of distributed systems. Here are the main features.

- **Config** - Load dynamic config from anywhere. The config interface provides a way to load application level config from any source such as env vars, file. You can merge the sources and even define fallbacks.

- **Service Discovery** - Automatic service registration and name resolution. Service discovery is at the core of Go Orb service development. When service A needs to speak to service B it needs the location of that service. The default discovery mechanism is multicast DNS (mdns), a zeroconf system.

- **Message Encoding** - Dynamic message encoding based on content-type. The client and server will use codecs along with content-type to seamlessly encode and decode Go types for you. Any variety of messages could be encoded and sent from different clients. The client and server handle this by default. This includes protobuf and json by default.

- **RPC Client/Server** - RPC based request/response with support for bidirectional streaming. We provide an abstraction for synchronous communication. A request made to a service will be automatically resolved, load balanced, dialled and streamed.

- **RPC over event topics** - RPC over event topics, makes RPC even easier.

- **Pluggable Interfaces** - Go Orb makes use of Go interfaces for each distributed system abstraction. Because of this these interfaces are pluggable and allows Go Orb to be runtime agnostic. You can plugin any underlying technology.

- **Strongly tested and linted** - We use golangci-lint to ensure code quality and we have a comprehensive test suite, all lints and tests are run on CI.

## Examples

Please see the [examples](https://github.com/go-orb/examples) repo.

## :rocket: What's new since forking from go-micro v4

### Use of [wire](https://github.com/go-orb/wire)

With wire we gain:

- compile-time safety
- no more globals

It was the main reason for starting orb, wire allows us to decouple the components and plugins.

### No more reflect

we have been working hard on removing all usage of reflect.

### Multiple Entrypoints

Orb allows you to listen on multiple port's with different protocols: gRPC, HTTP, HTTPS, DRPC, HTTP2, H2C, HTTP3.
See the config system entry on howto configure it.

### Advanced [config system](https://github.com/go-orb/go-orb/tree/main/config)

With orb you can configure your plugins with a config file or environment options.

```yaml
service1:
  server:
    handlers:
      - UserInfo
    middlewares:
      - middleware-1
    entrypoints:
      - name: grpc
        plugin: grpc
        insecure: true
        reflection: false
  registry:
    enabled: true
    plugin: mdns
    timeout: 350
    domain: svc.orb
```

```yaml
service1:
  server:
    handlers:
      - UserInfo
    middlewares:
      - middleware-1
      - middleware-2
    entrypoints:
      - name: hertzhttp
        plugin: hertz
        http2: false
        insecure: true

      - name: grpc
        plugin: grpc
        insecure: true
        reflection: false
        handlers:
          - ImOnlyOnGRPC
        middlewares:
          - ImAGRPCSpecificMiddlware

      - name: http
        plugin: http
        insecure: true

      - name: drpc
        plugin: drpc
  registry:
    plugin: consul
    address: consul:8500
```

These 2 config's with different options will both work, we first parse the config, get the "plugin" from it and pass a `map[any]any` with all config data to the plugin.

Both work with a single binary. :)

### Proto conform handlers

Return types as a result instead of HTTP req format.

New:

```go
req := HelloRequest{Name: "test"}

// Look at resp, it's now returned as a result.
resp , err := client.Call[HelloResponse](context.Background(), clientDi, "org.orb.svc.hello", "Say.Hello", &req)
```

Old:

```go
req := c.c.NewRequest(c.serviceName, "Greeter.Hello", in)
out := new(HelloResponse)
err := c.c.Call(ctx, req, out, opts...)
if err != nil {
  return nil, err
}
return out, nil
```

### Send json

go-orb has support for sending nearly anything you throw in as `application/json` to the Server.

#### pre encoded / proxy

```go
resp , err := client.Call[map[string]any](context.Background(), clientDi, "org.orb.svc.hello", "Say.Hello", `{"name": "Alex"}`, client.WithContentType(codecs.MimeJSON))
```

#### map[string]any{}

```go
req := make(map[string]any)
req["name"] = "Alex"

resp , err := client.Call[map[string]any](context.Background(), clientDi, "org.orb.svc.hello", "Say.Hello", req, client.WithContentType(codecs.MimeJSON))
```

#### Structured logging

We like structured logging, this is why we replaced all logging with one based on [slog](https://pkg.go.dev/log/slog).

#### go-orb/go-orb is just interfaces

We made sure that go-orb/go-orb (the core) is just a bunch of interfaces as well as some glue code, the most real code lives in [go-orb/plugins](https://github.com/go-orb/plugins).

#### Linted and analyzed

Everything is linted and staticaly analyzed by golangci-lint, enforced with CI pipelines on github.

## Community

Chat with us on [Discord](https://discord.gg/4n6E4NYjnR).

## Development

### golangci-lint

We use version v1.64.5 of golangci-lint.

```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.64.5
golangci-lint run --config .golangci.yaml
```

### Quirks

#### No go-orb/plugins imports here

To prevent import cycles it's not allowed to import [go-orb/plugins](https://github.com/go-orb/plugins) here.

#### Lint

We do not accept commits that fail to lint, use `./scripts/test.sh lint all`.

## Authors

### go-orb

- [David Brouwer](https://github.com/Davincible)
- [René Jochum](https://github.com/jochumdev)

### go-micro

A lot of this is copy&pasted from [go-micro](https://github.com/go-micro/go-micro/graphs/contributors), top contributors have been:

- [Asim Aslam](https://github.com/asim) the founder/creator of go-micro.
- [Milos Gajdos](https://github.com/milosgajdos)
- [ben-toogood](https://github.com/ben-toogood)
- [Vasiliy Tolstov](https://github.com/vtolstov)
- [Johnson C](https://github.com/xpunch)

## License

go-orb is Apache 2.0 licensed and is based on go-micro.
