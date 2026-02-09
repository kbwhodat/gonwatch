#!/usr/bin/env bash
#
# Gonwatch Universal Installer
# Cross-platform installation script for Linux, macOS, and Windows (via Git Bash/WSL/MSYS2)
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/kbwhodat/gonwatch/main/install.sh | bash
#   wget -qO- https://raw.githubusercontent.com/kbwhodat/gonwatch/main/install.sh | bash
#
# Or locally:
#   chmod +x install.sh && ./install.sh
#

set -euo pipefail

# ============================================================================
# Configuration
# ============================================================================

REPO_OWNER="${REPO_OWNER:-kbwhodat}"
REPO_NAME="gonwatch"
BINARY_NAME="gonwatch"
GITHUB_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}"
RELEASES_URL="${GITHUB_URL}/releases"
MIN_GO_VERSION="1.22"

# Python packages required for scraping
PYTHON_PACKAGES=(
    "selenium-driverless>=1.9.4"
    "requests"
    "numpy"
    "matplotlib"
    "scipy"
    "platformdirs"
    "jsondiff"
    "orjson"
    "langdetect"
    "beautifulsoup4"
    "tls_client"
    "selenium"
    "cdp-socket"
    "websockets"
    "aiofiles"
    "aiohttp"
)

# Colors (disabled if not a TTY)
if [[ -t 1 ]]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    MAGENTA='\033[0;35m'
    CYAN='\033[0;36m'
    BOLD='\033[1m'
    NC='\033[0m' # No Color
else
    RED='' GREEN='' YELLOW='' BLUE='' MAGENTA='' CYAN='' BOLD='' NC=''
fi

# ============================================================================
# Helper Functions
# ============================================================================

log_info()    { printf "${BLUE}[INFO]${NC} %s\n" "$1"; }
log_success() { printf "${GREEN}[OK]${NC} %s\n" "$1"; }
log_warn()    { printf "${YELLOW}[WARN]${NC} %s\n" "$1"; }
log_error()   { printf "${RED}[ERROR]${NC} %s\n" "$1" >&2; }
log_step()    { printf "${MAGENTA}[STEP]${NC} ${BOLD}%s${NC}\n" "$1"; }

die() {
    log_error "$1"
    exit 1
}

confirm() {
    local prompt="${1:-Continue?}"
    local default="${2:-y}"
    
    if [[ "$default" == "y" ]]; then
        prompt="$prompt [Y/n]: "
    else
        prompt="$prompt [y/N]: "
    fi
    
    printf "${CYAN}%s${NC}" "$prompt"
    read -r response
    response="${response:-$default}"
    
    [[ "$response" =~ ^[Yy]$ ]]
}

command_exists() {
    command -v "$1" &>/dev/null
}

version_gte() {
    # Returns 0 if $1 >= $2 (semantic version comparison)
    printf '%s\n%s\n' "$2" "$1" | sort -V -C
}

# ============================================================================
# OS/Arch Detection
# ============================================================================

detect_os() {
    local os
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          os="unknown" ;;
    esac
    echo "$os"
}

detect_arch() {
    local arch
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        armv7l|armhf)   arch="arm" ;;
        i386|i686)      arch="386" ;;
        *)              arch="unknown" ;;
    esac
    echo "$arch"
}

detect_linux_distro() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        echo "${ID:-unknown}"
    elif [[ -f /etc/debian_version ]]; then
        echo "debian"
    elif [[ -f /etc/redhat-release ]]; then
        echo "rhel"
    elif [[ -f /etc/arch-release ]]; then
        echo "arch"
    elif [[ -f /etc/alpine-release ]]; then
        echo "alpine"
    else
        echo "unknown"
    fi
}

detect_package_manager() {
    local os="$1"
    
    if [[ "$os" == "darwin" ]]; then
        if command_exists brew; then
            echo "brew"
        else
            echo "none"
        fi
        return
    fi
    
    if [[ "$os" == "windows" ]]; then
        if command_exists choco; then
            echo "choco"
        elif command_exists scoop; then
            echo "scoop"
        elif command_exists winget; then
            echo "winget"
        else
            echo "none"
        fi
        return
    fi
    
    # Linux package managers
    if command_exists apt-get; then
        echo "apt"
    elif command_exists dnf; then
        echo "dnf"
    elif command_exists yum; then
        echo "yum"
    elif command_exists pacman; then
        echo "pacman"
    elif command_exists apk; then
        echo "apk"
    elif command_exists zypper; then
        echo "zypper"
    elif command_exists emerge; then
        echo "emerge"
    elif command_exists nix-env; then
        echo "nix"
    elif command_exists xbps-install; then
        echo "xbps"
    else
        echo "none"
    fi
}

