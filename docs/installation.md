# Installation Guide

## Requirements

- Go 1.21 or later
- macOS, Linux, or Windows

## Installation Methods

### Using Go Install (Recommended)

```bash
go install github.com/shunichi-ikebuchi/rename-sony-photos-directories/cmd/rename-sony-photos-directories@latest
```

This will install the binary to `$GOPATH/bin` (usually `~/go/bin`).

### From Source

```bash
git clone https://github.com/shunichi-ikebuchi/rename-sony-photos-directories.git
cd rename-sony-photos-directories
make install
```

### Using Pre-built Binaries

Download the latest release from the [releases page](https://github.com/shunichi-ikebuchi/rename-sony-photos-directories/releases).

#### macOS

```bash
# Download and extract
curl -L https://github.com/shunichi-ikebuchi/rename-sony-photos-directories/releases/latest/download/rename-sony-photos-directories_Darwin_x86_64.tar.gz | tar xz

# Move to PATH
sudo mv rename-sony-photos-directories /usr/local/bin/
```

#### Linux

```bash
# Download and extract
curl -L https://github.com/shunichi-ikebuchi/rename-sony-photos-directories/releases/latest/download/rename-sony-photos-directories_Linux_x86_64.tar.gz | tar xz

# Move to PATH
sudo mv rename-sony-photos-directories /usr/local/bin/
```

#### Windows

Download the `.zip` file from the releases page and extract it to a directory in your PATH.

## Verification

Verify the installation:

```bash
rename-sony-photos-directories -help
```

## Updating

### Using Go Install

```bash
go install github.com/shunichi-ikebuchi/rename-sony-photos-directories/cmd/rename-sony-photos-directories@latest
```

### From Source

```bash
cd rename-sony-photos-directories
git pull
make install
```
