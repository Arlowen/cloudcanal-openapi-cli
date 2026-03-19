#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_DIR="$ROOT_DIR/bin"
BIN_PATH="$BIN_DIR/cloudcanal"
LOG_DIR="$(mktemp -d "${TMPDIR:-/tmp}/cloudcanal-openapi-cli-build.XXXXXX")"
START_TS="$(date +%s)"

cleanup() {
  rm -rf "$LOG_DIR"
}

trap cleanup EXIT

run_step() {
  local title="$1"
  local log_name="$2"
  shift 2

  local log_path="$LOG_DIR/$log_name.log"
  printf '[%s] %s\n' "$STEP_NO" "$title"

  if [[ "${VERBOSE:-0}" == "1" ]]; then
    "$@"
    return 0
  fi

  if "$@" >"$log_path" 2>&1; then
    return 0
  fi

  printf 'Failed: %s\n' "$title" >&2
  printf 'Log:\n' >&2
  cat "$log_path" >&2
  exit 1
}

cd "$ROOT_DIR"

printf 'CloudCanal OpenAPI CLI build\n'
printf 'Workspace: %s\n\n' "$ROOT_DIR"

STEP_NO="1/3"
printf '[%s] Clean build artifacts\n' "$STEP_NO"
if [[ -d "$BIN_DIR" ]]; then
  rm -rf "$BIN_DIR"
  printf 'Removed %s\n\n' "$BIN_DIR"
else
  printf 'Nothing to clean\n\n'
fi

STEP_NO="2/3"
run_step "Run tests" "test" go test ./...
printf 'Tests passed\n\n'

STEP_NO="3/3"
mkdir -p "$BIN_DIR"
run_step "Build CLI" "build" go build -o "$BIN_PATH" ./cmd/cloudcanal
printf 'Built %s\n\n' "$BIN_PATH"

ELAPSED="$(( $(date +%s) - START_TS ))"
printf 'Done in %ss\n' "$ELAPSED"
printf 'Tip: set VERBOSE=1 to show full command output\n'
