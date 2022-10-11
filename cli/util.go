package cli

import (
	"fmt"
	"strings"
)

func PrefixName(prefix, name string) string {
	if len(prefix) > 0 {
		return fmt.Sprintf("%s_%s", strings.ToLower(prefix), strings.ToLower(name))
	}

	return strings.ToLower(name)
}

func PrefixEnv(prefix, name string) string {
	if len(prefix) > 0 {
		return fmt.Sprintf("%s_%s", strings.ToUpper(prefix), strings.ToUpper(name))
	}

	return strings.ToUpper(name)
}
