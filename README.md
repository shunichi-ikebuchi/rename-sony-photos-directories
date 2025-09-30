# rename-sony-photos-directories

[![CI](https://github.com/shunichi-ikebuchi/rename-sony-photos-directories/workflows/CI/badge.svg)](https://github.com/shunichi-ikebuchi/rename-sony-photos-directories/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/shunichi-ikebuchi/rename-sony-photos-directories)](https://goreportcard.com/report/github.com/shunichi-ikebuchi/rename-sony-photos-directories)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/shunichi-ikebuchi/rename-sony-photos-directories)](https://github.com/shunichi-ikebuchi/rename-sony-photos-directories)

A Go program to convert Sony digital camera [date format folders](https://www.sony.jp/ServiceArea/impdf/pdf/44879440M.w-JP/jp/contents/TP0000220296.html) from the format `0YYMMDD0` to the more readable `yyyy-mm-dd` format.

## Features

- Convert Sony camera directory names (e.g., `02512310` → `2025-12-31`)
- Configurable paths via YAML configuration file
- Command-line flags for flexible usage
- Cross-platform support (Linux, macOS, Windows)
- Comprehensive test coverage

## Installation

### Using Go install

```bash
go install github.com/shunichi-ikebuchi/rename-sony-photos-directories/cmd/rename-sony-photos-directories@latest
```

### From source

```bash
git clone https://github.com/shunichi-ikebuchi/rename-sony-photos-directories.git
cd rename-sony-photos-directories
go build
```

### Using GoReleaser (for releases)

Pre-built binaries are available on the [releases page](https://github.com/shunichi-ikebuchi/rename-sony-photos-directories/releases).

## Configuration

### Configuration File

The program looks for a configuration file in the following order:

1. `./config.yaml` (current directory)
2. `~/.config/rename-sony-photos/config.yaml` (XDG config directory)

Create a configuration file from the example:

```bash
cp config.yaml.example config.yaml
```

Edit `config.yaml` with your paths:

```yaml
target_path: /Volumes/1-1          # Source SD card path
backup_path: /Volumes/1-2          # Backup SD card path
destination_path: /Volumes/a7iii   # Final destination for photos
tmp_dir: ~/Pictures/tmp            # Temporary directory for processing
```

Or create a default configuration file:

```bash
rename-sony-photos-directories -create-config
```

## Usage

### Basic Usage - Rename Only

Rename directories in the current directory:

```bash
rename-sony-photos-directories
```

Rename directories in a specific path:

```bash
rename-sony-photos-directories -path /Volumes/1-1/DCIM
```

### Full Workflow - Copy, Rename, Delete

Run the complete workflow (replaces the old shell script):

```bash
rename-sony-photos-directories -workflow
```

This will:
1. Copy photos from source SD card to temporary directory
2. Rename directories to `yyyy-mm-dd` format
3. Copy renamed directories to destination
4. Delete photos from source SD card
5. Eject source SD card

### Backup Cleanup

Delete all photos from backup SD card and eject it:

```bash
rename-sony-photos-directories -backup-cleanup
```

### Dry Run Mode

Preview what would be done without making any changes:

```bash
# Test the full workflow
rename-sony-photos-directories -workflow -dry-run

# Test backup cleanup
rename-sony-photos-directories -backup-cleanup -dry-run

# Test rename only
rename-sony-photos-directories -path /Volumes/1-1/DCIM -dry-run
```

### Using Configuration File

```bash
# Use config.yaml from current directory or ~/.config
rename-sony-photos-directories -workflow

# Specify a custom config file
rename-sony-photos-directories -config /path/to/config.yaml -workflow
```

### Command-Line Flags

- `-workflow` - Run full workflow: copy, rename, and delete
- `-backup-cleanup` - Delete files from backup SD card and eject
- `-path string` - Target path to rename directories (overrides config)
- `-config string` - Path to configuration file
- `-dry-run` - Show what would be done without making changes
- `-create-config` - Create a default configuration file

## Migration from Shell Scripts

If you were using the old shell scripts (`copy-and-delete.sh` and `delete.sh`), you can now use the integrated Go commands:

**Old way:**
```bash
./copy-and-delete.sh
```

**New way:**
```bash
rename-sony-photos-directories -workflow
```

**Old way:**
```bash
./delete.sh
```

**New way:**
```bash
rename-sony-photos-directories -backup-cleanup
```

The Go implementation provides:
- Cross-platform support (not just macOS)
- Dry-run mode for safety
- Better error handling
- Configurable paths via YAML
- No need for separate shell scripts

## Development

### Running Tests

```bash
go test -v -cover
```

### Building

```bash
go build
```

### Building with GoReleaser

```bash
goreleaser release --snapshot --clean
```

## Directory Format

Sony cameras use an 8-digit directory naming format:

- Format: `0YYMMDD0`
- Example: `02512310` represents December 31, 2025
- Positions:
  - `[0]` - Padding (position 0)
  - `[YY]` - Last 2 digits of year (positions 1-2)
  - `[MM]` - Month (positions 3-4)
  - `[DD]` - Day (positions 5-6)
  - `[0]` - Padding (position 7)

This program converts these to `YYYY-MM-DD` format:
- `02512310` → `2025-12-31`
- `02406150` → `2024-06-15`

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.