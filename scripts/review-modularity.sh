#!/bin/bash
# Modularity review - validates architectural boundaries

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=== Modularity Review ==="
echo ""

ERRORS=0

# Function to check for forbidden imports
check_forbidden_imports() {
    local dir="$1"
    local forbidden_pattern="$2"
    local description="$3"

    if [ ! -d "$dir" ]; then
        return 0
    fi

    # Find all TypeScript files and check for forbidden imports
    while IFS= read -r file; do
        if grep -qE "$forbidden_pattern" "$file" 2>/dev/null; then
            echo "ERROR: $description"
            echo "  File: $file"
            grep -nE "$forbidden_pattern" "$file" | head -5
            echo ""
            ERRORS=$((ERRORS + 1))
        fi
    done < <(find "$dir" -name "*.ts" -type f 2>/dev/null)
}

echo "Checking web client MVC boundaries..."
echo ""

WEB_SRC="$PROJECT_ROOT/web/src"

if [ -d "$WEB_SRC" ]; then
    # Models should not import from views/ or controllers/
    echo "Checking models/ imports..."
    check_forbidden_imports \
        "$WEB_SRC/models" \
        "from ['\"](\.\./views/|\.\./controllers/)" \
        "Model importing from views/ or controllers/ (Models should be independent)"

    # Views should not import from models/ or controllers/
    echo "Checking views/ imports..."
    check_forbidden_imports \
        "$WEB_SRC/views" \
        "from ['\"](\.\./models/|\.\./controllers/)" \
        "View importing from models/ or controllers/ (Views receive data via render(), not direct imports)"

    # lib/ should not import from models/, views/, or controllers/
    echo "Checking lib/ imports..."
    check_forbidden_imports \
        "$WEB_SRC/lib" \
        "from ['\"](\.\./models/|\.\./views/|\.\./controllers/)" \
        "lib/ importing from models/, views/, or controllers/ (lib should be independent)"
else
    echo "Web source directory not found: $WEB_SRC"
fi

echo ""
echo "Checking server modularity..."
echo ""

SERVER_SRC="$PROJECT_ROOT/server/internal"

if [ -d "$SERVER_SRC" ]; then
    # handlers/ should not import from cli/
    echo "Checking server/internal imports..."

    # Check that internal packages don't have circular dependencies
    # (basic check - handlers shouldn't import from cmd/)
    if grep -rE "project-template/server/cmd" "$SERVER_SRC" 2>/dev/null; then
        echo "ERROR: internal/ packages should not import from cmd/"
        ERRORS=$((ERRORS + 1))
    fi
fi

echo ""
echo "Checking cross-component boundaries..."
echo ""

# Server should not import from client, cli, or web
if [ -d "$PROJECT_ROOT/server" ]; then
    if grep -rE "project-template/(client|cli|web)" "$PROJECT_ROOT/server" 2>/dev/null; then
        echo "ERROR: server/ should not import from client/, cli/, or web/"
        ERRORS=$((ERRORS + 1))
    fi
fi

# Client should not import from server, cli, or web
if [ -d "$PROJECT_ROOT/client" ]; then
    if grep -rE "project-template/(server|cli|web)" "$PROJECT_ROOT/client/src" 2>/dev/null; then
        echo "ERROR: client/ should not import from server/, cli/, or web/"
        ERRORS=$((ERRORS + 1))
    fi
fi

echo ""

if [ $ERRORS -eq 0 ]; then
    echo "Modularity review passed - no violations found"
    exit 0
else
    echo "Modularity review failed - $ERRORS violation(s) found"
    exit 1
fi
