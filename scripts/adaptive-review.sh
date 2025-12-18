#!/bin/bash
# Adaptive review script - decides which reviews to run based on what changed

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

REVIEW_FAILED=0

# Helper to run a review and track failures
run_review() {
    local script="$1"
    if ! "$script"; then
        REVIEW_FAILED=1
    fi
    echo ""
}

# Get list of changed files
CHANGED_FILES=$(git diff --name-only HEAD 2>/dev/null || echo "")
STAGED_FILES=$(git diff --cached --name-only 2>/dev/null || echo "")
ALL_CHANGED="$CHANGED_FILES $STAGED_FILES"

# If no git or no changes, run all reviews
if [ -z "$ALL_CHANGED" ]; then
    echo "No changes detected, running all reviews..."
    echo ""
    run_review "$SCRIPT_DIR/review-changes.sh"
    run_review "$SCRIPT_DIR/review-modularity.sh"
    run_review "$SCRIPT_DIR/review-security.sh"
    run_review "$SCRIPT_DIR/review-api-spec.sh"
    exit $REVIEW_FAILED
fi

echo "Changed files:"
echo "$ALL_CHANGED" | tr ' ' '\n' | sort -u | grep -v '^$' || true
echo ""

# Always run reviewable moments check when there are changes
echo "Running reviewable moments check..."
run_review "$SCRIPT_DIR/review-changes.sh"

# Determine which reviews to run
RUN_MODULARITY=false
RUN_SECURITY=false
RUN_API_SPEC=false

for file in $ALL_CHANGED; do
    case "$file" in
        server/*|client/*|web/*)
            RUN_MODULARITY=true
            ;;
        specs/*)
            RUN_API_SPEC=true
            ;;
    esac

    # Security review for any code changes
    case "$file" in
        *.go|*.ts|*.js)
            RUN_SECURITY=true
            ;;
    esac
done

# Run affected reviews
if [ "$RUN_MODULARITY" = true ]; then
    echo "Running modularity review..."
    run_review "$SCRIPT_DIR/review-modularity.sh"
fi

if [ "$RUN_SECURITY" = true ]; then
    echo "Running security review..."
    run_review "$SCRIPT_DIR/review-security.sh"
fi

if [ "$RUN_API_SPEC" = true ]; then
    echo "Running API spec review..."
    run_review "$SCRIPT_DIR/review-api-spec.sh"
fi

echo "=== Adaptive Review Complete ==="
if [ $REVIEW_FAILED -eq 1 ]; then
    echo "Some reviews failed - see above for details"
    exit 1
else
    echo "All reviews passed"
    exit 0
fi
