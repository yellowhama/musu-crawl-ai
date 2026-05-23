package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yellowhama/musu-crawl-ai/internal/agent"
	"github.com/yellowhama/musu-crawl-ai/internal/processor"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Re-index the existing Wiki directory",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
		semantic, _ := cmd.Flags().GetBool("semantic")
		model, _ := cmd.Flags().GetString("model")

		proc := processor.NewWikiProcessor(out)

		fmt.Printf("Indexing directory: %s...\n", out)

		var err error
		if semantic {
			ollama := agent.NewOllamaClient(model)
			err = proc.UpdateIndexWithEmbedder(ollama.Embed)
		} else {
			err = proc.UpdateIndex()
		}

		if err != nil {
			fmt.Printf("Error during indexing: %v\n", err)
			return
		}
		fmt.Println("✅ Indexing completed (README, JSON, Bleve, and Vectors updated).")
	},
}

func init() {
	indexCmd.Flags().String("out", "./wiki", "Wiki directory to index")
	indexCmd.Flags().Bool("semantic", false, "Generate vector embeddings using Ollama")
	indexCmd.Flags().String("model", "nomic-embed-text", "Ollama model for embeddings")
	rootCmd.AddCommand(indexCmd)
}
