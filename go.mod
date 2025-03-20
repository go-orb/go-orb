module github.com/go-orb/go-orb

go 1.23.0

require (
	dario.cat/mergo v1.0.1
	github.com/cornelk/hashmap v1.0.8
	github.com/go-orb/wire v0.7.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/lithammer/shortuuid/v3 v3.0.7
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
)

// Fixing ambiguous import: found package google.golang.org/genproto/googleapis/api/annotations in multiple modules.
replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240924160255-9d4c2d233b61
