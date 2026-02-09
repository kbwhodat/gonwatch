#!/usr/bin/env bash

set -euo pipefail

BINARY_NAME="gonwatch"
VENV_DIR="${HOME}/.local/share/gonwatch/venv"
DATA_DIR="${HOME}/.local/share/gonwatch"

ASSUME_YES=false
for arg in "$@"; do
    case "$arg" in
        -y|--yes)
            ASSUME_YES=true
            ;;
    esac
done

log_info() { printf "[INFO] %s\n" "$1"; }
log_warn() { printf "[WARN] %s\n" "$1"; }
log_success() { printf "[OK] %s\n" "$1"; }

confirm() {
    local prompt="$1"
    if $ASSUME_YES; then
        return 0
    fi
    printf "%s [y/N]: " "$prompt"
    read -r reply
    [[ "$reply" =~ ^[Yy]$ ]]
}

remove_path() {
    local target="$1"
    if [[ -e "$target" ]]; then
        if [[ -w "$target" ]]; then
            rm -rf "$target"
        else
            sudo rm -rf "$target"
        fi
        log_success "Removed $target"
    fi
}

is_nix_managed() {
    local path="$1"
    [[ "$path" == /nix/* || "$path" == /run/current-system/* || "$path" == /etc/profiles/* ]]
}

targets=()
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    targets+=("$(command -v "$BINARY_NAME")")
fi
targets+=(
    "/usr/local/bin/${BINARY_NAME}"
    "${HOME}/.local/bin/${BINARY_NAME}"
    "${HOME}/bin/${BINARY_NAME}"
)

if ! confirm "This will remove $BINARY_NAME, its venv, and shell completions."; then
    log_info "Uninstall cancelled"
    exit 0
fi

for target in "${targets[@]}"; do
    if [[ -z "$target" || ! -e "$target" ]]; then
        continue
    fi
    if is_nix_managed "$target"; then
        log_warn "$target looks Nix-managed. Remove via nix profile or nixos-rebuild."
        continue
    fi
    remove_path "$target"
done

remove_path "/usr/local/share/gonwatch/scripts"
remove_path "${HOME}/.local/share/gonwatch/scripts"

remove_path "$VENV_DIR"

if [[ -d "$DATA_DIR" ]]; then
    if [[ -z "$(ls -A "$DATA_DIR" 2>/dev/null)" ]]; then
        rmdir "$DATA_DIR"
        log_success "Removed $DATA_DIR"
    fi
fi

if command -v brew >/dev/null 2>&1; then
    brew_prefix=$(brew --prefix 2>/dev/null || true)
    if [[ -n "$brew_prefix" ]]; then
        remove_path "$brew_prefix/etc/bash_completion.d/$BINARY_NAME"
    fi
fi

bash_completion_dir="${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions"
remove_path "$bash_completion_dir/$BINARY_NAME"

zsh_completion_dir="${ZDOTDIR:-$HOME}/.zfunc"
remove_path "$zsh_completion_dir/_$BINARY_NAME"

fish_completion_dir="${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions"
remove_path "$fish_completion_dir/$BINARY_NAME.fish"

log_success "Uninstall complete"
