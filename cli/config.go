package cli

type Config struct {
	Name        string
	Version     string
	Description string
	Usage       string
	NoFlags     *bool
	ArgPrefix   string
	Plugin      string
	ConfigFile  string

	// No config option but here as helper
	Flags []Flag
}

func NewConfig() *Config {
	return &Config{}
}

func (d *Config) Merge(src *Config) error {
	def := NewConfig()

	if src.NoFlags != def.NoFlags {
		d.NoFlags = src.NoFlags
	}
	if src.ArgPrefix != def.ArgPrefix {
		d.ArgPrefix = src.ArgPrefix
	}
	if src.Plugin != def.Plugin {
		d.Plugin = src.Plugin
	}
	if src.ConfigFile != def.ConfigFile {
		d.ConfigFile = src.ConfigFile
	}

	return nil
}
