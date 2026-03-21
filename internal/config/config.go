package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ModConfig holds all user-provided configuration for scaffolding a mod.
type ModConfig struct {
	ModName         string `json:"mod_name"`
	Author          string `json:"author"`
	Version         string `json:"version"`
	SptVersion      string `json:"spt_version"`       // NuGet package version (e.g. "4.0.13")
	SptVersionRange string `json:"spt_version_range"` // SemanticVersioning range for Mod.cs (e.g. "~4.0.0")
	Desc            string `json:"description"`
	RepoURL         string `json:"repository_url"`
	License         string `json:"license"`
	ModType         string `json:"mod_type"`
	ModTemplate     string `json:"mod_template"`
	// SptInstallPath is the absolute path to the SPT installation directory (client mods only).
	SptInstallPath string `json:"spt_install_path"`
	// ProjectGuid is the generated UUID for the .sln file.
	ProjectGuid string `json:"project_guid"`
}

var semverRe = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
var httpsRe = regexp.MustCompile(`^https://`)

// ValidateModName ensures the name is non-empty and has no spaces.
func ValidateModName(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("mod name is required")
	}
	if strings.ContainsAny(s, " \t") {
		return fmt.Errorf("mod name must not contain spaces")
	}
	return nil
}

// ValidateAuthor ensures the author is non-empty.
func ValidateAuthor(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("author is required")
	}
	return nil
}

// ValidateSemver ensures the version matches x.y.z format.
func ValidateSemver(s string) error {
	if !semverRe.MatchString(strings.TrimSpace(s)) {
		return fmt.Errorf("version must be in x.y.z format (e.g. 1.0.0)")
	}
	return nil
}

// ValidateSptVersion ensures the SPT version matches x.y.z format.
func ValidateSptVersion(s string) error {
	if !semverRe.MatchString(strings.TrimSpace(s)) {
		return fmt.Errorf("SPT version must be in x.y.z format (e.g. 4.0.3)")
	}
	return nil
}

// ValidateDescription ensures the description is at most 120 chars.
func ValidateDescription(s string) error {
	if len(s) > 120 {
		return fmt.Errorf("description must be 120 characters or fewer (%d/120)", len(s))
	}
	return nil
}

// ValidateRepoURL ensures the URL starts with https:// if not empty.
func ValidateRepoURL(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !httpsRe.MatchString(s) {
		return fmt.Errorf("repository URL must start with https://")
	}
	return nil
}

// ValidateSptInstallPath checks the path is non-empty and contains BepInEx/core/BepInEx.dll.
func ValidateSptInstallPath(p string) error {
	p = strings.TrimSpace(p)
	if p == "" {
		return fmt.Errorf("SPT install path is required")
	}
	check := filepath.Join(p, "BepInEx", "core", "BepInEx.dll")
	if _, err := os.Stat(check); err != nil {
		return fmt.Errorf("not a valid SPT install directory (BepInEx/core/BepInEx.dll not found)")
	}
	return nil
}

// LicenseEntry holds the display label and SPDX identifier for a license.
type LicenseEntry struct {
	Label string
	SPDX  string
}

// Licenses is the ordered list of available license choices.
var Licenses = []LicenseEntry{
	{"Apache License 2.0", "Apache-2.0"},
	{"Boost Software License 1.0", "BSL-1.0"},
	{"Creative Commons BY 3.0", "CC-BY-3.0"},
	{"Creative Commons BY-NC 3.0", "CC-BY-NC-3.0"},
	{"Creative Commons BY-NC-ND 3.0", "CC-BY-NC-ND-3.0"},
	{"Creative Commons BY-NC-ND 4.0", "CC-BY-NC-ND-4.0"},
	{"Creative Commons BY-NC-SA 3.0", "CC-BY-NC-SA-3.0"},
	{"Creative Commons BY-NC-SA 4.0", "CC-BY-NC-SA-4.0"},
	{"Creative Commons BY-ND 3.0", "CC-BY-ND-3.0"},
	{"Creative Commons BY-SA 3.0", "CC-BY-SA-3.0"},
	{"GNU AGPLv3", "AGPL-3.0"},
	{"GNU GPLv3", "GPL-3.0"},
	{"GNU LGPLv3", "LGPL-3.0"},
	{"MIT License", "MIT"},
	{"Mozilla Public License 2.0", "MPL-2.0"},
	{"The Unlicense", "Unlicense"},
	{"University of Illinois/NCSA Open Source License", "NCSA"},
}

// ModTypeServer and ModTypeClient are the canonical Value strings for mod types.
const (
	ModTypeServer = "server"
	ModTypeClient = "client"
)

// ModTypes lists the available mod types. Disabled items are shown but not selectable.
var ModTypes = []struct {
	Label    string
	Value    string
	Disabled bool
}{
	{"Server Mod", "server", false},
	{"Client Mod", "client", false},
}

// TemplateEntry describes a scaffold template option.
type TemplateEntry struct {
	Label string `json:"Label"`
	Value string `json:"-"` // We will set this from the directory name
	Desc  string `json:"Desc"`
}

// ServerTemplates lists the available server mod templates.
var ServerTemplates []TemplateEntry

// ClientTemplates lists the available client mod templates.
var ClientTemplates []TemplateEntry

// LoadTemplates given an fs.FS and a base directory ("server" or "client")
// reads the available templates.
func LoadTemplates(fsys interface{ ReadDir(string) ([]os.DirEntry, error); ReadFile(string) ([]byte, error) }, base string) ([]TemplateEntry, error) {
	entries, err := fsys.ReadDir(base)
	if err != nil {
		return nil, fmt.Errorf("reading template directory %s: %w", base, err)
	}

	var templates []TemplateEntry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		
		val := e.Name()
		jsonPath := filepath.ToSlash(filepath.Join(base, val, "template.json"))
		b, err := fsys.ReadFile(jsonPath)
		if err != nil {
			// skip if no template.json
			continue
		}

		var te TemplateEntry
		if err := json.Unmarshal(b, &te); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", jsonPath, err)
		}
		te.Value = val
		templates = append(templates, te)
	}
	
	// Ensure "empty" is always first if it exists
	sort.SliceStable(templates, func(i, j int) bool {
		return templates[i].Value == "empty"
	})
	
	return templates, nil
}

