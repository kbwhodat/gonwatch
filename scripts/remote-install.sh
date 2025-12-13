#!/bin/bash

# Gonwatch Remote Install Script
# Execute with: curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/gonwatch/main/scripts/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
        if [ -f /etc/debian_version ]; then
            DISTRO="debian"
        elif [ -f /etc/redhat-release ]; then
            DISTRO="redhat"
        elif [ -f /etc/arch-release ]; then
            DISTRO="arch"
        else
            DISTRO="unknown"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
        DISTRO="macos"
    else
        print_error "Unsupported OS: $OSTYPE"
        exit 1
    fi
    
    print_status "Detected OS: $OS ($DISTRO)"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Create temporary directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Version comparison function
version_compare() {
    local version1=$1
    local version2=$2
    
    # Split versions into arrays
    IFS='.' read -ra v1_parts <<< "$version1"
    IFS='.' read -ra v2_parts <<< "$version2"
    
    # Compare major, minor, patch versions
    for i in {0..2}; do
        local v1_part=${v1_parts[$i]:-0}
        local v2_part=${v2_parts[$i]:-0}
        
        if [ "$v1_part" -gt "$v2_part" ]; then
            return 0  # version1 > version2
        elif [ "$v1_part" -lt "$v2_part" ]; then
            return 1  # version1 < version2
        fi
    done
    return 0  # versions are equal
}

# Check Go version
check_go_version() {
    local go_version_str=$1
    local min_version="1.25.0"
    
    # Extract version number (remove 'go' prefix and any suffixes)
    local clean_version=$(echo "$go_version_str" | sed 's/^go//' | sed 's/[^0-9.].*//')
    
    if version_compare "$clean_version" "$min_version"; then
        return 0  # Version is acceptable
    else
        return 1  # Version is too old
    fi
}

# Install Go
install_go() {
    if command_exists go; then
        GO_VERSION_OUTPUT=$(go version)
        GO_VERSION=$(echo "$GO_VERSION_OUTPUT" | awk '{print $3}' | sed 's/go//')
        print_status "Go is already installed: $GO_VERSION"
        
        # Check if Go version meets minimum requirement (>= 1.25.0)
        if check_go_version "$GO_VERSION_OUTPUT"; then
            print_success "Go version meets requirements (>= 1.25.0)"
            return 0
        else
            print_warning "Go version $GO_VERSION is too old. Required: >= 1.25.0"
            print_status "Will install compatible Go version..."
            
            # Remove old installation and continue with installation
            if [ -d "/usr/local/go" ]; then
                print_status "Removing old Go installation..."
                sudo rm -rf /usr/local/go
            fi
        fi
    fi

    print_status "Installing Go (>= 1.25.0)..."
    
    GO_VERSION="1.25.0"
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) print_error "Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    
    if [[ "$OS" == "linux" ]]; then
        GO_TAR="go${GO_VERSION}.linux-${ARCH}.tar.gz"
    elif [[ "$OS" == "macos" ]]; then
        GO_TAR="go${GO_VERSION}.darwin-${ARCH}.tar.gz"
    fi
    
    # Download and install Go
    cd "$TEMP_DIR"
    if ! wget -q "https://golang.org/dl/${GO_TAR}"; then
        print_error "Failed to download Go"
        exit 1
    fi
    
    # Remove old installation if exists
    if [ -d "/usr/local/go" ]; then
        sudo rm -rf /usr/local/go
    fi
    
    # Extract Go
    sudo tar -C /usr/local -xzf "$GO_TAR"
    rm "$GO_TAR"
    
    # Add Go to PATH for current session
    export PATH=$PATH:/usr/local/go/bin
    
    # Add to shell profile
    SHELL_RC=""
    if [[ "$SHELL" == *"bash"* ]]; then
        SHELL_RC="$HOME/.bashrc"
    elif [[ "$SHELL" == *"zsh"* ]]; then
        SHELL_RC="$HOME/.zshrc"
    fi
    
    if [ -n "$SHELL_RC" ] && ! grep -q '/usr/local/go/bin' "$SHELL_RC" 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> "$SHELL_RC"
        print_status "Added Go to PATH in $SHELL_RC"
    fi
    
    # Verify the installation
    if command_exists /usr/local/go/bin/go; then
        INSTALLED_VERSION=$(/usr/local/go/bin/go version | awk '{print $3}')
        print_success "Go $INSTALLED_VERSION installed successfully"
        
        # Final version check
        if check_go_version "$INSTALLED_VERSION"; then
            print_success "Go version meets requirements (>= 1.25.0)"
        else
            print_error "Installed Go version is still incompatible: $INSTALLED_VERSION"
            exit 1
        fi
    else
        print_error "Go installation failed"
        exit 1
    fi
}

