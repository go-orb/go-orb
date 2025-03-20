# Go-Orb: From Zero to Hero in Go Development

Hey there! Go-Orb is a comprehensive, high-performance framework for building both monoliths and distributed systems in Go. We designed it as the successor to go-micro with tons of improvements in architecture, performance, and developer experience.

> **Beta Release Coming Soon!** We're actively developing Go-Orb and it'll be available as a beta release soon. Stay tuned for updates!

## Why Choose Go-Orb?

### Production-Ready Architecture

Go-Orb gives you a rock-solid foundation for building distributed systems with:

- **Near to zero Reflection**: Better type safety and faster performance by kicking runtime reflection to the curb
- **Wire-Based Dependency Injection**: Compile-time safety with no nasty globals or runtime surprises
- **Pluggable Architecture**: Swap components in and out without touching your core application code

### Clean, Interface-Driven Design

Go-Orb keeps things neat with a clean separation of concerns:

- **Core Interfaces Only**: [go-orb/go-orb](https://github.com/go-orb/go-orb) is just interfaces and minimal glue code
- **Plugins Do The Work**: The real implementations live in [go-orb/plugins](https://github.com/go-orb/plugins)
- **Simple Plugin System**: Just add a single blank import (`import _ "github.com/go-orb/plugins/..."`) and you're good to go
- **Mix and Match**: Pick exactly the plugins that fit your needs
- **Extensible**: Easily whip up your own plugins that implement core interfaces

### Simplified Distributed Systems Development

Focus on your business logic while Go-Orb handles all the distributed systems headaches:

- **Service Discovery**: Automatic service registration and name resolution
- **Load Balancing**: Smart request distribution across service instances
- **Fault Tolerance**: Built-in retries and circuit breaking to keep things running
- **Message Encoding**: Dynamic content-type based encoding and decoding

### Super-Fast In-Memory Communication

Go-Orb's got a high-performance in-memory adapter that lets you:

- **Direct Handler Calls**: Lightning-fast in-process communication with no serialization overhead
- **Same API as Network Calls**: Use the same client interface whether you're calling local or remote services
- **Perfect for Monoliths**: Start with everything in one process and split out services as needed
- **Seamless Testing**: Test your services in isolation without network dependencies

### Support for Modern Protocols

Communicate however you need to with support for:

- **gRPC**: High-performance RPC with bi-directional streaming
- **HTTP/HTTPS**: RESTful APIs with full support for HTTP/1.1 and HTTP/2
- **DRPC**: Ultra-fast RPC alternative with reduced overhead
- **HTTP/3**: Next-gen HTTP with QUIC
- **Event-Driven Communication**: Asynchronous messaging for decoupled architectures

Check out our [benchmarks](https://github.com/go-orb/go-orb/wiki/RPC-Benchmarks) for more details.

### Developer-Friendly Experience

Get up and running quickly and stay productive:

- **Intuitive APIs**: Clean, consistent interfaces that are easy to understand
- **Flexible Configuration**: Configure via files, environment variables, or code
- **Comprehensive Documentation**: Detailed guides and examples to get you started
- **Strong Test Support**: Built with testability in mind

### Quality-Focused Development

Go-Orb is all about code quality and reliability:

- **Comprehensive Static Analysis**: The entire codebase is validated with golangci-lint using strict rules
- **Comprehensive Test Suite**: High test coverage across all components and plugins
- **CI Enforcement**: Quality checks are automatically run for all pull requests
- **No Compromises**: Strict linting and testing requirements ensure consistent quality
- **Production-Ready**: Our strict development practices mean you can trust Go-Orb even after the beta release

## What Sets Go-Orb Apart?

### Start Small, Scale Big

Go-Orb is designed to grow with your application:

- **Start with a Monolith**: Begin development with all services in one process using the in-memory adapter
- **Transition to Microservices**: Gradually extract services without changing your business logic
- **Hybrid Architecture**: Run performance-critical components in-process while distributing others
- **Progressive Scaling**: Add more instances of specific services as your load increases

### Multiple Entry Points

Unlike traditional frameworks that lock you into a single protocol, Go-Orb lets you expose your services over multiple protocols simultaneously. Configure different handlers for different protocols in a single, cohesive service.

### Powerful, Multi-Stage Configuration System

Go-Orb's configuration system offers incredible flexibility:

- **Smart Configuration Loading**: Automatically merges configurations from multiple sources in a priority order:
  1. Predefined defaults
  2. User configuration files (local or remote)
  3. Environment variables
  4. Command-line flags
- **Section-Based Loading**: Load only the configuration sections you need, when you need them

### No Massive Structs Required

Define configuration types that match exactly what your component needs

### Format Agnostic

Support for YAML, TOML, JSON, and more

### Advanced Configuration

Define your entire service architecture in simple YAML, TOML or JSON:

```yaml
service1:
  server:
    logging:
        plugin: lumberjack
        level: INFO
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

      - name: http
        plugin: http
        insecure: true

      - name: drpc
        plugin: drpc
  client:
    middlewares:
      - name: log
      - name: retry
    logging:
        level: TRACE
  registry:
    plugin: kvstore
    kvstore:
        plugin: natsjs
        servers:
        - nats://localhost:9222
```

### Our Vision: A Unix-Like Service Ecosystem

We're building Go-Orb with a vision inspired by the Unix philosophy of "do one thing and do it well":

- **Focused Service Design**: Each service should have a clear, single responsibility
- **Composable Architecture**: Services working together through well-defined interfaces
- **Ready-to-Use Components**: We're building toward a library of pre-built services like API gateways and auth systems
- **Easy Integration**: Wire-based dependency injection makes service composition natural
- **Build Your Ecosystem**: Create your own tailored platform by mixing and matching exactly what you need

As we continue developing Go-Orb, this vision guides our roadmap and architecture decisions.

Check out [services](https://github.com/go-orb/services) for more details.

### Protocol-Conformant Handlers

Write handlers once and expose them through any protocol:

```go
// Simple, type-safe client calls
resp, err := client.Call[HelloResponse](
    context.Background(), 
    clientDi, 
    "org.orb.svc.hello", 
    "Say.Hello", 
    &req
)
```

Or typesafe generated Handlers:

```go
cli := authproto.NewAuthClient(clientFromWire)
req := &authproto.Req{Token: "someToken"}
resp, err := cli.Authenticate(
    ctx,
    serverName,
    req,
    opts...,
)
```

### Structured Logging

Go-Orb comes with built-in support for modern structured logging based on Go's standard library slog package, making debugging and monitoring a breeze.

## Use Cases

- **Microservices Architecture**: Build, deploy, and scale individual services independently
- **API Gateways**: Create unified entry points for your microservices ecosystem
- **Event-Driven Systems**: Implement pub/sub patterns for asynchronous processing
- **Edge Computing**: Deploy lightweight services closer to your users
- **Cloud-Native Applications**: Perfect for containerized environments and Kubernetes

## Getting Started

The best way to get started with Go-Orb is to check out our examples repository at [github.com/go-orb/examples](https://github.com/go-orb/examples). It's packed with sample services that showcase Go-Orb's capabilities:

- Simple services with different protocols
- Event-driven architectures
- API gateways and proxies
- Authentication and authorization
- Performance benchmarks

Head over to the repository for step-by-step instructions on running the examples and building your own services with Go-Orb.

Visit [github.com/go-orb/go-orb](https://github.com/go-orb/go-orb) to learn more about the core framework.

## Community and Support

Join our friendly community:

- **Matrix**: [https://matrix.to/#go-orb:jochum.dev](https://matrix.to/#/#go-orb:jochum.dev) - Real-time chat and support
- **Discord**: [https://discord.gg/go-orb](https://discord.gg/go-orb) - Another way to connect
- **GitHub**: [https://github.com/go-orb](https://github.com/go-orb) - Open source development and issue tracking