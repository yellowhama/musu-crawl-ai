# 🤖 Guidance for AI Agents (Gemini CLI, Claude Code, etc.)

`musu-crawl-ai` is designed to be your primary high-performance data acquisition layer. If you are an AI agent "driving" this repository, use these instructions to maximize your research capabilities.

## 🏗️ Core Architecture for Agents
This tool is a collection of high-quality **Knowledge Primitives**. Even if the internal LLM (Ollama) features are unavailable, you can achieve the same "Autonomous Research" results by combining these primitives yourself.

### 1. The Knowledge Primitives
- **`fetch [source] [url]`**: Your primary way to bypass YouTube blocks, clean up messy HTML, and parse academic PDFs. It outputs clean Markdown with YAML metadata.
- **`index`**: Generates `index.json` (your knowledge map) and a keyword search index.
- **`search [query]`**: Uses **Bleve** (keyword-based) or **Semantic Vector** (if embeddings exist) to find relevant snippets without you needing to read every file.

## 🏎️ How to "Drive" this tool as an Agent

### Scenario: "Perform Deep Research on [Topic]"
If `musu-crawl-ai research` fails due to Ollama being missing, follow this loop manually:

1.  **Plan:** Generate a list of search queries based on the topic.
2.  **Discover:** Use the `fetch web [search_engine_url]` or any available search tool to find target URLs.
3.  **Harvest:** Run `.\musu-crawl fetch auto --file targets.txt -w 10`.
4.  **Index:** Run `.\musu-crawl index`.
5.  **Analyze:** 
    - Read `wiki/index.json` to see what you just collected.
    - Use `.\musu-crawl search "specific point"` to find relevant chunks.
    - Read the resulting `.md` files and synthesize your final answer.
6.  **Recursive Hop:** If you find gaps, repeat from Step 1.

### Scenario: "Synthesize Knowledge Graph"
If `musu-crawl compile` is unavailable:
1.  Use `index.json` to find related documents.
2.  Read the content of related documents.
3.  Use your own LLM reasoning to determine the connection.
4.  Use `replace` or `append` tools to write `[[WikiLinks]]` directly into the `.md` files.

## 🛑 Critical Mandates for Agents
- **Always Verify:** Run `go build -o musu-crawl.exe main.go` before using CLI commands if you changed the code.
- **Privacy First:** Prefer the tool's local parsing/tagging logic over sending raw data to external APIs when possible.
- **Efficiency:** Use the worker pool (`-w` flag) for batch fetches to save time.
- **Structure:** Always ensure the YAML Frontmatter remains valid so the `index.json` generation doesn't break.
