module go-micro.dev/v5

go 1.19

require (
	github.com/go-micro/plugins/config/source/file v0.0.0-00010101000000-000000000000
	github.com/google/wire v0.5.0
	github.com/stretchr/testify v1.8.1
	golang.org/x/exp v0.0.0-20221109205753-fc8884afc316
)

require github.com/go-micro/plugins/codecs/yaml v0.0.0-00010101000000-000000000000 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/tools v0.3.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-micro/plugins/config/source/file => ../plugins/config/source/file

replace github.com/go-micro/plugins/codecs/yaml => ../plugins/codecs/yaml
