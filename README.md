# Project Template

A reference implementation for building web applications with:

- Go server with Google OAuth authentication
- TypeScript client library
- Go CLI
- Vanilla TypeScript web client (MVC architecture)

## Quick Start

```bash
# Start all services
make up

# Run tests
make test

# Run all tests (CI mode)
make test-ci

# Run reviews
make review
```

## Project Structure

```
project-template/
├── docs/           # PRDs, design docs, implementation plans, reviews
├── specs/          # OpenAPI and other machine-readable specs
├── generated/      # Generated code (from specs)
├── server/         # Go backend service
├── client/         # TypeScript client library
├── web/            # TypeScript web client (MVC)
├── cli/            # Go command-line interface
├── tests/          # Integration tests
├── scripts/        # Build and dev scripts
├── tools/          # Project-specific tooling
└── templates/      # Deployment config templates
```

## Development

See [CLAUDE.md](./CLAUDE.md) for development conventions and patterns.

## Documentation

- [PRD](./docs/prd.md) - Product requirements
- [Design docs](./docs/design/) - Technical design documents
- [Implementation plans](./docs/plans/) - Issue-linked implementation plans
- [Reviews](./docs/reviews/) - Post-implementation reviews and lessons learned
