package scaffold

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"spt-scaffold/internal/config"
	"spt-scaffold/internal/templates"
)

// getTemplateDir returns the path in the embed.FS for the specific mod type and template.
func getTemplateDir(cfg config.ModConfig) string {
	base := cfg.ModType
	if base != config.ModTypeServer && base != config.ModTypeClient {
		base = config.ModTypeServer
	}
	tmpl := cfg.ModTemplate
	if tmpl == "" {
		tmpl = "empty"
	}
	// e.g. "server/empty"
	return filepath.ToSlash(filepath.Join(base, tmpl))
}

// collectTemplateFiles walks the template directory recursively and returns
// relative paths to all .tmpl files (excluding template.json).
func collectTemplateFiles(dir string) ([]string, error) {
	var files []string
	err := fs.WalkDir(templates.FS, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "template.json" || !strings.HasSuffix(d.Name(), ".tmpl") {
			return nil
		}
		// Store relative path from the template root dir
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	return files, err
}

// renderName applies template variables to a file name string.
func renderName(nameTemplate string, cfg config.ModConfig) string {
	t := template.Must(template.New("").Parse(nameTemplate))
	var buf bytes.Buffer
	_ = t.Execute(&buf, cfg)
	return buf.String()
}

// FileNames returns the rendered file names (with relative paths) for the given config.
func FileNames(cfg config.ModConfig) []string {
	dir := getTemplateDir(cfg)
	tmplFiles, err := collectTemplateFiles(dir)
	if err != nil {
		return nil
	}

	var names []string
	for _, relPath := range tmplFiles {
		nameTemplate := strings.TrimSuffix(relPath, ".tmpl")
		names = append(names, renderName(nameTemplate, cfg))
	}
	return names
}

// GenerateFile generates the file at index idx inside the mod directory.
// Returns the file's relative path on success.
func GenerateFile(cfg config.ModConfig, idx int) (string, error) {
	dir := getTemplateDir(cfg)
	tmplFiles, err := collectTemplateFiles(dir)
	if err != nil {
		return "", fmt.Errorf("reading template directory: %w", err)
	}

	if idx >= len(tmplFiles) || idx < 0 {
		return "", fmt.Errorf("invalid file index %d", idx)
	}

	relPath := tmplFiles[idx]
	nameTemplate := strings.TrimSuffix(relPath, ".tmpl")

	// Read file contents from embed.FS
	fullPath := filepath.ToSlash(filepath.Join(dir, relPath))
	bodyTemplateBytes, err := templates.FS.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("reading template file: %w", err)
	}
	bodyTemplate := string(bodyTemplateBytes)

	// Render file name (including any subdirectory in the path).
	name := renderName(nameTemplate, cfg)

	// Render file body.
	bodyTmpl, err := template.New("body").Parse(bodyTemplate)
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

	outPath := filepath.Join(cwd, cfg.ModName, name)
	outDir := filepath.Dir(outPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("creating directory %s: %w", outDir, err)
	}

	if err := os.WriteFile(outPath, bodyBuf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("writing %s: %w", name, err)
	}

	return name, nil
}

// Generate generates all files and sends each name to the provided channel.
func Generate(cfg config.ModConfig, ch chan<- string) error {
	names := FileNames(cfg)
	for i := range names {
		name, err := GenerateFile(cfg, i)
		if err != nil {
			return err
		}
		if ch != nil {
			ch <- name
		}
	}
	return nil
}

