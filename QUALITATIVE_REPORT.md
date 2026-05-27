# Qualitative Report: musu-crawl-ai

## Grade
`A`

## Why It Improved
- the tool is now much easier for agents to preflight and recover locally
- `doctor --json` and `doctor --fix` make setup failures legible
- source setup metadata is now labeled honestly as static capability data
- local fetch CLI smoke now proves a real `web` harvest path against a controllable HTTP source
- URL normalization is safer for `http://` targets on Windows filesystem paths
- a gated real-endpoint integration harness now exists for fetch + semantic search command verification
- the real-integration runner now auto-diagnoses missing local Ollama/OpenAI-compatible runtime candidates instead of failing silently
- the runner now emits machine-readable JSON diagnostics (`-Json -ProbeOnly`) for CI or agent handoff
- the JSON diagnostics now carry stable `issue_codes` so automation can distinguish bind-address misconfiguration from missing installs or timeouts
- the real integration path is now model-configurable through `MUSU_CRAWL_INTEGRATION_EMBED_MODEL`
- a real Ollama-backed `fetch web` + semantic `search` integration pass was verified with `nomic-embed-text`
- a command-flag state leak was caught under integration load and fixed in `search_test.go`
- the integration doctor now treats `model` and `model:latest` as equivalent, avoiding false negatives against Ollama

## Strong Points
- broad harvesting surface
- clear wiki-first architecture
- deterministic machine-readable output path
- strong position as the upstream knowledge source for the other two repos
- local command and orchestrator smoke coverage now prove bootstrap, fetch, and search surfaces without external services

## Concerns
- AI endpoint probing is now correctly flag-aware and model-aware, but it is still limited to runtime/model contract checks rather than full source-contract validation
- source capability output is static, not a live probe
- doctor logic is starting to accumulate multiple roles in one command file

## Thermo Verdict
`PASS WITH CONCERNS`

## Immediate Priorities
1. decide whether live source probes should exist at all
2. add a real research-command smoke path with a controllable local search/AI harness
3. split doctor reporting helpers from static capability metadata
