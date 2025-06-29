package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
)

// ContainerTable manages the container table functionality
type ContainerTable struct {
	table           table.Model
	model           *Model
	containerStates []string // Track container states for styling
}

func NewContainerTable(m *Model) *ContainerTable {
	containerColumns := []table.Column{
		{Title: "ID", Width: 12},
		{Title: "Names", Width: 20},
		{Title: "Image", Width: 30},
		{Title: "Command", Width: 20},
		{Title: "Created", Width: 15},
		{Title: "Status", Width: 20},
		{Title: "Ports", Width: 25},
	}

	containerTable := table.New(
		table.WithColumns(containerColumns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	ct := &ContainerTable{
		table: containerTable,
		model: m,
	}

	ct.applyStyles()
	return ct
}

// GetTable returns the underlying table model
func (ct *ContainerTable) GetTable() table.Model {
	return ct.table
}

// SetTable updates the underlying table model
func (ct *ContainerTable) SetTable(t table.Model) {
	ct.table = t
}

// applyStyles applies the current theme styles to the table
func (ct *ContainerTable) applyStyles() {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(Colors.BorderNormal)).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color(Colors.TextHighlight))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(Colors.TextSecondary)).
		Background(lipgloss.Color(Colors.Primary)).
		Bold(true)

	ct.table.SetStyles(s)
}

// RefreshStyles reapplies the current styles (useful after theme changes)
func (ct *ContainerTable) RefreshStyles() {
	ct.applyStyles()
}

// SetHeight sets the table height
func (ct *ContainerTable) SetHeight(height int) {
	ct.table.SetHeight(height)
}

// SetColumns updates the table columns
func (ct *ContainerTable) SetColumns(columns []table.Column) {
	ct.table.SetColumns(columns)
}

// Update updates the table data based on current containers
func (ct *ContainerTable) Update() {
	if ct.model.groupByCompose {
		// Update container groups whenever we switch to grouped mode
		ct.UpdateContainerGroups()
		ct.updateGrouped()
	} else {
		ct.updateFlat()
	}
}

func (ct *ContainerTable) updateFlat() {
	rows := make([]table.Row, len(ct.model.containers))
	ct.containerStates = make([]string, len(ct.model.containers))

	for i, container := range ct.model.containers {
		created := time.Unix(container.Created, 0).Format("2006-01-02 15:04")

		status := StyleContainerStatusText(container.Status, container.State)
		ports := formatPorts(container.Ports)

		names := strings.Join(container.Names, ", ")
		if len(names) > 0 && names[0] == '/' {
			names = names[1:] // Remove leading slash
		}

		// Store plain text for proper layout
		rows[i] = table.Row{
			container.ID[:12],
			names,
			container.Image,
			container.Command,
			created,
			status,
			ports,
		}

		// Track container state for styling
		ct.containerStates[i] = container.State
	}
	ct.table.SetRows(rows)
}

// updateGrouped displays containers grouped by Docker Compose project
func (ct *ContainerTable) updateGrouped() {
	var rows []table.Row
	var states []string

	for _, group := range ct.model.containerGroups {
		// Add group header
		groupStatus := ct.getGroupStatus(group)
		groupPorts := ct.getGroupPorts(group)

		groupRow := table.Row{
			"", // No ID for group
			fmt.Sprintf("📁 %s (%d containers)", group.Name, len(group.Containers)),
			"",
			"",
			groupStatus,
			groupPorts,
			"",
		}
		rows = append(rows, groupRow)
		states = append(states, "group") // Special state for group headers

		// Add all containers in the group (expanded by default)
		for _, container := range group.Containers {
			created := time.Unix(container.Created, 0).Format("2006-01-02 15:04")
			status := StyleContainerStatusText(container.Status, container.State)
			ports := formatPorts(container.Ports)

			names := strings.Join(container.Names, ", ")
			if len(names) > 0 && names[0] == '/' {
				names = names[1:] // Remove leading slash
			}

			containerRow := table.Row{
				"  " + container.ID[:12],
				" |" + names,
				"  " + container.Image,
				"  " + container.Command,
				"  " + created,
				"  " + status,
				"  " + ports,
			}
			rows = append(rows, containerRow)
			states = append(states, container.State)
		}
	}

	ct.table.SetRows(rows)
	ct.containerStates = states
}

