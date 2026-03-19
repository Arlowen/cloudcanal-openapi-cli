#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
source "$SCRIPT_DIR/lib/log.sh"

default_shell_rc() {
  case "$(basename "${SHELL:-}")" in
    bash) printf '%s\n' "$HOME/.bashrc" ;;
    *) printf '%s\n' "$HOME/.zshrc" ;;
  esac
}

APP_NAME="${APP_NAME:-cloudcanal}"
BIN_PATH="$ROOT_DIR/bin/$APP_NAME"
INSTALL_ROOT="${INSTALL_ROOT:-$HOME/.cloudcanal-cli}"
INSTALL_BIN_DIR="${INSTALL_BIN_DIR:-$INSTALL_ROOT/bin}"
INSTALL_PATH="$INSTALL_BIN_DIR/$APP_NAME"
INSTALL_SHELL_RC="${INSTALL_SHELL_RC:-$(default_shell_rc)}"
INSTALL_COMPLETION_DIR="${INSTALL_COMPLETION_DIR:-$INSTALL_ROOT/completions}"
INSTALL_ZSH_COMPLETION_DIR="${INSTALL_ZSH_COMPLETION_DIR:-$INSTALL_COMPLETION_DIR/zsh}"
INSTALL_BASH_COMPLETION_DIR="${INSTALL_BASH_COMPLETION_DIR:-$INSTALL_COMPLETION_DIR/bash}"
ZSH_COMPLETION_PATH="$INSTALL_ZSH_COMPLETION_DIR/_$APP_NAME"
BASH_COMPLETION_PATH="$INSTALL_BASH_COMPLETION_DIR/$APP_NAME"
PATH_MARK_START="# >>> cloudcanal-openapi-cli >>>"
PATH_MARK_END="# <<< cloudcanal-openapi-cli <<<"
COMPLETION_MARK_START="# >>> cloudcanal-openapi-cli completion >>>"
COMPLETION_MARK_END="# <<< cloudcanal-openapi-cli completion <<<"

LEGACY_RELEASE_ROOT="${LEGACY_RELEASE_ROOT:-$HOME/.local/share/cloudcanal-openapi-cli}"
LEGACY_INSTALL_BIN_DIR="${LEGACY_INSTALL_BIN_DIR:-$HOME/bin}"
LEGACY_INSTALL_PATH="$LEGACY_INSTALL_BIN_DIR/$APP_NAME"
LEGACY_RELEASE_BIN_PATH="$LEGACY_RELEASE_ROOT/bin/$APP_NAME"
LEGACY_INSTALL_ZSH_COMPLETION_DIR="${LEGACY_INSTALL_ZSH_COMPLETION_DIR:-$HOME/.zsh/completions}"
LEGACY_INSTALL_BASH_COMPLETION_DIR="${LEGACY_INSTALL_BASH_COMPLETION_DIR:-$HOME/.local/share/bash-completion/completions}"
LEGACY_ZSH_COMPLETION_PATH="$LEGACY_INSTALL_ZSH_COMPLETION_DIR/_$APP_NAME"
LEGACY_BASH_COMPLETION_PATH="$LEGACY_INSTALL_BASH_COMPLETION_DIR/$APP_NAME"

remove_rc_block() {
  local start_mark="$1"
  local end_mark="$2"
  local description="$3"

  if [[ ! -f "$INSTALL_SHELL_RC" ]] || ! grep -Fq "$start_mark" "$INSTALL_SHELL_RC"; then
    log_info "No $description to remove from $INSTALL_SHELL_RC"
    return 0
  fi

  local tmp_file
  tmp_file="$(mktemp)"

  awk -v start="$start_mark" -v end="$end_mark" '
    $0 == start {skip = 1; next}
    $0 == end {skip = 0; next}
    !skip {print}
  ' "$INSTALL_SHELL_RC" > "$tmp_file"

  mv "$tmp_file" "$INSTALL_SHELL_RC"
  log_success "Updated $INSTALL_SHELL_RC"
}

