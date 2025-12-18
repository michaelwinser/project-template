#!/bin/bash
# Adaptive test script - decides which tests to run based on what changed

set -e

# Get list of changed files (staged + unstaged)
CHANGED_FILES=$(git diff --name-only HEAD 2>/dev/null || echo "")
STAGED_FILES=$(git diff --cached --name-only 2>/dev/null || echo "")
ALL_CHANGED="$CHANGED_FILES $STAGED_FILES"

# If no git or no changes, run all tests
if [ -z "$ALL_CHANGED" ]; then
    echo "No changes detected, running all tests..."
    make test-all
    exit 0
fi

echo "Changed files:"
echo "$ALL_CHANGED" | tr ' ' '\n' | sort -u | grep -v '^$' || true
echo ""

# Determine which components are affected
RUN_SERVER=false
RUN_CLIENT=false
RUN_CLI=false
RUN_WEB=false
RUN_INTEGRATION=false

for file in $ALL_CHANGED; do
    case "$file" in
        server/*)
            RUN_SERVER=true
            RUN_INTEGRATION=true
            ;;
        client/*)
            RUN_CLIENT=true
            RUN_INTEGRATION=true
            ;;
        cli/*)
            RUN_CLI=true
            ;;
        web/*)
            RUN_WEB=true
            ;;
        specs/*)
            # Spec changes affect everything
            RUN_SERVER=true
            RUN_CLIENT=true
            RUN_CLI=true
            RUN_WEB=true
            RUN_INTEGRATION=true
            ;;
        tests/integration/*)
            RUN_INTEGRATION=true
            ;;
    esac
done

# Run affected tests
if [ "$RUN_SERVER" = true ]; then
    echo "Running server tests..."
    make -C server test
fi

if [ "$RUN_CLIENT" = true ]; then
    echo "Running client tests..."
    make -C client test
fi

if [ "$RUN_CLI" = true ]; then
    echo "Running CLI tests..."
    make -C cli test
fi

if [ "$RUN_WEB" = true ]; then
    echo "Running web tests..."
    make -C web test
fi

if [ "$RUN_INTEGRATION" = true ]; then
    echo "Running integration tests..."
    make -C tests/integration test
fi

echo "Adaptive tests complete."
