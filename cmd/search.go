package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/yellowhama/musu-crawl-ai/internal/agent"
	"github.com/yellowhama/musu-crawl-ai/internal/processor"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the Wiki knowledge base locally",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
		semantic, _ := cmd.Flags().GetBool("semantic")
		model, _ := cmd.Flags().GetString("model")
		queryStr := args[0]

		if semantic {
			runSemanticSearch(out, queryStr, model)
		} else {
			runBleveSearch(out, queryStr)
		}
	},
}

func runBleveSearch(out, queryStr string) {
	blevePath := filepath.Join(out, "musu.bleve")
	if _, err := os.Stat(blevePath); os.IsNotExist(err) {
		fmt.Println("❌ Search index not found. Please run 'musu-crawl index' first.")
		return
	}

	index, err := bleve.Open(blevePath)
	if err != nil {
		fmt.Printf("❌ Error opening index: %v\n", err)
		return
	}
	defer index.Close()

	query := bleve.NewQueryStringQuery(queryStr)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"title", "source", "id", "summary"}

	searchResults, err := index.Search(searchRequest)
	if err != nil {
		fmt.Printf("❌ Error searching: %v\n", err)
		return
	}

	if searchResults.Total == 0 {
		fmt.Printf("No results found for %q\n", queryStr)
		return
	}

	fmt.Printf("Found %d matches for %q (Bleve Keyword):\n\n", searchResults.Total, queryStr)
	for i, hit := range searchResults.Hits {
		title := hit.Fields["title"]
		source := hit.Fields["source"]
		fmt.Printf("%d. [%s] %v (ID: %s)\n", i+1, source, title, hit.ID)
		if summary, ok := hit.Fields["summary"]; ok && summary != "" {
			fmt.Printf("   Summary: %v\n", summary)
		}
		fmt.Println()
	}
}

func runSemanticSearch(out, queryStr, model string) {
	vectorFile := filepath.Join(out, "musu.vectors.json")
	indexFile := filepath.Join(out, "index.json")
	if _, err := os.Stat(vectorFile); os.IsNotExist(err) {
		fmt.Println("❌ Vector index not found. Please run 'musu-crawl index --semantic' first.")
		return
	}

	vstore := processor.NewVectorStore()
	if err := vstore.Load(vectorFile); err != nil {
		fmt.Printf("❌ Error loading vectors: %v\n", err)
		return
	}

	// Load metadata for display
	metaData := make(map[string]processor.IndexEntry)
	if data, err := os.ReadFile(indexFile); err == nil {
		var entries []processor.IndexEntry
		json.Unmarshal(data, &entries)
		for _, e := range entries {
			metaData[e.ID] = e
		}
	}

	ollama := agent.NewOllamaClient(model)
	fmt.Printf("🧠 Generating query embedding for: %q...\n", queryStr)
	qVec, err := ollama.Embed(queryStr)
	if err != nil {
		fmt.Printf("❌ Embedding failed: %v\n", err)
		return
	}

	matches := vstore.Search(qVec, 5)
	if len(matches) == 0 {
		fmt.Println("No semantic matches found.")
		return
	}

	fmt.Printf("Found matches for %q (Semantic Vector):\n\n", queryStr)
	for i, match := range matches {
		entry, ok := metaData[match.ID]
		title := match.ID
		source := "?"
		summary := ""
		if ok {
			title = entry.Title
			source = entry.Source
			summary = entry.Summary
		}

		fmt.Printf("%d. [%s] %s (ID: %s, Score: %.4f)\n", i+1, source, title, match.ID, match.Score)
		if summary != "" {
			fmt.Printf("   Summary: %s\n", summary)
		}
		fmt.Println()
	}
}

func init() {
	searchCmd.Flags().String("out", "./wiki", "Wiki directory to search in")
	searchCmd.Flags().BoolP("semantic", "s", false, "Use vector semantic search")
	searchCmd.Flags().String("model", "nomic-embed-text", "Ollama model for query embedding")
	rootCmd.AddCommand(searchCmd)
}
