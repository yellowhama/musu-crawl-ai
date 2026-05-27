# Next Steps: musu-crawl-ai

## P1
- make `init` use the configured `--ai-url`
- keep `--capability-source` explicitly static unless a real live probe is implemented

## P2
- add optional live probes for source credentials/network health
- split doctor reporting helpers from static capability metadata

## P3
- optimize index/vector rewrite cost for larger bulk crawls
- improve ecosystem docs so the same entrypoint exists in all three repos
