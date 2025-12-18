# Claude Code Guidelines

This document captures conventions and patterns for working on this project with Claude Code.

## Project Overview

A reference template for building web applications with:
- Go server with Google OAuth authentication
- TypeScript client library
- Go CLI (thin wrapper around client functionality)
- Vanilla TypeScript web client (MVC architecture) - *Phase 2*

## Directory Structure

```
project-template/
├── docs/           # PRDs, design docs, plans, reviews
├── specs/          # OpenAPI specs (api.yaml)
├── generated/      # Generated code (from specs)
├── server/         # Go backend (port 8080)
├── client/         # TypeScript client library
├── cli/            # Go CLI
├── web/            # TypeScript web client (Phase 2)
├── tests/          # Integration tests
├── scripts/        # Build and review scripts
├── tools/          # Project-specific tooling
└── templates/      # Deployment config templates
```

## Development Workflow

### Building

```bash
make build          # Build all components
make server-build   # Build just the server
make cli-build      # Build just the CLI
```

### Testing

```bash
make test           # Adaptive tests (based on what changed)
make test-all       # Full test suite
make test-ci        # CI test suite
```

### Running

```bash
make up             # Start services via docker-compose
make down           # Stop services
make server-run     # Run server directly (dev mode)
```

### Reviews

```bash
make review         # Adaptive reviews (based on what changed)
make review-all     # All reviews
```

## Key Conventions

### Debugging

**Debug data, not code.** When investigating bugs:
1. Read logs first → understand what actually happened
2. Read relevant code → understand why
3. Back to logs → verify the fix

Logs are structured JSON. Look for `request_id` to trace requests.

### Modularity

Architectural boundaries are enforced:
- Components (server, client, cli, web) are isolated
- API spec (`specs/api.yaml`) is the contract between components
- Changes that cross component boundaries require explicit design

### Reviewable Moments

These changes require extra scrutiny:
- **Changing tests + code together** - Why did the contract change?
- **Modifying API spec** - Affects all consumers
- **Cross-component changes** - Plan first, don't drive-by

### Issue Lifecycle

1. Create Issue describing the work
2. Create implementation plan in `docs/plans/issue-NNN-description.md`
3. Link Issue ↔ Plan bidirectionally
4. Do the work
5. On completion, add closing comment summarizing:
   - What actually changed (commits)
   - How it differed from plan
   - Lessons learned
6. Add significant lessons to `docs/reviews/LESSONS_LEARNED.md`

## API Spec

The OpenAPI spec lives at `specs/api.yaml`. It defines:
- `/health` - Health check
- `/auth/login` - Initiate OAuth flow
- `/auth/callback` - OAuth callback
- `/auth/logout` - Log out
- `/auth/me` - Get current user

When modifying the API:
1. Update the spec first
2. Update server implementation
3. Update client library
4. Update CLI if needed
5. Update integration tests

## Environment Variables

See `templates/.env.example` for all configuration options.

Key variables:
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` - OAuth credentials
- `SESSION_SECRET` - Session signing key (change in production!)
- `ENV=development` - Enables dev mode (accepts any user)

## Testing

### Unit Tests
- Go: `*_test.go` files alongside code
- TypeScript: `*.test.ts` files alongside code

### Integration Tests
- Located in `tests/integration/`
- Scenarios defined in `tests/integration/scenarios.yaml`
- Require running server: `TEST_SERVER_URL=http://localhost:8080`

## Common Tasks

### Add a new API endpoint

1. Add to `specs/api.yaml`
2. Add handler in `server/internal/handlers/`
3. Wire up in `server/cmd/server/main.go`
4. Add to client library `client/src/client.ts`
5. Add to CLI if user-facing
6. Add integration test

### Debug an auth issue

1. Check server logs for the request (filter by `request_id`)
2. Look for session-related entries
3. Verify cookie is being set/sent correctly
4. Check if session exists in memory (dev mode doesn't persist)
