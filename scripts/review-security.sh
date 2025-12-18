#!/bin/bash
# Security review - checks for common security issues

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=== Security Review ==="
echo ""

WARNINGS=0
ERRORS=0

# Function to check for patterns in files
check_pattern() {
    local pattern="$1"
    local description="$2"
    local severity="$3"  # "error" or "warning"
    local exclude_pattern="${4:-}"

    local cmd="grep -rn --include='*.go' --include='*.ts' --include='*.js' -E \"$pattern\" \"$PROJECT_ROOT\""
    if [ -n "$exclude_pattern" ]; then
        cmd="$cmd | grep -v -E \"$exclude_pattern\""
    fi

    local results
    results=$(eval "$cmd" 2>/dev/null | grep -v node_modules | grep -v "\.test\." | grep -v "_test\.go" | head -10) || true

    if [ -n "$results" ]; then
        if [ "$severity" = "error" ]; then
            echo "ERROR: $description"
            ERRORS=$((ERRORS + 1))
        else
            echo "WARNING: $description"
            WARNINGS=$((WARNINGS + 1))
        fi
        echo "$results" | sed 's|^|  |'
        echo ""
    fi
}

echo "Checking for hardcoded secrets..."
echo ""

# Check for hardcoded passwords/secrets (but not in config structs or env lookups)
check_pattern \
    '(password|passwd|secret|api_key|apikey|private_key)\s*[:=]\s*["\x27][^"\x27]{8,}' \
    "Possible hardcoded secret found" \
    "error" \
    "(getEnv|os\.Getenv|\.example|config\.go|_test)"

# Check for hardcoded API keys (common patterns)
check_pattern \
    '(sk-[a-zA-Z0-9]{20,}|AIza[a-zA-Z0-9_-]{35}|ghp_[a-zA-Z0-9]{36})' \
    "Possible API key found" \
    "error"

# Check for private keys
check_pattern \
    'BEGIN (RSA |DSA |EC |OPENSSH )?PRIVATE KEY' \
    "Private key found in source code" \
    "error"

echo "Checking for sensitive files..."
echo ""

# Check for .env files that aren't examples
if find "$PROJECT_ROOT" -name ".env" -not -name ".env.example" -not -path "*/node_modules/*" 2>/dev/null | grep -q .; then
    echo "WARNING: .env file found (should not be committed)"
    find "$PROJECT_ROOT" -name ".env" -not -name ".env.example" -not -path "*/node_modules/*" 2>/dev/null | sed 's|^|  |'
    echo ""
    WARNINGS=$((WARNINGS + 1))
fi

# Check for credential files
for pattern in "credentials.json" "service-account*.json" "*.pem" "*.key"; do
    if find "$PROJECT_ROOT" -name "$pattern" -not -path "*/node_modules/*" 2>/dev/null | grep -q .; then
        echo "WARNING: Sensitive file pattern '$pattern' found"
        find "$PROJECT_ROOT" -name "$pattern" -not -path "*/node_modules/*" 2>/dev/null | sed 's|^|  |'
        echo ""
        WARNINGS=$((WARNINGS + 1))
    fi
done

echo "Checking for debug code with sensitive data..."
echo ""

# Check for console.log with sensitive variable names
check_pattern \
    'console\.log\([^)]*\b(password|secret|token|cookie|session)\b' \
    "Debug logging may expose sensitive data" \
    "warning"

# Check for fmt.Print with sensitive variable names
check_pattern \
    'fmt\.(Print|Println|Printf)\([^)]*\b(password|secret|token|cookie|session)\b' \
    "Debug printing may expose sensitive data" \
    "warning"

echo "Checking for SQL injection patterns..."
echo ""

# Check for string concatenation in SQL (basic pattern)
check_pattern \
    '(SELECT|INSERT|UPDATE|DELETE).*\+.*\$|fmt\.Sprintf\([^)]*SELECT' \
    "Possible SQL injection vulnerability (string concatenation)" \
    "warning"

echo ""
echo "=== Security Review Summary ==="

if [ $ERRORS -gt 0 ]; then
    echo "FAILED: $ERRORS error(s), $WARNINGS warning(s)"
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo "PASSED with $WARNINGS warning(s)"
    exit 0
else
    echo "PASSED: No security issues found"
    exit 0
fi
