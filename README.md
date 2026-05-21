# musu-crawl-ai

AI-Ready Knowledge Harvester & Wiki Generator.

`musu-crawl-ai` is a high-performance, concurrent web crawler built in Go. It extracts valuable information from various platforms and automatically organizes it into a structured, Wiki-like format optimized for AI consumption (RAG, Fine-tuning) and human readability.

## 🚀 Features

-   **Multi-Source Harvesters:**
    -   **YouTube (`yt`):** Captions/transcripts (with Innertube fallback to bypass restrictions).
    -   **GitHub (`gh`):** `README.md` and repository metadata.
    -   **Arxiv (`arxiv`):** Research paper metadata and PDF text extraction.
    -   **Hugging Face (`hf`):** Model Cards and metadata.
    -   **General Web (`web`):** Clean article content extraction (Readability).
    -   **Twitter/X (`x`):** Tweet content (via Syndication/OEmbed fallback).
-   **Auto-Wiki & Indexing:**
    -   Generates Markdown files with YAML Frontmatter.
    -   **Auto-Tagging:** Automatically extracts keywords from content.
    -   **Local Summarization:** Pure Go extractive summarizer (top 3 sentences) for quick context.
    -   **Local Search Engine:** Built-in high-performance text search across all documents using Bleve.
    -   **Machine-Readable Index:** Unified `index.json` catalog for RAG integration.
-   **High Performance:**
    -   **Concurrency:** Worker pool pattern (Goroutines) for massive parallel downloads.
    -   **Robustness:** Exponential backoff retry logic for all HTTP requests.

## 🛠️ Installation

```bash
# Clone the repository
git clone https://github.com/yellowhama/musu-crawl-ai
cd musu-crawl-ai

# Build the executable
go build -o musu-crawl
```

## 📖 Usage

### Single Fetch
```bash
./musu-crawl fetch yt [VIDEO_ID]
./musu-crawl fetch gh [OWNER/REPO]
./musu-crawl fetch arxiv [ARXIV_ID]
./musu-crawl fetch web [URL]
```

### Local Search
```bash
# Search for keywords across the entire harvested knowledge base
./musu-crawl search "machine learning"
```

### Batch Fetch (Parallel)
Create a `targets.txt` file:
```text
yt dQw4w9WgXcQ
gh spf13/cobra
web https://go.dev/blog/go1.22
```
Run with 10 concurrent workers:
```bash
./musu-crawl fetch --file targets.txt --workers 10
```

### Re-indexing
```bash
# Rebuild index.json, README.md, and Search Index
./musu-crawl index --out ./wiki
```

## 📂 Directory Structure
- `/wiki`: The generated knowledge base.
- `/wiki/index.json`: Centralized metadata catalog for RAG.
- `/wiki/musu.bleve`: Local search engine database.
- `/wiki/README.md`: Human-readable master index.
- `/internal/harvester`: Modular fetcher implementations.

## 📝 Roadmap
- [x] Concurrent Batch Processing
- [x] Robust Retry Logic (Exponential Backoff)
- [x] Auto-Tagging & Local Summarization
- [x] Built-in Local Search Engine
- [ ] Direct Vector DB Export (Ollama Embedding Integration)
- [ ] PDF OCR Support for scanned papers
