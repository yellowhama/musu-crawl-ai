# Project Handoff: musu-crawl-ai

## What This Repo Is
`musu-crawl-ai` is the Musu ecosystem's knowledge harvester and local wiki generator. It is the upstream source of truth for `musu-marketer` and `musu-nurikun`.

## Current Truth
- binary: `musu-crawl.exe`
- version constant: `v0.8.0`
- default wiki root: `./wiki`
- default AI contract: OpenAI-compatible endpoint at `--ai-url`
- key recovery path: `doctor --fix`

## What Changed In This Round
- added `doctor`
- added `doctor --fix`
- clarified static setup metadata via `--capability-source`
- added ecosystem workflow docs and machine-readable preflight guidance

## Operator Flow
1. `musu-crawl init --out ./wiki --project <name>`
2. `musu-crawl doctor --out ./wiki --project <name>`
3. `musu-crawl fetch <source> <id> --project <name>`
4. `musu-crawl index --out ./wiki --project <name>`
5. `musu-crawl research "<question>" --project <name>`

## Known Constraints
- `--capability-source` is static metadata, not a live reachability probe
- `init` still checks the default localhost AI port instead of respecting `--ai-url`
- very large bulk fetches still pay repeated index/vector rewrite cost

## Key Files
- `cmd/root.go`: global flags and JSON mode
- `cmd/init.go`: wiki/project scaffold bootstrap
- `cmd/doctor.go`: preflight and capability metadata output
- `internal/harvester/*`: source-specific fetchers
- `internal/processor/*`: wiki/index/vector processing
- `internal/agent/*`: orchestration and research logic
