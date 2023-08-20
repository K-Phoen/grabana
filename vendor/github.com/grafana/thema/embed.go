package thema

import (
	"embed"
)

// CueFS contains the raw .cue files that comprise the core thema system,
// making them available directly in Go.
//
// This virtual filesystem is relied on by the Go functions exported by this
// library, effectively co-versioning the Go logic with the CUE logic. It
// is exported such that other Go packages have the unfettered capability to
// create their own thema-based systems.
//
//go:embed *.cue crd/*.cue
var CueFS embed.FS

// CueJointFS contains the raw thema .cue files, as well as the cue.mod
// directory.
//
//go:embed *.cue crd/*.cue cue.mod
var CueJointFS embed.FS
