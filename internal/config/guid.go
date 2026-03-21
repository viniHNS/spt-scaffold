package config

import (
	"crypto/rand"
	"fmt"
)

// NewProjectGuid generates a random UUID v4 in the format {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}.
func NewProjectGuid() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	// Set version (4) and variant (RFC 4122).
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("{%08X-%04X-%04X-%04X-%012X}",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
