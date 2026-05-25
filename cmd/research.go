package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yellowhama/musu-crawl-ai/internal/agent"
	"github.com/yellowhama/musu-crawl-ai/internal/processor"
	"github.com/yellowhama/musu-crawl-ai/internal/utils"
)

var researchCmd = &cobra.Command{
	Use:   "research [question]",
	Short: "Perform autonomous recursive multi-agent research with skepticism",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := args[0]
		project, _ := cmd.Flags().GetString("project")
		conf, _ := utils.LoadConfig(project)

		model, _ := cmd.Flags().GetString("model")
		if !cmd.Flags().Changed("model") {
			model = conf.OllamaModel
		}

		limit, _ := cmd.Flags().GetInt("limit")

		out, _ := cmd.Flags().GetString("out")
		if !cmd.Flags().Changed("out") {
			out = conf.WikiDir
		}

		maxDepth, _ := cmd.Flags().GetInt("depth")

		// Load custom persona if exists
		customPersona := ""
		personaPath := filepath.Join(out, "projects", project, "PROMPT.md")
		if data, err := os.ReadFile(personaPath); err == nil {
			customPersona = string(data)
			fmt.Printf("🎭 Loaded custom persona for project '%s'\n", project)
		}

		ollama := agent.NewOllamaClient(model)
		planner := &agent.Planner{Client: ollama, CustomPersona: customPersona}
		searcher := &agent.Searcher{}
		analyst := &agent.Analyst{Client: ollama, CustomPersona: customPersona}
		proc := processor.NewWikiProcessor(out, project)

		fmt.Printf("🚀 Starting autonomous research for: %q (Project: %s)\n", question, project)

		seenURLs := make(map[string]bool)
		currentGoal := question
		var allSources []agent.SourceContent

		for depth := 1; depth <= maxDepth; depth++ {
			fmt.Printf("\n--- 📍 Research Depth %d/%d ---\n", depth, maxDepth)
			fmt.Printf("🎯 Current Target: %s\n", currentGoal)

			// 1. Plan
			fmt.Println("🧠 Planning search strategy (Socratic mode)...")
			plan, err := planner.CreatePlan(currentGoal)
			if err != nil {
				fmt.Printf("❌ Planning failed: %v\n", err)
				break
			}
			fmt.Printf("📋 Strategy: %s\n", plan.Reason)
			if len(plan.Hypotheses) > 0 {
				fmt.Printf("🔍 Testing Hypotheses: %s\n", strings.Join(plan.Hypotheses, ", "))
			}

			// 2. Discover
			fmt.Println("🌐 Discovering sources...")
			urls := searcher.DiscoverURLs(plan.Queries, limit)

			var newURLs []string
			for _, u := range urls {
				if !seenURLs[u] {
					seenURLs[u] = true
					newURLs = append(newURLs, u)
				}
			}

			if len(newURLs) == 0 {
				fmt.Println("ℹ️  No new sources discovered in this iteration.")
				if depth == 1 {
					break
				}
			} else {
				fmt.Printf("🔗 Found %d new sources.\n", len(newURLs))

				// 3. Harvest & Process
				fmt.Println("⛏️  Harvesting and evaluating content...")
				for _, url := range newURLs {
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
					
					reliability := GetReliability(source)
					allSources = append(allSources, agent.SourceContent{
						ID:          url,
						Content:     text,
						Reliability: reliability,
						SourceType:  source,
					})
				}
			}

			if len(allSources) == 0 {
				fmt.Println("❌ No content available to analyze.")
				break
			}

			// 4. Analyze
			fmt.Println("📊 Analyzing accumulated knowledge (Skeptical mode)...")
			report, err := analyst.Synthesize(question, allSources)
			if err != nil {
				fmt.Printf("❌ Analysis failed: %v\n", err)
				break
			}

			fmt.Println("\n--- CURRENT FINDINGS ---")
			fmt.Println(report)

			// Recursive logic
			recursiveQuery := ""
			
			// 1. Check for contradictions
			contraMatch := regexp.MustCompile(`(?s)CONTRADICTIONS:\s*(.*)`).FindStringSubmatch(report)
			if len(contraMatch) >= 2 {
				contraInfo := strings.TrimSpace(contraMatch[1])
				if strings.ToLower(contraInfo) != "none" && strings.ToLower(contraInfo) != "none." && len(contraInfo) > 5 {
					fmt.Printf("\n⚔️  Contradiction found! Initiating tie-break research for: %s\n", contraInfo)
					recursiveQuery += "Resolve this contradiction: " + contraInfo + ". "
				}
			}

			// 2. Check for missing info
			missingMatch := regexp.MustCompile(`(?s)MISSING:\s*(.*)`).FindStringSubmatch(report)
			if len(missingMatch) >= 2 {
				missingInfo := strings.TrimSpace(missingMatch[1])
				lowerMissing := strings.ToLower(missingInfo)
				if lowerMissing != "none" && lowerMissing != "none." && len(missingInfo) > 5 {
					fmt.Printf("\n🔁 Gaps identified. Initiating next hop for: %s\n", missingInfo)
					recursiveQuery += "Investigate these missing points: " + missingInfo
				}
			}

			if recursiveQuery == "" {
				fmt.Println("\n✨ Research goal fully satisfied. Stopping.")
				break
			}
			currentGoal = recursiveQuery
		}

		// Final Index
		proc.UpdateIndex()
		fmt.Println("\n✅ Autonomous research completed and indexed.")
	},
}

func init() {
	researchCmd.Flags().String("model", "llama3", "Local Ollama model to use")
	researchCmd.Flags().Int("limit", 5, "Maximum number of sources to fetch per hop")
	researchCmd.Flags().Int("depth", 2, "Maximum recursive research depth")
	researchCmd.Flags().String("out", "./wiki", "Output directory")
	researchCmd.Flags().StringP("project", "p", "default", "Project name to scope the research")
	rootCmd.AddCommand(researchCmd)
}
