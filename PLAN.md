# qr-cli — Development Plan

> A fast, beautiful CLI for generating QR codes from the terminal.
> "The curl of QR codes."

---

## A. Vision

### What It Is
A zero-dependency, single-binary CLI that generates QR codes instantly. No Python, no Node, no runtime—just one executable.

### Why It Matters
- Developers need QR codes constantly (WiFi sharing, URLs, app links, vCards)
- Existing tools require runtime dependencies or are bloated
- A clean, fast CLI that "just works" is missing from the ecosystem

### Design Principles
1. **Zero friction** — `go install` and you're done
2. **Sensible defaults** — works without flags, customizable when needed
3. **Unix philosophy** — does one thing well, composes with pipes
4. **Beautiful output** — terminal rendering that looks good

---

## B. Target User Experience

### Basic Usage
```bash
# Generate QR for any string
qr "https://realityshifting.tech"

# Save to file
qr "Hello world" -o hello.png

# Pipe from stdin
echo "secret message" | qr -o secret.png

# Render in terminal (ASCII/Unicode)
qr "https://example.com" --terminal

# Open immediately (macOS)
qr "https://example.com" --open
```

### Advanced Usage
```bash
# Custom size
qr "https://example.com" -s 1024 -o large.png

# Custom colors
qr "https://example.com" --fg "#1a1a1a" --bg "#ffffff"

# Error correction level
qr "https://example.com" --level H

# SVG output
qr "https://example.com" --format svg -o code.svg

# WiFi network
qr wifi --ssid "HomeNetwork" --pass "secret123" --security WPA

# Contact card
qr vcard --name "Elias Stevenson" --phone "+1234567890" --email "elias@example.com"

# Batch processing
qr batch -f urls.txt -d ./output/

# Logo overlay
qr "https://example.com" --logo logo.png -o branded.png
```

---

## C. Technology Stack

### Language: Go
- Single binary compilation
- Cross-platform (macOS, Linux, Windows)
- No runtime dependencies
- Fast startup time
- Easy distribution via `go install`

### Core Dependencies
| Package | Purpose |
|---------|---------|
| github.com/skip2/go-qrcode | QR code generation (pure Go, no CGO) |
| github.com/spf13/cobra | CLI framework |
| github.com/spf13/viper | Configuration (optional) |

### Optional Dependencies
| Package | Purpose |
|---------|---------|
| github.com/fatih/color | Colored terminal output |
| github.com/mattn/go-isatty | TTY detection |
| golang.design/x/clipboard | Clipboard access |

---

## D. Project Structure

```
qr-cli/
├── .github/
│   ├── workflows/
│   │   ├── ci.yml              # Build + test on PR
│   │   └── release.yml         # Build + release on tag
│   └── FUNDING.yml             # Sponsorship links
├── cmd/
│   ├── root.go                 # Main qr command
│   ├── wifi.go                 # WiFi subcommand
│   ├── vcard.go                # vCard subcommand
│   ├── batch.go                # Batch processing subcommand
│   └── version.go              # Version subcommand
├── internal/
│   ├── qr/
│   │   ├── generate.go         # Core QR generation
│   │   ├── generate_test.go    # Unit tests
│   │   ├── formats.go          # WiFi, vCard string builders
│   │   ├── formats_test.go     # Format tests
│   │   └── options.go          # Generation options struct
│   ├── output/
│   │   ├── file.go             # File output (PNG, SVG)
│   │   ├── terminal.go         # Terminal/ASCII rendering
│   │   ├── clipboard.go        # Clipboard output
│   │   └── viewer.go           # Open in system viewer
│   └── config/
│       └── config.go           # Default settings
├── scripts/
│   ├── install.sh              # curl-pipe installer
│   └── completions/
│       ├── qr.bash             # Bash completions
│       ├── qr.zsh              # Zsh completions
│       └── qr.fish             # Fish completions
├── testdata/
│   ├── urls.txt                # Test batch file
│   └── logo.png                # Test logo for overlay
├── .gitignore
├── .goreleaser.yml             # GoReleaser config
├── go.mod
├── go.sum
├── main.go                     # Entry point
├── LICENSE                     # MIT
├── README.md                   # User documentation
├── CHANGELOG.md                # Release notes
└── PLAN.md                     # This file
```

---

## E. Command & Flag Design

