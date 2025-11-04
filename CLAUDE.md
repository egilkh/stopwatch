# CLAUDE.md - AI Assistant Instructions

## Project Overview
This is a terminal-based stopwatch application written in Go using the Bubble Tea TUI framework. The application follows a component-based architecture for maintainability and testability.

## Architecture

### Components
- **timer** - Core timing logic, lap management, and time formatting
- **laps** - Scrollable lap list with viewport management
- **header** - Application header with version and start time
- **statusbar** - Help text, temporary messages, and scroll indicators

### Key Technologies
- **Bubble Tea** - Terminal UI framework based on Elm architecture
- **Lipgloss** - Styling library for terminal output
- **Bubbles** - Component library (viewport for scrolling)

## Development Guidelines

### Code Style
- Follow Go idioms and conventions
- Use `go fmt` for formatting
- Ensure all files end with a newline
- Keep components focused on single responsibilities

### Testing
- Run `go test ./...` before committing
- Maintain high test coverage (currently ~88%)
- Test files should be next to implementation files

### Building
```bash
go build -o sw .
```

### Running Tests
```bash
go test ./... -v        # Verbose output
go test ./... -cover    # With coverage
```

## Important Implementation Details

### Timer Update Rate
- Currently set to 25ms for smooth display updates
- Status bar messages use tick counts (80 ticks = 2 seconds)

### Export Functionality
- Exports to JSON with timestamp in filename
- Format: `stopwatch_YYYYMMDD_HHMMSS.json`
- Includes lap times and total elapsed time

### Component Communication
- Parent model routes messages to child components
- Components expose public methods for state updates
- Use `tea.Batch()` to combine multiple commands

## Common Tasks

### Adding a New Feature
1. Identify which component(s) need modification
2. Update the component's model and methods
3. Update the main model's Update() method if needed
4. Add tests for new functionality
5. Update README.md if user-facing

### Modifying Keyboard Shortcuts
1. Update the key handling in main.go Update()
2. Update help text in statusbar component
3. Update README.md documentation

### Changing Display Format
1. Modify the relevant component's View() method
2. Update tests to match new output
3. Consider impact on terminal width requirements

## File Structure
```
├── main.go                 # Main application and composition
├── components/
│   ├── timer/             # Timer logic
│   ├── laps/              # Lap list display
│   ├── header/            # Header display
│   └── statusbar/         # Status and help
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── .gitignore            # Git ignore rules
├── README.md             # User documentation
└── CLAUDE.md             # This file
```

## Known Constraints
- Requires terminal with ANSI escape sequence support
- Uses raw terminal mode (restores on exit)
- JSON export writes to current directory

## Future Considerations
- Could add configuration file support
- Could add different export formats (CSV, etc.)
- Could add lap statistics (average, best, worst)
- Could add split times between laps