#!/bin/bash

# Gonwatch Remove Script
# Removes Gonwatch and all its dependencies (optional)
# Execute with: curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/gonwatch/main/scripts/remove.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Remove Gonwatch installation
remove_gonwatch() {
    print_status "Removing Gonwatch installation..."
    
    # Remove binary from local bin
    if [ -f "$HOME/.local/bin/gonwatch" ]; then
        rm -f "$HOME/.local/bin/gonwatch"
        print_success "Removed gonwatch binary from ~/.local/bin/"
    fi
    
    # Remove from system bin
    if [ -f "/usr/local/bin/gonwatch" ]; then
        sudo rm -f /usr/local/bin/gonwatch
        print_success "Removed gonwatch binary from /usr/local/bin/"
    fi
    
    # Remove installation directory
    if [ -d "$HOME/.local/share/gonwatch" ]; then
        rm -rf "$HOME/.local/share/gonwatch"
        print_success "Removed installation directory ~/.local/share/gonwatch/"
    fi
    
    # Remove local build if exists
    if [ -f "gonwatch" ] && [ -f "go.mod" ]; then
        rm -f gonwatch
        print_success "Removed local gonwatch binary"
    fi
}

# Remove Python dependencies
remove_python_deps() {
    if [ "$1" = "--full-cleanup" ]; then
        print_status "Removing Python dependencies..."
        
        # List of packages to remove
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
            "selenium-driverless"
        )
        
        for package in "${PYTHON_PACKAGES[@]}"; do
            if pip3 show "$package" >/dev/null 2>&1; then
                print_status "Removing Python package: $package"
                pip3 uninstall -y "$package" --quiet
            fi
        done
        
        print_success "Python dependencies removed"
    else
        print_status "Skipping Python dependency removal (use --full-cleanup to remove)"
    fi
}

# Remove Go installation
remove_go() {
    if [ "$1" = "--full-cleanup" ]; then
        if command_exists go; then
            print_status "Removing Go installation..."
            
            # Remove Go directory
            if [ -d "/usr/local/go" ]; then
                sudo rm -rf /usr/local/go
                print_success "Removed Go from /usr/local/go/"
            fi
            
            # Remove Go from shell profiles
            for rc_file in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
                if [ -f "$rc_file" ]; then
                    if grep -q '/usr/local/go/bin' "$rc_file" 2>/dev/null; then
                        sed -i.bak '/export PATH=.*\/usr\/local\/go\/bin/d' "$rc_file" 2>/dev/null || \
                        sed -i '' '/export PATH=.*\/usr\/local\/go\/bin/d' "$rc_file" 2>/dev/null || true
                        print_status "Removed Go from PATH in $rc_file"
                    fi
                fi
            done
            
            # Remove ~/.local/bin from PATH if it was added by us
            for rc_file in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile"; do
                if [ -f "$rc_file" ]; then
                    if grep -q 'export PATH=.*~\/\.local\/bin' "$rc_file" 2>/dev/null; then
                        sed -i.bak '/export PATH=.*~\/\.local\/bin/d' "$rc_file" 2>/dev/null || \
                        sed -i '' '/export PATH=.*~\/\.local\/bin/d' "$rc_file" 2>/dev/null || true
                        print_status "Removed ~/.local/bin from PATH in $rc_file"
                    fi
                fi
            done
            
            print_success "Go removed"
        else
            print_status "Go not found, skipping removal"
        fi
    else
        print_status "Skipping Go removal (use --full-cleanup to remove)"
    fi
}

# Clean up temporary files
cleanup_temp() {
    print_status "Cleaning up temporary files..."
    
    # Remove any temporary directories
    if [ -d "/tmp/go*" ]; then
        rm -rf /tmp/go*
        print_status "Cleaned up temporary Go files"
    fi
    
    # Clean up pip cache
    if command_exists pip3; then
        pip3 cache purge --quiet 2>/dev/null || true
        print_status "Cleaned up pip cache"
    fi
    
    # Clean up Go modules cache if removing Go
    if [ "$1" = "--full-cleanup" ]; then
        if [ -d "$HOME/go" ]; then
            rm -rf "$HOME/go"
            print_status "Cleaned up Go modules cache"
        fi
        
        if [ -d "$HOME/.cache/go-build" ]; then
            rm -rf "$HOME/.cache/go-build"
            print_status "Cleaned up Go build cache"
        fi
    fi
}

# Verify removal
verify_removal() {
    print_status "Verifying removal..."
    
    # Check if gonwatch command still exists
    if command_exists gonwatch; then
        print_warning "Gonwatch command still found in PATH at: $(which gonwatch)"
        print_warning "You may need to restart your shell or manually remove it"
    else
        print_success "Gonwatch command removed from PATH"
    fi
    
    if [ "$1" = "--full-cleanup" ]; then
        # Check if Go still exists
        if command_exists go; then
            print_warning "Go command still found: $(go version)"
        else
            print_success "Go removed from system"
        fi
        
        # Check if Python packages still exist
        REMAINING_PACKAGES=()
        PYTHON_PACKAGES=("langdetect" "requests" "selenium")
        
        for package in "${PYTHON_PACKAGES[@]}"; do
            if pip3 show "$package" >/dev/null 2>&1; then
                REMAINING_PACKAGES+=("$package")
            fi
        done
        
        if [ ${#REMAINING_PACKAGES[@]} -eq 0 ]; then
            print_success "Key Python dependencies removed"
        else
            print_warning "Some Python packages still installed: ${REMAINING_PACKAGES[*]}"
        fi
    fi
}

# Show usage
show_usage() {
    echo "Gonwatch Remove Script"
    echo "====================="
    echo ""
    echo "Usage:"
    echo "  curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/gonwatch/main/scripts/remove.sh | bash"
    echo ""
    echo "For complete cleanup (removes Go and Python dependencies):"
    echo "  curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/gonwatch/main/scripts/remove.sh | bash -s -- --full-cleanup"
    echo ""
    echo "This script will:"
    echo "  - Remove gonwatch binary and installation directory"
    echo "  - Remove ~/.local/bin/gonwatch symlink"
    echo "  - Clean up temporary files"
    echo ""
    echo "With --full-cleanup:"
    echo "  - Also removes Go installation"
    echo "  - Also removes Python dependencies"
    echo "  - Also cleans up caches"
    echo ""
    echo "Note: You may need to restart your shell after removal to update PATH"
}

# Main removal function
main() {
    echo "Gonwatch Removal Script"
    echo "======================="
    echo ""
    
    # Parse arguments
    FULL_CLEANUP=""
    for arg in "$@"; do
        case $arg in
            --full-cleanup)
                FULL_CLEANUP="--full-cleanup"
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
    
    if [ -n "$FULL_CLEANUP" ]; then
        print_warning "Full cleanup mode enabled - This will remove Go and Python dependencies!"
        read -p "Are you sure? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Cleanup cancelled"
            exit 0
        fi
    fi
    
    # Run removal steps
    remove_gonwatch
    remove_python_deps $FULL_CLEANUP
    remove_go $FULL_CLEANUP
    cleanup_temp $FULL_CLEANUP
    verify_removal $FULL_CLEANUP
    
    echo ""
    print_success "Removal completed!"
    echo ""
    echo "Note: You may need to restart your shell or run:"
    echo "  source ~/.bashrc  # or ~/.zshrc"
    echo ""
    echo "To verify the removal, try: which gonwatch"
}

# Run main function
main "$@"