remove_if_empty_dir() {
  local dir="$1"
  if [[ -d "$dir" ]] && [[ -z "$(ls -A "$dir")" ]]; then
    rmdir "$dir"
  fi
}

prune_install_dirs() {
  remove_if_empty_dir "$INSTALL_ZSH_COMPLETION_DIR"
  remove_if_empty_dir "$INSTALL_BASH_COMPLETION_DIR"
  remove_if_empty_dir "$INSTALL_COMPLETION_DIR"
  remove_if_empty_dir "$INSTALL_BIN_DIR"
  remove_if_empty_dir "$INSTALL_ROOT"
}

remove_link() {
  if [[ -L "$INSTALL_PATH" ]]; then
    local target
    target="$(readlink "$INSTALL_PATH")"
    if [[ "$target" == "$BIN_PATH" ]]; then
      rm -f "$INSTALL_PATH"
      log_success "Removed $INSTALL_PATH"
      return 0
    fi
    log_info "Skipped $INSTALL_PATH because it is not managed by this project"
    return 0
  fi

  if [[ -e "$INSTALL_PATH" ]]; then
    log_info "Skipped $INSTALL_PATH because it is not a symlink created by this project"
    return 0
  fi

  log_info "No installed command found at $INSTALL_PATH"
}

remove_completion_files() {
  if [[ -f "$ZSH_COMPLETION_PATH" ]]; then
    rm -f "$ZSH_COMPLETION_PATH"
    log_success "Removed $ZSH_COMPLETION_PATH"
  else
    log_info "No zsh completion file found at $ZSH_COMPLETION_PATH"
  fi

  if [[ -f "$BASH_COMPLETION_PATH" ]]; then
    rm -f "$BASH_COMPLETION_PATH"
    log_success "Removed $BASH_COMPLETION_PATH"
  else
    log_info "No bash completion file found at $BASH_COMPLETION_PATH"
  fi
}

remove_legacy_install() {
  if [[ -L "$LEGACY_INSTALL_PATH" ]]; then
    local target
    target="$(readlink "$LEGACY_INSTALL_PATH")"
    if [[ "$target" == "$BIN_PATH" || "$target" == "$LEGACY_RELEASE_BIN_PATH" ]]; then
      rm -f "$LEGACY_INSTALL_PATH"
      log_success "Removed legacy managed symlink $LEGACY_INSTALL_PATH"
    else
      log_info "Skipped legacy command link $LEGACY_INSTALL_PATH because it is not managed by this project"
    fi
  elif [[ -e "$LEGACY_INSTALL_PATH" ]]; then
    log_info "Skipped legacy command link $LEGACY_INSTALL_PATH because it is not a managed symlink"
  else
    log_info "No legacy installed command found at $LEGACY_INSTALL_PATH"
  fi

  if [[ -f "$LEGACY_ZSH_COMPLETION_PATH" ]]; then
    rm -f "$LEGACY_ZSH_COMPLETION_PATH"
    log_success "Removed legacy zsh completion $LEGACY_ZSH_COMPLETION_PATH"
  fi

  if [[ -f "$LEGACY_BASH_COMPLETION_PATH" ]]; then
    rm -f "$LEGACY_BASH_COMPLETION_PATH"
    log_success "Removed legacy bash completion $LEGACY_BASH_COMPLETION_PATH"
  fi

  if [[ -d "$LEGACY_RELEASE_ROOT" ]]; then
    rm -rf "$LEGACY_RELEASE_ROOT"
    log_success "Removed legacy release install root $LEGACY_RELEASE_ROOT"
  fi
}

log_info "Uninstall $APP_NAME command"
remove_link
remove_rc_block "$PATH_MARK_START" "$PATH_MARK_END" "PATH configuration"
remove_completion_files
remove_rc_block "$COMPLETION_MARK_START" "$COMPLETION_MARK_END" "shell completion configuration"
prune_install_dirs
remove_legacy_install

log_info "Open a new shell or source $INSTALL_SHELL_RC to refresh PATH"
print_run_summary "Uninstall completed"