# ============================================================================
# Dependency Installation
# ============================================================================

install_go() {
    local os="$1"
    local pkg_manager="$2"
    
    log_step "Installing Go..."
    
    case "$pkg_manager" in
        apt)
            sudo apt-get update
            sudo apt-get install -y golang-go || {
                # Fallback: install from official tarball for newer version
                install_go_from_tarball "$os"
            }
            ;;
        dnf|yum)
            sudo $pkg_manager install -y golang
            ;;
        pacman)
            sudo pacman -Sy --noconfirm go
            ;;
        apk)
            sudo apk add --no-cache go
            ;;
        zypper)
            sudo zypper install -y go
            ;;
        brew)
            brew install go
            ;;
        choco)
            choco install golang -y
            ;;
        scoop)
            scoop install go
            ;;
        winget)
            winget install -e --id GoLang.Go
            ;;
        nix)
            nix-env -iA nixpkgs.go
            ;;
        xbps)
            sudo xbps-install -Sy go
            ;;
        *)
            install_go_from_tarball "$os"
            ;;
    esac
}

install_go_from_tarball() {
    local os="$1"
    local arch
    arch=$(detect_arch)
    
    log_info "Installing Go from official tarball..."
    
    # Get latest Go version
    local go_version
    go_version=$(curl -fsSL "https://go.dev/VERSION?m=text" | head -1)
    
    if [[ "$os" == "windows" ]]; then
        local go_url="https://go.dev/dl/${go_version}.windows-${arch}.zip"
        local install_dir="/c/Go"
        
        log_info "Downloading $go_url"
        curl -fsSL "$go_url" -o /tmp/go.zip
        unzip -q /tmp/go.zip -d /c/
        rm /tmp/go.zip
        
        log_warn "Add C:\\Go\\bin to your PATH manually"
    else
        local go_url="https://go.dev/dl/${go_version}.${os}-${arch}.tar.gz"
        local install_dir="/usr/local"
        
        log_info "Downloading $go_url"
        curl -fsSL "$go_url" | sudo tar -C "$install_dir" -xzf -
        
        # Add to PATH for current session
        export PATH="$PATH:/usr/local/go/bin"
        
        # Suggest adding to shell profile
        log_info "Add the following to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo '  export PATH="$PATH:/usr/local/go/bin"'
    fi
    
    log_success "Go ${go_version} installed successfully"
}

install_mpv() {
    local os="$1"
    local pkg_manager="$2"
    
    log_step "Installing mpv (video player)..."
    
    case "$pkg_manager" in
        apt)
            sudo apt-get update
            sudo apt-get install -y mpv
            ;;
        dnf|yum)
            # Enable RPM Fusion for Fedora/RHEL
            if ! rpm -q rpmfusion-free-release &>/dev/null; then
                log_info "Enabling RPM Fusion repository..."
                if command_exists dnf; then
                    sudo dnf install -y "https://mirrors.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm" 2>/dev/null || true
                fi
            fi
            sudo $pkg_manager install -y mpv
            ;;
        pacman)
            sudo pacman -Sy --noconfirm mpv
            ;;
        apk)
            sudo apk add --no-cache mpv
            ;;
        zypper)
            sudo zypper install -y mpv
            ;;
        brew)
            brew install mpv
            ;;
        choco)
            choco install mpv -y
            ;;
        scoop)
            scoop bucket add extras
            scoop install mpv
            ;;
        winget)
            winget install -e --id mpv.net
            ;;
        nix)
            nix-env -iA nixpkgs.mpv
            ;;
        xbps)
            sudo xbps-install -Sy mpv
            ;;
        emerge)
            sudo emerge --ask=n media-video/mpv
            ;;
        *)
            log_warn "Could not auto-install mpv. Please install it manually."
            log_info "mpv is required for video playback."
            return 1
            ;;
    esac
    
    log_success "mpv installed successfully"
}

