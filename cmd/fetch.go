package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/yellowhama/musu-crawl-ai/internal/harvester"
	"github.com/yellowhama/musu-crawl-ai/internal/processor"
	"github.com/yellowhama/musu-crawl-ai/internal/utils"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [source] [id]",
	Short: "Fetch content from a source (or multiple from a file)",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		workers, _ := cmd.Flags().GetInt("workers")
		lang, _ := cmd.Flags().GetString("lang")
		out, _ := cmd.Flags().GetString("out")

		proc := processor.NewWikiProcessor(out)

		if filePath != "" {
			// Batch processing mode
			runBatch(filePath, workers, lang, proc)
		} else {
			// Single item mode
			if len(args) < 2 {
				fmt.Println("Please provide source and id, or use --file")
				return
			}
			runSingle(args[0], args[1], lang, proc)
		}
	},
}

func runSingle(source, id, lang string, proc *processor.WikiProcessor) {
	title, text, err := dispatchFetch(source, id, lang)
	if err != nil {
		fmt.Printf("❌ Error [%s]: %v\n", id, err)
		return
	}

	// Auto-tagging
	tags := utils.ExtractKeywords(text, 5)

	// Local Summarization
	summary := utils.Summarize(text, 3)

	sourceDir := source
	if source == "yt" {
		sourceDir = "youtube"
	}
	if source == "gh" {
		sourceDir = "github"
	}
	if source == "arxiv" {
		sourceDir = "papers"
	}
	if source == "hf" {
		sourceDir = "huggingface"
	}
	if source == "x" {
		sourceDir = "twitter"
	}

	safeID := id
	if source == "gh" || source == "github" || source == "hf" || source == "huggingface" {
		safeID = strings.ReplaceAll(id, "/", "_")
	} else if source == "web" {
		safeID = strings.ReplaceAll(strings.ReplaceAll(id, "https://", ""), "/", "_")
		if len(safeID) > 100 {
			safeID = safeID[:100]
		}
	}

	fname, err := proc.SaveToWiki(sourceDir, safeID, title, text, tags, summary)
	if err != nil {
		fmt.Printf("❌ Error saving [%s]: %v\n", id, err)
		return
	}
	fmt.Printf("✅ Saved [%s] to Wiki: %s (Tags: %s)\n", id, fname, strings.Join(tags, ", "))
}

func dispatchFetch(source, id, lang string) (string, string, error) {
	source = strings.ToLower(source)
	switch source {
	case "yt", "youtube":
		f := &harvester.YouTubeFetcher{Language: lang}
		return f.Fetch(id)
	case "gh", "github":
		f := &harvester.GitHubFetcher{}
		return f.Fetch(id)
	case "web":
		f := &harvester.WebFetcher{}
		return f.Fetch(id)
	case "arxiv":
		f := &harvester.ArxivFetcher{}
		return f.Fetch(id)
	case "hf", "huggingface":
		f := &harvester.HuggingFaceFetcher{}
		return f.Fetch(id)
	case "x", "twitter":
		f := &harvester.TwitterFetcher{}
		return f.Fetch(id)
	default:
		return "", "", fmt.Errorf("unsupported source: %s", source)
	}
}

type job struct {
	source string
	id     string
}

func runBatch(filePath string, numWorkers int, lang string, proc *processor.WikiProcessor) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	jobs := make(chan job)
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range jobs {
				runSingle(j.source, j.id, lang, proc)
			}
		}(i)
	}

	// Feed jobs
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Expected format: "source id" or just "id" (if we could auto-detect, but let's stick to source id for now)
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			jobs <- job{source: parts[0], id: parts[1]}
		} else {
			// Auto-detection attempt
			id := parts[0]
			source := autoDetectSource(id)
			if source != "" {
				jobs <- job{source: source, id: id}
			} else {
				fmt.Printf("⚠️  Skipping line: %s (could not detect source)\n", line)
			}
		}
	}
	close(jobs)
	wg.Wait()
	fmt.Println("\nBatch processing completed.")
}

func autoDetectSource(input string) string {
	if strings.Contains(input, "youtube.com") || strings.Contains(input, "youtu.be") {
		return "yt"
	}
	if strings.Contains(input, "github.com") {
		return "gh"
	}
	if strings.Contains(input, "arxiv.org") {
		return "arxiv"
	}
	if strings.Contains(input, "huggingface.co") {
		return "hf"
	}
	if strings.Contains(input, "twitter.com") || strings.Contains(input, "x.com") {
		return "x"
	}
	if strings.HasPrefix(input, "http") {
		return "web"
	}
	return ""
}

func init() {
	fetchCmd.Flags().StringP("file", "f", "", "Input file with source and id/url per line")
	fetchCmd.Flags().IntP("workers", "w", 5, "Number of concurrent workers")
	fetchCmd.Flags().String("lang", "ko", "Preferred language")
	fetchCmd.Flags().String("out", "./wiki", "Output directory")
	rootCmd.AddCommand(fetchCmd)
}
