package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// FallbackSptVersion is used when the NuGet API is unreachable.
const FallbackSptVersion = "4.0.0"

type nugetIndex struct {
	Versions []string `json:"versions"`
}

// FetchSptVersions queries NuGet for all available SPTarkov.Server.Core versions.
// Returns versions in descending order (latest first).
func FetchSptVersions() ([]string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.nuget.org/v3-flatcontainer/sptarkov.server.core/index.json")
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	var idx nugetIndex
	if err := json.NewDecoder(resp.Body).Decode(&idx); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	if len(idx.Versions) == 0 {
		return nil, fmt.Errorf("no versions found")
	}

	// Keep only stable releases (x.y.z) and reverse to descending order (latest first).
	stable := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	var stables []string
	for _, v := range idx.Versions {
		if stable.MatchString(v) {
			stables = append(stables, v)
		}
	}
	if len(stables) == 0 {
		return nil, fmt.Errorf("no stable versions found")
	}
	for i, j := 0, len(stables)-1; i < j; i, j = i+1, j-1 {
		stables[i], stables[j] = stables[j], stables[i]
	}
	return stables, nil
}

// SptVersionRange returns the semver compatibility range for a given SPT version.
// E.g. "4.0.13" → "~4.0.0"
func SptVersionRange(version string) string {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return "~" + version
	}
	return "~" + parts[0] + "." + parts[1] + ".0"
}
