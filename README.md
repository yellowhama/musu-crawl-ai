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
    -   **Recursive Research Loop:** Fully autonomous multi-hop research. Decomposes questions, discovers sources, and repeats the process if information gaps are found.
    -   **Wiki Compiler Agent:** Automatically discovers and creates bidirectional links between related documents using local reasoning.
    -   **Semantic Vector Search:** Built-in meaning-based search using local Ollama embeddings (Hybrid search with Bleve keywords).
    -   **Auto-Tagging:** Automatically extracts keywords from content.
    -   **Local Summarization:** Pure Go extractive summarizer (top 3 sentences).
-   **High Performance:**
    -   **Concurrency:** Worker pool pattern (Goroutines).
    -   **Robustness:** Exponential backoff retry logic.

## 🛠️ Installation & Setup

### 1. Download
Download the latest binary for your OS from the [GitHub Releases](https://github.com/yellowhama/musu-crawl-ai/releases) page.

### 2. Initialize
Run the setup command to scaffold your local wiki and check for dependencies:
```bash
./musu-crawl init
```

### 3. Update (Anytime)
Keep your tool up-to-date with a single command:
```bash
./musu-crawl update
```

## 📖 Usage

### Autonomous Deep Research (Recursive)
```bash
# Decompose, Discover, Harvest, and recursively fill information gaps
./musu-crawl research "Compare the performance of Go 1.22 vs 1.21" --depth 2
```

### Semantic & Keyword Search
```bash
# Meaning-based search (Semantic Vector)
./musu-crawl search "machine learning concepts" --semantic

# Keyword-based search (Bleve)
./musu-crawl search "golang concurrency"
```

### Single Fetch with Auto-Linking
```bash
./musu-crawl fetch yt [VIDEO_ID] --compile
```

### Re-indexing (with Embeddings)
```bash
# Build keyword index and generate vector embeddings
./musu-crawl index --semantic
```

## 📂 Directory Structure
- `/wiki`: The generated knowledge base.
- `/wiki/index.json`: Machine-readable catalog.
- `/wiki/musu.bleve`: Keyword search database.
- `/wiki/musu.vectors.json`: Semantic vector store.
- `/wiki/README.md`: Human-readable master index.

## 📝 Roadmap
- [x] Semantic Vector Search (Local Ollama)
- [x] Recursive Multi-hop Research
- [x] Wiki Compiler Agent (Autonomous Cross-Linking)
- [ ] Direct Vector DB Export (Pinecone/Weaviate integration)
- [ ] PDF OCR Support