// getGroupStatus returns a summary status for a container group
func (ct *ContainerTable) getGroupStatus(group ContainerGroup) string {
	if len(group.Containers) == 0 {
		return ""
	}

	running := 0
	total := len(group.Containers)

	for _, container := range group.Containers {
		if container.State == "running" {
			running++
		}
	}

	if running == total {
		return fmt.Sprintf("%d/%d running", running, total)
	} else if running == 0 {
		return fmt.Sprintf("%d/%d stopped", total-running, total)
	} else {
		return fmt.Sprintf("%d/%d running", running, total)
	}
}

// getGroupPorts returns a summary of ports for a container group
func (ct *ContainerTable) getGroupPorts(group ContainerGroup) string {
	portSet := make(map[string]bool)

	for _, container := range group.Containers {
		for _, port := range container.Ports {
			if port.PublicPort != 0 {
				portStr := fmt.Sprintf("%d->%d", port.PublicPort, port.PrivatePort)
				portSet[portStr] = true
			}
		}
	}

	if len(portSet) == 0 {
		return ""
	}

	var ports []string
	for port := range portSet {
		ports = append(ports, port)
	}

	// Sort ports for consistent display
	sort.Strings(ports)

	if len(ports) > 3 {
		return fmt.Sprintf("%s... (%d ports)", strings.Join(ports[:3], ", "), len(ports))
	}

	return strings.Join(ports, ", ")
}

// GetSelectedContainer returns the currently selected container, if any
func (ct *ContainerTable) GetSelectedContainer() *container.Summary {
	if ct.model.groupByCompose {
		return ct.getSelectedContainerGrouped()
	}
	return ct.getSelectedContainerFlat()
}

// getSelectedContainerFlat returns selected container in flat mode
func (ct *ContainerTable) getSelectedContainerFlat() *container.Summary {
	cursor := ct.table.Cursor()
	if cursor >= 0 && cursor < len(ct.model.containers) {
		return &ct.model.containers[cursor]
	}
	return nil
}

// getSelectedContainerGrouped returns selected container in grouped mode
func (ct *ContainerTable) getSelectedContainerGrouped() *container.Summary {
	cursor := ct.table.Cursor()
	if cursor < 0 {
		return nil
	}

	currentRow := 0
	for _, group := range ct.model.containerGroups {
		// Group header row
		if cursor == currentRow {
			// Group header selected, no specific container
			return nil
		}
		currentRow++

		// Container rows (now all containers are shown)
		for containerIndex := range group.Containers {
			if cursor == currentRow {
				return &group.Containers[containerIndex]
			}
			currentRow++
		}
	}

	return nil
}

// View returns the rendered table view with custom row styling
func (ct *ContainerTable) View() string {
	// Get the base table view
	baseView := ct.table.View()

	// Split into lines for row-by-row styling
	lines := strings.Split(baseView, "\n")
	if len(lines) == 0 {
		return baseView
	}

	// Apply styling to data rows (skip header and borders)
	styledLines := make([]string, len(lines))
	dataRowIndex := 0

	for i, line := range lines {
		// Skip header rows and border lines
		if i == 0 || i == 1 || strings.TrimSpace(line) == "" ||
			strings.Contains(line, "─") || strings.Contains(line, "│") {
			styledLines[i] = line
			continue
		}

		// Check if we have state information for this row
		if dataRowIndex < len(ct.containerStates) {
			state := ct.containerStates[dataRowIndex]

			switch state {
			case "running":
				// Keep normal text for running containers
				styledLines[i] = line
			case "exited":
				// Muted text for exited containers
				styledLines[i] = AppStyles.TextMuted.Render(line)
			case "group":
				// Keep normal text for group headers
				styledLines[i] = line
			default:
				// Red text for error states
				styledLines[i] = AppStyles.TextError.Render(line)
			}
			dataRowIndex++
		} else {
			styledLines[i] = line
		}
	}

	return strings.Join(styledLines, "\n")
}

// Cursor returns the current cursor position
func (ct *ContainerTable) Cursor() int {
	return ct.table.Cursor()
}

// UpdateContainerGroups updates the container groups for grouped display
func (ct *ContainerTable) UpdateContainerGroups() {
	ct.model.containerGroups = ct.model.groupContainersByCompose()
}
