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

### Phase 6: Development Container
- Dockerfile.dev with Go 1.21, Node.js 20, Make, git
- VS Code devcontainer.json for seamless IDE integration
- Docker Compose dev service with volume mounts
- Make targets: docker-dev, docker-shell, docker-build, docker-clean

## Upcoming Phases

### Phase 7: Claude-Powered Review Agent
- Semantic code review using Claude
- Understands context, not just patterns
- Optional deep review mode (`make review-deep`)
- Complements fast bash scripts for CI
