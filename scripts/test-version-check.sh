#!/bin/bash

# Test script for Go version checking
# This validates the version comparison logic

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_pass() {
    echo -e "${GREEN}✅ PASS${NC}: $1"
}

print_fail() {
    echo -e "${RED}❌ FAIL${NC}: $1"
}

print_info() {
    echo -e "${YELLOW}ℹ️  INFO${NC}: $1"
}

# Version comparison function (copied from install script)
version_compare() {
    local version1=$1
    local version2=$2
    
    IFS='.' read -ra v1_parts <<< "$version1"
    IFS='.' read -ra v2_parts <<< "$version2"
    
    for i in {0..2}; do
        local v1_part=${v1_parts[$i]:-0}
        local v2_part=${v2_parts[$i]:-0}
        
        if [ "$v1_part" -gt "$v2_part" ]; then
            return 0
        elif [ "$v1_part" -lt "$v2_part" ]; then
            return 1
        fi
    done
    return 0
}

# Check Go version function (copied from install script)
check_go_version() {
    local go_version_str=$1
    local min_version="1.25.0"
    
    local clean_version=$(echo "$go_version_str" | sed 's/^go//' | sed 's/[^0-9.].*//')
    
    if version_compare "$clean_version" "$min_version"; then
        return 0
    else
        return 1
    fi
}

# Test cases
echo "Testing Go Version Validation"
echo "============================="

# Test version comparison function
echo -e "\n📋 Testing version_compare function:"
test_cases=(
    "1.25.0:1.25.0:0"
    "1.26.0:1.25.0:0" 
    "1.24.0:1.25.0:1"
    "2.0.0:1.25.0:0"
    "1.25.5:1.25.0:0"
    "1.20.0:1.25.0:1"
)

for case in "${test_cases[@]}"; do
    IFS=':' read -ra params <<< "$case"
    version1=${params[0]}
    version2=${params[1]}
    expected=${params[2]}
    
    if version_compare "$version1" "$version2"; then
        actual=0
    else
        actual=1
    fi
    
    if [ "$actual" -eq "$expected" ]; then
        print_pass "version_compare $version1 >= $version2"
    else
        print_fail "version_compare $version1 >= $version2 (expected $expected, got $actual)"
    fi
done

# Test Go version checking function
echo -e "\n📋 Testing check_go_version function:"
go_test_cases=(
    "go1.25.0:0"
    "go1.26.1:0"
    "go1.24.9:1"
    "go1.25.0rc1:0"
    "go1.24.0beta2:1"
    "go2.0.0-alpha1:0"
)

for case in "${go_test_cases[@]}"; do
    IFS=':' read -ra params <<< "$case"
    go_version=${params[0]}
    expected=${params[1]}
    
    if check_go_version "$go_version"; then
        actual=0
    else
        actual=1
    fi
    
    if [ "$actual" -eq "$expected" ]; then
        print_pass "check_go_version $go_version"
    else
        print_fail "check_go_version $go_version (expected $expected, got $actual)"
    fi
done

# Test with real system Go if available
echo -e "\n📋 Testing with system Go:"
if command -v go >/dev/null 2>&1; then
    system_go_version=$(go version)
    print_info "System Go: $system_go_version"
    
    if check_go_version "$system_go_version"; then
        print_pass "System Go version meets requirements"
    else
        print_fail "System Go version does NOT meet requirements"
    fi
else
    print_info "No system Go installation found"
fi

# Test the actual install script (dry run)
echo -e "\n📋 Testing install script logic:"
echo "Creating temporary test environment..."

TEMP_TEST_DIR=$(mktemp -d)
trap "rm -rf $TEMP_TEST_DIR" EXIT

# Create mock Go binary that returns old version
mkdir -p "$TEMP_TEST_DIR/bin"
cat > "$TEMP_TEST_DIR/bin/go" << 'EOF'
#!/bin/bash
echo "go version go1.24.0 linux/amd64"
EOF
chmod +x "$TEMP_TEST_DIR/bin/go"

# Add to PATH temporarily
export PATH="$TEMP_TEST_DIR/bin:$PATH"

print_info "Testing with mock Go 1.24.0..."
if command -v go >/dev/null 2>&1; then
    mock_version=$(go version)
    print_info "Mock Go: $mock_version"
    
    if check_go_version "$mock_version"; then
        print_fail "Mock old Go should fail validation"
    else
        print_pass "Mock old Go correctly fails validation"
    fi
fi

echo -e "\n✨ Version validation testing complete!"
echo ""
echo "Next steps:"
echo "1. Test actual installation in clean VM"
echo "2. Test with different pre-installed Go versions"
echo "3. Verify the download and installation process"