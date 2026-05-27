# musu-crawl-ai

> **The High-Bandwidth Data Supply Chain for AI Agents.**

`musu-crawl-ai` is a high-performance, agent-native knowledge harvester and LLM Wiki generator. It is the "Eye and Brain" of the Musu ecosystem, responsible for fetching, cleaning, and organizing world-wide knowledge into a structured, interlinked "AI Brain."

---

## 🚀 Key Features

### 🤖 Agent-Native Interface (v0.8.0)
- **MCP Server:** Native integration with Claude Desktop and Cursor. Use Musu as a built-in tool.
- **Machine-Readable:** Global `--json` mode for deterministic, noise-free output.
- **Agentic Recovery:** Error messages include `agent_actionable_fix` tips to guide LLMs.

### 🧠 Intelligence & Research
- **Researcher Mindset:** Socratic planning, hypothesis testing, and contradiction detection.
- **Universal Harvesters:** YouTube (Transcript), Arxiv (HTML-first + OCR), GitHub, Reddit, and Web.
- **Multi-Modal:** Automatic image harvesting and local description (via LLaVA).
- **Live Sync:** O(1) incremental indexing for instant search availability.

### 🦾 Production-Grade Foundation
- **Thread-Safety:** Mutex-protected indexing for high-concurrency parallel crawls.
- **Secure Secrets:** Project-scoped credential management.
- **Visual Galaxy:** Interactive D3.js dashboard for knowledge mapping.

---

## 🛠️ Installation & Setup

### 1. Prerequisite: [Ollama](https://ollama.com)
Local LLM intelligence is required for Research, Vision, and Semantic Search.

### 2. Quick Start
```bash
./musu-crawl init
./musu-crawl fetch web https://go.dev/blog/go1.22 --project my-research
```

### 3. MCP Integration (For Claude/Cursor)
Add the following to your `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "musu": {
      "command": "C:/path/to/musu-crawl-ai/musu-crawl.exe",
      "args": ["mcp"]
    }
  }
}
```

---

## 📖 Core Commands
- `fetch [source] [id]`: Harvest specific content.
- `research "[question]"`: Autonomous multi-agent deep research mission.
- `search "[query]"`: Local keyword and semantic vector search.
- `index --semantic`: Refresh the global knowledge graph and embeddings.
- `serve`: Launch the Galaxy Dashboard (Port 8080).

---

## 📂 Data Silos
Private research data is stored in the `wiki/` directory and is strictly excluded from Git tracking to ensure privacy.

---

## 🔗 The Ecosystem
- **musu-marketer:** The "Voice" that uses this knowledge for strategy.
- **musu-nurikun:** The "Hand" that creates identities and acts on the web.
