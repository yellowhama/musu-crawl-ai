# Next Steps: musu-crawl-ai

## P1
- keep `--capability-source` explicitly static unless a real live probe is implemented
- decide whether source readiness should ever become a live probe instead of documented static capability metadata

## P2
- add a real research-command smoke path with a controllable local search/AI harness
- add optional live probes for source credentials/network health
- split doctor reporting helpers from static capability metadata

## P3
- optimize index/vector rewrite cost for larger bulk crawls
- improve ecosystem docs so the same entrypoint exists in all three repos

## Verified Integration Harness
- set `MUSU_CRAWL_INTEGRATION_AI_URL`
- optionally set `MUSU_CRAWL_INTEGRATION_EMBED_MODEL` (verified locally with `nomic-embed-text`)
- run `go test -tags integration ./cmd`
- or run `powershell -ExecutionPolicy Bypass -File .\scripts\run-real-integration.ps1`
