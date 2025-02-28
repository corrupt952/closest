# Contributing to closest

Thank you for your interest in contributing to `closest`! This document provides guidelines and instructions for contributing to this project.

## Development Setup

### Prerequisites

- Git
- Go (we recommend using [aqua](https://github.com/aquaproj/aqua) for installation)

### Setting Up the Development Environment

1. **Clone the repository**

   ```sh
   git clone https://github.com/corrupt952/closest.git
   cd closest
   ```

2. **Install Go using aqua**

   If you have aqua installed:

   ```sh
   aqua i
   ```

   This will install the Go version specified in the `aqua.yaml` file.

3. **Install development dependencies**

   ```sh
   go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
   ```

## Development Workflow

### Making Changes

1. Create a new branch for your changes:

   ```sh
   git checkout -b feature/your-feature-name
   ```

2. Make your changes to the codebase.

3. Run tests to ensure your changes don't break existing functionality:

   ```sh
   go test -v ./...
   ```

4. Run the linter to ensure code quality:

   ```sh
   golangci-lint run
   ```

5. Build the binary locally to test your changes:

   ```sh
   go build
   ```

### Submitting Changes

1. Commit your changes with a descriptive commit message:

   ```sh
   git commit -m "feat: add your feature description"
   ```

   We follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.

2. Push your branch to GitHub:

   ```sh
   git push origin feature/your-feature-name
   ```

3. Create a Pull Request on GitHub.

## Release Process

Releases are automated using GitHub Actions and GoReleaser. Here's how to create a new release:

1. **Update version references** (if any) in the codebase.

2. **Create and push a new tag** with the version number:

   ```sh
   git tag v1.0.0
   git push origin v1.0.0
   ```

   This will automatically trigger the release workflow defined in `.github/workflows/release.yaml`.

3. **Verify the release** on the [GitHub Releases page](https://github.com/corrupt952/closest/releases).

### Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

## Project Structure

- `main.go` - Main application code
- `main_test.go` - Tests for the application
- `.github/workflows/` - GitHub Actions workflow definitions
- `.goreleaser.yml` - GoReleaser configuration

## Code Style

We follow standard Go code style and conventions:

- Use `gofmt` to format your code
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Write meaningful comments and documentation

## Getting Help

If you have questions or need help, please open an issue on GitHub.

Thank you for contributing to `closest`!