install_python() {
    local os="$1"
    local pkg_manager="$2"
    
    log_step "Installing Python 3..."
    
    case "$pkg_manager" in
        apt)
            sudo apt-get update
            sudo apt-get install -y python3 python3-pip python3-venv
            ;;
        dnf|yum)
            sudo $pkg_manager install -y python3 python3-pip
            ;;
        pacman)
            sudo pacman -Sy --noconfirm python python-pip
            ;;
        apk)
            sudo apk add --no-cache python3 py3-pip
            ;;
        zypper)
            sudo zypper install -y python3 python3-pip
            ;;
        brew)
            brew install python3
            ;;
        choco)
            choco install python -y
            ;;
        scoop)
            scoop install python
            ;;
        winget)
            winget install -e --id Python.Python.3.12
            ;;
        nix)
            nix-env -iA nixpkgs.python3 nixpkgs.python3Packages.pip
            ;;
        xbps)
            sudo xbps-install -Sy python3 python3-pip
            ;;
        emerge)
            sudo emerge --ask=n dev-lang/python
            ;;
        *)
            log_warn "Could not auto-install Python. Please install Python 3 manually."
            return 1
            ;;
    esac
    
    log_success "Python 3 installed successfully"
}

install_python_packages() {
    log_step "Installing Python packages..."

    local venv_dir="${HOME}/.local/share/gonwatch/venv"
    local venv_python="${venv_dir}/bin/python"

    mkdir -p "${HOME}/.local/share/gonwatch"

    if [[ ! -x "$venv_python" ]]; then
        log_info "Creating venv at ${venv_dir}"
        if command_exists python3; then
            python3 -m venv "$venv_dir" || return 1
        elif command_exists python; then
            python -m venv "$venv_dir" || return 1
        else
            log_warn "Python not found for venv creation"
            return 1
        fi
    fi

    log_info "Upgrading pip in venv"
    "$venv_python" -m pip install --upgrade pip >/dev/null 2>&1 || true

    # Install packages
    log_info "Installing: ${PYTHON_PACKAGES[*]}"

    if "$venv_python" -m pip install "${PYTHON_PACKAGES[@]}" 2>/dev/null; then
        log_success "Python packages installed successfully"
        log_info "Venv location: ${venv_dir}"
    else
        log_warn "Could not install Python packages automatically."
        log_info "Please run manually: ${venv_python} -m pip install ${PYTHON_PACKAGES[*]}"
        return 1
    fi
}

install_nodejs() {
    local os="$1"
    local pkg_manager="$2"
    
    log_step "Installing Node.js (for AnimePahe support)..."
    
    case "$pkg_manager" in
        apt)
            sudo apt-get update
            sudo apt-get install -y nodejs npm || {
                # Fallback to NodeSource
                log_info "Installing from NodeSource..."
                curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
                sudo apt-get install -y nodejs
            }
            ;;
        dnf|yum)
            sudo $pkg_manager install -y nodejs npm
            ;;
        pacman)
            sudo pacman -Sy --noconfirm nodejs npm
            ;;
        apk)
            sudo apk add --no-cache nodejs npm
            ;;
        zypper)
            sudo zypper install -y nodejs npm
            ;;
        brew)
            brew install node
            ;;
        choco)
            choco install nodejs -y
            ;;
        scoop)
            scoop install nodejs
            ;;
        winget)
            winget install -e --id OpenJS.NodeJS.LTS
            ;;
        nix)
            nix-env -iA nixpkgs.nodejs
            ;;
        xbps)
            sudo xbps-install -Sy nodejs
            ;;
        emerge)
            sudo emerge --ask=n net-libs/nodejs
            ;;
        *)
            log_warn "Could not auto-install Node.js."
            log_info "Node.js is required for AnimePahe streaming."
            log_info "Install from: https://nodejs.org/"
            return 1
            ;;
    esac
    
    log_success "Node.js installed successfully"
}

