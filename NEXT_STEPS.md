# Next Steps: musu-crawl-ai

## P1
- keep `--capability-source` explicitly static unless a real live probe is implemented
- decide whether source readiness should ever become a live probe instead of documented static capability metadata

## P2
- add optional live probes for source credentials/network health
- split doctor reporting helpers from static capability metadata
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
