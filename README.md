# ![go-orb Logo](docs/logo-header.png) [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/go-orb/go-orb?tab=doc) [![Go Report Card](https://goreportcard.com/badge/github.com/go-orb/go-orb)](https://goreportcard.com/report/github.com/go-orb/go-orb) [![Discord](https://dcbadge.vercel.app/api/server/sggGS389qb?style=flat-square&theme=default-inverted)](https://discord.gg/sggGS389qb)

Orb is a framework for distributed systems development, it can be seen as the successor of [go-micro.dev/v4](https://github.com/go-micro/go-micro).

## :warning: WIP

This project is a work in progress, please do not use yet!

We are not far from the releasing the first Alpha, see the [Roadmap](https://github.com/orgs/go-orb/projects/1/views/1).

## :rocket: What's new since forking from go-micro v4

### Use of [wire](https://github.com/google/wire)

With wire we gain:

- compile-time safety
- no more globals

It was the main reason for starting orb, wire allows us to decouple the components and plugins.

### No more reflect

we have been working hard on removing all usage of reflect.

### Multiple Entrypoints

Orb allows you to listen on multiple port's with different protocols: gRPC, HTTP, HTTPS, DRPC, HTTP2, H2C, HTTP3.
See the config system entry on howto configure it.

### Advanced [config system](config)

With orb you can configure your plugins with a config file or environment options.

```yaml
service1:
  server:
    handlers:
      - Streams
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
      - Streams
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
resp , err := client.Call[map[string]any](context.Background(), clientDi, "org.orb.svc.hello", "Say.Hello", `{"name": "Alex"}`, client.WithContentType("application/json"))
```

#### map[string]any{}

```go
req := make(map[string]any)
req["name"] = "Alex"

resp , err := client.Call[map[string]any](context.Background(), clientDi, "org.orb.svc.hello", "Say.Hello", req, client.WithContentType("application/json"))
```

#### Structured logging

We like structured logging, this is why we replaced all logging with one based on [slog](https://pkg.go.dev/log/slog).

#### go-orb/go-orb is just interfaces

We made sure that go-orb/go-orb (the core) is just a bunch of interfaces as well as some glue code, the most real code lives in [go-orb/plugins](https://github.com/go-orb/plugins).

#### Linted and analyzed

Everything is linted and staticaly analyzed by golangcli-lint, enforced with CI/CD pipelines on github.

## Community

Chat with us on [Discord](https://discord.gg/sggGS389qb).

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

We do not accept commits that fail to lint, either use `./scripts/test.sh lint all` or install [Trunk](https://trunk.io/) and it's extension for your editor.

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
