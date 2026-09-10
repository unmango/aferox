# GitHub Copilot Instructions for aferox

## Project Overview
This is aferox, a Go library that extends `github.com/spf13/afero` with additional `afero.Fs` implementations and utility functions. The project consists of a main module and several submodules located in separate directories.

## Project Structure
- **Main module**: Core utilities and base functionality
- **Submodules**: `containerregistry`, `docker`, `github`, `gitignore`, `protofs`
- Each submodule is a separate Go module with its own `go.mod` file

## Technology Stack
- **Language**: Go 1.25.5
- **Testing Framework**: Ginkgo/Gomega (BDD-style testing)
- **Build System**: Make and Nix
- **Package Manager**: Go modules

## Building the Project
- Use `make build` to build all modules using Nix
- Use `go build ./...` to build the main module with standard Go tools
- Use `nix build .#aferox` to build the main module
- Use `nix flake check --all-systems` to check all modules

## Testing
- **Test framework**: Ginkgo (BDD-style testing with Gomega matchers)
- Use `make test` to run tests for the main module (excludes E2E tests locally)
- Use `make test_all` to run all tests including submodules
- Use `go tool ginkgo run -r` to run all tests recursively
- In CI, tests run with `--github-output --race --trace --coverprofile=cover.profile`
- Test files follow the pattern `*_test.go`
- Test suites use `*_suite_test.go` files (Ginkgo convention)
- Generate new test files with `go tool ginkgo generate <name>`
- Bootstrap new test suites with `go tool ginkgo bootstrap`

## Coding Conventions
- Follow standard Go conventions and idioms
- Use `afero.Fs` interface for filesystem abstractions
- Implement filesystem operations through the `afero` interface
- Maintain compatibility with the `afero.Fs` interface
- Each submodule should be self-contained with minimal dependencies

## Dependencies
- Use `go mod tidy` to update dependencies in the main module
- Use `make tidy` to update all module dependencies
- Use `make import` to update `gomod2nix.toml` files for Nix builds
- Submodule dependencies are managed separately in their respective directories

## Development Workflow
1. Make changes to relevant Go files
2. Run `go mod tidy` if dependencies changed
3. Run tests with `make test` for quick validation
4. Run `make test_all` before submitting PR
5. Ensure all tests pass before committing

## Key Packages
- **Main**: Core utilities like `filter.Fs`, `ignore.Fs`, utility functions
- **containerregistry**: Wrap container images/layers as filesystems
- **docker**: Docker container filesystem access
- **github**: GitHub API as filesystem interface
- **gitignore**: Gitignore-based filtering filesystem
- **protofs**: gRPC-based filesystem protocol
- **testing**: Mock/stub utilities for testing
- **writer**: Write-only filesystem that dumps to `io.Writer`
- **context**: Context-aware filesystem operations

## Important Notes
- The main module excludes submodule tests in CI (`--skip-package` flag)
- Each submodule has its own CI workflow
- Use `GOWORK=off` when building to avoid workspace mode issues
- E2E tests are excluded in local development (use `--label-filter !E2E`)