# Install Python and dependencies
install_python_deps() {
    print_status "Setting up Python dependencies..."
    
    # Check Python version
    if command_exists python3; then
        PYTHON_VERSION=$(python3 --version | awk '{print $2}')
        print_status "Python3 found: $PYTHON_VERSION"
    else
        print_error "Python3 is required but not found"
        exit 1
    fi
    
    # Install pip if not present
    if ! command_exists pip3; then
        print_status "Installing pip3..."
        if [[ "$OS" == "linux" ]]; then
            if [[ "$DISTRO" == "debian" ]]; then
                sudo apt update && sudo apt install -y python3-pip
            elif [[ "$DISTRO" == "redhat" ]]; then
                sudo yum install -y python3-pip
            elif [[ "$DISTRO" == "arch" ]]; then
                sudo pacman -S --noconfirm python-pip
            fi
        elif [[ "$OS" == "macos" ]]; then
            if command_exists brew; then
                brew install python3
            else
                print_error "Please install Homebrew first: https://brew.sh/"
                exit 1
            fi
        fi
    fi
    
    # Install Python dependencies
    print_status "Installing Python packages..."
    
    # Create ~/.local directory if it doesn't exist
    mkdir -p "$HOME/.local/bin"
    mkdir -p "$HOME/.local/lib"
    
    # Required packages
    PYTHON_PACKAGES=(
        "langdetect"
        "requests"
        "websockets"
        "numpy"
        "aiofiles"
        "matplotlib"
        "scipy"
        "platformdirs"
        "aiohttp"
        "jsondiff"
        "orjson"
        "selenium"
        "selenium-driverless==1.9.4"
    )
    
    for package in "${PYTHON_PACKAGES[@]}"; do
        print_status "Installing Python package: $package"
        pip3 install --prefix="$HOME/.local" --user "$package" --quiet
    done
    
    print_success "Python dependencies installed"
}

# Install system dependencies
install_system_deps() {
    print_status "Installing system dependencies..."
    
    if [[ "$OS" == "linux" ]]; then
        if [[ "$DISTRO" == "debian" ]]; then
            sudo apt update
            sudo apt install -y build-essential wget curl gcc git
        elif [[ "$DISTRO" == "redhat" ]]; then
            sudo yum groupinstall -y "Development Tools"
            sudo yum install -y wget curl gcc git
        elif [[ "$DISTRO" == "arch" ]]; then
            sudo pacman -S --noconfirm base-devel wget curl gcc git
        fi
    elif [[ "$OS" == "macos" ]]; then
        if ! command_exists brew; then
            print_warning "Homebrew not found. Installing Homebrew..."
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi
        brew install gcc git
    fi
    
    print_success "System dependencies installed"
}

