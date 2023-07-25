# go-orb/config

Package go-orb/config is a pluggable config provider for loosely coupled components.

It provides 2 main functions and 2 helper

## Functions

### config.Read

Read reads urls into []source.Data where source.Data is basicaly a map[string]any.

This is done over config/source.Plugins, currently there are 3 Plugins for config.source:

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
    _ "github.com/go-orb/plugins/codecs/json"
    _ "github.com/go-orb/plugins/codecs/yaml"
    _ "github.com/go-orb/plugins/config/source/file"
    _ "github.com/go-orb/plugins/config/source/http"
)

func main() {
    // https://raw.githubusercontent.com/go-orb/plugins/main/config/tests/data/set1/registry1.yaml
    u1, err := url.Parse("./data/set1/registry1.yaml")
    if err != nil {
        log.Fatal(err)
    }

    u2, err := url.Parse("https://raw.githubusercontent.com/go-orb/plugins/main/config/tests/data/set1/registry2.json")
    if err != nil {
        log.Fatal(err)
    }

    datas, err := config.Read([]*url.URL{u1, u2}, []string{"app"})
    if err != nil {
        log.Fatal(err)
    }
}
```

### config.Parse

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

### config.ParseStruct

ParseStruct is a helper to make any struct with `json` tags a source.Data (map[string]any{} with some more fields) with sections.

Example:

```go
func main() {
    cfg := log.NewConfig(log.WithLevel(log.LevelTrace), log.WithPlugin("slog"))

    data, err := config.ParseStruct([]string{"com", "example", "app", "registry", "logger"}, &cfg)
    if err != nil {
        return l, nil //nolint:nilerr
    }

    datas := []source.Data{data}
}
```

### config.HasKey

HasKey returns a boolean which indidcates if the given sections and key exists in the configs.

Example:

```go
func main() {
    test := config.HasKey([]string{"com", "example", "app", "registry", "logger"}, "plugin", configs)
}
```

## Authors

- [Asim Aslam](https://github.com/asim/) - Author of [go-micro/config](https://github.com/go-micro/go-micro/tree/master/config) on which this is based on.
- [David Brouwer](https://github.com/Davincible/) - Ideas and reviews.
- [René Jochum](https://github.com/jochumdev) - Developer.

## License

go-orb is Apache 2.0 licensed and is based on go-micro.
