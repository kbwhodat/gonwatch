#!/bin/bash

# Upload helper for Gonwatch installation scripts
# This script helps you upload the remote install script to GitHub

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

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if git repo exists
check_git_repo() {
    if [ ! -d ".git" ]; then
        print_error "Not in a git repository"
        exit 1
    fi
    
    # Check remote
    if ! git remote get-url origin >/dev/null 2>&1; then
        print_error "No origin remote found"
        exit 1
    fi
    
    REMOTE_URL=$(git remote get-url origin)
    print_status "Git remote: $REMOTE_URL"
}

# Upload files to GitHub
upload_files() {
    print_status "Staging installation scripts..."
    
    # Add scripts directory
    git add scripts/
    
    # Add local install script
    git add install.sh
    
    # Commit changes
    print_status "Committing changes..."
    git commit -m "Add installation scripts for Linux and macOS
    
    - Added local install.sh for direct execution
    - Added remote-install.sh for curl installation
    - Scripts support Linux (Debian/RedHat/Arch) and macOS
    - Auto-install Go, Python dependencies, and build Gonwatch"
    
    # Push to GitHub
    print_status "Pushing to GitHub..."
    git push origin main
    
    print_success "Installation scripts uploaded to GitHub!"
}

# Show usage instructions
show_usage() {
    echo "Installation Instructions:"
    echo "========================="
    echo ""
    echo "For users to install Gonwatch, they can run:"
    echo ""
    echo "1. Remote installation (recommended):"
    echo "   curl -fsSL https://raw.githubusercontent.com/$(git config remote.origin.url | sed 's/.*github.com[:/]\([^/]*\)\/.*/\1/')/gonwatch/main/scripts/remote-install.sh | bash"
    echo ""
    echo "2. Local installation (after cloning):"
    echo "   git clone $(git config remote.origin.url)"
    echo "   cd gonwatch"
    echo "   ./install.sh"
    echo ""
    echo "3. System-wide installation:"
    echo "   curl -fsSL https://raw.githubusercontent.com/$(git config remote.origin.url | sed 's/.*github.com[:/]\([^/]*\)\/.*/\1/')/gonwatch/main/scripts/remote-install.sh | bash -s -- --help"
    echo ""
}

main() {
    print_status "Gonwatch Upload Helper"
    echo "==========================="
    echo ""
    
    check_git_repo
    
    # Ask for confirmation
    echo "This will upload the installation scripts to GitHub."
    echo "The following files will be added:"
    echo "  - scripts/remote-install.sh"
    echo "  - install.sh"
    echo ""
    read -p "Continue? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        upload_files
        show_usage
    else
        print_status "Upload cancelled"
    fi
}

main "$@"