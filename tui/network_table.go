package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/network"
)

type NetworkTable struct {
	table table.Model
	model *Model
}

func NewNetworkTable(m *Model) *NetworkTable {
	networkColumns := []table.Column{
		{Title: "ID", Width: 12},
		{Title: "Name", Width: 20},
		{Title: "Driver", Width: 15},
		{Title: "Scope", Width: 10},
		{Title: "Created", Width: 15},
	}

	networkTable := table.New(
		table.WithColumns(networkColumns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	nt := &NetworkTable{
		table: networkTable,
		model: m,
	}

	nt.applyStyles()
	return nt
}

func (nt *NetworkTable) GetTable() table.Model {
	return nt.table
}

func (nt *NetworkTable) SetTable(t table.Model) {
	nt.table = t
}

func (nt *NetworkTable) applyStyles() {
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

	nt.table.SetStyles(s)
}

func (nt *NetworkTable) RefreshStyles() {
	nt.applyStyles()
}

func (nt *NetworkTable) SetHeight(height int) {
	nt.table.SetHeight(height)
}

func (nt *NetworkTable) SetColumns(columns []table.Column) {
	nt.table.SetColumns(columns)
}

func (nt *NetworkTable) Update() {
	rows := make([]table.Row, len(nt.model.networks))
	for i, net := range nt.model.networks {
		created := ""
		if !net.Created.IsZero() {
			created = net.Created.Format("2006-01-02 15:04")
		}

		rows[i] = table.Row{
			net.ID[:12],
			net.Name,
			net.Driver,
			net.Scope,
			created,
		}
	}
	nt.table.SetRows(rows)
}

func (nt *NetworkTable) GetSelectedNetwork() *network.Summary {
	cursor := nt.table.Cursor()
	if cursor >= 0 && cursor < len(nt.model.networks) {
		return &nt.model.networks[cursor]
	}
	return nil
}

func (nt *NetworkTable) View() string {
	return nt.table.View()
}

func (nt *NetworkTable) Cursor() int {
	return nt.table.Cursor()
}
