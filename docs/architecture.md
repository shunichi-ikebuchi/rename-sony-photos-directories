# Architecture

## Overview

This project follows standard Go project layout with clear separation of concerns.

## Package Structure

```
rename-sony-photos-directories/
├── cmd/
│   └── rename-sony-photos-directories/  # Application entry point
├── internal/                            # Private application code
│   ├── config/                          # Configuration management
│   ├── rename/                          # Core renaming logic
│   └── workflow/                        # Workflow operations
└── examples/                            # Usage examples
```

## Core Components

### 1. Configuration (`internal/config`)

**Responsibility**: Manage application configuration

**Key Functions**:
- `Default()` - Returns default configuration
- `Load(path)` - Loads configuration from file
- `Save(config, path)` - Saves configuration to file
- `LoadOrDefault()` - Loads config or returns default

**Design Decisions**:
- YAML format for human-readability
- Multiple search paths (current dir, ~/.config)
- Immutable config struct

### 2. Rename (`internal/rename`)

**Responsibility**: Directory renaming logic

**Key Functions**:
- `IsValidDateDir(name)` - Validates directory name format
- `ConvertDirName(name, century)` - Converts Sony format to yyyy-mm-dd
- `Directories(path)` - Renames all valid directories in path

**Design Decisions**:
- Pure functions where possible
- No side effects in validation functions
- Detailed logging for user feedback

**Sony Camera Format**:
```
Format: 0YYMMDD0
Example: 02512310 → 2025-12-31

Positions:
[0] Padding
[YY] Last 2 digits of year
[MM] Month
[DD] Day
[0] Padding
```

### 3. Workflow (`internal/workflow`)

**Responsibility**: Complex multi-step operations

**Key Functions**:
- `Run(config, dryRun)` - Full workflow (copy, rename, delete, eject)
- `RunBackupCleanup(config, dryRun)` - Backup cleanup workflow
- `CopyDir(src, dst, dryRun)` - Recursive directory copy
- `RemoveContents(dir, dryRun)` - Safe directory cleanup
- `EjectVolume(name, dryRun)` - Volume ejection (macOS)

**Design Decisions**:
- Dry-run support for safety
- Platform-specific code isolated
- Detailed error messages with context

### 4. Main (`cmd/rename-sony-photos-directories`)

**Responsibility**: CLI interface and orchestration

**Key Functions**:
- Parse command-line flags
- Load configuration
- Execute requested operation
- Handle errors and logging

**Design Decisions**:
- Thin main package (orchestration only)
- Business logic in internal packages
- Clear error messages for users

## Data Flow

### Rename-Only Mode

```
User Input (path)
    ↓
Load Config
    ↓
Validate Path
    ↓
Read Directories
    ↓
Filter Valid Names → IsValidDateDir()
    ↓
Convert Names → ConvertDirName()
    ↓
Rename Directories
    ↓
Log Results
```

### Full Workflow Mode

```
User Input (-workflow flag)
    ↓
Load Config
    ↓
Validate Paths (source, destination)
    ↓
Copy Source → Temp
    ↓
Rename in Temp
    ↓
Copy Temp → Destination
    ↓
Delete Source
    ↓
Clean Temp
    ↓
Eject Volume
    ↓
Log Completion
```

## Error Handling

### Strategy

1. **Validation**: Check inputs before operations
2. **Wrapping**: Use `fmt.Errorf` with `%w` for context
3. **Logging**: Log non-fatal errors, return fatal errors
4. **Recovery**: Continue processing other items on error

### Example

```go
if err := someOperation(); err != nil {
    return fmt.Errorf("failed to do operation: %w", err)
}
```

## Concurrency

**Current State**: Sequential processing

**Future Considerations**:
- Parallel directory processing
- Concurrent file copying
- Progress reporting

## Testing Strategy

### Unit Tests
- Test each function in isolation
- Mock external dependencies
- Cover edge cases and errors

### Integration Tests
- Test package interactions
- Use temporary directories
- Verify file system operations

### Test Coverage
- Target: >80% coverage
- Focus on business logic
- Test error paths

## Dependencies

### External
- `gopkg.in/yaml.v3` - YAML parsing

### Standard Library
- `os` - File system operations
- `path/filepath` - Path manipulation
- `log` - Logging
- `flag` - Command-line parsing
- `os/exec` - External commands (eject)

## Design Principles

### 1. Separation of Concerns
Each package has a single responsibility

### 2. Testability
- Pure functions where possible
- Dependencies can be mocked
- Small, focused functions

### 3. Explicit Error Handling
- No silent failures
- Detailed error messages
- Context in error wrapping

### 4. User Safety
- Dry-run mode for testing
- Validation before destructive operations
- Clear logging of actions

### 5. Cross-Platform Support
- Platform-specific code isolated
- Graceful degradation (e.g., eject on non-macOS)
- Standard Go conventions

## Future Enhancements

### Potential Improvements
1. **Progress Reporting**: Real-time progress bars
2. **Parallel Processing**: Concurrent operations
3. **Undo Support**: Backup before operations
4. **Web UI**: Browser-based interface
5. **Plugin System**: Custom rename strategies
6. **Database**: Track processed files
7. **Cloud Integration**: Upload to cloud storage

### Performance Optimizations
1. Buffer sizes for file copying
2. Parallel directory scanning
3. Incremental processing
4. Caching file stats

## Security Considerations

1. **File Permissions**: Config files set to 0600
2. **Path Validation**: Prevent directory traversal
3. **Input Sanitization**: Validate user inputs
4. **No Secrets**: Don't store sensitive data
5. **Audit Logging**: Log important operations

## Scalability

Current design handles:
- Thousands of directories
- Gigabytes of data
- Multiple concurrent users (separate configs)

Limitations:
- Sequential processing (can be slow for large datasets)
- Memory usage scales with directory count
- No distributed operation support
