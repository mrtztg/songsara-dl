# SongSara Downloader

A fast CLI tool to download entire albums and playlists from SongSara with concurrent downloads.

## Features

- ⚡ Concurrent downloads (configurable, max 10 by default)
- 📁 Organized folder structure
- 🔄 Skip existing files (resume downloads)
- 📊 Progress tracking
- 🛡️ Anti-bot protection bypass
- 🧹 Filename sanitization

## Quick Start

### Download Pre-built Binaries

Download the latest release for your platform:

**macOS:**
- Intel: [songsara-dl_Darwin_x86_64](https://github.com/mrtztg/songsara-dl/releases/latest/download/songsara-dl_Darwin_x86_64)
- Apple Silicon: [songsara-dl_Darwin_arm64](https://github.com/mrtztg/songsara-dl/releases/latest/download/songsara-dl_Darwin_arm64)

**Linux:**
- x86_64: [songsara-dl_Linux_x86_64](https://github.com/mrtztg/songsara-dl/releases/latest/download/songsara-dl_Linux_x86_64)
- ARM64: [songsara-dl_Linux_arm64](https://github.com/mrtztg/songsara-dl/releases/latest/download/songsara-dl_Linux_arm64)

**Windows:**
- x86_64: [songsara-dl_Windows_x86_64](https://github.com/mrtztg/songsara-dl/releases/latest/download/songsara-dl_Windows_x86_64.exe)

### Build from Source

```bash
git clone <repository-url>
cd songsara-dl
go mod tidy
go build -o songsara-dl
```

### Building Releases

To build executables for all platforms:

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Build snapshot (local testing)
make snapshot

# Build release (requires git tag)
make release
```

📖 **For GitHub Releases**: See [RELEASE_GUIDE.md](RELEASE_GUIDE.md) for step-by-step instructions.

## Usage

```bash
# Download a single album
./songsara-dl "https://songsara.net/59021/"

# Download multiple albums
./songsara-dl "https://songsara.net/59021/" "https://songsara.net/12345/"

# With options
./songsara-dl -c 5 -o /music -v "https://songsara.net/59021/"
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `-c, --concurrency` | Max concurrent downloads | `10` |
| `-o, --output` | Output directory | `downloads` |
| `-v, --verbose` | Verbose output | `false` |
| `-n, --dry-run` | Preview downloads | `false` |
| `-s, --skip-existing` | Skip existing files | `true` |
| `-t, --timeout` | HTTP timeout (seconds) | `30` |

## Examples

```bash
# Basic download
./songsara-dl "https://songsara.net/59021/"

# Multiple albums with custom settings
./songsara-dl -c 8 -o /music -v "https://songsara.net/59021/" "https://songsara.net/12345/"

# Dry run to preview
./songsara-dl -n -v "https://songsara.net/59021/"

# Custom timeout
./songsara-dl -t 60 "https://songsara.net/59021/"
```

## Output Structure

```
downloads/
├── Album Name 1/
│   ├── 01 - Song Title 1.mp3
│   ├── 02 - Song Title 2.mp3
│   └── ...
└── Album Name 2/
    ├── 01 - Song Title 1.mp3
    └── ...
```

## Supported Platforms

- macOS (Intel & Apple Silicon)
- Linux (x86_64 & ARM64)
- Windows (x86_64)

## Build Status

✅ Successfully tested builds for all platforms:
- macOS Intel (8.3MB)
- macOS Apple Silicon (7.8MB)
- Linux x86_64 (8.0MB)
- Linux ARM64 (7.6MB)
- Windows x86_64 (8.4MB)

## License

This project is for educational purposes. Please respect SongSara's terms of service. 