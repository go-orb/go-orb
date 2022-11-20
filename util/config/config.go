// Package config provides config utilities.
package config

import (
	"fmt"

	"go-micro.dev/v5/codecs"
)

var codec codecs.Marshaler //nolint:gochecknoglobals

// OverlayMap takes in a map[string]any and will overlay the file contents on
// a struct. It does this through the yaml/json structs. This function thus
// relies on both json/yaml struct tags, and that it's values are the same.
//
// This method is often used as dirty trick around https://github.com/golang/go/issues/48522
func OverlayMap(data map[string]any, target any) error {
	if target == nil {
		return nil
	}

	var err error
	if codec == nil {
		codec, err = codecs.GetCodec([]string{"yaml", "json"})
		if err != nil {
			return fmt.Errorf("parse entrypoint config: %w", err)
		}
	}

	b, err := codec.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := codec.Unmarshal(b, target); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}
