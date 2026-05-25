package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const Version = "v0.7.1"

var rootCmd = &cobra.Command{
	Use:     "musu-crawl-ai",
	Short:   "AI-Ready Knowledge Harvester & Wiki Generator",
	Long:    `A high-performance crawler to fetch and organize content for AI.`,
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Silent background check for updates
		if cmd.Name() != "update" && cmd.Name() != "help" {
			go checkNewVersion()
		}
	},
}

func checkNewVersion() {
	latest, _, err := GetLatestRelease("yellowhama", "musu-crawl-ai")
	if err == nil && latest != Version {
		fmt.Fprintf(os.Stderr, "\n💡 New version available: %s (Current: %s)\n", latest, Version)
		fmt.Fprintf(os.Stderr, "👉 Run 'musu-crawl update' to upgrade.\n\n")
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Root flags if needed
}
