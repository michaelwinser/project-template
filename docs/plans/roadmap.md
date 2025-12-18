# Project Template Roadmap

## Completed Phases

### Phase 1: Skeleton + Auth Vertical
- Go server with Google OAuth
- TypeScript client library
- Go CLI with --config flag
- OpenAPI spec
- Docker setup

### Phase 2: Web Client + MVC Template
- Vanilla TypeScript MVC framework
- Auth feature (model, view, controller)
- Modularity linter enforcing boundaries
- Static file serving from Go server

### Phase 3: Adaptive Testing
- Unit tests for server (handlers, auth, middleware, session)
- Unit tests for CLI (config)
- Adaptive test runner based on git changes

### Phase 4: Review Automation
- Security review (hardcoded secrets, sensitive files)
- API spec review (endpoint alignment)
- Reviewable moments detection (test+code, cross-component)

### Phase 5: Logging Infrastructure
- Structured JSON logging with env configuration
- Client log upload endpoint (POST /api/logs)
- Enhanced request logging (size, user agent, duration)

### Phase 6: Docker Infrastructure
- Unified production Dockerfile (multi-stage: Node → Go → Alpine runtime)
- Self-contained 24MB production image with baked-in web assets
- Development container (Go 1.23, Node.js 20, Make, git)
- VS Code devcontainer.json for seamless IDE integration
- Docker Compose with production and dev services
- Make targets: docker-dev, docker-shell, docker-build, docker-clean

### Phase 7: Claude-Powered Code Review Agent
- `/review` slash command for Claude Code
- Two-phase approach: triage (is it significant?) then deep review
- Semantic analysis: security, architecture, logic, quality
- Complements fast bash scripts (`make review`) for CI

## Upcoming Phases

(None currently planned)
