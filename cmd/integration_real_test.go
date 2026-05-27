//go:build integration

package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func captureCommandStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	fn()

	_ = w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	_ = r.Close()
	return buf.String()
}

func TestFetchAndSemanticSearchRealAIIntegration(t *testing.T) {
	aiURL := strings.TrimSpace(os.Getenv("MUSU_CRAWL_INTEGRATION_AI_URL"))
	if aiURL == "" {
		t.Skip("set MUSU_CRAWL_INTEGRATION_AI_URL to run real integration")
	}
	embedModel := strings.TrimSpace(os.Getenv("MUSU_CRAWL_INTEGRATION_EMBED_MODEL"))
	if embedModel == "" {
		embedModel = "nomic-embed-text"
	}

	wd, _ := os.Getwd()
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	page := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html><html><head><title>Scheduler Integration Page</title></head><body><main><h1>Scheduler Reliability</h1><p>Schedulers need reliability, timing, and operator trust.</p></main></body></html>`))
	}))
	defer page.Close()

	wikiDir := "./wiki"
	if err := bootstrapProjectDirs(wikiDir, "integration-demo", "openai", aiURL, false); err != nil {
		t.Fatalf("bootstrapProjectDirs failed: %v", err)
	}
	if err := os.Setenv("MUSU_OUT", wikiDir); err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("MUSU_OUT")

	if err := fetchCmd.Flags().Set("project", "integration-demo"); err != nil {
		t.Fatal(err)
	}
	if err := fetchCmd.Flags().Set("out", wikiDir); err != nil {
		t.Fatal(err)
	}
	if err := fetchCmd.Flags().Set("compile", "false"); err != nil {
		t.Fatal(err)
	}
	if err := fetchCmd.Flags().Set("model", embedModel); err != nil {
		t.Fatal(err)
	}
	fetchCmd.Run(fetchCmd, []string{"web", page.URL})

	if _, err := os.Stat(filepath.Join(wikiDir, "musu.vectors.json")); err != nil {
		t.Fatalf("expected vector store after real fetch integration: %v", err)
	}

	if err := searchCmd.Flags().Set("out", wikiDir); err != nil {
		t.Fatal(err)
	}
	if err := searchCmd.Flags().Set("project", "integration-demo"); err != nil {
		t.Fatal(err)
	}
	if err := searchCmd.Flags().Set("semantic", "true"); err != nil {
		t.Fatal(err)
	}
	if err := searchCmd.Flags().Set("model", embedModel); err != nil {
		t.Fatal(err)
	}
	if err := searchCmd.Flags().Set("json", "true"); err != nil {
		t.Fatal(err)
	}
	viper.Set("json", true)

	out := captureCommandStdout(t, func() {
		searchCmd.Run(searchCmd, []string{"scheduler"})
	})
	if !strings.Contains(out, `"status": "success"`) {
		t.Fatalf("expected successful semantic search JSON, got:\n%s", out)
	}
	if !strings.Contains(out, `Scheduler Integration Page`) {
		t.Fatalf("expected semantic search result to mention fetched page, got:\n%s", out)
	}
}
