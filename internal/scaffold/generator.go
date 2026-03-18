package scaffold

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"spt-scaffold/internal/config"
)

// fileEntry describes a file to generate.
type fileEntry struct {
	// nameTemplate is a Go template for the file name (e.g. "{{.ModName}}.csproj").
	nameTemplate string
	bodyTemplate string
}

var entries = []fileEntry{
	{"{{.ModName}}.csproj", CsprojTemplate},
	{"Mod.cs", ModCSTemplate},
	{"README.md", ReadmeTemplate},
	{".gitignore", GitignoreTemplate},
}

// FileNames returns the rendered file names for the given config.
func FileNames(cfg config.ModConfig) []string {
	names := make([]string, len(entries))
	for i, e := range entries {
		t := template.Must(template.New("").Parse(e.nameTemplate))
		var buf bytes.Buffer
		_ = t.Execute(&buf, cfg)
		names[i] = buf.String()
	}
	return names
}

// GenerateFile generates the file at index idx inside the mod directory.
// Returns the file's base name on success.
func GenerateFile(cfg config.ModConfig, idx int) (string, error) {
	if idx >= len(entries) {
		return "", fmt.Errorf("invalid file index %d", idx)
	}
	e := entries[idx]

	// Render file name.
	nameTmpl := template.Must(template.New("name").Parse(e.nameTemplate))
	var nameBuf bytes.Buffer
	if err := nameTmpl.Execute(&nameBuf, cfg); err != nil {
		return "", fmt.Errorf("rendering file name: %w", err)
	}
	name := nameBuf.String()

	// Render file body.
	bodyTmpl, err := template.New("body").Parse(e.bodyTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing template for %s: %w", name, err)
	}
	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, cfg); err != nil {
		return "", fmt.Errorf("rendering template for %s: %w", name, err)
	}

	// Determine output directory.
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working directory: %w", err)
	}
	outDir := filepath.Join(cwd, cfg.ModName)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("creating directory %s: %w", outDir, err)
	}

	outPath := filepath.Join(outDir, name)
	if err := os.WriteFile(outPath, bodyBuf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("writing %s: %w", name, err)
	}

	return name, nil
}

// Generate generates all files and sends each name to the provided channel.
// It is used by the non-streaming code path.
func Generate(cfg config.ModConfig, ch chan<- string) error {
	for idx := range entries {
		name, err := GenerateFile(cfg, idx)
		if err != nil {
			return err
		}
		if ch != nil {
			ch <- name
		}
	}
	return nil
}
