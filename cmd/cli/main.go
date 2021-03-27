package main

import (
	"github.com/K-Phoen/grabana/cmd/cli/cmd"
	"github.com/spf13/cobra"
)

var version = "SNAPSHOT"

func main() {
	root := &cobra.Command{Use: "grabana"}
	root.Version = version
	root.SetVersionTemplate(version)

	root.AddCommand(cmd.Apply())
	root.AddCommand(cmd.Validate())

	_ = root.Execute()
}
