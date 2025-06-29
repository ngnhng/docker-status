package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type ColorPalette struct {
	// Primary colors
	Primary    string
	Secondary  string
	Accent     string
	Background string
	Surface    string

	// Text colors
	TextPrimary   string
	TextSecondary string
	TextMuted     string
	TextHighlight string

	// Status colors
	Success string
	Warning string
	Error   string
	Info    string

	// Border colors
	BorderNormal   string
	BorderActive   string
	BorderInactive string
}

var DefaultColors = ColorPalette{
	// Primary colors
	Primary:    "69",
	Secondary:  "81",
	Accent:     "39",
	Background: "232",
	Surface:    "240",

	// Text colors
	TextPrimary:   "255", // White
	TextSecondary: "229", // Light yellow
	TextMuted:     "241", // Gray
	TextHighlight: "81",

	// Status colors
	Success: "46",  // Green
	Warning: "226", // Yellow
	Error:   "124", // Red
	Info:    "81",  // Light blue

	// Border colors
	BorderNormal:   "240", // Gray
	BorderActive:   "57",  // Blue
	BorderInactive: "241", // Darker gray
}

var DarkColors = ColorPalette{
	Primary:    "25",
	Secondary:  "165",
	Accent:     "33",
	Background: "232",
	Surface:    "236",

	TextPrimary:   "252",
	TextSecondary: "220",
	TextMuted:     "244",
	TextHighlight: "165",

	Success: "40",
	Warning: "214",
	Error:   "160",
	Info:    "33",

	BorderNormal:   "236",
	BorderActive:   "25",
	BorderInactive: "244",
}

var LightColors = ColorPalette{
	Primary:    "27",
	Secondary:  "125",
	Accent:     "32",
	Background: "255",
	Surface:    "254",

	TextPrimary:   "16",
	TextSecondary: "17",
	TextMuted:     "240",
	TextHighlight: "125",

	Success: "28",
	Warning: "172",
	Error:   "160",
	Info:    "27",

	BorderNormal:   "240",
	BorderActive:   "27",
	BorderInactive: "250",
}

var Colors = DefaultColors

type Styles struct {
	// Base styles
	Base      lipgloss.Style
	Container lipgloss.Style

	// Typography
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Text        lipgloss.Style
	TextMuted   lipgloss.Style
	TextSuccess lipgloss.Style
	TextWarning lipgloss.Style
	TextError   lipgloss.Style

	// Navigation
	TabActive   lipgloss.Style
	TabInactive lipgloss.Style

	// Tables
	TableHeader   lipgloss.Style
	TableRow      lipgloss.Style
	TableBorder   lipgloss.Style
	TableSelected lipgloss.Style

	// Dialogs
	Dialog       lipgloss.Style
	DialogBorder lipgloss.Style

	// Help
	HelpTitle   lipgloss.Style
	HelpSection lipgloss.Style
	HelpFooter  lipgloss.Style

	// Status
	StatusBar   lipgloss.Style
	StatusError lipgloss.Style
	StatusInfo  lipgloss.Style
}

func NewStyles() *Styles {
	return &Styles{
		// Base styles
		Base: lipgloss.NewStyle(),
		Container: lipgloss.NewStyle().
			Padding(0, 1),

		// Typography
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(Colors.Secondary)),

		Subtitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(Colors.Accent)),

		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextPrimary)),

		TextMuted: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextMuted)),

		TextSuccess: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.Success)),

		TextWarning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.Warning)),

		TextError: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.Error)),

		// Navigation
		TabActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextSecondary)).
			Background(lipgloss.Color(Colors.Primary)).
			Padding(0, 1),

		TabInactive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextMuted)).
			Padding(0, 1),

		// Tables
		TableHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextHighlight)).
			Bold(true),

		TableRow: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextPrimary)),

		TableBorder: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(Colors.BorderNormal)),

		TableSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextSecondary)).
			Background(lipgloss.Color(Colors.Primary)).
			Bold(true),

		// Dialogs
		Dialog: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Background(lipgloss.Color(Colors.Background)).
			Foreground(lipgloss.Color(Colors.TextPrimary)),

		DialogBorder: lipgloss.NewStyle().
			BorderForeground(lipgloss.Color(Colors.Error)),

		// Help
		HelpTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(Colors.Secondary)).
			Align(lipgloss.Center),

		HelpSection: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(Colors.Accent)),

		HelpFooter: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextMuted)).
			Align(lipgloss.Center),

		// Status
		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.TextMuted)),

		StatusError: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.Error)),

		StatusInfo: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.Info)),
	}
}

// Global styles instance
var AppStyles = NewStyles()

// Theme switching functions
func SetTheme(colors ColorPalette) {
	Colors = colors
	AppStyles = NewStyles()
}

func SetDefaultTheme() {
	SetTheme(DefaultColors)
}

func SetDarkTheme() {
	SetTheme(DarkColors)
}

func SetLightTheme() {
	SetTheme(LightColors)
}

func StyleContainer(content string) string {
	return AppStyles.Container.Render(content)
}

func StyleTitle(text string) string {
	return AppStyles.Title.Render(text)
}

func StyleSubtitle(text string) string {
	return AppStyles.Subtitle.Render(text)
}

func StyleError(text string) string {
	return AppStyles.TextError.Render(text)
}

func StyleSuccess(text string) string {
	return AppStyles.TextSuccess.Render(text)
}

func StyleMuted(text string) string {
	return AppStyles.TextMuted.Render(text)
}

func StyleContainerStatus(status, state string) string {
	if state == "running" {
		return AppStyles.TextSuccess.Render("✓ " + status)
	}
	if state == "exited" {
		return AppStyles.TextMuted.Render("✗ " + status)
	}
	return AppStyles.TextError.Render("✗ " + status)
}

// StyleContainerStatusText returns just the status text with icon (no additional styling)
func StyleContainerStatusText(status, state string) string {
	if state == "running" {
		return "✓ " + status
	}
	return "✗ " + status
}

func StyleTab(text string, isActive bool) string {
	if isActive {
		return AppStyles.TabActive.Render(text)
	}
	return AppStyles.TabInactive.Render(text)
}

func StyleTitleWithWidth(text string, width int) string {
	return AppStyles.Title.Width(width).Render(text)
}

func StyleHelpTitleWithWidth(text string, width int) string {
	return AppStyles.HelpTitle.Width(width).Render(text)
}
