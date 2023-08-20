package util

import (
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue/load"
)

// ToOverlay converts an fs.FS into a CUE loader overlay.
func ToOverlay(prefix string, vfs fs.FS, overlay map[string]load.Source) error {
	// TODO why not just stick the prefix on automatically...?
	if !filepath.IsAbs(prefix) {
		return fmt.Errorf("must provide absolute path prefix when generating cue overlay, got %q", prefix)
	}
	err := fs.WalkDir(vfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := vfs.Open(path)
		if err != nil {
			return err
		}
		defer f.Close() // nolint: errcheck

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		overlay[filepath.Join(prefix, path)] = load.FromBytes(b)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq produces random (basic, not crypto) letters of a given length.
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// TODO make a better list - rely on CUE Go API somehow? Tokens?
func mustQuote(n string) bool {
	quoteneed := []string{
		"string",
		"number",
		"int",
		"uint",
		"float",
		"byte",
	}

	for _, s := range quoteneed {
		if n == s {
			return true
		}
	}
	return false
}

// SanitizeLabelString strips characters from a string that are not allowed for
// use in a CUE label.
func SanitizeLabelString(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			fallthrough
		case r >= 'A' && r <= 'Z':
			fallthrough
		case r >= '0' && r <= '9':
			fallthrough
		case r == '_':
			return r
		default:
			return -1
		}
	}, s)
}
