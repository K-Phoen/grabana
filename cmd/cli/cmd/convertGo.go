package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/K-Phoen/grabana/encoder"
	"github.com/K-Phoen/sdk"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type convertGoOpts struct {
	inputJSON string
}

func ConvertGo(logger *zap.Logger) *cobra.Command {
	opts := convertGoOpts{}

	cmd := &cobra.Command{
		Use:   "convert-go",
		Short: "Converts a JSON dashboard to Golang",
		RunE: func(cmd *cobra.Command, args []string) error {
			return convertGo(logger, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.inputJSON, "input", "i", "", "JSON file used as input")
	_ = cmd.MarkFlagFilename("input", "json")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}

func convertGo(logger *zap.Logger, opts convertGoOpts) error {
	file, err := os.Open(opts.inputJSON)
	if err != nil {
		return fmt.Errorf("could not open input file '%s': %w", opts.inputJSON, err)
	}
	defer func() { _ = file.Close() }()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	dashboard := &sdk.Board{}
	if err := json.Unmarshal(content, dashboard); err != nil {
		return fmt.Errorf("could not unmarshall dashboard from JSON: %w", err)
	}

	golangDashboard, err := encoder.ToGolang(logger, *dashboard)
	if err != nil {
		return fmt.Errorf("could not encode input file '%s' to Go: %w", opts.inputJSON, err)
	}

	fmt.Println(golangDashboard)

	return nil
}
