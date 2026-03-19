#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/log.sh"

APP_NAME="${APP_NAME:-cloudcanal}"
BIN_PATH="$ROOT_DIR/bin/$APP_NAME"
INSTALL_BIN_DIR="${INSTALL_BIN_DIR:-$HOME/bin}"
INSTALL_PATH="$INSTALL_BIN_DIR/$APP_NAME"
INSTALL_SHELL_RC="${INSTALL_SHELL_RC:-$HOME/.zshrc}"
INSTALL_ZSH_COMPLETION_DIR="${INSTALL_ZSH_COMPLETION_DIR:-$HOME/.zsh/completions}"
INSTALL_BASH_COMPLETION_DIR="${INSTALL_BASH_COMPLETION_DIR:-$HOME/.local/share/bash-completion/completions}"
ZSH_COMPLETION_PATH="$INSTALL_ZSH_COMPLETION_DIR/_$APP_NAME"
BASH_COMPLETION_PATH="$INSTALL_BASH_COMPLETION_DIR/$APP_NAME"
PATH_MARK_START="# >>> cloudcanal-openapi-cli >>>"
PATH_MARK_END="# <<< cloudcanal-openapi-cli <<<"
COMPLETION_MARK_START="# >>> cloudcanal-openapi-cli completion >>>"
COMPLETION_MARK_END="# <<< cloudcanal-openapi-cli completion <<<"

ensure_binary() {
  if [[ -x "$BIN_PATH" ]]; then
    log_success "Found binary at $BIN_PATH"
    return 0
  fi

  log_info "Binary not found, running all_build.sh first"
  "$SCRIPT_DIR/all_build.sh"
}

ensure_path_block() {
  mkdir -p "$(dirname "$INSTALL_SHELL_RC")"
  touch "$INSTALL_SHELL_RC"

  if grep -Fq "$PATH_MARK_START" "$INSTALL_SHELL_RC"; then
    log_success "PATH configuration already present in $INSTALL_SHELL_RC"
    return 0
  fi

  {
    printf '\n%s\n' "$PATH_MARK_START"
    printf 'export PATH="%s:$PATH"\n' "$INSTALL_BIN_DIR"
    printf '%s\n' "$PATH_MARK_END"
  } >> "$INSTALL_SHELL_RC"

  log_success "Updated $INSTALL_SHELL_RC"
}

ensure_completion_files() {
  mkdir -p "$INSTALL_ZSH_COMPLETION_DIR" "$INSTALL_BASH_COMPLETION_DIR"

  "$BIN_PATH" completion zsh "$APP_NAME" > "$ZSH_COMPLETION_PATH"
  log_success "Installed zsh completion to $ZSH_COMPLETION_PATH"

  "$BIN_PATH" completion bash "$APP_NAME" > "$BASH_COMPLETION_PATH"
  log_success "Installed bash completion to $BASH_COMPLETION_PATH"
}

ensure_completion_block() {
  mkdir -p "$(dirname "$INSTALL_SHELL_RC")"
  touch "$INSTALL_SHELL_RC"

  if grep -Fq "$COMPLETION_MARK_START" "$INSTALL_SHELL_RC"; then
    log_success "Shell completion configuration already present in $INSTALL_SHELL_RC"
    return 0
  fi

  {
    printf '\n%s\n' "$COMPLETION_MARK_START"
    printf 'if [[ -d "%s" ]]; then\n' "$INSTALL_ZSH_COMPLETION_DIR"
    printf '  fpath=("%s" $fpath)\n' "$INSTALL_ZSH_COMPLETION_DIR"
    printf '  autoload -Uz compinit\n'
    printf '  compinit\n'
    printf 'fi\n'
    printf '%s\n' "$COMPLETION_MARK_END"
  } >> "$INSTALL_SHELL_RC"

  log_success "Updated $INSTALL_SHELL_RC"
}

log_info "Install $APP_NAME command"
ensure_binary
mkdir -p "$INSTALL_BIN_DIR"
ln -sfn "$BIN_PATH" "$INSTALL_PATH"
log_success "Installed $INSTALL_PATH"
ensure_path_block
ensure_completion_files
ensure_completion_block

log_info "Open a new shell or source $INSTALL_SHELL_RC, then run: $APP_NAME jobs list"
print_run_summary "Install completed"
