package scaffold

import (
	"embed"
	"io/fs"
)

//go:embed all:templates
var templatesRoot embed.FS

// Templates returns the embedded template tree rooted at "templates/".
// Consumers pass a kind (e.g. "cli-go") + name to Render along with this FS.
func Templates() fs.FS {
	sub, err := fs.Sub(templatesRoot, "templates")
	if err != nil {
		panic(err)
	}
	return sub
}
