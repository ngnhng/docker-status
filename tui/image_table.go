package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/image"
)

// ImageTable manages the image table functionality
type ImageTable struct {
	table table.Model
	model *Model
}

// NewImageTable creates a new image table
func NewImageTable(m *Model) *ImageTable {
	imageColumns := []table.Column{
		{Title: "Repository", Width: 30},
		{Title: "Tag", Width: 15},
		{Title: "Image ID", Width: 12},
		{Title: "Created", Width: 15},
		{Title: "Size", Width: 12},
	}

	imageTable := table.New(
		table.WithColumns(imageColumns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	it := &ImageTable{
		table: imageTable,
		model: m,
	}

	it.applyStyles()
	return it
}

// GetTable returns the underlying table model
func (it *ImageTable) GetTable() table.Model {
	return it.table
}

// SetTable updates the underlying table model
func (it *ImageTable) SetTable(t table.Model) {
	it.table = t
}

// applyStyles applies the current theme styles to the table
func (it *ImageTable) applyStyles() {
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

	it.table.SetStyles(s)
}

// RefreshStyles reapplies the current styles (useful after theme changes)
func (it *ImageTable) RefreshStyles() {
	it.applyStyles()
}

// SetHeight sets the table height
func (it *ImageTable) SetHeight(height int) {
	it.table.SetHeight(height)
}

// SetColumns updates the table columns
func (it *ImageTable) SetColumns(columns []table.Column) {
	it.table.SetColumns(columns)
}

// Update updates the table data based on current images
func (it *ImageTable) Update() {
	rows := make([]table.Row, len(it.model.images))
	for i, img := range it.model.images {
		created := time.Unix(img.Created, 0).Format("2006-01-02 15:04")

		// Format size in human readable format
		size := formatSize(img.Size)

		// Handle repository and tag
		repo := "<none>"
		tag := "<none>"
		if len(img.RepoTags) > 0 && img.RepoTags[0] != "<none>:<none>" {
			parts := parseRepoTag(img.RepoTags[0])
			repo = parts[0]
			tag = parts[1]
		}

		rows[i] = table.Row{
			repo,
			tag,
			img.ID[7:19], // Remove "sha256:" prefix and show first 12 chars
			created,
			size,
		}
	}
	it.table.SetRows(rows)
}

// GetSelectedImage returns the currently selected image, if any
func (it *ImageTable) GetSelectedImage() *image.Summary {
	cursor := it.table.Cursor()
	if cursor >= 0 && cursor < len(it.model.images) {
		return &it.model.images[cursor]
	}
	return nil
}

// View returns the rendered table view
func (it *ImageTable) View() string {
	return it.table.View()
}

// Cursor returns the current cursor position
func (it *ImageTable) Cursor() int {
	return it.table.Cursor()
}

// parseRepoTag splits a repo:tag string into repository and tag parts
func parseRepoTag(repoTag string) [2]string {
	if repoTag == "" {
		return [2]string{"<none>", "<none>"}
	}

	// Find the last colon to separate repo from tag
	lastColon := -1
	for i := len(repoTag) - 1; i >= 0; i-- {
		if repoTag[i] == ':' {
			lastColon = i
			break
		}
	}

	if lastColon == -1 {
		return [2]string{repoTag, "<none>"}
	}

	repo := repoTag[:lastColon]
	tag := repoTag[lastColon+1:]

	if repo == "" {
		repo = "<none>"
	}
	if tag == "" {
		tag = "<none>"
	}

	return [2]string{repo, tag}
}
