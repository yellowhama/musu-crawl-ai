# Next Steps: musu-crawl-ai

> See `C:\Users\empty\MUSU_MCP_AUDIT_2026-05-28.md` for the 2026-05-28 real-usage audit that produced the MCP-related items below.

## P1
- keep `--capability-source` explicitly static unless a real live probe is implemented
- decide whether source readiness should ever become a live probe instead of documented static capability metadata
- declare MCP tool parameter schemas in `cmd/mcp.go` for `fetch` / `search` / `research` (currently `WithDescription` only → empty JSON schema → MCP clients cannot pass args); mirror the `WithString`/`WithNumber`/`Required` pattern from `mcp-go`
- guard empty input in `handleResearch` (currently a no-arg call returns silently with no output instead of erroring on missing `question`)

## P2
- add optional live probes for source credentials/network health
- split doctor reporting helpers from static capability metadata
- in `handleSearch`, distinguish "no index" from "0 matches" so an unbuilt index does not look like an empty corpus
- add `json:"…"` tags to `preflight.DoctorResult` (Report/Blocking/ActionableFix) so the MCP envelope is snake_case-consistent with the inner Report
- decide whether the local search harness seam (`MUSU_SEARCH_BASE_URL`) should stay integration-focused or become a documented alternative provider hook

## P3
- optimize index/vector rewrite cost for larger bulk crawls
- improve ecosystem docs so the same entrypoint exists in all three repos
- extract shared module(s) for `AgentClient` + `preflight/doctor` + env-loader to remove triple-duplicated logic across the three repos

## Verified Integration Harness
- set `MUSU_CRAWL_INTEGRATION_AI_URL`
- optionally set `MUSU_CRAWL_INTEGRATION_EMBED_MODEL` (verified locally with `nomic-embed-text`)
- optionally set `MUSU_CRAWL_INTEGRATION_CHAT_MODEL` (verified locally with `llama3.2:1b`)
- run `go test -tags integration ./cmd`
- or run `powershell -ExecutionPolicy Bypass -File .\scripts\run-real-integration.ps1`