# Clone and build the application
clone_and_build() {
    print_status "Downloading and building Gonwatch..."
    
    # Set installation directory
    INSTALL_DIR="$HOME/.local/share/gonwatch"
    BIN_DIR="$HOME/.local/bin"
    
    # Create directories
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$BIN_DIR"
    
    # Clone repository
    cd "$TEMP_DIR"
    if ! git clone https://github.com/katob/gonwatch.git; then
        print_error "Failed to clone repository"
        exit 1
    fi
    
    cd gonwatch
    
    # Ensure Go is in PATH
    if [ "$PATH" != *"/usr/local/go/bin"* ]; then
        export PATH=$PATH:/usr/local/go/bin
    fi
    
    # Download Go modules and build
    go mod download
    go mod tidy
    
    # Build the application
    go build -o "$INSTALL_DIR/gonwatch" main.go
    
    if [ -f "$INSTALL_DIR/gonwatch" ]; then
        print_success "Gonwatch built successfully"
        
        # Make it executable
        chmod +x "$INSTALL_DIR/gonwatch"
        
        # Create symlink in local bin
        ln -sf "$INSTALL_DIR/gonwatch" "$BIN_DIR/gonwatch"
        
        # Add local bin to PATH if not already
        if [ "$PATH" != *"$HOME/.local/bin"* ]; then
            SHELL_RC=""
            if [[ "$SHELL" == *"bash"* ]]; then
                SHELL_RC="$HOME/.bashrc"
            elif [[ "$SHELL" == *"zsh"* ]]; then
                SHELL_RC="$HOME/.zshrc"
            fi
            
            if [ -n "$SHELL_RC" ]; then
                if ! grep -q 'PATH=.*\$HOME/.local/bin' "$SHELL_RC" 2>/dev/null; then
                    echo 'export PATH=$PATH:$HOME/.local/bin' >> "$SHELL_RC"
                    print_status "Added ~/.local/bin to PATH in $SHELL_RC"
                fi
                # Add PYTHONPATH for Python packages
                if ! grep -q 'PYTHONPATH=.*\$HOME/.local/lib' "$SHELL_RC" 2>/dev/null; then
                    echo 'export PYTHONPATH=$HOME/.local/lib/python3*/site-packages:$PYTHONPATH' >> "$SHELL_RC"
                    print_status "Added PYTHONPATH for user packages in $SHELL_RC"
                fi
            fi
        fi
        
        print_success "Gonwatch installed to ~/.local/bin/gonwatch"
    else
        print_error "Build failed"
        exit 1
    fi
    
    # Copy the Python scripts
    mkdir -p "$INSTALL_DIR/scripts"
    cp scripts/*.py "$INSTALL_DIR/scripts/" 2>/dev/null || true
}

# Verify installation
verify_installation() {
    print_status "Verifying installation..."
    
    # Check Go
    if command_exists go; then
        GO_VERSION_OUTPUT=$(go version)
        GO_VERSION=$(echo "$GO_VERSION_OUTPUT" | awk '{print $3}')
        print_status "Go found: $GO_VERSION"
        
        if check_go_version "$GO_VERSION_OUTPUT"; then
            print_success "Go version meets requirements (>= 1.25.0)"
        else
            print_error "Go version is too old: $GO_VERSION (required: >= 1.25.0)"
            return 1
        fi
    else
        print_error "Go not found in PATH"
        return 1
    fi
    
    # Check Python packages
    PYTHONPATH="$HOME/.local/lib/python3*/site-packages:$PYTHONPATH" python3 -c "import langdetect, requests, selenium" 2>/dev/null
    if [ $? -eq 0 ]; then
        print_success "Python dependencies verified"
    else
        print_warning "Some Python dependencies may be missing - trying with PYTHONPATH"
        # Try with explicit PYTHONPATH
        export PYTHONPATH="$HOME/.local/lib/python3*/site-packages:$PYTHONPATH"
        python3 -c "import langdetect, requests, selenium" 2>/dev/null
        if [ $? -eq 0 ]; then
            print_success "Python dependencies verified with PYTHONPATH"
        else
            print_warning "Some Python dependencies may be missing"
        fi
    fi
    
    # Check binary
    if command_exists gonwatch; then
        GONWATCH_PATH=$(which gonwatch)
        print_success "Gonwatch installed: $GONWATCH_PATH"
    else
        print_error "Gonwatch binary not found in PATH"
        return 1
    fi
    
    print_success "Installation completed successfully!"
}

# Show usage
show_usage() {
    echo "Gonwatch Remote Installation"
    echo "==========================="
    echo ""
    echo "Execute with:"
    echo "  curl -fsSL https://raw.githubusercontent.com/katob/gonwatch/main/scripts/install.sh | bash"
    echo ""
    echo "Options (pipe to bash with additional args):"
    echo "  bash <(curl -fsSL https://raw.githubusercontent.com/katob/gonwatch/main/scripts/install.sh) --help"
    echo ""
    echo "This script will install:"
    echo "  - Go (if not present)"
    echo "  - Python3 and required packages"
    echo "  - System dependencies"
    echo "  - Gonwatch to ~/.local/bin/"
}

# Main installation function
main() {
    echo "Gonwatch Remote Installation"
    echo "==========================="
    echo ""
    
    # Parse arguments
    for arg in "$@"; do
        case $arg in
            --help|-h)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $arg"
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Check if git is available
    if ! command_exists git; then
        print_error "Git is required but not found"
        exit 1
    fi
    
    # Run installation steps
    detect_os
    install_system_deps
    install_go
    install_python_deps
    clone_and_build
    verify_installation
    
    echo ""
    print_success "Installation complete!"
    echo ""
    echo "To run gonwatch:"
    echo "  gonwatch"
    echo ""
    echo "For debug mode:"
    echo "  gonwatch debug"
    echo ""
    echo "Note: You may need to restart your shell or run:"
    echo "  source ~/.bashrc  # or ~/.zshrc"
    echo ""
    echo "If you encounter any issues, please check:"
    echo "  - ~/.local/bin is in your PATH"
    echo "  - Go is in your PATH"
    echo "  - Python dependencies are installed"
}

# Run main function
main "$@"