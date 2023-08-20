//go:build mage

package main

import (
	"golang.org/x/tools/txtar"
	"os"
	"path/filepath"
	"strings"
)

// Clean cleans all the output from corpus txtar files
func Clean() error {
	return filepath.Walk("./testdata", cleanTxtarOut)
}

func cleanTxtarOut(path string, _ os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if filepath.Ext(path) != ".txtar" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	archive := txtar.Parse(data)
	archive.Files = filter(
		archive.Files,
		func(t txtar.File) bool {
			return !strings.HasPrefix(t.Name, "out")
		},
	)

	return os.WriteFile(path, txtar.Format(archive), 0666)

}

func filter(ff []txtar.File, criteria func(f txtar.File) bool) []txtar.File {
	filtered := make([]txtar.File, 0, len(ff))
	for _, f := range ff {
		if criteria(f) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}
