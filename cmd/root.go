package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "musu-crawl-ai",
	Short: "AI-Ready Knowledge Harvester & Wiki Generator",
	Long:  `A high-performance crawler to fetch and organize content for AI.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Root flags if needed
}
