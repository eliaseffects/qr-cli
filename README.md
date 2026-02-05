# qr-cli

[![CI](https://github.com/eliaseffects/qr-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/eliaseffects/qr-cli/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Generate QR codes from the terminal. Fast, minimal, and easy to script.

## Features
- Single binary, no runtime dependencies
- PNG, SVG, or terminal output
- WiFi and vCard helpers
- Batch generation
- Logo overlays (PNG/SVG)
- Decode QR codes from images
- Clipboard copy + open in viewer
- Config file and `QR_*` env support

## Installation

### Homebrew (macOS/Linux)
```bash
brew install eliaseffects/tap/qr-cli
```

### Go Install
```bash
go install github.com/eliaseffects/qr-cli@latest
```

### Curl Installer
```bash
curl -sSL https://raw.githubusercontent.com/eliaseffects/qr-cli/main/scripts/install.sh | bash
```

### Download Binary
Download from [GitHub Releases](https://github.com/eliaseffects/qr-cli/releases).

## Quick Start
```bash
qr "https://example.com"
qr "Hello world" -o hello.png
qr "https://example.com" --terminal
```

## Usage

### Basic
```bash
# Generate QR code
qr "https://example.com"

# Save to file
qr "Hello world" -o hello.png

# Read from stdin
echo "secret message" | qr -o secret.png

# Render in terminal
qr "https://example.com" --terminal

# Open in viewer
qr "https://example.com" --open

# Colored terminal output
qr "https://example.com" --terminal --terminal-color

# Invert terminal output
qr "https://example.com" --terminal --invert
```

### WiFi Network
```bash
qr wifi --ssid "MyNetwork" --pass "secret123"
```

### Contact Card
```bash
qr vcard --name "John Doe" --phone "+1234567890" --email "john@example.com"
```

### Customization
```bash
# Custom size
qr "https://example.com" -s 512

# Custom colors
qr "https://example.com" --fg "#1a1a1a" --bg "#f0f0f0"

# High error correction
qr "https://example.com" --level H

# Add logo overlay
qr "https://example.com" --logo ./logo.png
```

### Batch Processing
```bash
qr batch -f urls.txt -d ./output/
```

### Decode
```bash
qr decode ./code.png
```

## Commands
- `qr [data]` Generate a QR code from a string or stdin
- `qr wifi` Generate a WiFi QR
- `qr vcard` Generate a vCard QR
- `qr batch` Generate multiple QR codes from a file
- `qr decode` Decode QR codes from an image
- `qr version` Print version info

## Common Flags
- `-o, --output` Output file path (default: `qr.png`/`qr.svg`)
- `-s, --size` Image size in pixels (default: 256)
- `-f, --format` `png`, `svg`, or `terminal` (default: `png`)
- `-l, --level` Error correction `L|M|Q|H` (default: `M`)
- `--fg`, `--bg` Hex colors (default: `#000000`, `#ffffff`)
- `--border` Border size in modules (default: 4)
- `--logo` Logo file path (PNG/JPEG/GIF)
- `--logo-scale` Logo fraction of QR (default: 0.2)
- `-t, --terminal` Render in terminal
- `--terminal-color` Use ANSI colors in terminal output
- `--invert` Invert terminal output
- `--copy` Copy output to clipboard (PNG/SVG only)
- `--open` Open output in default viewer
- `--config` Path to config file

## Config File + Env

Config file is optional. By default `qr-cli` looks for `qr-cli.{yaml|yml|json|toml}` in:
- Current directory
- `~/.config/qr-cli/`
- `~` (home directory)

You can also pass an explicit path with `--config /path/to/config.yaml`.

Example `~/.config/qr-cli/qr-cli.yaml`:
```yaml
size: 512
format: png
fg: "#111111"
bg: "#ffffff"
border: 4
logo: "/absolute/path/to/logo.png"
logo-scale: 0.22

wifi:
  ssid: "MyNetwork"
  security: "WPA"

vcard:
  name: "John Doe"
  email: "john@example.com"
```

Environment variables use the `QR_` prefix (dots and dashes become underscores):
```bash
export QR_SIZE=512
export QR_WIFI_SSID="MyNetwork"
```

## Shell Completions
Completion scripts are in `scripts/completions/`.

## Development
- Requires Go 1.22+
- Run tests: `go test ./...`
- Format: `gofmt -w ./...`

## Release
See `docs/RELEASE.md` for the Homebrew tap setup and release smoke checklist.

## Contributing
Contributions are welcome. See `CONTRIBUTING.md`.

## License
MIT
