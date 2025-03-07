package cli

// ParserFunc is a provider function type for parsing flags.
type ParserFunc func(appContext *AppContext, args []string) ([]*Flag, error)
