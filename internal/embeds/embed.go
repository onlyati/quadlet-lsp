// Packkage embdes
//
// Package to store embedded files for templates
package embeds

import "embed"

//go:embed *.tpl
var TemplateFs embed.FS