### Root Command: `qr`

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| --output | -o | string | stdout/qr.png | Output file path |
| --size | -s | int | 256 | Image size in pixels |
| --format | -f | string | png | Output format: png, svg, terminal |
| --level | -l | string | M | Error correction: L, M, Q, H |
| --fg | | string | #000000 | Foreground color (hex) |
| --bg | | string | #ffffff | Background color (hex) |
| --border | | int | 4 | Border size in modules |
| --terminal | -t | bool | false | Render in terminal |
| --open | | bool | false | Open in system viewer |
| --copy | | bool | false | Copy to clipboard |
| --quiet | -q | bool | false | Suppress non-error output |
| --version | -v | bool | | Print version |
| --help | -h | bool | | Print help |

### Subcommand: `qr wifi`

Generate QR code for WiFi network connection.

| Flag | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| --ssid | string | yes | | Network name |
| --pass | string | no | | Network password |
| --security | string | no | WPA | Security type: WPA, WEP, nopass |
| --hidden | bool | no | false | Network is hidden |
| --output | string | no | | Output file path |
| --terminal | bool | no | false | Render in terminal |

### Subcommand: `qr vcard`

Generate QR code for contact card.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| --name | string | yes | Full name |
| --phone | string | no | Phone number |
| --email | string | no | Email address |
| --org | string | no | Organization |
| --title | string | no | Job title |
| --url | string | no | Website URL |
| --address | string | no | Street address |
| --output | string | no | Output file path |
| --terminal | bool | no | Render in terminal |

### Subcommand: `qr batch`

Generate QR codes from a file of data (one per line).

| Flag | Short | Type | Required | Default | Description |
|------|-------|------|----------|---------|-------------|
| --file | -f | string | yes | | Input file path |
| --dir | -d | string | no | ./qr-output | Output directory |
| --size | -s | int | no | 256 | Image size |
| --format | | string | no | png | Output format |
| --prefix | | string | no | qr- | Filename prefix |

---

## F. Feature Breakdown

### Phase 1: MVP
| Feature | Description |
|---------|-------------|
| Basic generation | QR from string argument |
| File output | Save to PNG |
| Custom output path | -o flag |
| Custom size | -s flag |
| Stdin support | Pipe data in |
| Help/version | Standard CLI help |

### Phase 2: Polish
| Feature | Description |
|---------|-------------|
| Terminal render | ASCII/Unicode output with --terminal |
| Open in viewer | --open flag (macOS/Linux/Windows) |
| Error correction | --level L/M/Q/H |
| Custom colors | --fg and --bg hex colors |
| SVG output | --format svg |
| WiFi subcommand | qr wifi --ssid --pass |
| vCard subcommand | qr vcard --name --phone --email |

### Phase 3: Advanced
| Feature | Description |
|---------|-------------|
| Batch processing | qr batch -f urls.txt |
| Clipboard copy | --copy flag |
| Logo overlay | --logo flag |
| Shell completions | bash, zsh, fish |
| Custom border | --border flag |

---

## G. Implementation Details

### Core Logic: QR Generation
- Use go-qrcode library (pure Go, no CGO required)
- Support all 4 error correction levels (L=7%, M=15%, Q=25%, H=30%)
- Default to Medium (M) for balance of size and reliability
- Parse hex colors for foreground/background customization

### WiFi Format
Standard WiFi QR format: `WIFI:T:<security>;S:<ssid>;P:<password>;H:<hidden>;;`
- Escape special characters: `\`, `;`, `,`, `"`
- Security types: WPA, WEP, nopass

### vCard Format
Standard vCard 3.0 format with fields:
- FN (full name), N (structured name)
- TEL, EMAIL, ORG, TITLE, URL, ADR
- Properly escape special characters

### Terminal Rendering
- Use Unicode block characters: ▀ (upper half), ▄ (lower half), █ (full), space
- Combine 2 vertical pixels into 1 character for better aspect ratio
- Support inverted colors for dark terminals

### System Viewer Integration
- macOS: `open <path>`
- Linux: `xdg-open <path>`
- Windows: `cmd /c start <path>`

---

## H. Testing Strategy

### Unit Tests
| Component | Test Coverage |
|-----------|---------------|
| QR generation | Valid input, empty input, unicode, long text |
| WiFi format | Standard, special chars, hidden network |
| vCard format | All fields, partial fields, escaping |
| Color parsing | Valid hex, invalid hex, defaults |
| Terminal render | Output structure, character mapping |

### Integration Tests
| Scenario | Description |
|----------|-------------|
| CLI end-to-end | Generate file, verify PNG magic bytes |
| Stdin pipe | Echo data, verify output |
| WiFi command | Full flow with all flags |
| Batch processing | Multiple inputs, verify all outputs |

### Manual QA Checklist
- [ ] Scan generated QR with phone camera
- [ ] Test WiFi QR connects to network
- [ ] Test vCard QR adds contact
- [ ] Verify terminal output renders correctly in iTerm, Terminal.app, Alacritty
- [ ] Test --open on macOS, Linux, Windows

