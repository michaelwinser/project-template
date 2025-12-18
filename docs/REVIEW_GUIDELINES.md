# Review Guidelines

This document defines review processes that can be run manually or via `make review-*`.

## Quick Reference

| Review Type | Command | When to Run |
|-------------|---------|-------------|
| Modularity | `make review-modularity` | Changes to server, client, or web code |
| Security | `make review-security` | Any code changes |
| API Spec | `make review-api-spec` | Changes to specs/ or API implementation |
| All | `make review-all` | Before major releases or PRs |
| Adaptive | `make review` | Automatically selects based on changes |

---

## Modularity Review

**Purpose:** Ensure architectural boundaries are respected.

### Checklist

- [ ] **Server isolation** - Server code doesn't import client or web code
- [ ] **Client isolation** - Client library doesn't import server internals
- [ ] **CLI thin wrapper** - CLI only calls client library, no direct API calls
- [ ] **No circular dependencies** - Components depend in one direction only
- [ ] **MVC boundaries** (web) - View doesn't import Model internals, Controller bridges

### Automated Checks

The modularity linter (`scripts/review-modularity.sh`) checks:
- Import patterns across component boundaries
- Circular dependency detection
- MVC layer violations

---

## Security Review

**Purpose:** Catch common security vulnerabilities.

### Checklist

#### Input Validation
- [ ] All user input is validated before use
- [ ] No SQL injection (parameterized queries only)
- [ ] No command injection (avoid shell execution with user input)
- [ ] No path traversal (validate file paths)

#### Authentication & Sessions
- [ ] Session cookies are HttpOnly and Secure (in production)
- [ ] Session tokens are cryptographically random
- [ ] Session expiration is enforced
- [ ] Logout invalidates session server-side

#### Secrets
- [ ] No hardcoded secrets in code
- [ ] Secrets loaded from environment variables
- [ ] No secrets in logs
- [ ] `.env` files are gitignored

#### Output Encoding
- [ ] HTML output is escaped (prevent XSS)
- [ ] JSON responses use proper content-type
- [ ] Error messages don't leak internal details

### Automated Checks

The security review (`scripts/review-security.sh`) checks:
- Hardcoded secret patterns
- SQL string concatenation
- Shell command construction with variables
- Missing input validation on handlers

---

## API Spec Review

**Purpose:** Ensure implementation matches the OpenAPI spec.

### Checklist

- [ ] All spec endpoints are implemented
- [ ] Request validation matches spec schemas
- [ ] Response formats match spec schemas
- [ ] Error responses use defined error schema
- [ ] No undocumented endpoints exist

### Automated Checks

The API spec review (`scripts/review-api-spec.sh`) checks:
- Endpoint coverage (spec vs implementation)
- Schema validation
- Response type correctness

---

## Pre-PR Review

Before creating a pull request, run:

```bash
make test-all
make review-all
```

### Additional PR Checklist

- [ ] Tests pass
- [ ] No linter warnings
- [ ] CHANGELOG updated (if applicable)
- [ ] Documentation updated (if API changed)
- [ ] Commit messages are clear and descriptive

---

## Post-Implementation Review

After completing an Issue, document:

1. **What changed** - List of commits with brief descriptions
2. **Deviations from plan** - What was different from the implementation plan?
3. **Lessons learned** - What would you do differently?

Add significant lessons to `docs/reviews/LESSONS_LEARNED.md`.
