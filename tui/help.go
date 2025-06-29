package tui

import (
	"strings"
)

// HelpView represents the help screen
type HelpView struct {
	width  int
	height int
}

// NewHelpView creates a new help view
func NewHelpView() *HelpView {
	return &HelpView{}
}

// SetSize sets the help view dimensions
func (h *HelpView) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// Render renders the help screen
func (h *HelpView) Render() string {
	var content strings.Builder

	title := StyleHelpTitleWithWidth("Docker TUI - Help", h.width)

	content.WriteString(title)
	content.WriteString("\n\n")

	// General navigation
	content.WriteString(h.renderSection("General Navigation", []string{
		"↑/k, ↓/j         Navigate up/down in tables",
		"←/h, →/l         Navigate left/right (future use)",
		"Tab              Switch between views",
		"1-4              Jump directly to view (1=Containers, 2=Images, 3=Networks, 4=Volumes)",
		"r, Ctrl+R        Refresh data",
		"q, Ctrl+C        Quit application",
		"?                Show/hide this help",
	}))

	// Container specific
	content.WriteString(h.renderSection("Container Management", []string{
		"d                Delete selected container (with confirmation)",
		"g                Toggle grouping by Docker Compose project",
		"L                View logs for selected container (coming soon)",
		"Enter            Inspect selected container (coming soon)",
		"s                Start/stop selected container (coming soon)",
	}))

	// Image specific
	content.WriteString(h.renderSection("Image Management", []string{
		"d                Delete selected image (with confirmation)",
		"Enter            Inspect selected image (coming soon)",
		"p                Pull new image (coming soon)",
	}))

	// Network specific
	content.WriteString(h.renderSection("Network Management", []string{
		"d                Delete selected network (with confirmation)",
		"Enter            Inspect selected network (coming soon)",
		"n                Create new network (coming soon)",
	}))

	// Volume specific
	content.WriteString(h.renderSection("Volume Management", []string{
		"d                Delete selected volume (with confirmation)",
		"Enter            Inspect selected volume (coming soon)",
		"v                Create new volume (coming soon)",
	}))

	// Features
	content.WriteString(h.renderSection("Features", []string{
		"• Real-time updates every 5 seconds",
		"• Color-coded status indicators",
		"• Keyboard-driven navigation (like k9s)",
		"• Resource management (containers, images, networks, volumes)",
		"• Safety confirmations for destructive operations",
	}))

	footer := AppStyles.HelpFooter.
		Width(h.width).
		Render("Press ? to close help")

	content.WriteString("\n")
	content.WriteString(footer)

	return content.String()
}

// renderSection renders a help section with title and items
func (h *HelpView) renderSection(title string, items []string) string {
	var section strings.Builder

	sectionTitle := StyleSubtitle(title)

	section.WriteString(sectionTitle)
	section.WriteString("\n")

	for _, item := range items {
		section.WriteString("  ")
		section.WriteString(item)
		section.WriteString("\n")
	}
	section.WriteString("\n")

	return section.String()
}
