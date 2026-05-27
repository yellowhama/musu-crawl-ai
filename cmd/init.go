package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func bootstrapProjectDirs(out, project string, verbose bool) error {
	dirs := []string{
		out,
		filepath.Join(out, "projects"),
		filepath.Join(out, "projects", project),
		filepath.Join(out, "projects", project, "youtube"),
		filepath.Join(out, "projects", project, "github"),
		filepath.Join(out, "projects", project, "papers"),
		filepath.Join(out, "projects", project, "web"),
		filepath.Join(out, "projects", project, "twitter"),
		filepath.Join(out, "projects", project, "huggingface"),
		filepath.Join(out, "projects", project, "reddit"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", d, err)
		}
		if verbose {
			fmt.Printf("✅ Directory ready: %s\n", d)
		}
	}
	return nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Wiki directory and check environment",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			project = "default"
		}

		fmt.Printf("🚀 Initializing musu-crawl-ai (Version %s)...\n", Version)

		// 1. Create directory structure
		if err := bootstrapProjectDirs(out, project, true); err != nil {
			fmt.Printf("❌ Failed to initialize wiki structure: %v\n", err)
			return
		}

		// 2. Check for AI Service
		fmt.Print("🧠 Checking for AI service (Ollama/SGLang)... ")
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get("http://localhost:11434/v1/models") // Use standard OpenAI endpoint
		if err != nil || resp.StatusCode != 200 {
			fmt.Println("⚠️  Local AI service not detected on default port. AI features will be limited.")
			fmt.Println("   Tip: Start Ollama or SGLang to enable full intelligence.")
		} else {
			fmt.Println("✅ AI service detected and active.")
			resp.Body.Close()
		}

		// 3. Check for Search Index
		blevePath := filepath.Join(out, "musu.bleve")
		if _, err := os.Stat(blevePath); os.IsNotExist(err) {
			fmt.Println("ℹ️  Search index not found. You may want to run 'musu-crawl index' once you have fetched content.")
		} else {
			fmt.Println("✅ Search index found.")
		}

		fmt.Println("\n✨ Initialization complete! You are ready to crawl.")
	},
}

func init() {
	initCmd.Flags().String("out", "./wiki", "Wiki directory to initialize")
	initCmd.Flags().String("project", "default", "Project directory to initialize under wiki/projects")
	rootCmd.AddCommand(initCmd)
}
