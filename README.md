# musu-crawl-ai

> **The High-Bandwidth Data Supply Chain for AI Researcher Agents.**

`musu-crawl-ai` is a high-performance, production-hardened Go-based knowledge harvester and LLM Wiki generator. It empowers AI agents (and humans) to discover, fetch, organize, and synthesize knowledge into a structured, searchable, and interlinked "AI Brain" with a built-in **Researcher Mindset**.

---

## 🚀 Key Features

### 🧠 The Researcher Mindset (New in v0.6.0)
- **Socratic Planning:** The agent deconstructs queries into hypotheses and designs search queries to *falsify* assumptions.
- **Evidence Discrimination:** Every source is assigned a reliability score (Arxiv=0.9, Reddit=0.5), weighting the final analysis.
- **Cross-Verification:** The Analyst explicitly detects contradictions between sources and triggers "tie-breaker" research hops.
- **Recursive Multi-hop:** Fully autonomous loops that recursively fill information gaps identified by the Analyst.

### 📡 Universal Knowledge Harvesters
- **YouTube / Academic (Arxiv) / Code (GitHub) / Social (Twitter/Reddit) / General Web.**
- **Perfect Layout:** HTML-first Arxiv parsing and OCR fallback for scanned PDFs.
- **Multi-Modal:** Localized image harvesting and re-linking.

### 🦾 Production-Hardened Architecture
- **Thread-Safe Core:** sync.Mutex protected indexing.
- **Multi-Project Scoping:** Isolate missions in `wiki/projects/`.
- **Standalone UX:** `init`, `update`, and D3.js Galaxy Dashboard.

---

## 📖 User Manual (v0.6.0 Highlights)

### 1. Socratic Research
```bash
./musu-crawl research "Does nuclear fusion provide a net energy gain in 2024?" --depth 3
```
*The agent will set hypotheses, look for refuting evidence, and resolve contradictions between different news sources.*

### 2. Knowledge Galaxy
Visualize how the agent's skepticism has linked different viewpoints:
```bash
./musu-crawl serve
```

---

## 🤖 For AI Agents
AI Agents should refer to [**AGENTS.md**](./AGENTS.md) to understand how to leverage the new **Reliability Scores** and **Contradiction Detection** primitives.

---

## 📝 Roadmap
- [x] v0.5.0: Multi-Modal (Images & OCR)
- [x] v0.6.0: The Researcher Mindset (Hypothesis & Cross-Verification)
- [ ] v0.7.0: Local Vision Analysis (LLaVA Integration)
