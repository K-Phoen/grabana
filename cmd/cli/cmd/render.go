package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/K-Phoen/grabana/decoder"
	"io/ioutil"

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
	content, err := ioutil.ReadFile(opts.inputYAML)
	if err != nil {
		return fmt.Errorf("could not read input file '%s': %w", opts.inputYAML, err)
	}

	dashboard, err := decoder.UnmarshalYAML(bytes.NewBuffer(content))
	if err != nil {
		return fmt.Errorf("could not decode input file '%s': %w", opts.inputYAML, err)
	}

	buf, err := json.MarshalIndent(dashboard.Internal(), "", "  ")
	if err != nil {
		return err
	}

	fmt.Print(string(buf))

	return nil
}
