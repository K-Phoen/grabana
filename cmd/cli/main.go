package main

import (
	"os"

	"github.com/K-Phoen/grabana/cmd/cli/cmd"
	"github.com/spf13/cobra"
)

var version = "SNAPSHOT"

func main() {
	root := &cobra.Command{Use: "grabana"}
	root.Version = version

	root.AddCommand(cmd.Apply())
	root.AddCommand(cmd.Validate())
	root.AddCommand(cmd.SelfUpdate(version))
	root.AddCommand(cmd.Render())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
