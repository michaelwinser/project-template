# Development Environment

This project includes a containerized development environment with all required build tools.

## Quick Start

```bash
# Build the dev container
make docker-dev

# Shell into the dev container
make docker-shell

# Inside the container, you can run all make commands:
make build
make test
```

## What's Included

The dev container includes:
- Go 1.21
- Node.js 20 LTS + npm
- Make
- git, curl, wget
- Non-root user (dev) with sudo access

## VS Code Dev Containers

This project supports VS Code Dev Containers (and DevPod):

1. Install the "Dev Containers" extension in VS Code
2. Open the project folder
3. Click "Reopen in Container" when prompted (or use Command Palette → "Dev Containers: Reopen in Container")

The container will:
- Mount your workspace at `/workspace`
- Forward port 8080
- Install Go, ESLint, Prettier, and Docker extensions
- Enable format-on-save

## Manual Docker Commands

If you prefer not to use make or VS Code:

```bash
# Build the dev container
docker-compose build dev

# Run a shell in the dev container
docker-compose run --rm dev

# Run a specific command
docker-compose run --rm dev make test

# Build production server image
docker-compose build server

# Start all services
docker-compose up -d

# Stop services and clean up volumes
docker-compose down -v
```

## Environment Variables

The dev container inherits environment variables from your shell. For OAuth to work:

```bash
export GOOGLE_CLIENT_ID=your-client-id
export GOOGLE_CLIENT_SECRET=your-client-secret
```

Or create a `.env` file in the project root (see `templates/.env.example`).

## Volume Mounts

The dev service mounts:
- `.:/workspace` - Your code (live edits)
- `go-mod-cache:/home/dev/go/pkg/mod` - Go module cache (persisted)
- `node-modules-client:/workspace/client/node_modules` - Client deps (persisted)
- `node-modules-web:/workspace/web/node_modules` - Web deps (persisted)

This means:
- Code changes are immediately visible in the container
- Go and npm dependencies are cached between container restarts
- First `npm install` or `go mod download` is slow, subsequent ones are fast

## Troubleshooting

### Permission issues with mounted files

The container runs as user `dev` (UID 1000). If your host user has a different UID, you may see permission issues. Solutions:

1. Run the container as root: `docker-compose run --rm -u root dev`
2. Rebuild the image with your UID: Edit `Dockerfile.dev` and change `USER_UID=1000` to your UID

### Clean rebuild

```bash
# Remove containers, volumes, and images
make docker-clean

# Rebuild from scratch
make docker-dev
```

### Port 8080 already in use

Stop any local server first, or change the port mapping in `docker-compose.yml`.
