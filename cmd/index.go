package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yellowhama/musu-crawl-ai/internal/processor"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Re-index the existing Wiki directory",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
		proc := processor.NewWikiProcessor(out)

		fmt.Printf("Indexing directory: %s...\n", out)
		err := proc.UpdateIndex()
		if err != nil {
			fmt.Printf("Error during indexing: %v\n", err)
			return
		}
		fmt.Println("✅ Indexing completed (README.md and index.json updated).")
	},
}

func init() {
	indexCmd.Flags().String("out", "./wiki", "Wiki directory to index")
	rootCmd.AddCommand(indexCmd)
}
