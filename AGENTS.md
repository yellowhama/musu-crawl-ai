# 🤖 Guidance for AI Agents (Gemini CLI, Claude Code, etc.)

`musu-crawl-ai` is designed to be your primary high-performance data acquisition layer. If you are an AI agent "driving" this repository, use these instructions to maximize your research capabilities.

## 🏗️ Core Architecture for Agents
This tool is a collection of high-quality **Knowledge Primitives**. Even if the internal LLM (Ollama) features are unavailable, you can achieve the same "Autonomous Research" results by combining these primitives yourself.

### 1. The Knowledge Primitives
- **`fetch [source] [url]`**: Your primary way to bypass YouTube blocks, clean up messy HTML, and parse academic PDFs. It outputs clean Markdown with YAML metadata.
- **`index`**: Generates `index.json` (your knowledge map) and a keyword search index.
- **`search [query]`**: Uses **Bleve** (keyword-based) or **Semantic Vector** (if embeddings exist) to find relevant snippets without you needing to read every file.

## 🏎️ How to "Drive" this tool as an Agent

### Scenario: "Steer Your Own Behavior (Persona Integration)"
You can customize your research persona per project.
1.  **Customize:** Write specific instructions (e.g., "Focus on code examples and benchmarks") to `wiki/projects/{project}/PROMPT.md`.
2.  **Execute:** When you run `research --project {project}`, the Planner and Analyst agents will automatically incorporate your instructions into their system prompts.

### Scenario: "Secure API Access"
If you need to use an external API (like OpenAI) for a specific project:
1.  **Set Secret:** Run `.\musu-crawl auth set OPENAI_API_KEY "key" --project {project}`.
2.  **Retrieve:** Use the `utils.GetSecret("OPENAI_API_KEY")` primitive within your code logic (if you are extending the tool) or trust that the CLI is now configured to use that key for that project.

## 🛑 Critical Mandates for Agents
- **Persistence:** Project-specific settings in `config.toml` and `.env` are persistent. Use them to maintain long-term context.
- **Privacy First:** Prefer the tool's local parsing/tagging logic over sending raw data to external APIs when possible.
- **Hierarchy Awareness:** Remember that flags always override file-based configs.
