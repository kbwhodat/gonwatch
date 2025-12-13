#!/bin/bash

# Gonwatch Install Script
# Supports Linux and macOS

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

# Install Go
install_go() {
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_status "Go is already installed: $GO_VERSION"
        
        # Check if Go version meets minimum requirement (1.25.0)
        if go version | grep -q "go1\.[2-9][5-9]\|go1\.[3-9][0-9]\|go[2-9]\."; then
            print_success "Go version meets requirements"
        else
            print_warning "Go version might be too old. Recommended: 1.25.0 or later"
        fi
        return 0
    fi

    print_status "Installing Go..."
    
    GO_VERSION="1.23.5"  # Using stable version as of late 2024
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
    cd /tmp
    if ! wget -q "https://golang.org/dl/${GO_TAR}"; then
        print_error "Failed to download Go"
        exit 1
    fi
    
    # Remove old installation if exists
    sudo rm -rf /usr/local/go
    
    # Extract Go
    sudo tar -C /usr/local -xzf "$GO_TAR"
    rm "$GO_TAR"
    
    # Add Go to PATH
    export PATH=$PATH:/usr/local/go/bin
    
    # Add to shell profile
    SHELL_RC=""
    if [[ "$SHELL" == *"bash"* ]]; then
        SHELL_RC="$HOME/.bashrc"
    elif [[ "$SHELL" == *"zsh"* ]]; then
        SHELL_RC="$HOME/.zshrc"
    fi
    
    if [ -n "$SHELL_RC" ] && ! grep -q '/usr/local/go/bin' "$SHELL_RC"; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> "$SHELL_RC"
        print_status "Added Go to PATH in $SHELL_RC"
    fi
    
    print_success "Go installed successfully"
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
    
    # Required packages based on shell.nix
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
    )
    
    for package in "${PYTHON_PACKAGES[@]}"; do
        print_status "Installing Python package: $package"
        pip3 install --user "$package"
    done
    
    # Install selenium-driverless and cdp-socket from source
    print_status "Installing selenium-driverless..."
    pip3 install --user selenium-driverless==1.9.4
    
    print_success "Python dependencies installed"
}

# Install system dependencies
install_system_deps() {
    print_status "Installing system dependencies..."
    
    if [[ "$OS" == "linux" ]]; then
        if [[ "$DISTRO" == "debian" ]]; then
            sudo apt update
            sudo apt install -y build-essential wget curl gcc
        elif [[ "$DISTRO" == "redhat" ]]; then
            sudo yum groupinstall -y "Development Tools"
            sudo yum install -y wget curl gcc
        elif [[ "$DISTRO" == "arch" ]]; then
            sudo pacman -S --noconfirm base-devel wget curl gcc
        fi
    elif [[ "$OS" == "macos" ]]; then
        if ! command_exists brew; then
            print_warning "Homebrew not found. Installing Homebrew..."
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi
        brew install gcc
    fi
    
    print_success "System dependencies installed"
}

# Build the Go application
build_app() {
    print_status "Building Gonwatch..."
    
    # Navigate to project directory
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please run this script from the project root."
        exit 1
    fi
    
    # Download Go modules
    if [ "$PATH" != *"/usr/local/go/bin"* ]; then
        export PATH=$PATH:/usr/local/go/bin
    fi
    
    go mod download
    go mod tidy
    
    # Build the application
    go build -o gonwatch main.go
    
    if [ -f "gonwatch" ]; then
        print_success "Gonwatch built successfully"
        
        # Make it executable
        chmod +x gonwatch
        
        # Optionally move to system PATH
        if [ "$1" = "--system-install" ]; then
            sudo mv gonwatch /usr/local/bin/
            print_success "Gonwatch installed to /usr/local/bin/gonwatch"
        else
            print_status "Gonwatch binary created in current directory"
            print_status "Run: ./gonwatch"
        fi
    else
        print_error "Build failed"
        exit 1
    fi
}

# Verify installation
verify_installation() {
    print_status "Verifying installation..."
    
    # Check Go
    if command_exists go; then
        GO_VERSION=$(go version)
        print_success "Go: $GO_VERSION"
    else
        print_error "Go not found in PATH"
        return 1
    fi
    
    # Check Python packages
    python3 -c "import langdetect, requests, selenium" 2>/dev/null
    if [ $? -eq 0 ]; then
        print_success "Python dependencies verified"
    else
        print_warning "Some Python dependencies may be missing"
    fi
    
    # Check binary
    if [ -f "/usr/local/bin/gonwatch" ]; then
        print_success "Gonwatch installed system-wide"
    elif [ -f "gonwatch" ]; then
        print_success "Gonwatch binary ready"
    else
        print_error "Gonwatch binary not found"
        return 1
    fi
    
    print_success "Installation completed successfully!"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --system-install    Install gonwatch to /usr/local/bin"
    echo "  --help, -h         Show this help message"
    echo ""
    echo "This script will install:"
    echo "  - Go (if not present)"
    echo "  - Python3 and required packages"
    echo "  - System dependencies"
    echo "  - Build and install gonwatch"
}

# Main installation function
main() {
    echo "Gonwatch Installation Script"
    echo "============================="
    echo ""
    
    # Parse arguments
    SYSTEM_INSTALL=""
    for arg in "$@"; do
        case $arg in
            --system-install)
                SYSTEM_INSTALL="--system-install"
                ;;
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
    
    # Check if running from project root
    if [ ! -f "go.mod" ] || [ ! -f "main.go" ]; then
        print_error "Please run this script from the gonwatch project root directory"
        exit 1
    fi
    
    # Run installation steps
    detect_os
    install_system_deps
    install_go
    install_python_deps
    build_app $SYSTEM_INSTALL
    verify_installation
    
    echo ""
    print_success "Installation complete!"
    echo ""
    echo "To run gonwatch:"
    if [ -f "/usr/local/bin/gonwatch" ]; then
        echo "  gonwatch"
    else
        echo "  ./gonwatch"
    fi
    echo ""
    echo "For debug mode:"
    echo "  gonwatch debug"
    echo ""
    echo "If you encounter any issues, please check:"
    echo "  - Go is in your PATH"
    echo "  - Python dependencies are installed"
    echo "  - The binary has execute permissions"
}

# Run main function
main "$@"