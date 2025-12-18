# Project Template

A reference architecture for building web applications with Go, TypeScript, and Claude Code. Designed as a living template that demonstrates best practices for authentication, testing, code review automation, and containerized deployment.

## Features

- **Go Server** - HTTP API with Google OAuth, session management, structured logging
- **TypeScript Client** - Type-safe API client library
- **Go CLI** - Command-line interface wrapping client functionality
- **Web Client** - Vanilla TypeScript MVC framework with enforced boundaries
- **Docker** - Multi-stage production builds, dev containers, VS Code integration
- **Adaptive Testing** - Git-aware test runner that only tests what changed
- **Review Automation** - Security scanning, API spec validation, modularity checks
- **Claude Code Integration** - Semantic code review with `/code-review` command

## Quick Start

### Option 1: Docker (Recommended)

```bash
# Build and run production server
make docker-build
docker-compose up server

# Open http://localhost:8080
```

### Option 2: Development Container

```bash
# Build dev container with all tools
make docker-dev

# Shell into container
make docker-shell

# Inside container: build and run
make build
cd server && go run ./cmd/server
```

### Option 3: Local Development

Requirements: Go 1.23+, Node.js 20+, npm

```bash
# Install dependencies
(cd client && npm install)
(cd web && npm install)

# Build everything
make build

# Run server (dev mode - no OAuth required)
cd server && go run ./cmd/server

# Open http://localhost:8080
```

## Project Structure

```
project-template/
├── server/           # Go HTTP server (port 8080)
│   ├── cmd/server/   # Entry point
│   └── internal/     # Handlers, auth, middleware, logging
├── client/           # TypeScript API client library
├── cli/              # Go CLI (wraps client)
├── web/              # Browser MVC client
│   ├── src/lib/      # Framework base classes
│   ├── src/models/   # State management
│   ├── src/views/    # UI rendering
│   └── src/controllers/
├── specs/            # OpenAPI specification
├── scripts/          # Test and review automation
├── .claude/          # Claude Code commands
└── docs/             # Documentation
```

## Make Targets

### Building

```bash
make build          # Build all components
make server-build   # Build server only
make client-build   # Build client library only
make web-build      # Build web client only
make cli-build      # Build CLI only
```

### Testing

```bash
make test           # Adaptive tests (based on git changes)
make test-all       # Run all tests
make test-ci        # CI test suite
```

### Code Review

```bash
make review         # Adaptive reviews (based on git changes)
make review-all     # Run all reviews
```

Reviews include:
- **Security** - Hardcoded secrets, API keys, sensitive files
- **API Spec** - Endpoint implementation matches OpenAPI spec
- **Modularity** - Architectural boundaries enforced
- **Changes** - Detects "reviewable moments" (test+code, cross-component)

### Docker

```bash
make docker-build   # Build production image (~24MB)
make docker-dev     # Build development container
make docker-shell   # Shell into dev container
make docker-clean   # Remove images and volumes

make up             # Start services (docker-compose up -d)
make down           # Stop services
```

## Claude Code Integration

This project includes a `/code-review` slash command for semantic code review:

```bash
/code-review           # Review last commit
/code-review staged    # Review staged changes
/code-review HEAD~3    # Review last 3 commits
```

The review performs:
1. **Triage** - Decides if changes warrant deep review
2. **Deep Review** - Analyzes security, architecture, logic, quality
3. **Structured Output** - Critical/Concerns/Suggestions with summary

## Configuration

### Environment Variables

```bash
# OAuth (optional in dev mode)
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/callback

# Server
PORT=8080
HOST=0.0.0.0
ENV=development          # or "production"
SESSION_SECRET=change-me # Required in production

# Logging
LOG_LEVEL=debug          # debug, info, warn, error
LOG_FORMAT=json          # json or text
```

### Dev Mode

When `ENV=development` (default), the server:
- Accepts any email for OAuth (no real Google auth needed)
- Enables verbose logging
- Serves from local filesystem

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/auth/login` | Initiate OAuth flow |
| GET | `/auth/callback` | OAuth callback |
| POST | `/auth/logout` | Log out |
| GET | `/auth/me` | Get current user |
| POST | `/api/logs` | Client log upload |

See [specs/api.yaml](./specs/api.yaml) for full OpenAPI specification.

## Web Client MVC

The web client uses a vanilla TypeScript MVC framework with strict boundaries:

| Layer | Can Import | Cannot Import |
|-------|------------|---------------|
| Models | `lib/` | views, controllers |
| Views | `lib/` | models, controllers |
| Controllers | `lib/`, `models/`, `views/` | - |

Boundaries are enforced by `make review` (modularity check).

## Documentation

- [CLAUDE.md](./CLAUDE.md) - Development conventions for Claude Code
- [docs/DEVELOPMENT.md](./docs/DEVELOPMENT.md) - Docker and dev environment setup
- [docs/IMPLEMENTATION_LOG.md](./docs/IMPLEMENTATION_LOG.md) - Build history and decisions
- [docs/plans/roadmap.md](./docs/plans/roadmap.md) - Project phases and status

## License

MIT
