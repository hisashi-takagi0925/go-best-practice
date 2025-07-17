# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a Go best practices project that is currently in its initial setup phase. The repository is intended to demonstrate Go development best practices and patterns.

## Current State

The repository is newly initialized with:
- Basic git setup (main branch)
- Minimal README.md
- Claude AI integration configured in `.claude/settings.local.json`

## Setting Up the Go Project

Since this is a new Go project, you'll need to initialize it properly:

```bash
# Initialize Go module (replace with actual module path)
go mod init github.com/username/go-best-practice

# Create standard Go project structure
mkdir -p cmd internal pkg api scripts configs test docs examples
```

## Common Development Tasks

Once the project is set up with Go code:

```bash
# Download dependencies
go mod download

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test
go test -run TestName ./path/to/package

# Build the project
go build ./...

# Format code
go fmt ./...

# Vet code for common mistakes
go vet ./...
```

## Recommended Project Structure

When implementing features, follow this standard Go project layout:

```
/cmd          - Main applications (each subdirectory is one application)
/internal     - Private application and library code
/pkg          - Library code that can be used by external applications
/api          - API protocol definitions (OpenAPI specs, Protocol Buffers)
/scripts      - Scripts for build, install, analysis operations
/configs      - Configuration file templates or default configs
/test         - Additional external test apps and test data
/docs         - Design and user documents
/examples     - Examples for using the project
```

## Code Architecture Guidelines

1. **Package Design**: Keep packages small and focused on a single responsibility
2. **Interfaces**: Define interfaces in the package that uses them, not the package that implements them
3. **Error Handling**: Always handle errors explicitly; wrap errors with context
4. **Testing**: Write table-driven tests; use testify for assertions if needed
5. **Dependencies**: Use interfaces for external dependencies to enable testing

## Important Files to Create

When setting up this project, ensure these files exist:
- `.gitignore` - Use GitHub's Go gitignore template
- `go.mod` - Module definition
- `Makefile` - Build automation tasks
- `.golangci.yml` - Linting configuration

## Notes

- The repository currently has no Go code, so standard Go commands won't work until the module is initialized
- Follow Go's official style guide and effective Go principles
- Use conventional commit messages for version history