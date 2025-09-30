# Development Guide

## Prerequisites

- Go 1.25 or later
- golangci-lint
- make (optional, but recommended)

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/shunichi-ikebuchi/rename-sony-photos-directories.git
cd rename-sony-photos-directories
```

### Install Dependencies

```bash
go mod download
```

### Build

```bash
# Using make
make build

# Or directly with go
go build -o rename-sony-photos-directories ./cmd/rename-sony-photos-directories
```

## Project Structure

```
rename-sony-photos-directories/
├── cmd/
│   └── rename-sony-photos-directories/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go                  # Configuration management
│   │   └── config_test.go
│   ├── rename/
│   │   ├── rename.go                  # Directory renaming logic
│   │   └── rename_test.go
│   └── workflow/
│       ├── workflow.go                # Workflow operations
│       └── workflow_test.go
├── examples/                          # Usage examples
├── docs/                              # Documentation
├── Makefile                           # Build automation
├── go.mod                             # Go module definition
└── README.md                          # Project overview
```

## Running Tests

### All Tests

```bash
# Using make
make test

# Or directly with go
go test -v -cover ./...
```

### Specific Package

```bash
go test -v ./internal/config
go test -v ./internal/rename
go test -v ./internal/workflow
```

### With Coverage

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Code Quality

### Running Linter

```bash
# Using make
make lint

# Or directly
golangci-lint run
```

### Formatting Code

```bash
go fmt ./...
gofmt -w .
```

## Making Changes

### Development Workflow

1. Create a feature branch
   ```bash
   git checkout -b feature/my-feature
   ```

2. Make your changes

3. Run tests
   ```bash
   make test
   ```

4. Run linter
   ```bash
   make lint
   ```

5. Commit your changes
   ```bash
   git add .
   git commit -m "feat: add my feature"
   ```

6. Push and create pull request
   ```bash
   git push origin feature/my-feature
   ```

### Commit Message Convention

Follow conventional commits:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

## Adding New Features

### Adding a New Package

1. Create directory under `internal/`
2. Add Go files with package documentation
3. Write tests
4. Update imports in other packages

### Adding Configuration Options

1. Update `Config` struct in `internal/config/config.go`
2. Add YAML tags
3. Update `config.yaml.example`
4. Update documentation in `docs/configuration.md`

## Testing

### Writing Tests

```go
func TestMyFunction(t *testing.T) {
    // Arrange
    input := "test"

    // Act
    result := MyFunction(input)

    // Assert
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### Test Coverage Goals

- Aim for >80% coverage
- All public functions should have tests
- Include edge cases and error conditions

## Debugging

### Enable Verbose Logging

```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### Debugging Tests

```bash
go test -v -run TestSpecificTest ./internal/rename
```

### Using Delve Debugger

```bash
dlv test ./internal/rename
```

## Release Process

### Creating a Release

1. Update version in relevant files
2. Update CHANGELOG (if using)
3. Commit changes
4. Create and push tag
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

5. GitHub Actions will automatically:
   - Run tests
   - Run linter
   - Build binaries
   - Create GitHub release

### Manual Release

```bash
goreleaser release --clean
```

## Continuous Integration

The project uses GitHub Actions for CI/CD:

- `.github/workflows/ci.yml` - Runs on every push/PR
- `.github/workflows/release.yml` - Runs on tag push

## Tips and Best Practices

1. **Keep Functions Small**: Each function should do one thing well
2. **Write Tests First**: TDD approach helps design better APIs
3. **Document Public APIs**: All exported functions need documentation
4. **Handle Errors**: Always check and handle errors appropriately
5. **Use Dry Run**: Implement dry-run mode for destructive operations

## Common Tasks

### Adding a New Command-Line Flag

1. Add flag in `cmd/rename-sony-photos-directories/main.go`
2. Update help text
3. Update `docs/usage.md`
4. Update README

### Updating Dependencies

```bash
go get -u ./...
go mod tidy
```

### Benchmarking

```go
func BenchmarkMyFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        MyFunction()
    }
}
```

Run benchmarks:
```bash
go test -bench=. ./...
```

## Getting Help

- Check existing issues on GitHub
- Review documentation in `docs/`
- Ask questions in GitHub Discussions
