#!/usr/bin/env bash
set -euo pipefail
(cd "$(dirname "$0")/.." && \
  GOWORK=off go run ./cmd/visualizer see --name "Go" >/dev/null && \
  GOWORK=off go run ./cmd/visualizer see --reverse=codepoints --name "U+0041 0x0042 67" >/dev/null && \
  GOWORK=off go run ./cmd/visualizer see --reverse=bytes --name "0xF0 0x9F 0x98 0x8A" >/dev/null)
echo "Visualizer sample commands completed successfully."
