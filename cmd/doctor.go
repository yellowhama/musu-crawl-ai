package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yellowhama/musu-crawl-ai/internal/utils"
)

var doctorCapabilitySources []string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check wiki, index, project config, and AI connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, _ := cmd.Flags().GetString("out")
		project, _ := cmd.Flags().GetString("project")
		conf, err := utils.LoadConfig(project)
		if err != nil {
			return err
		}
		if out == "" {
			out = conf.WikiDir
		}

		jsonMode := viper.GetBool("json")
		if !jsonMode {
			fmt.Println("==> musu-crawl-ai doctor")
			fmt.Printf("Wiki Dir     : %s\n", out)
			fmt.Printf("Project      : %s\n", project)
			fmt.Printf("AI Provider  : %s\n", conf.AIProvider)
			fmt.Printf("AI URL       : %s\n", conf.AIBaseURL)
		}

		hasError := false
		report := map[string]interface{}{
			"wiki_dir":             out,
			"project":              project,
			"ai_provider":          conf.AIProvider,
			"ai_url":               conf.AIBaseURL,
			"wiki_exists":          false,
			"wiki_markdown_count":  0,
			"index_exists":         false,
			"project_exists":       false,
			"ai_reachable":         false,
		}
		if len(doctorCapabilitySources) > 0 {
			report["source_capabilities"] = sourceCapabilityReport(doctorCapabilitySources)
		}
		autoFix, _ := cmd.Flags().GetBool("fix")

		if info, statErr := os.Stat(out); statErr != nil || !info.IsDir() {
			if autoFix {
				if !jsonMode {
					fmt.Printf("🛠️  Auto-fixing missing wiki scaffold at %s\n", out)
				}
				if fixErr := bootstrapProjectDirs(out, projectName(project), !jsonMode); fixErr == nil {
					report["wiki_exists"] = true
					report["project_exists"] = true
				} else {
					report["wiki_fix_error"] = fixErr.Error()
					if !jsonMode {
						fmt.Printf("❌ Auto-fix failed: %v\n", fixErr)
					}
					hasError = true
				}
			} else if !jsonMode {
				fmt.Printf("❌ Wiki directory missing: %s\n", out)
				fmt.Println("   Fix: run 'musu-crawl init --out <dir>' first, or use doctor --fix.")
				hasError = true
			}
		} else {
			mdCount := countMarkdownFiles(out)
			if !jsonMode {
				fmt.Printf("✅ Wiki directory exists (%d markdown files)\n", mdCount)
			}
			report["wiki_exists"] = true
			report["wiki_markdown_count"] = mdCount
		}

		blevePath := filepath.Join(out, "musu.bleve")
		if _, statErr := os.Stat(blevePath); statErr != nil {
			if !jsonMode {
				fmt.Printf("⚠️  Search index missing: %s\n", blevePath)
				fmt.Println("   Run: musu-crawl index --out", out)
			}
		} else {
			if !jsonMode {
				fmt.Println("✅ Search index exists")
			}
			report["index_exists"] = true
		}

		if project != "" && project != "all" && project != "default" {
			projectDir := filepath.Join(out, "projects", project)
			if _, statErr := os.Stat(projectDir); statErr != nil {
				if !jsonMode {
					fmt.Printf("⚠️  Project dir missing: %s\n", projectDir)
				}
			} else {
				if !jsonMode {
					fmt.Println("✅ Project directory exists")
				}
				report["project_exists"] = true
			}
		}

		if err := probeModels(conf.AIBaseURL); err != nil {
			if !jsonMode {
				fmt.Printf("❌ AI endpoint probe failed: %v\n", err)
			}
			report["ai_error"] = err.Error()
			hasError = true
		} else {
			if !jsonMode {
				fmt.Println("✅ AI endpoint reachable")
			}
			report["ai_reachable"] = true
		}

		if hasError {
			if jsonMode {
				encoded, _ := json.MarshalIndent(utils.JSONResponse{
					Status:    "error",
					Message:   "doctor found blocking issues",
					Data:      report,
					ActionFix: "Initialize the wiki/project or start the configured AI endpoint.",
				}, "", "  ")
				fmt.Println(string(encoded))
			}
			return fmt.Errorf("doctor found blocking issues")
		}
		if !jsonMode {
			fmt.Println("✅ Doctor passed")
		}
		utils.PrintJSON("Doctor passed", report)
		return nil
	},
}

func init() {
	doctorCmd.Flags().String("out", "./wiki", "Wiki directory to inspect")
	doctorCmd.Flags().StringP("project", "p", "default", "Project to inspect")
	doctorCmd.Flags().Bool("fix", false, "Auto-create missing local wiki scaffold when safe")
	doctorCmd.Flags().StringSliceVar(&doctorCapabilitySources, "capability-source", nil, "Optional source(s) to describe capability metadata for (web, yt, gh, arxiv, reddit, hf, x)")
	rootCmd.AddCommand(doctorCmd)
}

func projectName(project string) string {
	if strings.TrimSpace(project) == "" || project == "all" {
		return "default"
	}
	return project
}

func sourceCapabilityReport(sources []string) []map[string]interface{} {
	var report []map[string]interface{}
	for _, source := range sources {
		key := strings.ToLower(strings.TrimSpace(source))
		switch key {
		case "web", "yt", "gh", "arxiv", "reddit", "hf", "x":
			report = append(report, map[string]interface{}{
				"source":            key,
				"supported":         true,
				"auth_required":     false,
				"capability_static": true,
				"recommended_mode":  "public-fetch",
			})
		default:
			report = append(report, map[string]interface{}{
				"source":            key,
				"supported":         false,
				"auth_required":     false,
				"capability_static": true,
			})
		}
	}
	return report
}

func probeModels(baseURL string) error {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return fmt.Errorf("empty ai-url")
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(baseURL + "/models")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected status %s from %s/models", resp.Status, baseURL)
}

func countMarkdownFiles(root string) int {
	count := 0
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".md") {
			count++
		}
		return nil
	})
	return count
}
