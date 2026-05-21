package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the Wiki knowledge base locally",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
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

		queryStr := args[0]
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

		fmt.Printf("Found %d matches for %q:\n\n", searchResults.Total, queryStr)
		for i, hit := range searchResults.Hits {
			title := hit.Fields["title"]
			source := hit.Fields["source"]
			fmt.Printf("%d. [%s] %v (ID: %s)\n", i+1, source, title, hit.ID)
			if summary, ok := hit.Fields["summary"]; ok && summary != "" {
				fmt.Printf("   Summary: %v\n", summary)
			}
			fmt.Println()
		}
	},
}

func init() {
	searchCmd.Flags().String("out", "./wiki", "Wiki directory to search in")
	rootCmd.AddCommand(searchCmd)
}
