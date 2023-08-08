package packages

import (
	"encoding/json"
	"io"
	"os"
)

func Load(input io.Reader) (NormalizedPackage, error) {
	normalized := NormalizedPackage{}
	if err := json.NewDecoder(input).Decode(&normalized); err != nil {
		return normalized, err
	}

	return normalized, nil
}

func LoadFile(filepath string) (NormalizedPackage, error) {
	pkgHandle, err := os.Open(filepath)
	if err != nil {
		return NormalizedPackage{}, err
	}

	return Load(pkgHandle)
}
