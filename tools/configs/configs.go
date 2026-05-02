// Package configs embeds static configuration files bundled with the piltover
// engine binary.
package configs

import _ "embed"

// DefaultsYAML is the content of defaults.yaml, which holds the default
// lint/test/build commands for each supported language.
//
//go:embed defaults.yaml
var DefaultsYAML []byte
