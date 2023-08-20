package kindsys

import (
	"embed"
)

// CueSchemaFS embeds all CUE files in the Kindsys project.
//
//go:embed cue.mod/module.cue *.cue
var CueSchemaFS embed.FS
