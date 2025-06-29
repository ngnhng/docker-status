package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type ViewMode int

const (
	ContainersView ViewMode = iota
	ImagesView
	NetworksView
	VolumesView
	LogsView
	HelpViewMode
)

var viewNames = map[ViewMode]string{
	ContainersView: "Containers",
	ImagesView:     "Images",
	NetworksView:   "Networks",
	VolumesView:    "Volumes",
	LogsView:       "Logs",
	HelpViewMode:   "Help",
}

// ContainerGroup represents a group of containers (e.g., from the same Docker Compose project)
type ContainerGroup struct {
	Name       string
	Containers []container.Summary
}

type Model struct {
	// Docker client
	dockerClient *client.Client
	dockerCli    command.Cli
	ctx          context.Context

	// Current view
	currentView ViewMode

	// Tables for different views
	containerTable *ContainerTable
	imageTable     *ImageTable
	networkTable   *NetworkTable
	volumeTable    *VolumeTable

	// Data
	containers []container.Summary
	images     []image.Summary
	networks   []network.Summary
	volumes    []*volume.Volume

	// Container grouping
	containerGroups []ContainerGroup
	groupByCompose  bool
	selectedGroup   int // Index of currently selected group when grouped
	selectedInGroup int // Index of container within the group (-1 for group header)

	// UI state
	width  int
	height int
	ready  bool

	// Status and error handling
	status string
	err    error

	// Refresh ticker
	ticker *time.Ticker

	// Key bindings
	keys KeyMap

	// Help view
	helpView *HelpView
	showHelp bool

	// Confirmation dialog
	confirmDialog *ConfirmationDialog
	pendingAction func() tea.Cmd

	// Styles
	styles *Styles
}

// Compile time check to ensure Model implements tea.Model interface
var _ tea.Model = (*Model)(nil)

// KeyMap defines the key bindings
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	Tab          key.Binding
	Enter        key.Binding
	Refresh      key.Binding
	Quit         key.Binding
	Help         key.Binding
	Delete       key.Binding
	Logs         key.Binding
	Containers   key.Binding
	Images       key.Binding
	Networks     key.Binding
	Volumes      key.Binding
	ThemeDefault key.Binding
	ThemeDark    key.Binding
	ThemeLight   key.Binding
	GroupToggle  key.Binding
	GroupStop    key.Binding
	GroupStart   key.Binding
	GroupDelete  key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next view"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "ctrl+r"),
			key.WithHelp("r", "refresh"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d", "delete"),
			key.WithHelp("d", "delete"),
		),
		Logs: key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "logs"),
		),
		Containers: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "containers"),
		),
		Images: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "images"),
		),
		Networks: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "networks"),
		),
		Volumes: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "volumes"),
		),
		ThemeDefault: key.NewBinding(
			key.WithKeys("t", "1"),
			key.WithHelp("t+1", "default theme"),
		),
		ThemeDark: key.NewBinding(
			key.WithKeys("t", "2"),
			key.WithHelp("t+2", "dark theme"),
		),
		ThemeLight: key.NewBinding(
			key.WithKeys("t", "3"),
			key.WithHelp("t+3", "light theme"),
		),
		GroupToggle: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "group by compose"),
		),
		GroupStop: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stop group"),
		),
		GroupStart: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "start group"),
		),
		GroupDelete: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "delete group"),
		),
	}
}

// tickMsg is sent periodically to refresh data
type tickMsg time.Time

func NewModel(dockerClient *client.Client, dockerCli command.Cli) *Model {
	m := &Model{
		dockerClient: dockerClient,
		dockerCli:    dockerCli,
		ctx:          context.Background(),
		currentView:  ContainersView,
		keys:         DefaultKeyMap(),
		ticker:       time.NewTicker(5 * time.Second), // TODO: allow configurable interval
		styles:       NewStyles(),
	}

	m.initTables()

	m.helpView = NewHelpView()
	m.showHelp = false

	return m
}

// initTables initializes all the table models
func (m *Model) initTables() {
	m.containerTable = NewContainerTable(m)
	m.imageTable = NewImageTable(m)
	m.networkTable = NewNetworkTable(m)
	m.volumeTable = NewVolumeTable(m)
}

func (m *Model) SetTheme(colors ColorPalette) {
	SetTheme(colors)
	m.styles = NewStyles()
	m.refreshTableStyles()
}

func (m *Model) SetDefaultTheme() {
	m.SetTheme(DefaultColors)
}

func (m *Model) SetDarkTheme() {
	m.SetTheme(DarkColors)
}

func (m *Model) SetLightTheme() {
	m.SetTheme(LightColors)
}

// refreshTableStyles applies the current styles to all tables
func (m *Model) refreshTableStyles() {
	m.containerTable.RefreshStyles()
	m.imageTable.RefreshStyles()
	m.networkTable.RefreshStyles()
	m.volumeTable.RefreshStyles()
}
