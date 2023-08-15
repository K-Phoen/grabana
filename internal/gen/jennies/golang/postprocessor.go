package golang

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"golang.org/x/tools/imports"
)

func PostProcessFile(file codejen.File) (codejen.File, error) {
	output, err := imports.Process(filepath.Base(file.RelativePath), file.Data, nil)
	if err != nil {
		return codejen.File{}, fmt.Errorf("goimports processing of generated file failed: %w", err)
	}

	return codejen.File{
		RelativePath: file.RelativePath,
		Data:         output,
		From:         file.From,
	}, nil
}
