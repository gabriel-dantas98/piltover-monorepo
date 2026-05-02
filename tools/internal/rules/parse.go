// Package rules provides parsing and validation for Kody Custom Rules.
package rules

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rule is a parsed Kody Custom Rule.
type Rule struct {
	Title       string   `yaml:"title"        json:"title"`
	Scope       string   `yaml:"scope"        json:"scope"`
	Path        []string `yaml:"path"         json:"path"`
	SeverityMin string   `yaml:"severity_min" json:"severity_min"`
	Languages   []string `yaml:"languages"    json:"languages"`
	Buckets     []string `yaml:"buckets"      json:"buckets"`
	Enabled     bool     `yaml:"enabled"      json:"enabled"`

	// Slug is set by the caller from the filename (without .md).
	Slug string `yaml:"-" json:"slug"`
	// Body is the markdown body after the frontmatter.
	Body string `yaml:"-" json:"-"`
}

// Parse splits YAML frontmatter from markdown body and returns a populated Rule.
// The caller fills Slug afterwards.
func Parse(data []byte) (*Rule, error) {
	const fence = "---"
	rest := strings.TrimLeft(string(data), " \t\n\r")
	if !strings.HasPrefix(rest, fence+"\n") && !strings.HasPrefix(rest, fence+"\r\n") {
		return nil, fmt.Errorf("rules.Parse: missing frontmatter fence")
	}
	rest = strings.TrimPrefix(rest, fence+"\n")
	rest = strings.TrimPrefix(rest, fence+"\r\n")
	idx := strings.Index(rest, "\n"+fence)
	if idx < 0 {
		return nil, fmt.Errorf("rules.Parse: unterminated frontmatter")
	}
	front := rest[:idx]
	body := strings.TrimLeft(rest[idx+len("\n"+fence):], "\n\r")

	var r Rule
	if err := yaml.Unmarshal([]byte(front), &r); err != nil {
		return nil, fmt.Errorf("rules.Parse: yaml: %w", err)
	}
	r.Body = body
	return &r, nil
}

// Validate returns an error if the rule is missing required fields.
func (r *Rule) Validate() error {
	var missing []string
	if r.Title == "" {
		missing = append(missing, "title")
	}
	if r.Scope == "" {
		missing = append(missing, "scope")
	}
	if len(r.Path) == 0 {
		missing = append(missing, "path")
	}
	if r.SeverityMin == "" {
		missing = append(missing, "severity_min")
	}
	if len(r.Languages) == 0 {
		missing = append(missing, "languages")
	}
	if len(r.Buckets) == 0 {
		missing = append(missing, "buckets")
	}
	if len(missing) > 0 {
		return fmt.Errorf("rule %q missing fields: %s", r.Slug, strings.Join(missing, ", "))
	}
	return nil
}
