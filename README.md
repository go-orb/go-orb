# go-orb

Orb is a framework for distributed systems development, it can be seen as the successor of [go-micro.dev/v4](https://github.com/go-micro/go-micro).

## :warning: WIP

This project is a work in progress, please do not use yet!

## What's new since v4

### Use of [wire](https://github.com/google/wire)

It was the main reason for starting orb, wire allows us to decouple the components and plugins.

### No more reflect

we have been working hard on removing all usage of reflect.

### Multiple Entrypoints

Orb allows you to listen on multiple port's with different protocols: gRPC, HTTP, HTTP2, H2C, HTTP3.
See the config system entry on howto configure it.

### Advanced [config system](config)

With orb you can configure your plugins with a config file or environment options.

```yaml
service1:
  server:
    http:
      gzip: true
      handlers:
        - Streams
      # middleware:
      #   - middleware-1
      entrypoints:
        - name: ep1
          address: :4512
          insecure: true
          h2c: true
  registry:
    enabled: true
    plugin: mdns
    timeout: 350
    domain: svc.orb
```

```yaml
service1:
  server:
    grpc:
      insecure: true
      handlers:
        - Streams
      # middleware:
      #   - middleware-1
      # streamMiddleware:
      #   - middleware-S1
      entrypoints:
        - name: ep1
          address: :4512
          health: false
          reflection: false
          # handlers:
          #   - handler-1
          #   - handler-2
          # middleware:
          #   - middleware-1
          #   - middleware-4
  registry:
    plugin: nats
    address: nats://10.0.0.1:4222
    quorum: false
```

These 2 config's with different options will both work, we first parse the config, get the "plugin" from it and pass a `map[any]any` with all config data to the plugin.

Both work with a single binary. :)

## Community

- Chat with us on [Discord](https://discord.gg/sggGS389qb)

## Development

### Quirks

#### No go-orb/plugins imports here

To prevent import cycles it's not allowed to import github.com/go-orb/plugins here.

## Authors

- [David Brouwer](https://github.com/Davincible/)
- [René Jochum](https://github.com/jochumdev)

## License

go-orb is Apache 2.0 licensed.
