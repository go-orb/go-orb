module github.com/go-orb/go-orb

go 1.23

toolchain go1.23.0

require (
	github.com/cornelk/hashmap v1.0.8
	github.com/google/wire v0.6.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/sanity-io/litter v1.5.5
	github.com/stretchr/testify v1.9.0
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/tools v0.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// Fixing ambiguous import: found package google.golang.org/genproto/googleapis/api/annotations in multiple modules.
replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240924160255-9d4c2d233b61
