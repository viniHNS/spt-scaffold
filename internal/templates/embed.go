package templates

import (
	"embed"
	"spt-scaffold/internal/config"
)

//go:embed all:server all:client
var FS embed.FS

func init() {
	var err error
	config.ServerTemplates, err = config.LoadTemplates(FS, "server")
	if err != nil {
		panic("failed to load server templates: " + err.Error())
	}
	config.ClientTemplates, err = config.LoadTemplates(FS, "client")
	if err != nil {
		panic("failed to load client templates: " + err.Error())
	}
}
