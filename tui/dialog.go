package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type ConfirmationDialog struct {
	message string
	width   int
	height  int
	visible bool
}

func NewConfirmationDialog(message string) *ConfirmationDialog {
	return &ConfirmationDialog{
		message: message,
		visible: true,
	}
}

func (c *ConfirmationDialog) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c *ConfirmationDialog) Show() {
	c.visible = true
}

func (c *ConfirmationDialog) Hide() {
	c.visible = false
}

func (c *ConfirmationDialog) IsVisible() bool {
	return c.visible
}

func (c *ConfirmationDialog) Render() string {
	if !c.visible {
		return ""
	}

	// Dialog content styling with enhanced appearance
	dialogStyle := AppStyles.Dialog.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(Colors.Error)).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center)

	content := fmt.Sprintf("%s\n\nPress 'y' to confirm, 'n' to cancel", c.message)
	dialog := dialogStyle.Render(content)

	// Create overlay if we have dimensions
	if c.width > 0 && c.height > 0 {
		// Center the dialog with overlay background
		return lipgloss.Place(
			c.width, c.height,
			lipgloss.Center, lipgloss.Center,
			dialog,
			lipgloss.WithWhitespaceForeground(lipgloss.Color("238")),
		)
	}

	return dialog
}
