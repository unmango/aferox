# aferox - Agent Instructions

## Project Overview
aferox is a Go library that extends `github.com/spf13/afero` with additional filesystem implementations and utilities. The project uses a multi-module structure with a main module and several specialized submodules.

## Module Structure
The repository contains multiple Go modules:
- **Main module** (`github.com/unmango/aferox`): Core utilities at the repository root
- **containerregistry** (`github.com/unmango/aferox/containerregistry`): Container image/layer filesystem wrappers
- **docker** (`github.com/unmango/aferox/docker`): Docker container filesystem access
- **github** (`github.com/unmango/aferox/github`): GitHub API as filesystem interface
- **gitignore** (`github.com/unmango/aferox/gitignore`): Gitignore-based filtering
- **protofs** (`github.com/unmango/aferox/protofs`): gRPC filesystem protocol

Each submodule has its own `go.mod`, `go.sum`, and separate CI workflow.

## Building
```bash
# Build everything with Nix (preferred)
make build

# Build with standard Go tools
go build ./...

# Build specific module with Nix
nix build .#aferox-docker
nix build .#aferox-github

# Check all Nix flake outputs
nix flake check --all-systems
```

## Testing
The project uses **Ginkgo** (BDD-style testing) with **Gomega** matchers.

### Running Tests
```bash
# Quick test (main module only, excludes E2E tests)
make test

# Test all modules including submodules
make test_all

# Run all tests with Ginkgo directly
go tool ginkgo run -r

# Run tests for specific package
go tool ginkgo run ./filter

# Generate new test file
cd <package-dir>
go tool ginkgo generate <filename>

# Bootstrap new test suite
cd <package-dir>
go tool ginkgo bootstrap
```

### Test Conventions
- Test files: `*_test.go`
- Test suites: `*_suite_test.go` (bootstrapped with Ginkgo)
- Use Ginkgo's BDD-style `Describe`, `Context`, `It` blocks
- Use Gomega matchers: `Expect()`, `To()`, `BeNil()`, etc.
- E2E tests are labeled and excluded locally: `--label-filter !E2E`
- CI runs with: `--github-output --race --trace --coverprofile=cover.profile`

## Dependency Management
```bash
# Update dependencies for main module
go mod tidy

# Update all module dependencies
make tidy

# Update Nix dependencies (after go.mod changes)
make import

# Or update specific module
go tool gomod2nix generate --dir containerregistry
```

## Development Workflow
1. Make code changes in relevant Go files
2. Write or update tests following Ginkgo/Gomega conventions
3. Run `go mod tidy` if you added/removed dependencies
4. Run `make test` for quick validation
5. Run `make test_all` before creating PR
6. Update `gomod2nix.toml` with `make import` if dependencies changed
7. Ensure all tests pass

## Key Patterns
- All filesystem implementations should satisfy `afero.Fs` interface
- Use `afero.Fs` for filesystem abstractions, not `os` package directly
- Prefer composition and wrapping of existing `afero.Fs` implementations
- Each submodule is self-contained to avoid dependency bloat in main module
- Readonly filesystems should implement only read operations
- Context-aware operations use the `context.Fs` interface from `context` package

## Common Commands
```bash
# Build all modules
make build

# Run tests (main module, excludes E2E)
make test

# Run all tests including submodules
make test_all

# Update all dependencies
make tidy

# Update Nix configs after dependency changes
make import

# Create new test suite
cd <directory>
go tool ginkgo bootstrap

# Generate test for specific file
cd <directory>
go tool ginkgo generate <name>
```

## Important Environment Variables
- `GOWORK=off`: Disabled in Makefile to avoid workspace mode issues
- `CI`: Automatically set in GitHub Actions, changes test behavior

## CI/CD
- Main CI workflow: `.github/workflows/ci.yml`
- Each submodule has its own workflow: `containerregistry.yml`, `docker.yml`, `github.yml`, `protofs.yml`
- Tests in CI include race detection, coverage reporting, and GitHub Actions integration
- Codecov uploads coverage reports from main module

## Testing Quick Reference
- **Create suite**: `go tool ginkgo bootstrap` in target directory
- **Create test**: `go tool ginkgo generate <name>` in target directory
- **Run tests**: `make test` (fast, local), `make test_all` (complete)
- **Test one package**: `go tool ginkgo run ./path/to/package`
- **Watch mode**: `go tool ginkgo watch -r`
