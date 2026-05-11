// Package scaffold renders an embedded fs.FS of templates into a destination
// directory, substituting {{.Name}}/{{.KebabName}} placeholders in both file
// contents and path segments.
package scaffold

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Vars are the variables exposed to templates.
type Vars struct {
	Name      string
	KebabName string
}

// Render copies every file under `root` of the templates FS into destDir,
// applying text/template substitution on contents (when the filename ends
// with .tmpl) and on path segments (always).
//
// Returns an error if the root subtree does not exist or if any write fails.
func Render(templates fs.FS, root string, destDir string, vars Vars) error {
	// Verify the root subtree exists.
	if _, err := fs.Stat(templates, root); err != nil {
		return fmt.Errorf("template tree %q: %w", root, err)
	}
	return fs.WalkDir(templates, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Compute the destination path relative to root.
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// Apply template substitution to each path segment.
		segments := strings.Split(filepath.ToSlash(rel), "/")
		for i, seg := range segments {
			rendered, err := renderString(seg, vars)
			if err != nil {
				return fmt.Errorf("render path segment %q: %w", seg, err)
			}
			segments[i] = rendered
		}
		rel = filepath.Join(segments...)

		// Strip the .tmpl suffix if present.
		isTemplate := strings.HasSuffix(rel, ".tmpl")
		if isTemplate {
			rel = strings.TrimSuffix(rel, ".tmpl")
		}

		destPath := filepath.Join(destDir, rel)
		if err := os.MkdirAll(filepath.Dir(destPath), 0o750); err != nil {
			return err
		}

		// Read the source content from the embedded FS.
		src, err := fs.ReadFile(templates, path)
		if err != nil {
			return err
		}

		var out []byte
		if isTemplate {
			rendered, err := renderBytes(src, vars)
			if err != nil {
				return fmt.Errorf("render %q: %w", path, err)
			}
			out = rendered
		} else {
			out = src
		}
		return os.WriteFile(destPath, out, 0o600)
	})
}

func renderString(s string, vars Vars) (string, error) {
	tpl, err := template.New("seg").Option("missingkey=error").Parse(s)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	if err := tpl.Execute(&b, vars); err != nil {
		return "", err
	}
	return b.String(), nil
}

func renderBytes(src []byte, vars Vars) ([]byte, error) {
	tpl, err := template.New("file").Option("missingkey=error").Parse(string(src))
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	if err := tpl.Execute(&b, vars); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
