package cmd

import (
	"fmt"
	"os"

	"github.com/K-Phoen/grabana/decoder"
	"github.com/spf13/cobra"
)

type renderOpts struct {
	inputYAML string
}

func Render() *cobra.Command {
	opts := renderOpts{}

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render a YAML dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return renderYAML(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.inputYAML, "input", "i", "", "YAML file used as input")

	_ = cmd.MarkFlagFilename("input", "yaml", "yml")

	_ = cmd.MarkFlagRequired("input")

	return cmd
}

func renderYAML(opts renderOpts) error {
	file, err := os.Open(opts.inputYAML)
	if err != nil {
		return fmt.Errorf("could not open input file '%s': %w", opts.inputYAML, err)
	}

	dashboard, err := decoder.UnmarshalYAML(file)
	if err != nil {
		return fmt.Errorf("could not decode input file '%s': %w", opts.inputYAML, err)
	}

	buf, err := dashboard.MarshalIndentJSON()
	if err != nil {
		return err
	}

	fmt.Println(string(buf))

	return nil
}
