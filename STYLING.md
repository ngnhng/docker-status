# Styling System Documentation

## Overview

The Docker TUI application now uses a centralized styling system that makes it easy to:

- Maintain consistent colors and styles across the application
- Switch between different themes (Default, Dark, Light)
- Add new themes easily
- Modify colors globally from a single location

## Architecture

### Core Files

- `tui/styles.go` - Contains all color palettes, style definitions, and helper functions
- Individual TUI files use the centralized styles instead of hardcoded values

### Key Components

1. **ColorPalette** - Defines a set of colors for a theme
2. **Styles** - Contains all styled components (titles, tabs, tables, etc.)
3. **Helper Functions** - Convenience functions for common styling operations

## Usage Examples

### Using Predefined Styles

```go
// Instead of:
title := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("212")).
    Render("Docker Status")

// Use:
title := StyleTitle("Docker Status")
```

### Using Style Instances

```go
// Access styles through the model
errorText := m.styles.TextError.Render("Error message")
```

### Styling with Helpers

```go
// Container status with automatic color coding
status := StyleContainerStatus(container.Status, container.State)

// Tab styling with active/inactive states
tab := StyleTab("Containers", isActive)
```

## Theme Switching

### Available Themes

1. **Default Theme** - Blue/pink color scheme (original)
2. **Dark Theme** - Darker, more muted colors
3. **Light Theme** - Light background with dark text

### Switching Themes at Runtime

```go
// Switch to dark theme
m.SetDarkTheme()

// Switch to light theme
m.SetLightTheme()

// Switch to default theme
m.SetDefaultTheme()
```

### Key Bindings for Theme Switching

- `t+1` - Switch to default theme
- `t+2` - Switch to dark theme
- `t+3` - Switch to light theme

## Creating Custom Themes

### Define a New Color Palette

```go
var MyCustomColors = ColorPalette{
    // Primary colors
    Primary:    "27",  // Your preferred accent color
    Secondary:  "165", // Secondary accent
    Accent:     "33",  // Highlight color
    Background: "232", // Background color
    Surface:    "236", // Surface elements

    // Text colors
    TextPrimary:   "252", // Main text
    TextSecondary: "220", // Secondary text
    TextMuted:     "244", // Muted text
    TextHighlight: "165", // Highlighted text

    // Status colors
    Success: "40",  // Success/running states
    Warning: "214", // Warning states
    Error:   "160", // Error/stopped states
    Info:    "33",  // Info messages

    // Border colors
    BorderNormal:   "236", // Default borders
    BorderActive:   "25",  // Active element borders
    BorderInactive: "244", // Inactive borders
}
```

### Apply Custom Theme

```go
// Apply your custom theme
SetTheme(MyCustomColors)

// Or through the model
m.SetTheme(MyCustomColors)
```

## Style Components

### Typography

- `Title` - Main headings
- `Subtitle` - Section headings
- `Text` - Regular text
- `TextMuted` - Secondary text
- `TextSuccess/Warning/Error` - Status text

### Navigation

- `TabActive` - Active navigation tabs
- `TabInactive` - Inactive navigation tabs

### Tables

- `TableHeader` - Table column headers
- `TableRow` - Table row content
- `TableBorder` - Table borders

### Dialogs

- `Dialog` - Dialog content styling
- `DialogBorder` - Dialog border styling

### Status

- `StatusBar` - Status bar text
- `StatusError` - Error messages
- `StatusInfo` - Info messages

## Helper Functions

### Quick Styling

- `StyleTitle(text)` - Apply title styling
- `StyleSubtitle(text)` - Apply subtitle styling
- `StyleError(text)` - Apply error styling
- `StyleSuccess(text)` - Apply success styling
- `StyleMuted(text)` - Apply muted styling

### Conditional Styling

- `StyleContainerStatus(status, state)` - Auto-style based on container state
- `StyleTab(text, isActive)` - Style tabs based on active state

### Dynamic Sizing

- `StyleTitleWithWidth(text, width)` - Title with specific width
- `StyleHelpTitleWithWidth(text, width)` - Help title with width

## Benefits

1. **Consistency** - All components use the same color palette
2. **Maintainability** - Change colors in one place
3. **Flexibility** - Easy theme switching
4. **Extensibility** - Simple to add new themes or modify existing ones
5. **Performance** - Styles are created once and reused

## Migration Notes

The refactoring replaces hardcoded `lipgloss.Color()` calls with centralized color constants and helper functions. This makes the codebase more maintainable while preserving all existing functionality.

### Before

```go
lipgloss.NewStyle().
    Foreground(lipgloss.Color("212")).
    Bold(true).
    Render("Title")
```

### After

```go
StyleTitle("Title")
```

This approach makes the code cleaner, more consistent, and easier to maintain.

```go
package main

import (
	"context"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/client"
	"github.com/ngnhng/docker-status/tui"
)

// Example of creating a custom theme
var GreenTheme = tui.ColorPalette{
	// Primary colors - Green theme
	Primary:    "28",  // Dark green
	Secondary:  "46",  // Bright green
	Accent:     "34",  // Forest green
	Background: "232", // Very dark background
	Surface:    "235", // Dark surface

	// Text colors
	TextPrimary:   "252", // Light gray
	TextSecondary: "154", // Light green
	TextMuted:     "242", // Medium gray
	TextHighlight: "46",  // Bright green

	// Status colors
	Success: "46",  // Bright green
	Warning: "214", // Orange
	Error:   "196", // Red
	Info:    "39",  // Blue

	// Border colors
	BorderNormal:   "235", // Dark gray
	BorderActive:   "28",  // Dark green
	BorderInactive: "240", // Gray
}

func main() {
	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	// Test connection
	ctx := context.Background()
	_, err = cli.Info(ctx)
	if err != nil {
		log.Fatal("Docker daemon not available:", err)
	}

	// Create TUI model
	model := tui.NewModel(cli)

	// Apply custom green theme
	model.SetTheme(GreenTheme)

	// You could also do:
	// model.SetDarkTheme()    // For built-in dark theme
	// model.SetLightTheme()   // For built-in light theme
	// model.SetDefaultTheme() // For original theme

	// Run the TUI
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}

```