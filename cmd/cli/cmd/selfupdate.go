package cmd

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func SelfUpdate(currentVersion string) *cobra.Command {
	return &cobra.Command{
		Use:   "self-update",
		Short: "Self-update",
		RunE: func(cmd *cobra.Command, args []string) error {
			return selfUpdate(currentVersion)
		},
	}
}

func selfUpdate(currentVersion string) error {
	current, err := semver.ParseTolerant(currentVersion)
	if err != nil {
		return fmt.Errorf("could not parse current version: %w", err)
	}

	latest, err := selfupdate.UpdateSelf(current, "K-Phoen/grabana")
	if err != nil {
		return fmt.Errorf("binary update failed: %w", err)
	}

	if latest.Version.Equals(current) {
		fmt.Printf("Current binary is the latest version: %s\n", currentVersion)
	} else {
		fmt.Printf("Successfully updated to version %s\n", latest.Version)
	}

	return nil
}
