module github.com/go-orb/go-orb

go 1.20

require (
	github.com/go-orb/plugins/codecs/yaml v0.0.0-00010101000000-000000000000
	github.com/go-orb/plugins/config/source/file v0.0.1-00010101000000-000000000000
	github.com/google/wire v0.5.0
	github.com/sanity-io/litter v1.5.5
	github.com/stretchr/testify v1.8.3
	golang.org/x/exp v0.0.0-20230626212559-97b1e661b5df
)

require github.com/kr/pretty v0.3.1 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/tools v0.11.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/go-orb/plugins/config/source/file => ../plugins/config/source/file

replace github.com/go-orb/plugins/codecs/yaml => ../plugins/codecs/yaml
