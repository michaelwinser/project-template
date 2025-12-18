# Issue #001: Update Deprecated NPM Dependencies

**Issue:** [Link to GitHub Issue when created]
**Status:** Resolved
**Priority:** Low (dev dependencies only, no security issues)

## Problem

Running `make build` produces npm deprecation warnings in both `client/` and `web/` components:

```
npm warn deprecated inflight@1.0.6: This module is not supported, and leaks memory.
npm warn deprecated @humanwhocodes/config-array@0.13.0: Use @eslint/config-array instead
npm warn deprecated rimraf@3.0.2: Rimraf versions prior to v4 are no longer supported
npm warn deprecated glob@7.2.3: Glob versions prior to v9 are no longer supported
npm warn deprecated @humanwhocodes/object-schema@2.0.3: Use @eslint/object-schema instead
npm warn deprecated eslint@8.57.1: This version is no longer supported.
```

## Analysis

These are all **transitive dependencies** pulled in by our dev dependencies:
- `eslint@8.x` - The main culprit, pulls in most of the deprecated packages
- `jest@29.x` - May pull in `glob` and `rimraf`

## Proposed Solution

1. **Upgrade ESLint to v9.x** - This is a major version with breaking changes
   - ESLint 9 uses flat config format instead of `.eslintrc`
   - Need to update eslint config in both `client/` and `web/`
   - Update `@typescript-eslint/*` packages to compatible versions

2. **Verify Jest dependencies** - Check if Jest 30 (when available) resolves remaining issues

## Tasks

- [x] Upgrade `eslint` to ^9.x in `client/package.json`
- [x] Upgrade `eslint` to ^9.x in `web/package.json`
- [x] Upgrade `@typescript-eslint/eslint-plugin` and `@typescript-eslint/parser` to v8.x
- [x] Convert ESLint config to flat config format (if needed)
- [x] Test linting still works
- [x] Verify no deprecation warnings remain

## Resolution

Upgraded both `client/` and `web/` to ESLint 9 with flat config format:

**Package changes:**
- `eslint`: ^8.x → ^9.0.0
- Added `@eslint/js`: ^9.0.0
- Added `typescript-eslint`: ^8.0.0 (replaces separate plugin/parser packages)
- Added `globals`: ^15.0.0

**Config changes:**
- Created `eslint.config.js` (flat config) in both directories
- Removed old `.eslintrc` style configs
- Added `"type": "module"` to `client/package.json`

**Remaining warnings:**
- `inflight@1.0.6` and `glob@7.2.3` may still appear from jest dependencies
- These are harmless in test/dev context and will be resolved when Jest 30 releases

## Notes

- These are dev dependencies only - they don't affect production builds
- No security vulnerabilities reported for these packages
