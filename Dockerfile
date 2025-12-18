# Production Dockerfile - unified build for server + web assets
#
# Build: docker build -t project-template .
# Run:   docker run -p 8080:8080 project-template

# =============================================================================
# Stage 1: Build web assets
# =============================================================================
FROM node:20-alpine AS web-builder

WORKDIR /app/web

# Install dependencies first (better caching)
COPY web/package.json web/package-lock.json ./
RUN npm ci --production=false

# Copy source and build
COPY web/ ./
RUN npm run build

# =============================================================================
# Stage 2: Build Go server
# =============================================================================
FROM golang:1.23-alpine AS server-builder

WORKDIR /app

# Install git for go mod (some deps may need it)
RUN apk add --no-cache git

# Download dependencies first (better caching)
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy source and build
COPY server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# =============================================================================
# Stage 3: Production runtime
# =============================================================================
FROM alpine:3.19

# Install ca-certificates for HTTPS and wget for healthcheck
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Create non-root user
RUN adduser -D -g '' appuser

# Copy server binary
COPY --from=server-builder /server .

# Copy web assets and fix permissions
COPY --from=web-builder /app/web/dist ./web/dist
COPY --from=web-builder /app/web/public ./web/public
RUN chown -R appuser:appuser /app && chmod -R a+r /app/web

USER appuser

# Server configuration
ENV PORT=8080
ENV HOST=0.0.0.0
ENV WEB_DIR=/app/web
ENV ENV=production
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/health || exit 1

CMD ["./server"]