install_chromium() {
    local os="$1"
    local pkg_manager="$2"
    
    log_step "Installing Chromium (for headless browsing)..."
    
    # Check if Chrome or Chromium already exists
    if command_exists google-chrome || command_exists google-chrome-stable || \
       command_exists chromium || command_exists chromium-browser; then
        log_success "Chrome/Chromium already installed"
        return 0
    fi
    
    case "$pkg_manager" in
        apt)
            sudo apt-get update
            sudo apt-get install -y chromium-browser || sudo apt-get install -y chromium
            ;;
        dnf)
            sudo dnf install -y chromium
            ;;
        yum)
            sudo yum install -y chromium
            ;;
        pacman)
            sudo pacman -Sy --noconfirm chromium
            ;;
        apk)
            sudo apk add --no-cache chromium
            ;;
        zypper)
            sudo zypper install -y chromium
            ;;
        brew)
            brew install --cask chromium || brew install chromium
            ;;
        choco)
            choco install chromium -y
            ;;
        scoop)
            scoop bucket add extras
            scoop install chromium
            ;;
        winget)
            winget install -e --id Chromium.Chromium
            ;;
        nix)
            nix-env -iA nixpkgs.chromium
            ;;
        xbps)
            sudo xbps-install -Sy chromium
            ;;
        emerge)
            sudo emerge --ask=n www-client/chromium
            ;;
        *)
            log_warn "Could not auto-install Chromium."
            log_info "Chrome or Chromium is required for video source scraping."
            log_info "Install Google Chrome or Chromium manually."
            return 1
            ;;
    esac
    
    log_success "Chromium installed successfully"
}

install_git() {
    local pkg_manager="$1"
    
    log_step "Installing git..."
    
    case "$pkg_manager" in
        apt)
            sudo apt-get update
            sudo apt-get install -y git
            ;;
        dnf|yum)
            sudo $pkg_manager install -y git
            ;;
        pacman)
            sudo pacman -Sy --noconfirm git
            ;;
        apk)
            sudo apk add --no-cache git
            ;;
        zypper)
            sudo zypper install -y git
            ;;
        brew)
            brew install git
            ;;
        choco)
            choco install git -y
            ;;
        scoop)
            scoop install git
            ;;
        winget)
            winget install -e --id Git.Git
            ;;
        nix)
            nix-env -iA nixpkgs.git
            ;;
        xbps)
            sudo xbps-install -Sy git
            ;;
        *)
            die "Could not install git. Please install it manually."
            ;;
    esac
    
    log_success "git installed successfully"
}

# ============================================================================
# Installation Methods
# ============================================================================

get_install_dir() {
    local os="$1"
    
    case "$os" in
        linux|darwin)
            if [[ -w "/usr/local/bin" ]]; then
                echo "/usr/local/bin"
            elif [[ -d "$HOME/.local/bin" ]]; then
                echo "$HOME/.local/bin"
            else
                mkdir -p "$HOME/.local/bin"
                echo "$HOME/.local/bin"
            fi
            ;;
        windows)
            # For Windows, use user's local bin or Program Files
            if [[ -d "$HOME/bin" ]]; then
                echo "$HOME/bin"
            else
                mkdir -p "$HOME/bin"
                echo "$HOME/bin"
            fi
            ;;
    esac
}

install_from_release() {
    local os="$1"
    local arch="$2"
    local install_dir="$3"
    
    log_step "Installing from GitHub releases..."
    
    # Get latest release
    local release_info
    release_info=$(curl -fsSL "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null) || {
        log_warn "No releases found. Falling back to building from source."
        return 1
    }
    
    local version
    version=$(echo "$release_info" | grep -o '"tag_name": *"[^"]*"' | head -1 | cut -d'"' -f4)
    
    if [[ -z "$version" ]]; then
        log_warn "Could not determine latest version. Falling back to building from source."
        return 1
    fi
    
    log_info "Latest version: $version"
    
    # Construct expected asset name
    local ext="tar.gz"
    [[ "$os" == "windows" ]] && ext="zip"
    
    local asset_name="${BINARY_NAME}_${version#v}_${os}_${arch}.${ext}"
    local download_url="${RELEASES_URL}/download/${version}/${asset_name}"
    
    log_info "Downloading $download_url"
    
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT
    
    if ! curl -fsSL "$download_url" -o "$tmp_dir/archive.${ext}"; then
        log_warn "Binary release not found for ${os}/${arch}. Falling back to building from source."
        return 1
    fi
    
    # Extract
    cd "$tmp_dir"
    if [[ "$ext" == "zip" ]]; then
        unzip -q "archive.zip"
    else
        tar -xzf "archive.tar.gz"
    fi
    
    # Find and install binary
    local binary_file
    binary_file=$(find . -name "$BINARY_NAME" -o -name "${BINARY_NAME}.exe" | head -1)
    
    if [[ -z "$binary_file" ]]; then
        log_warn "Binary not found in release archive. Falling back to building from source."
        return 1
    fi
    
    chmod +x "$binary_file"
    
    if [[ -w "$install_dir" ]]; then
        mv "$binary_file" "$install_dir/"
    else
        sudo mv "$binary_file" "$install_dir/"
    fi
    
    # Also install scripts directory if present in archive
    if [[ -d "scripts" ]]; then
        local scripts_dest="$install_dir/../share/gonwatch/scripts"
        mkdir -p "$(dirname "$scripts_dest")"
        if [[ -w "$(dirname "$scripts_dest")" ]]; then
            cp -r scripts "$scripts_dest"
        else
            sudo mkdir -p "$(dirname "$scripts_dest")"
            sudo cp -r scripts "$scripts_dest"
        fi
        log_info "Scripts installed to $scripts_dest"
    fi
    
    log_success "Installed $BINARY_NAME $version to $install_dir"
    return 0
}

