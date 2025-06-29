package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/volume"
)

// VolumeTable manages the volume table functionality
type VolumeTable struct {
	table table.Model
	model *Model
}

// NewVolumeTable creates a new volume table
func NewVolumeTable(m *Model) *VolumeTable {
	volumeColumns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Driver", Width: 15},
		{Title: "Mountpoint", Width: 50},
		{Title: "Created", Width: 15},
	}

	volumeTable := table.New(
		table.WithColumns(volumeColumns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	vt := &VolumeTable{
		table: volumeTable,
		model: m,
	}

	vt.applyStyles()
	return vt
}

// GetTable returns the underlying table model
func (vt *VolumeTable) GetTable() table.Model {
	return vt.table
}

// SetTable updates the underlying table model
func (vt *VolumeTable) SetTable(t table.Model) {
	vt.table = t
}

// applyStyles applies the current theme styles to the table
func (vt *VolumeTable) applyStyles() {
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

	vt.table.SetStyles(s)
}

// RefreshStyles reapplies the current styles (useful after theme changes)
func (vt *VolumeTable) RefreshStyles() {
	vt.applyStyles()
}

// SetHeight sets the table height
func (vt *VolumeTable) SetHeight(height int) {
	vt.table.SetHeight(height)
}

// SetColumns updates the table columns
func (vt *VolumeTable) SetColumns(columns []table.Column) {
	vt.table.SetColumns(columns)
}

// Update updates the table data based on current volumes
func (vt *VolumeTable) Update() {
	rows := make([]table.Row, len(vt.model.volumes))
	for i, vol := range vt.model.volumes {
		created := ""
		if vol.CreatedAt != "" {
			created = vol.CreatedAt[:16] // Show date and time portion
		}

		rows[i] = table.Row{
			vol.Name,
			vol.Driver,
			vol.Mountpoint,
			created,
		}
	}
	vt.table.SetRows(rows)
}

// GetSelectedVolume returns the currently selected volume, if any
func (vt *VolumeTable) GetSelectedVolume() *volume.Volume {
	cursor := vt.table.Cursor()
	if cursor >= 0 && cursor < len(vt.model.volumes) {
		return vt.model.volumes[cursor]
	}
	return nil
}

// View returns the rendered table view
func (vt *VolumeTable) View() string {
	return vt.table.View()
}

// Cursor returns the current cursor position
func (vt *VolumeTable) Cursor() int {
	return vt.table.Cursor()
}
