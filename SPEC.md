# musu-crawl-ai Spec (STATUS: v0.8.0 AGENT-NATIVE)

## Goal
`musu-crawl-ai` is the Musu ecosystem's "Brain." It harvests source material, compiles it into a local wiki, and exposes search/research/fetch primitives to both humans and other agents.

## Current Product Truth

### Harvesting
- supported source families include `web`, `yt`, `gh`, `arxiv`, `reddit`, `hf`, and `x`
- harvested material lands in a local markdown wiki rooted at `./wiki` by default
- project scoping is handled under `wiki/projects/<project>`

### Research And Search
- `fetch` acquires raw source material
- `compile` links harvested markdown into a more useful wiki graph
- `search` supports local keyword and semantic retrieval
- `research` orchestrates a multi-step planning/search/harvest/analyze loop

### Preflight / Recovery
- `doctor` verifies:
  - wiki presence
  - search index presence
  - project directory presence
  - AI endpoint reachability
- `doctor --fix` can safely create a missing local wiki scaffold
- `doctor --capability-source` exposes static source capability metadata during setup
- `--json` provides deterministic machine-readable output

### Important Contract Detail
`--capability-source` is intentionally static metadata. It is not a live probe of credentials, network health, or external service reachability.

## Completed Milestones
- [x] multi-source harvesting
- [x] project-scoped wiki layout
- [x] JSON output mode
- [x] MCP server surface
- [x] doctor preflight with `--fix`
- [x] capability metadata output for setup-time automation
- [x] index and telemetry I/O errors are now surfaced (no silent index corruption, no silently-lost traces)
- [x] compiled `musu-crawl.exe` binary is no longer tracked in git

## Known Constraints
- `init` and `doctor` now probe the configured AI endpoint, but they still do not validate per-model compatibility beyond simple reachability
- source capability metadata is not a live readiness probe
- bulk crawl indexing still rewrites index/vector artifacts repeatedly and is not yet optimized for large N

## Next Work
1. Add optional live source probes for credentials or network reachability.
2. Separate static capability metadata from live doctor logic more cleanly.
3. Reduce repeated index/vector rewrite cost for large bulk crawls.

---
**Build Date:** 2026-05-27
**Status:** 🧠 BRAIN READY
