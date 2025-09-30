# Configuration Guide

## Configuration File

The application uses YAML format for configuration.

### Default Locations

The program searches for `config.yaml` in the following order:

1. `./config.yaml` - Current directory
2. `~/.config/rename-sony-photos/config.yaml` - User config directory

### Configuration Structure

```yaml
target_path: /Volumes/1-1          # Source SD card path
backup_path: /Volumes/1-2          # Backup SD card path
destination_path: /Volumes/a7iii   # Final destination for photos
tmp_dir: ~/Pictures/tmp            # Temporary directory for processing
```

### Configuration Options

#### `target_path`
- **Type**: String
- **Required**: Yes (for workflow mode)
- **Description**: Path to the source SD card containing photos
- **Default**: `/Volumes/1-1`
- **Example**: `/Volumes/CAMERA-SD`

#### `backup_path`
- **Type**: String
- **Required**: Yes (for backup cleanup)
- **Description**: Path to the backup SD card
- **Default**: `/Volumes/1-2`
- **Example**: `/Volumes/BACKUP-SD`

#### `destination_path`
- **Type**: String
- **Required**: Yes (for workflow mode)
- **Description**: Final destination directory for processed photos
- **Default**: `/Volumes/a7iii`
- **Example**: `/Users/username/Photos/Camera`

#### `tmp_dir`
- **Type**: String
- **Required**: Yes (for workflow mode)
- **Description**: Temporary directory for processing photos
- **Default**: `~/Pictures/tmp`
- **Example**: `/tmp/photo-processing`

## Creating Configuration

### Method 1: Auto-generate

```bash
rename-sony-photos-directories -create-config
```

This creates `config.yaml` in the current directory with default values.

### Method 2: Manual Creation

Create `config.yaml`:

```yaml
target_path: /Volumes/MY-CAMERA
backup_path: /Volumes/MY-BACKUP
destination_path: /Users/myname/Photos/Imported
tmp_dir: /Users/myname/tmp
```

### Method 3: Copy Example

```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your values
```

## Configuration Examples

### Example 1: Basic Setup

```yaml
target_path: /Volumes/1-1
backup_path: /Volumes/1-2
destination_path: /Users/photographer/Photos/Camera
tmp_dir: /Users/photographer/tmp
```

### Example 2: Network Storage

```yaml
target_path: /Volumes/SD-CARD
backup_path: /Volumes/BACKUP
destination_path: /Volumes/NAS/Photos/Camera
tmp_dir: ~/Pictures/temp
```

### Example 3: Windows Paths

```yaml
target_path: D:\DCIM
backup_path: E:\DCIM
destination_path: C:\Users\Username\Photos\Camera
tmp_dir: C:\Temp\PhotoProcessing
```

### Example 4: Linux Paths

```yaml
target_path: /media/username/SD-CARD
backup_path: /media/username/BACKUP
destination_path: /home/username/Photos/Camera
tmp_dir: /tmp/photo-processing
```

## Multiple Configurations

You can maintain multiple configuration files for different cameras or workflows:

```bash
# For Camera A
rename-sony-photos-directories -config camera-a.yaml -workflow

# For Camera B
rename-sony-photos-directories -config camera-b.yaml -workflow
```

## Environment-Specific Configuration

### Development

```yaml
target_path: ./test-data/source
backup_path: ./test-data/backup
destination_path: ./test-data/dest
tmp_dir: ./test-data/tmp
```

### Production

```yaml
target_path: /Volumes/CAMERA-SD
backup_path: /Volumes/BACKUP-SD
destination_path: /Volumes/Photos/Camera
tmp_dir: ~/Pictures/tmp
```

## Troubleshooting

### Config File Not Found

If you see "Failed to load config", check:
1. File exists at specified path
2. File has correct YAML syntax
3. File has read permissions

### Invalid Paths

If paths don't exist, the program will fail with an error. Ensure:
- SD cards are mounted
- Directories exist and are accessible
- You have write permissions

### Permissions

The config file should be readable:
```bash
chmod 600 config.yaml
```
