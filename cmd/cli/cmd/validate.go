package cmd

import (
	"fmt"
	"os"

	"github.com/K-Phoen/grabana/decoder"
	"github.com/spf13/cobra"
)

type validateOpts struct {
	inputYAML string
}

func Validate() *cobra.Command {
	opts := validateOpts{}

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a YAML dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validateYAML(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.inputYAML, "input", "i", "", "YAML file used as input")

	_ = cmd.MarkFlagFilename("input", "yaml", "yml")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}

func validateYAML(opts validateOpts) error {
	file, err := os.Open(opts.inputYAML)
	if err != nil {
		return fmt.Errorf("could not open input file '%s': %w", opts.inputYAML, err)
	}

	if _, err := decoder.UnmarshalYAML(file); err != nil {
		return fmt.Errorf("could not decode input file '%s': %w", opts.inputYAML, err)
	}

	return nil
}
