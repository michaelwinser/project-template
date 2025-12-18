#!/bin/bash
# Review changes - detects "reviewable moments" that need human attention

set -e

echo "=== Reviewable Moments Check ==="
echo ""

MOMENTS=0

# Get list of changed files
CHANGED_FILES=$(git diff --name-only HEAD 2>/dev/null || echo "")
STAGED_FILES=$(git diff --cached --name-only 2>/dev/null || echo "")
ALL_CHANGED=$(echo "$CHANGED_FILES $STAGED_FILES" | tr ' ' '\n' | sort -u | grep -v '^$') || true

if [ -z "$ALL_CHANGED" ]; then
    echo "No changes detected."
    exit 0
fi

# Helper to check if any file matches pattern
has_changes_matching() {
    echo "$ALL_CHANGED" | grep -qE "$1"
}

# Helper to list matching files
list_matching() {
    echo "$ALL_CHANGED" | grep -E "$1" | sed 's/^/  /'
}

# 1. Test + code changed together
echo "Checking for test + code changes..."

TEST_FILES=$(echo "$ALL_CHANGED" | grep -E '(_test\.go|\.test\.ts|\.test\.js|\.spec\.ts)$') || true
CODE_FILES=$(echo "$ALL_CHANGED" | grep -E '\.(go|ts|js)$' | grep -vE '(_test\.go|\.test\.ts|\.test\.js|\.spec\.ts)$') || true

if [ -n "$TEST_FILES" ] && [ -n "$CODE_FILES" ]; then
    echo ""
    echo "REVIEWABLE MOMENT: Tests and implementation changed together"
    echo "This requires clear explanation of why tests needed to change."
    echo ""
    echo "Test files changed:"
    echo "$TEST_FILES" | sed 's/^/  /'
    echo ""
    echo "Implementation files changed:"
    echo "$CODE_FILES" | sed 's/^/  /'
    echo ""
    MOMENTS=$((MOMENTS + 1))
fi

# 2. API spec modified
echo "Checking for API spec changes..."

if has_changes_matching '^specs/'; then
    echo ""
    echo "REVIEWABLE MOMENT: API specification changed"
    echo "Changes to the API contract affect all consumers."
    echo ""
    echo "Changed spec files:"
    list_matching '^specs/'
    echo ""
    MOMENTS=$((MOMENTS + 1))
fi

# 3. Cross-component changes
echo "Checking for cross-component changes..."

COMPONENTS_CHANGED=0
[ -n "$(echo "$ALL_CHANGED" | grep -E '^server/')" ] && COMPONENTS_CHANGED=$((COMPONENTS_CHANGED + 1))
[ -n "$(echo "$ALL_CHANGED" | grep -E '^client/')" ] && COMPONENTS_CHANGED=$((COMPONENTS_CHANGED + 1))
[ -n "$(echo "$ALL_CHANGED" | grep -E '^cli/')" ] && COMPONENTS_CHANGED=$((COMPONENTS_CHANGED + 1))
[ -n "$(echo "$ALL_CHANGED" | grep -E '^web/')" ] && COMPONENTS_CHANGED=$((COMPONENTS_CHANGED + 1))

if [ $COMPONENTS_CHANGED -gt 1 ]; then
    echo ""
    echo "REVIEWABLE MOMENT: Multiple components changed ($COMPONENTS_CHANGED)"
    echo "Cross-component changes may indicate an API change or tight coupling."
    echo ""
    for component in server client cli web; do
        matching=$(echo "$ALL_CHANGED" | grep -E "^$component/") || true
        if [ -n "$matching" ]; then
            echo "$component/:"
            echo "$matching" | sed 's/^/  /'
        fi
    done
    echo ""
    MOMENTS=$((MOMENTS + 1))
fi

# 4. Security-sensitive files
echo "Checking for security-sensitive changes..."

SECURITY_PATTERNS='(auth|session|oauth|login|password|secret|credential|\.env)'
SECURITY_FILES=$(echo "$ALL_CHANGED" | grep -iE "$SECURITY_PATTERNS") || true

if [ -n "$SECURITY_FILES" ]; then
    echo ""
    echo "REVIEWABLE MOMENT: Security-sensitive files changed"
    echo "Extra scrutiny required for authentication/authorization changes."
    echo ""
    echo "Security-related files:"
    echo "$SECURITY_FILES" | sed 's/^/  /'
    echo ""
    MOMENTS=$((MOMENTS + 1))
fi

# 5. Configuration changes
echo "Checking for configuration changes..."

CONFIG_PATTERNS='(Makefile|package\.json|go\.mod|tsconfig|\.yaml$|\.yml$|config)'
CONFIG_FILES=$(echo "$ALL_CHANGED" | grep -E "$CONFIG_PATTERNS") || true

if [ -n "$CONFIG_FILES" ]; then
    echo ""
    echo "NOTE: Configuration files changed"
    echo "$CONFIG_FILES" | sed 's/^/  /'
    echo ""
    # Not counting as a "moment" but worth noting
fi

echo ""
echo "=== Reviewable Moments Summary ==="

if [ $MOMENTS -gt 0 ]; then
    echo "Found $MOMENTS reviewable moment(s) requiring attention"
    echo ""
    echo "Please ensure each moment has appropriate review/explanation."
    exit 0  # Exit 0 since these are prompts for review, not failures
else
    echo "No special reviewable moments detected"
    exit 0
fi
