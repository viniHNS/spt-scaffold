package scaffold

import (
	"testing"

	"spt-scaffold/internal/config"
)

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func TestFileNames_ServerMod(t *testing.T) {
	cfg := config.ModConfig{ModType: "server", ModName: "MyMod"}
	names := FileNames(cfg)
	if len(names) != 5 {
		t.Fatalf("expected 5 files for server mod, got %d", len(names))
	}
	if !contains(names, "MyMod.sln") {
		t.Errorf("missing MyMod.sln")
	}
	if !contains(names, "MyMod.csproj") {
		t.Errorf("missing MyMod.csproj")
	}
	if !contains(names, "Mod.cs") {
		t.Errorf("missing Mod.cs")
	}
}

func TestFileNames_ClientMod(t *testing.T) {
	cfg := config.ModConfig{ModType: "client", ModName: "MyMod"}
	names := FileNames(cfg)
	if len(names) != 5 {
		t.Fatalf("expected 5 files for client mod, got %d", len(names))
	}
	if !contains(names, "MyMod.sln") {
		t.Errorf("missing MyMod.sln")
	}
	if !contains(names, "MyMod.csproj") {
		t.Errorf("missing MyMod.csproj")
	}
	if !contains(names, "Plugin.cs") {
		t.Errorf("missing Plugin.cs")
	}
}

func TestFileNames_UnknownModTypeFallsBackToServer(t *testing.T) {
	cfg := config.ModConfig{ModType: "unknown", ModName: "MyMod"}
	names := FileNames(cfg)
	if len(names) < 3 {
		t.Fatalf("expected at least 3 files for fallback, got %d", len(names))
	}
	if !contains(names, "MyMod.sln") {
		t.Errorf("expected server fallback (MyMod.sln) to be present")
	}
	if !contains(names, "Mod.cs") {
		t.Errorf("expected server fallback (Mod.cs) to be present")
	}
}

func TestFileNames_ClientHarmonyPatch(t *testing.T) {
	cfg := config.ModConfig{ModType: "client", ModName: "MyMod", ModTemplate: "harmonyPatch"}
	names := FileNames(cfg)
	if len(names) != 6 {
		t.Fatalf("expected 6 files for client harmonyPatch, got %d: %v", len(names), names)
	}
	if !contains(names, "MyMod.sln") {
		t.Errorf("missing MyMod.sln")
	}
	if !contains(names, "Plugin.cs") {
		t.Errorf("missing Plugin.cs")
	}
	if !contains(names, "Patches/ExamplePatch.cs") {
		t.Errorf("missing Patches/ExamplePatch.cs")
	}
}

