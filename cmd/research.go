package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yellowhama/musu-crawl-ai/internal/agent"
	"github.com/yellowhama/musu-crawl-ai/internal/processor"
)

var researchCmd = &cobra.Command{
	Use:   "research [question]",
	Short: "Perform autonomous multi-agent research",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := args[0]
		model, _ := cmd.Flags().GetString("model")
		limit, _ := cmd.Flags().GetInt("limit")
		out, _ := cmd.Flags().GetString("out")

		ollama := agent.NewOllamaClient(model)
		planner := &agent.Planner{Client: ollama}
		searcher := &agent.Searcher{}
		analyst := &agent.Analyst{Client: ollama}
		proc := processor.NewWikiProcessor(out)

		fmt.Printf("🚀 Starting research for: %q\n", question)

		// 1. Plan
		fmt.Println("🧠 Planning search strategy...")
		plan, err := planner.CreatePlan(question)
		if err != nil {
			fmt.Printf("❌ Planning failed: %v\n", err)
			return
		}
		fmt.Printf("📋 Plan: %s\n", plan.Reason)

		// 2. Discover
		fmt.Println("🌐 Discovering sources...")
		urls := searcher.DiscoverURLs(plan.Queries, limit)
		if len(urls) == 0 {
			fmt.Println("❌ No sources discovered.")
			return
		}
		fmt.Printf("🔗 Discovered %d unique sources.\n", len(urls))

		// 3. Harvest & Process (Using unified RunSingle)
		fmt.Println("⛏️  Harvesting and processing content...")
		var contents []string
		for _, url := range urls {
			fmt.Printf("   Processing: %s...\n", url)
			source := autoDetectSource(url)
			if source == "" {
				source = "web"
			}
			
			text, err := RunSingle(source, url, "en", proc, false, "")
			if err != nil {
				fmt.Printf("   ⚠️  Skip [%s]: %v\n", url, err)
				continue
			}
			contents = append(contents, text)
		}

		if len(contents) == 0 {
			fmt.Println("❌ No content could be harvested.")
			return
		}

		// 4. Analyze
		fmt.Println("📊 Analyzing and synthesizing results...")
		result, err := analyst.Synthesize(question, contents)
		if err != nil {
			fmt.Printf("❌ Analysis failed: %v\n", err)
			return
		}

		fmt.Println("\n--- FINAL RESEARCH REPORT ---")
		fmt.Println(result)
		
		// Final Index
		proc.UpdateIndex()
		fmt.Println("\n✅ Research completed and indexed.")
	},
}

func init() {
	researchCmd.Flags().String("model", "llama3", "Local Ollama model to use")
	researchCmd.Flags().Int("limit", 5, "Maximum number of sources to fetch")
	researchCmd.Flags().String("out", "./wiki", "Output directory")
	rootCmd.AddCommand(researchCmd)
}
