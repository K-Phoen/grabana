package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/decoder"
	"github.com/spf13/cobra"
)

type applyOpts struct {
	inputYAML         string
	destinationFolder string
	grafanaHost       string
	grafanaToken      string
}

func Apply() *cobra.Command {
	opts := applyOpts{}

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply a YAML dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return applyYAML(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.inputYAML, "input", "i", "", "YAML file used as input")
	cmd.Flags().StringVarP(&opts.destinationFolder, "folder", "f", "", "Folder in which the dashboard will be created")
	cmd.Flags().StringVarP(&opts.grafanaHost, "grafana", "g", "", "Grafana host. Example: http://grafana-host:3000")
	cmd.Flags().StringVarP(&opts.grafanaToken, "token", "t", "", "Grafana API token")

	_ = cmd.MarkFlagFilename("input", "yaml", "yml")

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("folder")
	_ = cmd.MarkFlagRequired("grafana")

	return cmd
}

func applyYAML(opts applyOpts) error {
	ctx := context.Background()
	client := grabanaClient(opts)

	file, err := os.Open(opts.inputYAML)
	if err != nil {
		return fmt.Errorf("could not open input file '%s': %w", opts.inputYAML, err)
	}

	dashboard, err := decoder.UnmarshalYAML(file)
	if err != nil {
		return fmt.Errorf("could not decode input file '%s': %w", opts.inputYAML, err)
	}

	folder, err := client.FindOrCreateFolder(ctx, opts.destinationFolder)
	if err != nil {
		return fmt.Errorf("could not find or create folder '%s': %w", opts.destinationFolder, err)
	}

	if _, err := client.UpsertDashboard(ctx, folder, dashboard); err != nil {
		return fmt.Errorf("could not apply dashboard: %w", err)
	}

	return nil
}

func grabanaClient(opts applyOpts) *grabana.Client {
	var clientOpts []grabana.Option
	if len(opts.grafanaToken) != 0 {
		clientOpts = append(clientOpts, grabana.WithAPIToken(opts.grafanaToken))
	}

	return grabana.NewClient(&http.Client{}, opts.grafanaHost, clientOpts...)
}