---

## I. Build & Release

### Build Configuration
- Use GoReleaser for cross-platform builds
- CGO_ENABLED=0 for static binaries
- Build targets: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64
- Inject version, commit, date via ldflags

### CI/CD Pipeline

**On Pull Request:**
- Run tests
- Run linter (golangci-lint)
- Build all platforms
- Upload coverage

**On Tag (v*):**
- Run GoReleaser
- Create GitHub release with binaries
- Update Homebrew formula
- Build deb/rpm packages

### Versioning
- Semantic versioning (MAJOR.MINOR.PATCH)
- v0.1.0 = MVP
- v0.2.0 = WiFi + vCard subcommands
- v1.0.0 = Stable with all Phase 2 features

---

## J. Documentation

### README Sections
1. Project description (one-liner)
2. Installation (Homebrew, go install, binary download)
3. Usage examples (basic, wifi, vcard)
4. All flags and options
5. License

### Man Page
- Generate from Cobra docs
- Include in Homebrew formula

---

## K. Distribution Channels

| Channel | Command |
|---------|---------|
| go install | `go install github.com/eliaseffects/qr-cli@latest` |
| Homebrew | `brew install eliaseffects/tap/qr-cli` |
| GitHub Releases | Pre-built binaries for all platforms |
| curl installer | `curl -sSL https://.../install.sh \| bash` |
| npm wrapper | `npm install -g @eliaseffects/qr-cli` (optional) |

---

## L. Implementation Checklist

### Setup
- [ ] Create GitHub repository `eliaseffects/qr-cli`
- [ ] Initialize Go module (`go mod init`)
- [ ] Add .gitignore
- [ ] Add MIT LICENSE
- [ ] Create project structure (cmd/, internal/, scripts/)

### Core Implementation
- [ ] Implement `internal/qr/options.go` — Options struct
- [ ] Implement `internal/qr/generate.go` — Core generation
- [ ] Implement `internal/output/file.go` — PNG/SVG output
- [ ] Implement `internal/output/terminal.go` — ASCII render
- [ ] Implement `internal/output/viewer.go` — System open
- [ ] Implement `cmd/root.go` — Main command with all flags
- [ ] Implement `main.go` — Entry point

### Subcommands
- [ ] Implement `internal/qr/formats.go` — WiFi and vCard formatters
- [ ] Implement `cmd/wifi.go` — WiFi subcommand
- [ ] Implement `cmd/vcard.go` — vCard subcommand
- [ ] Implement `cmd/batch.go` — Batch processing
- [ ] Implement `cmd/version.go` — Version info

### Testing
- [ ] Write unit tests for QR generation
- [ ] Write unit tests for formats (WiFi, vCard)
- [ ] Write unit tests for terminal rendering
- [ ] Write integration tests for CLI
- [ ] Achieve >80% code coverage

### Documentation
- [ ] Write comprehensive README.md
- [ ] Add usage examples for all commands
- [ ] Document all flags and options
- [ ] Create CHANGELOG.md

### CI/CD
- [ ] Set up GitHub Actions CI (.github/workflows/ci.yml)
- [ ] Set up release workflow (.github/workflows/release.yml)
- [ ] Configure GoReleaser (.goreleaser.yml)
- [ ] Test release with v0.0.1-rc1 tag

### Distribution
- [ ] Create `homebrew-tap` repository
- [ ] Push v0.1.0 release
- [ ] Verify Homebrew installation
- [ ] Create curl installer script

### Polish
- [ ] Add shell completions (bash, zsh, fish)
- [ ] Add clipboard support (--copy)
- [ ] Add SVG output format
- [ ] Add logo overlay feature

### Launch
- [ ] Tweet announcement ("just shipped: qr-cli")
- [ ] Submit to Hacker News
- [ ] Post on r/golang, r/commandline
- [ ] Add to awesome-go list

---

## M. Post-Launch Improvements

### Additional Formats
- SMS QR (`sms:+1234567890?body=Hello`)
- Email QR (`mailto:email@example.com`)
- Phone QR (`tel:+1234567890`)
- Event QR (iCalendar format)
- Bitcoin/Ethereum address QR

### Advanced Features
- QR code decoder (`qr decode image.png`)
- Web UI mode (`qr serve --port 8080`)
- Custom shapes (rounded corners)
- Animated QR (GIF with color transition)

---

## N. Success Metrics

| Metric | Target (30 days) |
|--------|------------------|
| GitHub stars | 100+ |
| Homebrew installs | 50+ |
| Tweet impressions | 10K+ |
| Hacker News | Front page |
| Contributors | 3+ |

---

*Status: Ready for Implementation*
