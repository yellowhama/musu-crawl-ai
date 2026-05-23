package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/yellowhama/musu-crawl-ai/internal/processor"
	"github.com/yuin/goldmark"
)

type Server struct {
	WikiDir string
	Port    int
}

func NewServer(wikiDir string, port int) *Server {
	return &Server{WikiDir: wikiDir, Port: port}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/search", s.handleSearch)
	http.HandleFunc("/view", s.handleView)

	fmt.Printf("🌐 musu-crawl dashboard starting at http://localhost:%d\n", s.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	indexFile := filepath.Join(s.WikiDir, "index.json")
	var entries []processor.IndexEntry
	
	if data, err := os.ReadFile(indexFile); err == nil {
		json.Unmarshal(data, &entries)
	}

	tmpl := template.Must(template.New("layout").Parse(layoutHTML))
	template.Must(tmpl.New("content").Parse(indexHTML))
	
	tmpl.Execute(w, map[string]interface{}{
		"Entries": entries,
	})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	queryStr := r.URL.Query().Get("q")
	blevePath := filepath.Join(s.WikiDir, "musu.bleve")
	
	var results []processor.IndexEntry
	
	if queryStr != "" {
		index, err := bleve.Open(blevePath)
		if err == nil {
			defer index.Close()
			query := bleve.NewQueryStringQuery(queryStr)
			searchRequest := bleve.NewSearchRequest(query)
			searchRequest.Fields = []string{"title", "source", "id", "summary", "path"}
			searchRes, err := index.Search(searchRequest)
			if err == nil {
				for _, hit := range searchRes.Hits {
					results = append(results, processor.IndexEntry{
						ID:      hit.ID,
						Title:   hit.Fields["title"].(string),
						Source:  hit.Fields["source"].(string),
						Summary: hit.Fields["summary"].(string),
						Path:    hit.Fields["path"].(string),
					})
				}
			}
		}
	}

	tmpl := template.Must(template.New("layout").Parse(layoutHTML))
	template.Must(tmpl.New("content").Parse(searchHTML))
	
	tmpl.Execute(w, map[string]interface{}{
		"Query":   queryStr,
		"Results": results,
	})
}

func (s *Server) handleView(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	fullPath := filepath.Join(s.WikiDir, path)
	
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}

	// Simple frontmatter strip
	content := string(data)
	if idx := strings.Index(content, "---"); idx == 0 {
		if nextIdx := strings.Index(content[3:], "---"); nextIdx != -1 {
			content = content[nextIdx+6:]
		}
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(content), &buf); err != nil {
		http.Error(w, "Markdown conversion failed", 500)
		return
	}

	tmpl := template.Must(template.New("layout").Parse(layoutHTML))
	template.Must(tmpl.New("content").Parse(viewHTML))
	
	tmpl.Execute(w, map[string]interface{}{
		"Title":   filepath.Base(path),
		"Content": template.HTML(buf.String()),
	})
}

// --- Embedded Templates ---

const layoutHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MUSU Knowledge Portal</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        'musu-accent': '#ffa602',
                        'musu-deep': '#432c1c',
                        'musu-cream': '#FDF8F1',
                    }
                }
            }
        }
    </script>
    <style>
        .prose h1 { color: #432c1c; font-size: 2rem; font-weight: bold; margin-bottom: 1rem; }
        .prose h2 { color: #432c1c; font-size: 1.5rem; font-weight: bold; margin-top: 1.5rem; }
        .prose p { margin-bottom: 1rem; line-height: 1.6; }
        .prose a { color: #ffa602; text-decoration: underline; }
    </style>
</head>
<body class="bg-musu-cream text-musu-deep min-h-screen">
    <nav class="bg-musu-deep text-white p-4 shadow-lg">
        <div class="container mx-auto flex justify-between items-center">
            <a href="/" class="text-2xl font-bold tracking-tighter text-musu-accent">MUSU <span class="text-white font-light">CRAWL</span></a>
            <form action="/search" method="GET" class="flex gap-2">
                <input type="text" name="q" placeholder="Search knowledge..." class="rounded px-3 py-1 text-musu-deep focus:ring-2 focus:ring-musu-accent outline-none w-64">
                <button type="submit" class="bg-musu-accent text-musu-deep font-bold px-4 py-1 rounded hover:opacity-90 transition">SEARCH</button>
            </form>
        </div>
    </nav>
    <main class="container mx-auto py-8 px-4">
        {{template "content" .}}
    </main>
    <footer class="text-center py-8 text-gray-500 text-sm">
        Powered by musu-crawl-ai • Autonomous Researcher Agent
    </footer>
</body>
</html>
`

const indexHTML = `
<div class="flex justify-between items-center mb-8">
    <h1 class="text-3xl font-bold">Knowledge Dashboard</h1>
    <span class="bg-musu-deep text-musu-accent px-4 py-1 rounded-full font-mono text-sm">{{len .Entries}} Documents</span>
</div>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    {{range .Entries}}
    <div class="bg-white p-6 rounded-xl shadow-sm border-l-4 border-musu-accent hover:shadow-md transition group">
        <div class="flex justify-between items-start mb-2">
            <span class="text-xs font-bold uppercase tracking-widest text-musu-accent">{{.Source}}</span>
            <span class="text-xs text-gray-400">{{.Date}}</span>
        </div>
        <h2 class="text-xl font-bold mb-3 line-clamp-2 group-hover:text-musu-accent transition">{{.Title}}</h2>
        <p class="text-gray-600 text-sm mb-4 line-clamp-3 leading-relaxed">{{.Summary}}</p>
        <a href="/view?path={{.Path}}" class="inline-block text-musu-deep font-bold text-sm border-b-2 border-musu-accent pb-0.5 hover:opacity-70 transition">READ MORE →</a>
    </div>
    {{end}}
</div>
`

const searchHTML = `
<h1 class="text-3xl font-bold mb-8">Results for "{{.Query}}"</h1>

<div class="space-y-6">
    {{range .Results}}
    <div class="bg-white p-6 rounded-xl shadow-sm border-l-4 border-musu-accent">
        <span class="text-xs font-bold uppercase tracking-widest text-musu-accent mb-2 block">{{.Source}}</span>
        <h2 class="text-2xl font-bold mb-2">{{.Title}}</h2>
        <p class="text-gray-600 mb-4">{{.Summary}}</p>
        <a href="/view?path={{.Path}}" class="text-musu-accent font-bold hover:underline">View Document →</a>
    </div>
    {{else}}
    <p class="text-gray-500 italic">No documents matched your query.</p>
    {{end}}
</div>
`

const viewHTML = `
<div class="max-w-4xl mx-auto">
    <a href="/" class="text-musu-accent font-bold mb-8 inline-block hover:underline">← Back to Dashboard</a>
    <article class="bg-white p-10 rounded-2xl shadow-xl prose prose-slate lg:prose-xl">
        {{.Content}}
    </article>
</div>
`
