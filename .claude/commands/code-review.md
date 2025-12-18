# Code Review Agent

Perform intelligent code review with triage and deep analysis.

## Arguments
- No argument: Review the last commit (`HEAD~1..HEAD`)
- `staged`: Review staged but uncommitted changes
- `HEAD~N`: Review the last N commits (e.g., `HEAD~3`)
- Any git revision range: Review that specific range

## Instructions

### Step 1: Gather Context

First, determine what to review based on the argument:

1. **Parse the argument** from this command invocation
2. **Get the diff**:
   - If "staged": Run `git diff --cached`
   - If no argument: Run `git diff HEAD~1..HEAD`
   - Otherwise: Run `git diff {argument}`
3. **Get changed files**: Run `git diff --name-only` with the same range
4. **Get commit message(s)**: Run `git log --oneline` for the range

If there are no changes, report that and stop.

### Step 2: Triage

Analyze the changes and decide if they warrant deep review. Consider:

- **Security sensitivity**: Auth, sessions, passwords, tokens, encryption, user data
- **API changes**: Endpoints, request/response formats, OpenAPI spec
- **Cross-component impact**: Changes spanning server/, client/, cli/, web/
- **Complexity**: Non-trivial logic, error handling, concurrency
- **Size**: Large refactors or many files changed

**Triage decision:**
- **SKIP**: Minor changes (typos, formatting, comments, simple renames, dependency updates)
- **REVIEW**: Anything that could affect behavior, security, or architecture

Output your triage decision with a brief reason before proceeding.

### Step 3: Deep Review (if REVIEW decision)

Analyze the code changes across these dimensions:

#### Security
- Authentication/authorization changes
- Input validation and sanitization
- Secrets or credentials in code
- SQL injection, XSS, command injection risks
- Secure defaults and error messages

#### Architecture
- Module boundary violations (see CLAUDE.md for project boundaries)
- API contract changes without spec updates
- Tight coupling or dependency issues
- Patterns inconsistent with codebase

#### Logic
- Edge cases not handled
- Error handling gaps
- Race conditions or concurrency issues
- Off-by-one errors, null checks
- State management issues

#### Quality
- Code clarity and maintainability
- Test coverage implications
- Documentation needs
- Performance concerns

### Step 4: Output Format

Structure your review as follows:

```
## Code Review: [one-line summary of changes]

### Triage
**[SKIP/REVIEW]** - [brief reason]

---

### Findings

#### Critical
[Issues that must be fixed before merge. Security vulnerabilities, data loss risks, breaking changes.]

#### Concerns
[Issues that should be addressed. Logic errors, missing edge cases, architectural issues.]

#### Suggestions
[Nice-to-have improvements. Style, clarity, minor optimizations.]

---

### Summary
[2-3 sentence overall assessment. Is this ready to merge? What's the risk level?]
```

If SKIP was chosen, only output the Triage section.

## Examples

**Minor change (SKIP):**
```
## Code Review: Fix typo in README

### Triage
**SKIP** - Documentation-only change with no code impact.
```

**Significant change (REVIEW):**
```
## Code Review: Add session timeout handling

### Triage
**REVIEW** - Security-sensitive change affecting authentication flow.

---

### Findings

#### Critical
None identified.

#### Concerns
- `session.go:45`: Timeout value is hardcoded. Consider making this configurable.
- Missing test for the timeout edge case when session expires mid-request.

#### Suggestions
- Add debug logging when session timeout occurs for easier troubleshooting.

---

### Summary
This change correctly implements session timeouts but needs a test for the edge case.
Low risk overall - the timeout logic is straightforward. Ready to merge after adding the test.
```
