package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateSptInstallPath_EmptyPath(t *testing.T) {
	err := ValidateSptInstallPath("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestValidateSptInstallPath_NonExistentPath(t *testing.T) {
	// Use a path guaranteed not to exist (subdirectory of a fresh temp dir)
	dir := t.TempDir()
	nonExistent := filepath.Join(dir, "does-not-exist")
	err := ValidateSptInstallPath(nonExistent)
	if err == nil {
		t.Fatal("expected error for non-existent path, got nil")
	}
}

func TestValidateSptInstallPath_ValidSptDir(t *testing.T) {
	// Create a temp dir that looks like an SPT install
	dir := t.TempDir()
	coreDir := filepath.Join(dir, "BepInEx", "core")
	if err := os.MkdirAll(coreDir, 0755); err != nil {
		t.Fatal(err)
	}
	bepInExDll := filepath.Join(coreDir, "BepInEx.dll")
	if err := os.WriteFile(bepInExDll, []byte("stub"), 0644); err != nil {
		t.Fatal(err)
	}

	err := ValidateSptInstallPath(dir)
	if err != nil {
		t.Fatalf("expected nil for valid SPT dir, got: %v", err)
	}
}

func TestValidateSptInstallPath_DirExistsButMissingBepInEx(t *testing.T) {
	dir := t.TempDir() // exists but has no BepInEx structure
	err := ValidateSptInstallPath(dir)
	if err == nil {
		t.Fatal("expected error for dir without BepInEx.dll, got nil")
	}
}