install_from_source() {
    local install_dir="$1"
    
    log_step "Building from source..."
    
    # Check Go version
    if command_exists go; then
        local go_version
        go_version=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
        
        if ! version_gte "$go_version" "$MIN_GO_VERSION"; then
            log_warn "Go version $go_version is too old. Minimum required: $MIN_GO_VERSION"
            return 1
        fi
        log_info "Using Go version $go_version"
    else
        die "Go is required to build from source"
    fi
    
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT
    
    log_info "Cloning repository..."
    git clone --depth 1 "$GITHUB_URL" "$tmp_dir/$REPO_NAME"
    
    cd "$tmp_dir/$REPO_NAME"
    
    log_info "Building..."
    go build -ldflags="-s -w" -o "$BINARY_NAME" .
    
    if [[ ! -f "$BINARY_NAME" ]]; then
        die "Build failed: binary not created"
    fi
    
    chmod +x "$BINARY_NAME"
    
    if [[ -w "$install_dir" ]]; then
        mv "$BINARY_NAME" "$install_dir/"
    else
        sudo mv "$BINARY_NAME" "$install_dir/"
    fi
    
    # Install scripts directory
    local scripts_dest="$install_dir/../share/gonwatch/scripts"
    if [[ -d "scripts" ]]; then
        if [[ -w "$(dirname "$scripts_dest")" ]] || mkdir -p "$(dirname "$scripts_dest")" 2>/dev/null; then
            mkdir -p "$scripts_dest"
            cp -r scripts/* "$scripts_dest/"
        else
            sudo mkdir -p "$scripts_dest"
            sudo cp -r scripts/* "$scripts_dest/"
        fi
        log_info "Scripts installed to $scripts_dest"
    fi
    
    log_success "Built and installed $BINARY_NAME to $install_dir"
}

install_with_go_install() {
    local install_dir="$1"
    
    log_step "Installing with 'go install'..."
    
    if ! command_exists go; then
        die "Go is required for 'go install'"
    fi
    
    go install "${GITHUB_URL}@latest"
    
    # Move from GOPATH/bin to install_dir if different
    local gopath_bin="${GOPATH:-$HOME/go}/bin"
    
    if [[ "$gopath_bin" != "$install_dir" ]] && [[ -f "$gopath_bin/$BINARY_NAME" ]]; then
        if [[ -w "$install_dir" ]]; then
            mv "$gopath_bin/$BINARY_NAME" "$install_dir/"
        else
            sudo mv "$gopath_bin/$BINARY_NAME" "$install_dir/"
        fi
    fi
    
    log_success "Installed $BINARY_NAME to $install_dir"
}

# ============================================================================
# Setup Shell Completion
# ============================================================================

setup_completions() {
    local os="$1"
    local shell_name
    shell_name=$(basename "${SHELL:-/bin/bash}")
    
    log_step "Setting up shell completions..."
    
    # Check if gonwatch supports completions
    if ! "$BINARY_NAME" completion "$shell_name" &>/dev/null 2>&1; then
        log_info "Shell completions not available for this build"
        return 0
    fi
    
    case "$shell_name" in
        bash)
            local completion_dir
            if [[ "$os" == "darwin" ]]; then
                completion_dir="$(brew --prefix 2>/dev/null)/etc/bash_completion.d" || completion_dir="$HOME/.bash_completion.d"
            else
                completion_dir="${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions"
            fi
            mkdir -p "$completion_dir"
            "$BINARY_NAME" completion bash > "$completion_dir/$BINARY_NAME"
            log_success "Bash completions installed to $completion_dir"
            ;;
        zsh)
            local completion_dir="${ZDOTDIR:-$HOME}/.zfunc"
            mkdir -p "$completion_dir"
            "$BINARY_NAME" completion zsh > "$completion_dir/_$BINARY_NAME"
            log_success "Zsh completions installed to $completion_dir"
            log_info "Add 'fpath+=~/.zfunc' to your .zshrc if not already present"
            ;;
        fish)
            local completion_dir="${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions"
            mkdir -p "$completion_dir"
            "$BINARY_NAME" completion fish > "$completion_dir/$BINARY_NAME.fish"
            log_success "Fish completions installed to $completion_dir"
            ;;
        *)
            log_info "Shell completions not configured for $shell_name"
            ;;
    esac
}

# ============================================================================
# Post-Installation
# ============================================================================

verify_installation() {
    local install_dir="$1"
    local binary_path="$install_dir/$BINARY_NAME"
    
    log_step "Verifying installation..."
    
    if [[ ! -x "$binary_path" ]] && ! command_exists "$BINARY_NAME"; then
        die "Installation verification failed: $BINARY_NAME not found or not executable"
    fi
    
    # Try to run version command
    if "$BINARY_NAME" --version &>/dev/null 2>&1; then
        local version
        version=$("$BINARY_NAME" --version 2>&1 | head -1)
        log_success "Installed: $version"
    else
        log_success "$BINARY_NAME installed successfully"
    fi
}

print_post_install_info() {
    local install_dir="$1"
    local os="$2"
    
    echo ""
    echo "=============================================="
    printf "${GREEN}${BOLD}Installation Complete!${NC}\n"
    echo "=============================================="
    echo ""
    
    # Check if install_dir is in PATH
    if [[ ":$PATH:" != *":$install_dir:"* ]]; then
        log_warn "$install_dir is not in your PATH"
        echo ""
        echo "Add the following to your shell profile:"
        echo ""
        echo "  # For bash (~/.bashrc or ~/.bash_profile):"
        echo "  export PATH=\"\$PATH:$install_dir\""
        echo ""
        echo "  # For zsh (~/.zshrc):"
        echo "  export PATH=\"\$PATH:$install_dir\""
        echo ""
        echo "  # For fish (~/.config/fish/config.fish):"
        echo "  fish_add_path $install_dir"
        echo ""
        echo "Then reload your shell or run: source ~/.bashrc (or equivalent)"
        echo ""
    fi
    
    echo "Usage:"
    echo "  $BINARY_NAME              # Launch the TUI"
    echo "  $BINARY_NAME debug        # Launch with debug logging"
    echo ""
    echo "For more information, visit: $GITHUB_URL"
    echo ""
}

# ============================================================================
# Main Installation Flow
# ============================================================================

print_banner() {
    cat << 'EOF'

   ██████╗  ██████╗ ███╗   ██╗██╗    ██╗ █████╗ ████████╗ ██████╗██╗  ██╗
  ██╔════╝ ██╔═══██╗████╗  ██║██║    ██║██╔══██╗╚══██╔══╝██╔════╝██║  ██║
  ██║  ███╗██║   ██║██╔██╗ ██║██║ █╗ ██║███████║   ██║   ██║     ███████║
  ██║   ██║██║   ██║██║╚██╗██║██║███╗██║██╔══██║   ██║   ██║     ██╔══██║
  ╚██████╔╝╚██████╔╝██║ ╚████║╚███╔███╔╝██║  ██║   ██║   ╚██████╗██║  ██║
   ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝ ╚══╝╚══╝ ╚═╝  ╚═╝   ╚═╝    ╚═════╝╚═╝  ╚═╝

                    Go and Watch - Universal Installer
EOF
    echo ""
}

main() {
    print_banner
    
    # Detect environment
    local os arch distro pkg_manager install_dir
    os=$(detect_os)
    arch=$(detect_arch)
    pkg_manager=$(detect_package_manager "$os")
    install_dir=$(get_install_dir "$os")
    
    if [[ "$os" == "linux" ]]; then
        distro=$(detect_linux_distro)
        log_info "Detected: Linux ($distro) on $arch"
    else
        log_info "Detected: $os on $arch"
    fi
    
    log_info "Package manager: $pkg_manager"
    log_info "Install directory: $install_dir"
    echo ""
    
    # Check for unsupported platforms
    if [[ "$os" == "unknown" ]] || [[ "$arch" == "unknown" ]]; then
        die "Unsupported platform: $(uname -s) / $(uname -m)"
    fi
    
    # =========== Dependency Checks ===========
    log_step "Checking dependencies..."
    
    local need_go=false
    local need_mpv=false
    local need_git=false
    local need_python=false
    local need_nodejs=false
    local need_chromium=false
    local need_pip_packages=false
    
    # Check git
    if ! command_exists git; then
        need_git=true
        log_warn "git not found"
    else
        log_success "git found: $(git --version | head -1)"
    fi
    
    # Check Go
    if ! command_exists go; then
        need_go=true
        log_warn "Go not found"
    else
        local go_version
        go_version=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
        if version_gte "$go_version" "$MIN_GO_VERSION"; then
            log_success "Go found: version $go_version"
        else
            log_warn "Go version $go_version is too old (need >= $MIN_GO_VERSION)"
            need_go=true
        fi
    fi
    
    # Check Python
    if ! command_exists python3 && ! command_exists python; then
        need_python=true
        log_warn "Python 3 not found (required for video source scraping)"
    else
        local py_cmd="python3"
        command_exists python3 || py_cmd="python"
        local py_version
        py_version=$($py_cmd --version 2>&1 | grep -oE '[0-9]+\.[0-9]+')
        log_success "Python found: version $py_version"
        
        # Check for required packages
        if ! $py_cmd -c "import selenium_driverless" 2>/dev/null; then
            need_pip_packages=true
            log_warn "Python packages missing (selenium-driverless, etc.)"
        else
            log_success "Python packages found"
        fi
    fi
    
    # Check Node.js
    if ! command_exists node; then
        need_nodejs=true
        log_warn "Node.js not found (required for AnimePahe support)"
    else
        local node_version
        node_version=$(node --version 2>&1)
        log_success "Node.js found: $node_version"
    fi
    
    # Check Chrome/Chromium
    if ! command_exists google-chrome && ! command_exists google-chrome-stable && \
       ! command_exists chromium && ! command_exists chromium-browser; then
        need_chromium=true
        log_warn "Chrome/Chromium not found (required for video source scraping)"
    else
        log_success "Chrome/Chromium found"
    fi
    
    # Check mpv
    if ! command_exists mpv; then
        need_mpv=true
        log_warn "mpv not found (required for video playback)"
    else
        log_success "mpv found: $(mpv --version | head -1)"
    fi
    
    echo ""
    
    # =========== Install Missing Dependencies ===========
    if $need_git || $need_go || $need_python || $need_nodejs || $need_chromium || $need_mpv || $need_pip_packages; then
        log_step "Installing missing dependencies..."
        
        if [[ "$pkg_manager" == "none" ]]; then
            log_warn "No supported package manager found."
            log_info "Please install the following manually:"
            $need_git && echo "  - git"
            $need_go && echo "  - Go (>= $MIN_GO_VERSION)"
            $need_python && echo "  - Python 3 with pip"
            $need_nodejs && echo "  - Node.js"
            $need_chromium && echo "  - Google Chrome or Chromium"
            $need_mpv && echo "  - mpv"
            $need_pip_packages && echo "  - Python packages: ${PYTHON_PACKAGES[*]}"
            echo ""
            
            if ! confirm "Continue anyway? (may fail)"; then
                exit 1
            fi
        else
            if $need_git; then
                install_git "$pkg_manager"
            fi
            
            if $need_go; then
                install_go "$os" "$pkg_manager"
                # Refresh PATH for newly installed Go
                export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"
            fi
            
            if $need_python; then
                install_python "$os" "$pkg_manager"
                need_pip_packages=true  # Always install packages after fresh Python install
            fi
            
            if $need_pip_packages; then
                install_python_packages || true  # Don't fail entire install
            fi
            
            if $need_nodejs; then
                install_nodejs "$os" "$pkg_manager" || true  # Optional, don't fail
            fi
            
            if $need_chromium; then
                install_chromium "$os" "$pkg_manager" || true  # Don't fail if already have Chrome
            fi
            
            if $need_mpv; then
                install_mpv "$os" "$pkg_manager" || true  # Don't fail if mpv install fails
            fi
        fi
        
        echo ""
    fi
    
    # =========== Install Gonwatch ===========
    log_step "Installing $BINARY_NAME..."
    
    # Try release binary first, then fall back to source
    if ! install_from_release "$os" "$arch" "$install_dir"; then
        log_info "Building from source instead..."
        
        if ! command_exists go; then
            die "Go is required to build from source. Please install Go >= $MIN_GO_VERSION and try again."
        fi
        
        if ! command_exists git; then
            die "git is required to build from source. Please install git and try again."
        fi
        
        install_from_source "$install_dir"
    fi
    
    # =========== Post-Installation ===========
    
    # Add install_dir to PATH for verification
    export PATH="$PATH:$install_dir"
    
    verify_installation "$install_dir"
    
    # Setup completions (optional, don't fail if it doesn't work)
    setup_completions "$os" 2>/dev/null || true
    
    print_post_install_info "$install_dir" "$os"
}

# ============================================================================
# Entry Point
# ============================================================================

# Handle arguments
case "${1:-}" in
    --help|-h)
        cat << EOF
Gonwatch Universal Installer

USAGE:
    $0 [OPTIONS]

OPTIONS:
    -h, --help      Show this help message
    --deps-only     Only install dependencies (go, mpv, python, nodejs, chromium)
    --no-deps       Skip dependency installation
    --source        Force build from source (skip release check)

DEPENDENCIES:
    Required:
      - mpv           Video player for playback
      - python3       Runs scraping scripts
      - Chrome/Chromium  Headless browser for scraping
      - pip packages  selenium-driverless, requests, langdetect, beautifulsoup4, tls_client
    
    For AnimePahe:
      - nodejs        Resolves obfuscated stream URLs
    
    Build-time only:
      - go (>= $MIN_GO_VERSION)  Compiles the binary
      - git           Clones the repository

ENVIRONMENT VARIABLES:
    REPO_OWNER      GitHub repository owner (default: kbwhodat)
    
EXAMPLES:
    # Standard installation
    ./install.sh
    
    # Install from URL
    curl -fsSL https://raw.githubusercontent.com/kbwhodat/gonwatch/main/install.sh | bash
    
    # Only install dependencies
    ./install.sh --deps-only
EOF
        exit 0
        ;;
    --deps-only)
        print_banner
        os=$(detect_os)
        pkg_manager=$(detect_package_manager "$os")
        
        log_step "Installing dependencies only..."
        
        command_exists git || install_git "$pkg_manager"
        command_exists go || install_go "$os" "$pkg_manager"
        command_exists python3 || command_exists python || install_python "$os" "$pkg_manager"
        install_python_packages || true
        command_exists node || install_nodejs "$os" "$pkg_manager" || true
        command_exists google-chrome || command_exists chromium || command_exists chromium-browser || \
            install_chromium "$os" "$pkg_manager" || true
        command_exists mpv || install_mpv "$os" "$pkg_manager"
        
        log_success "Dependencies installed"
        exit 0
        ;;
    --no-deps)
        print_banner
        os=$(detect_os)
        arch=$(detect_arch)
        install_dir=$(get_install_dir "$os")
        
        if ! install_from_release "$os" "$arch" "$install_dir"; then
            install_from_source "$install_dir"
        fi
        
        export PATH="$PATH:$install_dir"
        verify_installation "$install_dir"
        print_post_install_info "$install_dir" "$os"
        exit 0
        ;;
    --source)
        print_banner
        os=$(detect_os)
        pkg_manager=$(detect_package_manager "$os")
        install_dir=$(get_install_dir "$os")
        
        # Quick dependency check
        command_exists git || install_git "$pkg_manager"
        command_exists go || install_go "$os" "$pkg_manager"
        command_exists python3 || command_exists python || install_python "$os" "$pkg_manager"
        install_python_packages || true
        command_exists node || install_nodejs "$os" "$pkg_manager" || true
        command_exists google-chrome || command_exists chromium || command_exists chromium-browser || \
            install_chromium "$os" "$pkg_manager" || true
        command_exists mpv || install_mpv "$os" "$pkg_manager" || true
        
        export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"
        
        install_from_source "$install_dir"
        
        export PATH="$PATH:$install_dir"
        verify_installation "$install_dir"
        print_post_install_info "$install_dir" "$os"
        exit 0
        ;;
    "")
        main
        ;;
    *)
        die "Unknown option: $1. Use --help for usage information."
        ;;
esac
