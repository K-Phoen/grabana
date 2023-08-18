package jennies

import (
	"path/filepath"

	"github.com/grafana/codejen"
)

// Prefixer returns a FileMapper that injects the provided path prefix to files
// passed through it.
func Prefixer(prefix string) codejen.FileMapper {
	return func(f codejen.File) (codejen.File, error) {
		f.RelativePath = filepath.Join(prefix, f.RelativePath)
		return f, nil
	}
}
