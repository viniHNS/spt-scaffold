package config

import (
	"fmt"
	"regexp"
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
