package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/container"
)

// formatPorts formats container ports for display
func formatPorts(ports []container.Port) string {
	if len(ports) == 0 {
		return ""
	}

	var formatted []string
	for _, port := range ports {
		if port.PublicPort > 0 {
			formatted = append(formatted, fmt.Sprintf("%d->%d/%s", port.PublicPort, port.PrivatePort, port.Type))
		} else {
			formatted = append(formatted, fmt.Sprintf("%d/%s", port.PrivatePort, port.Type))
		}
	}

	result := strings.Join(formatted, ", ")
	if len(result) > 23 {
		return result[:20] + "..."
	}
	return result
}

// formatSize formats a byte size into human readable format
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fGB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1fMB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1fKB", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}

// groupContainersByCompose groups containers by their Docker Compose project
func (m *Model) groupContainersByCompose() []ContainerGroup {
	groups := make(map[string][]container.Summary)

	for _, container := range m.containers {
		projectName := getComposeProjectName(container)
		groups[projectName] = append(groups[projectName], container)
	}

	var result []ContainerGroup
	for name, containers := range groups {
		// Sort containers within each group by name
		sort.Slice(containers, func(i, j int) bool {
			return containers[i].Names[0] < containers[j].Names[0]
		})

		result = append(result, ContainerGroup{
			Name:       name,
			Containers: containers,
		})
	}

	// Sort groups by name for consistent display
	sort.Slice(result, func(i, j int) bool {
		// Put "Standalone Containers" at the end
		if result[i].Name == "Standalone Containers" {
			return false
		}
		if result[j].Name == "Standalone Containers" {
			return true
		}
		return result[i].Name < result[j].Name
	})

	return result
}

// getComposeProjectName extracts the Docker Compose project name from container labels
func getComposeProjectName(container container.Summary) string {
	// Try Docker Compose project label
	if projectName, exists := container.Labels["com.docker.compose.project"]; exists {
		return projectName
	}

	// Try Docker Stack label (for swarm mode)
	if stackName, exists := container.Labels["com.docker.stack.namespace"]; exists {
		return stackName + " (stack)"
	}

	// Check if it's a standalone container (no compose labels)
	hasComposeLabels := false
	for key := range container.Labels {
		if strings.HasPrefix(key, "com.docker.compose.") {
			hasComposeLabels = true
			break
		}
	}

	if !hasComposeLabels {
		return "Standalone Containers"
	}

	// Fallback to "Unknown Project" if we can't determine the project
	return "Unknown Project"
}

// getSelectedItem returns information about the currently selected item in grouped mode
func (m *Model) getSelectedItem() (groupIndex int, containerIndex int, isGroupHeader bool) {
	if !m.groupByCompose || len(m.containerGroups) == 0 {
		return -1, -1, false
	}

	cursor := m.containerTable.Cursor()
	currentRow := 0

	for groupIdx, group := range m.containerGroups {
		// Check if cursor is on group header
		if currentRow == cursor {
			return groupIdx, -1, true
		}
		currentRow++

		// Check if cursor is on any container in this group
		for containerIdx := range group.Containers {
			if currentRow == cursor {
				return groupIdx, containerIdx, false
			}
			currentRow++
		}
	}

	return -1, -1, false
}

func (m *Model) getSelectedGroup() *ContainerGroup {
	groupIdx, _, _ := m.getSelectedItem()
	if groupIdx >= 0 && groupIdx < len(m.containerGroups) {
		return &m.containerGroups[groupIdx]
	}
	return nil
}

func (m *Model) getSelectedContainer() *container.Summary {
	groupIdx, containerIdx, isGroupHeader := m.getSelectedItem()
	if isGroupHeader || groupIdx < 0 || containerIdx < 0 {
		return nil
	}

	if groupIdx < len(m.containerGroups) && containerIdx < len(m.containerGroups[groupIdx].Containers) {
		return &m.containerGroups[groupIdx].Containers[containerIdx]
	}
	return nil
}
