# Usage Guide

## Basic Usage

### Rename Directories Only

Rename directories in the current directory:

```bash
rename-sony-photos-directories
```

Rename directories in a specific path:

```bash
rename-sony-photos-directories -path /Volumes/1-1/DCIM
```

### Full Workflow

Run the complete workflow (copy, rename, delete, eject):

```bash
rename-sony-photos-directories -workflow
```

This will:
1. Copy photos from source SD card to temporary directory
2. Rename directories to `yyyy-mm-dd` format
3. Copy renamed directories to destination
4. Delete photos from source SD card
5. Eject source SD card (macOS only)

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

## Configuration

### Using Configuration File

The program looks for `config.yaml` in:
1. Current directory (`./config.yaml`)
2. `~/.config/rename-sony-photos/config.yaml`

Specify a custom config file:

```bash
rename-sony-photos-directories -config /path/to/config.yaml -workflow
```

### Creating Default Configuration

Create a default configuration file:

```bash
rename-sony-photos-directories -create-config
```

This creates `config.yaml` in the current directory with default values.

## Command-Line Flags

- `-workflow` - Run full workflow: copy, rename, and delete
- `-backup-cleanup` - Delete files from backup SD card and eject
- `-path string` - Target path to rename directories (overrides config)
- `-config string` - Path to configuration file
- `-dry-run` - Show what would be done without making changes
- `-create-config` - Create a default configuration file
- `-help` - Show help message

## Examples

### Example 1: Basic Renaming

```bash
cd /Volumes/1-1/DCIM
rename-sony-photos-directories
```

### Example 2: Full Workflow with Custom Config

```bash
rename-sony-photos-directories -config ~/my-camera-config.yaml -workflow
```

### Example 3: Dry Run Before Actual Execution

```bash
# First, check what will be done
rename-sony-photos-directories -workflow -dry-run

# If everything looks good, run for real
rename-sony-photos-directories -workflow
```

### Example 4: Backup Workflow

```bash
# Process main SD card
rename-sony-photos-directories -workflow

# Clean up backup SD card
rename-sony-photos-directories -backup-cleanup
```

## Common Workflows

### Photography Workflow

1. Insert main SD card into card reader
2. Run: `rename-sony-photos-directories -workflow`
3. Insert backup SD card
4. Run: `rename-sony-photos-directories -backup-cleanup`

### Safe Testing

1. Create test config: `rename-sony-photos-directories -create-config`
2. Edit config to point to test directories
3. Test with dry-run: `rename-sony-photos-directories -workflow -dry-run`
4. Execute: `rename-sony-photos-directories -workflow`
