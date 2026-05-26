package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yellowhama/musu-crawl-ai/internal/agent"
	"github.com/yellowhama/musu-crawl-ai/internal/processor"
	"github.com/yellowhama/musu-crawl-ai/internal/utils"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch [source] [id]",
	Short: "Fetch content from a source (or multiple from a file)",
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("json", cmd.Flags().Lookup("json"))
		project, _ := cmd.Flags().GetString("project")
		conf, _ := utils.LoadConfig(project)

		filePath, _ := cmd.Flags().GetString("file")
		workers, _ := cmd.Flags().GetInt("workers")

		lang, _ := cmd.Flags().GetString("lang")
		if !cmd.Flags().Changed("lang") {
			lang = conf.Language
		}

		out, _ := cmd.Flags().GetString("out")
		if !cmd.Flags().Changed("out") {
			out = conf.WikiDir
		}

		compile, _ := cmd.Flags().GetBool("compile")

		model, _ := cmd.Flags().GetString("model")
		if !cmd.Flags().Changed("model") {
			model = conf.OllamaModel
		}

		proc := processor.NewWikiProcessor(out, project)

		if filePath != "" {
			runBatch(filePath, workers, lang, proc, compile, model)
		} else {
			if len(args) < 2 {
				utils.PrintError(fmt.Errorf("missing arguments"), "Provide [source] and [id] or use --file")
				return
			}
			RunSingle(args[0], args[1], lang, proc, compile, model)
		}
	},
}

func RunSingle(source, id, lang string, proc *processor.WikiProcessor, compile bool, model string) (string, error) {
	// Use agent.FetchAndSave instead of processor.FetchAndSave
	fname, text, reliability, tags, err := agent.FetchAndSave(source, id, lang, proc, model)
	if err != nil {
		utils.PrintError(err, "")
		return "", err
	}

	utils.PrintSuccess("Saved [%s] to Wiki project '%s': %s (Reliability: %.1f)", id, proc.Project, fname, reliability)

	if viper.GetBool("json") {
		utils.PrintJSON("Fetch successful", map[string]interface{}{
			"id":          id,
			"project":     proc.Project,
			"filename":    fname,
			"reliability": reliability,
			"tags":        tags,
		})
	}

	return text, nil
}

type job struct {
	source string
	id     string
}

func runBatch(filePath string, numWorkers int, lang string, proc *processor.WikiProcessor, compile bool, model string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	jobs := make(chan job)
	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range jobs {
				RunSingle(j.source, j.id, lang, proc, compile, model)
			}
		}(i)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") { continue }
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			jobs <- job{source: parts[0], id: parts[1]}
		} else {
			id := parts[0]
			source := autoDetectSource(id)
			if source != "" {
				jobs <- job{source: source, id: id}
			}
		}
	}
	close(jobs)
	wg.Wait()
}

func autoDetectSource(input string) string {
	if strings.Contains(input, "youtube.com") || strings.Contains(input, "youtu.be") { return "yt" }
	if strings.Contains(input, "github.com") { return "gh" }
	if strings.Contains(input, "arxiv.org") { return "arxiv" }
	if strings.Contains(input, "huggingface.co") { return "hf" }
	if strings.Contains(input, "twitter.com") || strings.Contains(input, "x.com") { return "x" }
	if strings.Contains(input, "reddit.com") { return "reddit" }
	if strings.HasPrefix(input, "http") { return "web" }
	return ""
}

func init() {
	fetchCmd.Flags().Bool("json", false, "Output in machine-readable JSON format")
	fetchCmd.Flags().StringP("file", "f", "", "Input file with source and id/url per line")
	fetchCmd.Flags().IntP("workers", "w", 5, "Number of concurrent workers")
	fetchCmd.Flags().String("lang", "ko", "Preferred language")
	fetchCmd.Flags().String("out", "./wiki", "Output directory")
	fetchCmd.Flags().Bool("compile", false, "Automatically compile knowledge links after fetch")
	fetchCmd.Flags().String("model", "llama3", "Ollama model for compilation reasoning")
	fetchCmd.Flags().StringP("project", "p", "default", "Project name to scope the knowledge")
	rootCmd.AddCommand(fetchCmd)
}
