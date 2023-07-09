# go-micro/config

Package go-micro/config is a pluggable config provider for loosely coupled components.

It provides 2 main functions

## config.Read

Read reads urls into []source.Data where Data is basicaly map[string]any.

This is done over config/source.Plugins and currently there are 3 Plugins for config.source:

- `cli` provides config from cli/env sources.
- `file` provides config from file sources.
- `http` provides config from http sources.

It's straight forward to write Plugins for config.source and we will provide more:

- nats
- etcd
- consul

An example:

```go
package main

import (
    "log"
    "net/url"

    "github.com/go-orb/go-orb/config"
    _ "github.com/go-micro/plugins/codecs/json/v5"
    _ "github.com/go-micro/plugins/codecs/yaml/v5"
    _ "github.com/go-micro/plugins/config/source/file/v5"
    _ "github.com/go-micro/plugins/config/source/http/v5"
)

func main() {
    // https://raw.githubusercontent.com/go-micro/plugins/main/config/tests/data/set1/registry1.yaml
    u1, err := url.Parse("./data/set1/registry1.yaml")
    if err != nil {
        log.Fatal(err)
    }

    u2, err := url.Parse("https://raw.githubusercontent.com/go-micro/plugins/main/config/tests/data/set1/registry2.json")
    if err != nil {
        log.Fatal(err)
    }

    datas, err := config.Read([]*url.URL{u1, u2}, []string{"app"})
    if err != nil {
        log.Fatal(err)
    }
}
```

## config.Parse

Parse parses the config from config.Read into the given struct.

Example:

```go
    // extend the config.Read example here

    //
    // All from here is in the plugin itself.
    //
    cfg := newRegistryMdnsConfig()
    err := config.Parse([]string{"app", "registry"}, datas, cfg)
    if err != nil {
        log.Fatal(err)
    }
```

## Authors

- [Asim Aslam](https://github.com/asim/) - Author of [go-micro/config](https://github.com/go-micro/go-micro/tree/master/config) on which this is based on.
- [David Brouwer](https://github.com/Davincible/) - Ideas
- [René Jochum](https://github.com/jochumdev) - Developer

## License

go-micro is Apache 2.0 licensed.
