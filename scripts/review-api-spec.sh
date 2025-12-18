#!/bin/bash
# API spec review - validates implementation matches spec

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=== API Spec Review ==="
echo ""

ERRORS=0
WARNINGS=0

SPEC_FILE="$PROJECT_ROOT/specs/api.yaml"
SERVER_MAIN="$PROJECT_ROOT/server/cmd/server/main.go"

if [ ! -f "$SPEC_FILE" ]; then
    echo "WARNING: API spec not found at $SPEC_FILE"
    exit 0
fi

if [ ! -f "$SERVER_MAIN" ]; then
    echo "WARNING: Server main.go not found at $SERVER_MAIN"
    exit 0
fi

echo "Checking API spec endpoints..."
echo ""

# Extract endpoints from OpenAPI spec (paths that start with /)
SPEC_ENDPOINTS=$(grep -E '^\s+/[a-zA-Z]' "$SPEC_FILE" | sed 's/:$//' | tr -d ' ' | sort -u)

echo "Endpoints defined in spec:"
echo "$SPEC_ENDPOINTS" | sed 's/^/  /'
echo ""

# Check each endpoint has a handler registered
echo "Checking server handlers..."
echo ""

for endpoint in $SPEC_ENDPOINTS; do
    # Look for HandleFunc or Handle with this path
    if ! grep -q "\"$endpoint\"" "$SERVER_MAIN" 2>/dev/null; then
        echo "WARNING: Endpoint '$endpoint' not found in server registration"
        WARNINGS=$((WARNINGS + 1))
    fi
done

# Check for handlers not in spec (potential undocumented endpoints)
echo "Checking for undocumented endpoints..."
echo ""

# Extract registered routes from main.go
SERVER_ROUTES=$(grep -oE 'HandleFunc\("[^"]+"|Handle\("[^"]+"|mux\.Handle[^(]*\("[^"]+"' "$SERVER_MAIN" 2>/dev/null | grep -oE '"/[^"]+"' | tr -d '"' | sort -u) || true

for route in $SERVER_ROUTES; do
    # Skip static routes
    if [[ "$route" == "/static/"* ]] || [[ "$route" == "/" ]]; then
        continue
    fi

    if ! echo "$SPEC_ENDPOINTS" | grep -qF "$route"; then
        echo "WARNING: Route '$route' in server but not in API spec"
        WARNINGS=$((WARNINGS + 1))
    fi
done

# Check that spec file is valid YAML (basic check)
echo "Validating spec format..."
echo ""

if ! grep -q "^openapi:" "$SPEC_FILE"; then
    echo "ERROR: Missing 'openapi:' version declaration"
    ERRORS=$((ERRORS + 1))
fi

if ! grep -q "^paths:" "$SPEC_FILE"; then
    echo "ERROR: Missing 'paths:' section"
    ERRORS=$((ERRORS + 1))
fi

if ! grep -q "^components:" "$SPEC_FILE"; then
    echo "WARNING: Missing 'components:' section (schemas)"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "=== API Spec Review Summary ==="

if [ $ERRORS -gt 0 ]; then
    echo "FAILED: $ERRORS error(s), $WARNINGS warning(s)"
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo "PASSED with $WARNINGS warning(s)"
    exit 0
else
    echo "PASSED: API spec and implementation are aligned"
    exit 0
fi
