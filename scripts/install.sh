#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
APP_NAME="${APP_NAME:-cloudcanal}"
BIN_PATH="$ROOT_DIR/bin/$APP_NAME"
INSTALL_BIN_DIR="${INSTALL_BIN_DIR:-$HOME/bin}"
INSTALL_PATH="$INSTALL_BIN_DIR/$APP_NAME"
INSTALL_SHELL_RC="${INSTALL_SHELL_RC:-$HOME/.zshrc}"
PATH_MARK_START="# >>> cloudcanal-openapi-cli >>>"
PATH_MARK_END="# <<< cloudcanal-openapi-cli <<<"

ensure_binary() {
  if [[ -x "$BIN_PATH" ]]; then
    return 0
  fi

  printf 'Binary not found, running all_build.sh first...\n'
  "$SCRIPT_DIR/all_build.sh"
}

ensure_path_block() {
  mkdir -p "$(dirname "$INSTALL_SHELL_RC")"
  touch "$INSTALL_SHELL_RC"

  if grep -Fq "$PATH_MARK_START" "$INSTALL_SHELL_RC"; then
    return 0
  fi

  {
    printf '\n%s\n' "$PATH_MARK_START"
    printf 'export PATH="%s:$PATH"\n' "$INSTALL_BIN_DIR"
    printf '%s\n' "$PATH_MARK_END"
  } >> "$INSTALL_SHELL_RC"
}

ensure_binary
mkdir -p "$INSTALL_BIN_DIR"
ln -sfn "$BIN_PATH" "$INSTALL_PATH"
ensure_path_block

printf 'Installed %s\n' "$INSTALL_PATH"
printf 'Source %s or open a new shell, then run: %s jobs list\n' "$INSTALL_SHELL_RC" "$APP_NAME"
