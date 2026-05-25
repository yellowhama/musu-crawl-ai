# musu-crawl-ai

> **The High-Bandwidth Data Supply Chain for AI Researcher Agents.**

`musu-crawl-ai` is a high-performance, production-hardened Go-based knowledge harvester and LLM Wiki generator. It empowers AI agents (and humans) to discover, fetch, organize, and synthesize knowledge from across the web into a structured, searchable, and interlinked "AI Brain."

---

## 🚀 Key Features

### 📡 Universal Knowledge Harvesters
- **YouTube:** High-fidelity transcript extraction with Innertube fallback to bypass PO-Token/403 blocks.
- **Academic (Arxiv):** Perfect layout preservation using HTML-first parsing and **OCR fallback** for scanned documents.
- **Code (GitHub):** Repository metadata and documentation extraction.
- **Social (Twitter/Reddit):** Resilient post harvesting via Syndication and OEmbed bypasses.
- **General Web:** Noise-free content extraction using Readability algorithms.

### 🧠 Autonomous Intelligence (The Brain)
- **Recursive Research:** Multi-agent loop (Planner -> Searcher -> Harvester -> Analyst) that recursively fills information gaps.
- **Wiki Compiler:** Autonomous document cross-linking and relationship explanation using local LLM reasoning.
- **Hybrid Search:** Instant keyword (Bleve) and semantic vector (Ollama) search.
- **Multi-Modal Ready:** Automatically downloads and re-links images for future vision-based AI analysis.

### 🦾 Production-Hardened Architecture
- **Thread-Safe Core:** Robust locking mechanism to prevent data corruption during massive parallel fetches.
- **Multi-Project Scoping:** Isolate research missions into dedicated project silos (`wiki/projects/`).
- **Hierarchical Config:** Professional configuration precedence (Flags > Env > Project > Global).
- **Secure Secrets:** Built-in `auth` manager to keep API keys in project-specific hidden stores.

---

## 🛠️ Installation & Setup

### 1. Download
Get the latest binary (v0.5.1) from [GitHub Releases](https://github.com/yellowhama/musu-crawl-ai/releases).

### 2. Initialize
```bash
./musu-crawl init
```

### 3. External Dependencies (Optional but Recommended)
- **Intelligence:** Install [Ollama](https://ollama.com) to enable Research, Compilation, and Semantic Search.
- **OCR:** Install [Tesseract OCR](https://tesseract-ocr.github.io/tessdoc/Installation.html) (`winget install tesseract` on Windows) to extract text from scanned PDFs.

---

## 📖 User Manual

### 1. Project Management
Scope your work to maintain context and isolation:
```bash
./musu-crawl fetch web https://example.com --project my-new-mission
```

### 2. Autonomous Deep Research
```bash
./musu-crawl research "Explain the current state of Solid State Batteries" --project tech-review --depth 3
```

### 3. Knowledge Galaxy Dashboard
Visualize your interlinked knowledge base:
```bash
./musu-crawl serve --port 8080
# Visit http://localhost:8080/galaxy
```

### 4. Self-Maintenance
```bash
# Update to the latest version automatically
./musu-crawl update
```

---

## 📂 Directory Structure
- `/wiki/projects/{name}`: Project-specific knowledge silos.
- `/wiki/projects/{name}/images`: Locally harvested visual assets.
- `/wiki/index.json`: Global machine-readable knowledge map.
- `/wiki/musu.bleve`: Local search database.

---

## 🤖 For AI Agents
AI Agents "driving" this tool should follow the [**AGENTS.md**](./AGENTS.md) protocol to utilize knowledge primitives and steer their own research personas.

---

## 📝 Roadmap & Status
- [x] v0.3.0: AI Brain (Multi-Project & Galaxy)
- [x] v0.4.0: Configuration & Secrets
- [x] v0.5.0: Multi-Modal (Images & OCR)
- [x] v0.5.1: Production Hardening (Mutex & Encapsulation)
- [ ] v0.6.0: Web-based Agent Control Panel & Cloud Vector Sync
