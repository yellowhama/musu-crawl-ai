# musu-crawl-ai

AI-Ready Knowledge Harvester & Wiki Generator.

`musu-crawl-ai` is a high-performance, concurrent web crawler built in Go. It extracts valuable information from various platforms and automatically organizes it into a structured, Wiki-like format optimized for AI consumption (RAG, Fine-tuning) and human readability.

## 🚀 Features

-   **Multi-Source Harvesters:**
    -   **YouTube (`yt`):** Captions/transcripts (with Innertube fallback).
    -   **GitHub (`gh`):** `README.md` and repository metadata.
    -   **Arxiv (`arxiv`):** Research paper metadata and PDF text extraction.
    -   **Hugging Face (`hf`):** Model Cards and metadata.
    -   **General Web (`web`):** Clean article content extraction (Readability).
    -   **Twitter/X (`x`):** Tweet content (via Syndication/OEmbed fallback).
    -   **Reddit:** Post and subreddit listing extraction.
-   **Auto-Wiki & Intelligence:**
    -   **Wiki Compiler Agent:** Automatically discovers and creates bidirectional links between related documents using local reasoning.
    -   **Auto-Tagging:** Automatically extracts keywords from content.
    -   **Local Summarization:** Pure Go extractive summarizer (top 3 sentences).
    -   **Local Search Engine:** Built-in text search using Bleve.
-   **High Performance:**
    -   **Concurrency:** Worker pool pattern (Goroutines).
    -   **Robustness:** Exponential backoff retry logic.

## 📖 Usage

### Single Fetch with Auto-Linking
```bash
./musu-crawl fetch yt [VIDEO_ID] --compile
```

### Standalone Wiki Compilation
```bash
# Analyze all documents and forge knowledge links
./musu-crawl compile
```

### Autonomous Research
```bash
# Decompose, Discover, Harvest, and Synthesize
./musu-crawl research "What are the latest AI trends?"
```

### Local Search
```bash
./musu-crawl search "machine learning"
```

### Re-indexing
```bash
./musu-crawl index
```

## 📂 Directory Structure
- `/wiki`: The generated knowledge base.
- `/wiki/index.json`: Machine-readable catalog.
- `/wiki/musu.bleve`: Search database.
- `/wiki/README.md`: Human-readable index.

## 📝 Roadmap
- [x] Wiki Compiler Agent (Autonomous Cross-Linking)
- [x] Multi-Agent Research Loop
- [ ] Direct Vector DB Export (Ollama Embedding Integration)
- [ ] PDF OCR Support
