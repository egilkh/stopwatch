# Stopwatch

A feature-rich terminal stopwatch application written in Go with a beautiful TUI (Terminal User Interface).

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

## Features

- ⏱️ **Precise timing** with millisecond accuracy
- 🏃 **Lap tracking** with individual and total times
- 📊 **Scrollable lap list** for unlimited laps
- 💾 **Export to JSON** with detailed timing data
- 🎨 **Colorful interface** with intuitive status indicators
- ⚡ **Responsive updates** at 40 FPS
- 🖥️ **Cross-platform** support (macOS, Linux, Windows)

## Installation

### From Source

Requires Go 1.24 or later:

```bash
git clone https://github.com/egilkh/stopwatch.git
cd stopwatch
go build -o sw .
```

### Using Go Install

```bash
go install github.com/egilkh/stopwatch@latest
```

## Usage

Start the stopwatch:
```bash
./sw
```

The stopwatch starts running immediately when launched.

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `S` | Start/Stop the timer |
| `L` | Record a lap (while running) |
| `R` | Reset the timer |
| `F` | Export laps to JSON file |
| `Q` | Quit the application |
| `↑/↓` | Scroll through laps (when list is long) |

### Display Layout

```
Stopwatch v1.0.0 - Started: 2024-12-25T10:30:00Z
[S]tart/Stop [L]ap [R]eset [F]ile/Export [Q]uit
00:15.234 [RUNNING]

Lap 01: 00:05.123 (Total: 00:05.123)
Lap 02: 00:04.567 (Total: 00:09.690)
Lap 03: 00:05.544 (Total: 00:15.234)
```

## Export Format

Pressing `F` exports lap data to a JSON file with timestamp:

```json
{
  "export_time": "2024-12-25T10:35:30Z",
  "start_time": "2024-12-25T10:30:00Z",
  "total_time": "05:30.234",
  "total_time_ms": 330234,
  "lap_count": 10,
  "laps": [
    {
      "number": 1,
      "lap_time": "00:32.123",
      "lap_time_ms": 32123,
      "total_time": "00:32.123",
      "total_time_ms": 32123
    },
    ...
  ]
}
```

## Building

### Prerequisites

- Go 1.24 or later
- Terminal with ANSI escape sequence support

### Build Commands

```bash
# Build for current platform
go build -o sw .

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o sw-linux .
GOOS=darwin GOARCH=amd64 go build -o sw-mac .
GOOS=windows GOARCH=amd64 go build -o sw.exe .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run with verbose output
go test ./... -v
```

## Architecture

The application uses a component-based architecture:

- **Timer Component**: Core timing logic and lap management
- **Lap List Component**: Scrollable viewport for lap display
- **Header Component**: Version and start time display
- **Status Bar Component**: Help text and temporary messages

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`go test ./...`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with the excellent [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework by Charm
- Inspired by the need for a simple, beautiful terminal stopwatch