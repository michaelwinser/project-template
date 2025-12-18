# Implementation Log

This document chronicles the implementation of project-template, a reference architecture for building web applications with Go, TypeScript, and Claude Code.

## Overview

| Phase | Description | Commit |
|-------|-------------|--------|
| 1 | Skeleton + Auth Vertical | `310c053` |
| 2 | Web Client MVC + ESLint 9 | `4f8179a` |
| 3 | Adaptive Testing | `ff27be2` |
| 4 | Review Automation | `8297bae` |
| 5 | Logging Infrastructure | `359951e` |
| 6 | Docker Infrastructure | `baa913e`, `9c7b886` |
| 7 | Claude-Powered Code Review | `186d1ae` |

---

## Phase 1: Skeleton + Auth Vertical

**Goal:** Establish the foundational architecture with a working authentication flow.

**Components Created:**
- **server/** - Go HTTP server with Google OAuth
  - `cmd/server/main.go` - Entry point, route registration
  - `internal/handlers/` - Health, auth, static file handlers
  - `internal/auth/` - OAuth provider, session management
  - `internal/middleware/` - Request ID, logging middleware
- **client/** - TypeScript client library
  - `src/client.ts` - HTTP client for server API
  - `src/types.ts` - Shared type definitions
- **cli/** - Go CLI wrapping client functionality
  - `cmd/cli/main.go` - Command dispatch
  - `internal/config/` - Configuration management
- **specs/api.yaml** - OpenAPI specification

**Key Decisions:**
- Session-based auth (not JWT) for simplicity
- Dev mode accepts any user without real OAuth
- TypeScript client is the source of truth for API types

---

## Phase 2: Web Client MVC + ESLint 9

**Goal:** Add browser-based UI with clean architectural boundaries.

**Components Created:**
- **web/** - Vanilla TypeScript MVC framework
  - `src/lib/` - Base Model, View, Controller classes
  - `src/models/auth.ts` - Auth state management
  - `src/views/auth.ts` - Login/logout UI rendering
  - `src/controllers/auth.ts` - Wires model and view
  - `src/app.ts` - Application composition root

**Architectural Boundaries:**
- Models: Can only import from `lib/`
- Views: Can only import from `lib/`
- Controllers: Can import from `lib/`, `models/`, `views/`
- Enforced by `scripts/review-modularity.sh`

**ESLint 9 Upgrade:**
- Migrated from `.eslintrc` to flat config (`eslint.config.js`)
- Created configs for both `client/` and `web/`

---

## Phase 3: Adaptive Testing

**Goal:** Comprehensive test coverage with smart test selection.

**Tests Created:**
- **server/**
  - `internal/handlers/health_test.go`
  - `internal/handlers/auth_test.go`
  - `internal/middleware/request_id_test.go`
  - `internal/auth/session_test.go`
- **cli/**
  - `internal/config/config_test.go`
- **client/**
  - `src/client.test.ts`

**Adaptive Test Runner:**
- `scripts/adaptive-test.sh` - Analyzes git changes
- Only runs tests for modified components
- Falls back to full suite if no git history

**Fixes Applied:**
- Added missing `devMode` parameter to auth tests
- Renamed `jest.config.js` → `jest.config.cjs` for ESM compatibility
- Fixed double mock issue in client tests

---

## Phase 4: Review Automation

**Goal:** Automated code review for common issues and architectural violations.

**Scripts Created:**
- `scripts/review-security.sh` - Scans for:
  - Hardcoded secrets and API keys
  - Private keys in source
  - Debug logging with sensitive data
  - SQL injection patterns
- `scripts/review-api-spec.sh` - Validates:
  - All spec endpoints are implemented
  - No undocumented endpoints exist
  - Spec file structure is valid
- `scripts/review-changes.sh` - Detects reviewable moments:
  - Tests + code changed together
  - API spec modifications
  - Cross-component changes
  - Security-sensitive file changes
- `scripts/adaptive-review.sh` - Orchestrates reviews based on changes

**Makefile Targets:**
- `make review` - Adaptive reviews (CI-friendly)
- `make review-all` - Full review suite

---

## Phase 5: Logging Infrastructure

**Goal:** Production-ready structured logging with client log aggregation.

**Enhancements:**
- `internal/logging/logger.go`:
  - `NewFromEnv()` - Configure from LOG_LEVEL environment
  - `ParseLevel()` - String to level conversion
  - `LogEntry()` - Forward pre-structured entries
  - Added `Source` field to Entry struct
- `internal/handlers/logs.go`:
  - `POST /api/logs` - Client log upload endpoint
  - Validates and forwards client logs with source attribution
- Enhanced request middleware:
  - Response size tracking
  - User agent logging
  - Duration in milliseconds

**Tests Added:**
- `internal/logging/logger_test.go`
- `internal/handlers/logs_test.go`

---

## Phase 6: Docker Infrastructure

**Goal:** Containerized development and production environments.

### Part 1: Development Container (`baa913e`)

**Files Created:**
- `Dockerfile.dev` - Ubuntu with Go 1.23, Node.js 20, npm, make, git
- `.devcontainer/devcontainer.json` - VS Code/DevPod integration

**docker-compose.yml Updates:**
- Added `dev` service with volume mounts
- Persistent caches for Go modules and node_modules

**Makefile Targets:**
- `make docker-dev` - Build dev container
- `make docker-shell` - Shell into dev container
- `make docker-build` - Build production images
- `make docker-clean` - Remove images and volumes

### Part 2: Production Build (`9c7b886`)

**Files Created:**
- `Dockerfile` - Unified multi-stage production build:
  1. Stage 1 (Node): Build web assets
  2. Stage 2 (Go): Compile server binary
  3. Stage 3 (Alpine): Minimal 24MB runtime

**Key Features:**
- Self-contained image (no volume mounts needed)
- Non-root user for security
- Health check included
- Stripped binary for size optimization

**Go Version Upgrade:**
- Updated all Dockerfiles from Go 1.21 → 1.23
- Required for `http.ServeFileFS` function

---

## Phase 7: Claude-Powered Code Review

**Goal:** Semantic code review using Claude's understanding of context.

**Files Created:**
- `.claude/commands/code-review.md` - Slash command for Claude Code

**Two-Phase Approach:**
1. **Triage**: Claude analyzes changes and decides significance
   - Security implications
   - API changes
   - Cross-component impact
   - Complexity and size
2. **Deep Review** (if significant):
   - Security analysis
   - Architecture evaluation
   - Logic review
   - Quality assessment

**Output Format:**
- Triage decision (SKIP/REVIEW)
- Findings by severity (Critical/Concerns/Suggestions)
- Summary with merge recommendation

**Usage:**
- `/code-review` - Review last commit
- `/code-review staged` - Review staged changes
- `/code-review HEAD~3` - Review last 3 commits
- Or just ask: "review the last commit"

---

## Architecture Summary

```
project-template/
├── server/          # Go HTTP server (port 8080)
│   ├── cmd/server/  # Entry point
│   └── internal/    # Handlers, auth, middleware, logging
├── client/          # TypeScript API client library
├── cli/             # Go CLI (wraps client)
├── web/             # Browser MVC client
│   ├── src/lib/     # Framework base classes
│   ├── src/models/  # State management
│   ├── src/views/   # UI rendering
│   └── src/controllers/  # Orchestration
├── specs/           # OpenAPI specification
├── scripts/         # Test and review automation
├── .claude/         # Claude Code commands
└── docs/            # Documentation and plans
```

## Key Commands

```bash
# Development
make build          # Build all components
make test           # Adaptive tests
make review         # Adaptive reviews

# Docker
make docker-dev     # Build dev container
make docker-shell   # Shell into dev container
make docker-build   # Build production image
docker-compose up   # Run production server

# Claude Code
/code-review        # Semantic code review
```

## Lessons Learned

1. **Go version compatibility** - `http.ServeFileFS` requires Go 1.22+. Always check function availability against target Go version.

2. **ESM + Jest** - Projects with `"type": "module"` need `jest.config.cjs` (CommonJS) for Jest configuration.

3. **Docker file permissions** - Files copied into Docker images retain host permissions. Use `chmod` in Dockerfile for consistent permissions.

4. **Adaptive tooling** - Git-aware test/review selection significantly speeds up development cycles while maintaining coverage.

5. **Slash command naming** - Check for existing commands before creating new ones. Renamed `/review` → `/code-review` to avoid conflict with built-in PR review command.
