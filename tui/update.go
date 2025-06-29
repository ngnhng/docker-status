package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
)

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.refreshData(),
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.updateTableSizes()
		m.helpView.SetSize(msg.Width, msg.Height)

	case tickMsg:
		cmds = append(cmds, m.refreshData())
		cmds = append(cmds, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}))

	case tea.KeyMsg:
		// Handle confirmation dialog first, if it's visible
		if m.confirmDialog != nil && m.confirmDialog.IsVisible() {
			switch msg.String() {
			case "y", "Y":
				if m.pendingAction != nil {
					cmds = append(cmds, m.pendingAction())
					m.pendingAction = nil
				}
				m.confirmDialog.Hide()
			case "n", "N", "esc":
				m.confirmDialog.Hide()
				m.pendingAction = nil
			}
			// Don't process other keys when dialog is visible
			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			m.ticker.Stop()
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			m.nextView()

		case key.Matches(msg, m.keys.Containers):
			m.currentView = ContainersView

		case key.Matches(msg, m.keys.Images):
			m.currentView = ImagesView

		case key.Matches(msg, m.keys.Networks):
			m.currentView = NetworksView

		case key.Matches(msg, m.keys.Volumes):
			m.currentView = VolumesView

		case key.Matches(msg, m.keys.Refresh):
			cmds = append(cmds, m.refreshData())

		case key.Matches(msg, m.keys.Delete):
			// Show confirmation dialog
			m.showDeleteConfirmation()

		case key.Matches(msg, m.keys.Stop):
			if m.currentView == ContainersView {
				m.showStopConfirmation()
			}

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp

		case key.Matches(msg, m.keys.Logs):
			if m.currentView == ContainersView {
				cmds = append(cmds, m.showLogs())
			}

		case key.Matches(msg, m.keys.ThemeDefault):
			m.SetDefaultTheme()

		case key.Matches(msg, m.keys.ThemeDark):
			m.SetDarkTheme()

		case key.Matches(msg, m.keys.ThemeLight):
			m.SetLightTheme()

		case key.Matches(msg, m.keys.GroupToggle):
			if m.currentView == ContainersView {
				m.groupByCompose = !m.groupByCompose
				m.containerTable.Update()
			}

		case key.Matches(msg, m.keys.GroupStop):
			if m.currentView == ContainersView && m.groupByCompose {
				if group := m.getSelectedGroup(); group != nil {
					cmds = append(cmds, m.stopGroup(*group))
				}
			}

		case key.Matches(msg, m.keys.GroupStart):
			if m.currentView == ContainersView && m.groupByCompose {
				if group := m.getSelectedGroup(); group != nil {
					cmds = append(cmds, m.startGroup(*group))
				}
			}

		case key.Matches(msg, m.keys.GroupDelete):
			if m.currentView == ContainersView && m.groupByCompose {
				if group := m.getSelectedGroup(); group != nil {
					// Show confirmation dialog for group deletion
					message := fmt.Sprintf("Are you sure you want to delete the entire compose stack '%s'? This will remove all containers, networks, volumes, and images.", group.Name)
					m.confirmDialog = NewConfirmationDialog(message)
					m.confirmDialog.SetSize(m.width, m.height)
					m.confirmDialog.Show()
					m.pendingAction = func() tea.Cmd {
						return m.deleteGroup(*group)
					}
				}
			}
		}

	case dataRefreshedMsg:
		m.handleDataRefresh(msg)

	case errorMsg:
		m.err = msg.error
		m.status = ""

	case statusMsg:
		m.status = string(msg)
		m.err = nil
	}

	// Only update tables if confirmation dialog is not visible
	if m.confirmDialog == nil || !m.confirmDialog.IsVisible() {
		switch m.currentView {
		case ContainersView:
			table := m.containerTable.GetTable()
			table, cmd = table.Update(msg)
			m.containerTable.SetTable(table)
		case ImagesView:
			table := m.imageTable.GetTable()
			table, cmd = table.Update(msg)
			m.imageTable.SetTable(table)
		case NetworksView:
			table := m.networkTable.GetTable()
			table, cmd = table.Update(msg)
			m.networkTable.SetTable(table)
		case VolumesView:
			table := m.volumeTable.GetTable()
			table, cmd = table.Update(msg)
			m.volumeTable.SetTable(table)
		}
	}

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.showHelp {
		return m.helpView.Render()
	}

	var content strings.Builder

	// Header
	content.WriteString(m.renderHeader())
	content.WriteString("\n\n")

	// Main content
	switch m.currentView {
	case ContainersView:
		content.WriteString(m.containerTable.View())
	case ImagesView:
		content.WriteString(m.imageTable.View())
	case NetworksView:
		content.WriteString(m.networkTable.View())
	case VolumesView:
		content.WriteString(m.volumeTable.View())
	case LogsView:
		content.WriteString("Container logs view (coming soon)")
	}

	content.WriteString("\n\n")

	// Footer
	content.WriteString(m.renderFooter())

	view := content.String()

	// Overlay confirmation dialog if visible
	if m.confirmDialog != nil && m.confirmDialog.IsVisible() {
		dialog := m.confirmDialog.Render()
		return dialog
	}

	return view
}

func (m *Model) nextView() {
	switch m.currentView {
	case ContainersView:
		m.currentView = ImagesView
	case ImagesView:
		m.currentView = NetworksView
	case NetworksView:
		m.currentView = VolumesView
	case VolumesView:
		m.currentView = ContainersView
	}
}

func (m *Model) renderHeader() string {
	var tabs []string

	for view := ContainersView; view <= VolumesView; view++ {
		name := viewNames[view]
		tabs = append(tabs, StyleTab(name, view == m.currentView))
	}

	title := StyleTitle("Docker TUI (Beta v0.1.0 - Bugs are expected) - " + "Made by github.com/ngnhng, ©MIT License")

	return fmt.Sprintf("%s\n\n%s", title, strings.Join(tabs, " │ "))
}

func (m *Model) renderFooter() string {
	var help []string

	switch m.currentView {
	case ContainersView:
		help = []string{
			"1-4: switch views",
			"↑/↓: navigate",
			"r: refresh",
			"g: group toggle",
			"ctrl+s: stop",
			"d: delete",
			"L: logs",
			"q: quit",
		}
		if m.groupByCompose {
			help = append(help, "[grouped by compose]")
			help = append(help, "s: stop group", "S: start group", "D: delete group")
		}
	default:
		help = []string{
			"1-4: switch views",
			"↑/↓: navigate",
			"r: refresh",
			"d: delete",
			"q: quit",
		}
	}

	helpText := StyleMuted(strings.Join(help, " • "))

	status := ""
	if m.status != "" {
		status = m.styles.StatusInfo.Render(m.status)
	}

	if m.err != nil {
		status = StyleError(fmt.Sprintf("Error: %v", m.err))
	}

	return fmt.Sprintf("%s\n%s", helpText, status)
}

func (m *Model) updateTableSizes() {
	if m.height <= 0 || m.width <= 0 {
		return
	}

	tableHeight := m.height - 10 // Reserve space for header and footer
	if tableHeight < 5 {
		tableHeight = 5
	}

	m.containerTable.SetHeight(tableHeight)
	m.imageTable.SetHeight(tableHeight)
	m.networkTable.SetHeight(tableHeight)
	m.volumeTable.SetHeight(tableHeight)

	m.updateColumnWidths()
}

func (m *Model) updateColumnWidths() {
	// Reserve some space for borders and padding
	availableWidth := m.width - 4

	containerColumns := m.calculateContainerColumnWidths(availableWidth)
	m.containerTable.SetColumns(containerColumns)

	imageColumns := m.calculateImageColumnWidths(availableWidth)
	m.imageTable.SetColumns(imageColumns)

	networkColumns := m.calculateNetworkColumnWidths(availableWidth)
	m.networkTable.SetColumns(networkColumns)

	volumeColumns := m.calculateVolumeColumnWidths(availableWidth)
	m.volumeTable.SetColumns(volumeColumns)
}

// TODO: move these to a separate file for better organization
func (m *Model) calculateContainerColumnWidths(availableWidth int) []table.Column {
	minWidths := []int{12, 15, 15, 15, 15, 15, 15}
	preferredWidths := []int{12, 30, 25, 16, 20, 25, 20}
	titles := []string{"ID", "Names", "Image", "Command", "Created", "Status", "Ports"}

	return m.distributeColumnWidths(titles, minWidths, preferredWidths, availableWidth)
}

func (m *Model) calculateImageColumnWidths(availableWidth int) []table.Column {
	minWidths := []int{15, 10, 12, 15, 10} // Repository, Tag, Image ID, Created, Size
	preferredWidths := []int{30, 15, 12, 16, 12}
	titles := []string{"Repository", "Tag", "Image ID", "Created", "Size"}

	return m.distributeColumnWidths(titles, minWidths, preferredWidths, availableWidth)
}

func (m *Model) calculateNetworkColumnWidths(availableWidth int) []table.Column {
	minWidths := []int{12, 15, 10, 8, 15} // ID, Name, Driver, Scope, Created
	preferredWidths := []int{12, 25, 15, 10, 16}
	titles := []string{"ID", "Name", "Driver", "Scope", "Created"}

	return m.distributeColumnWidths(titles, minWidths, preferredWidths, availableWidth)
}

func (m *Model) calculateVolumeColumnWidths(availableWidth int) []table.Column {
	minWidths := []int{15, 10, 20, 15} // Name, Driver, Mountpoint, Created
	preferredWidths := []int{25, 15, 40, 16}
	titles := []string{"Name", "Driver", "Mountpoint", "Created"}

	return m.distributeColumnWidths(titles, minWidths, preferredWidths, availableWidth)
}

func (m *Model) distributeColumnWidths(titles []string, minWidths, preferredWidths []int, availableWidth int) []table.Column {
	if len(titles) != len(minWidths) || len(titles) != len(preferredWidths) {
		columns := make([]table.Column, len(titles))
		for i, title := range titles {
			width := 10
			if i < len(minWidths) {
				width = minWidths[i]
			}
			columns[i] = table.Column{Title: title, Width: width}
		}
		return columns
	}

	columns := make([]table.Column, len(titles))

	totalMinWidth := 0
	for _, w := range minWidths {
		totalMinWidth += w
	}

	// If we don't have enough space for minimum widths, scale everything down proportionally
	if availableWidth < totalMinWidth {
		scale := float64(availableWidth) / float64(totalMinWidth)
		for i, title := range titles {
			width := int(float64(minWidths[i]) * scale)
			if width < 8 {
				width = 8 // Absolute minimum
			}
			columns[i] = table.Column{Title: title, Width: width}
		}
		return columns
	}

	totalPreferredWidth := 0
	for _, w := range preferredWidths {
		totalPreferredWidth += w
	}

	if availableWidth >= totalPreferredWidth {
		extraSpace := availableWidth - totalPreferredWidth
		extraPerColumn := extraSpace / len(titles)
		remainder := extraSpace % len(titles)

		for i, title := range titles {
			width := preferredWidths[i] + extraPerColumn
			if i < remainder {
				width++ // Distribute remainder to first few columns
			}
			columns[i] = table.Column{Title: title, Width: width}
		}
		return columns
	}

	scale := float64(availableWidth) / float64(totalPreferredWidth)
	for i, title := range titles {
		width := int(float64(preferredWidths[i]) * scale)
		if width < minWidths[i] {
			width = minWidths[i]
		}
		columns[i] = table.Column{Title: title, Width: width}
	}

	return columns
}

func (m *Model) refreshData() tea.Cmd {
	return func() tea.Msg {
		containers, err := m.dockerClient.ContainerList(m.ctx, container.ListOptions{All: true})
		if err != nil {
			return errorMsg{err}
		}

		images, err := m.dockerClient.ImageList(m.ctx, image.ListOptions{})
		if err != nil {
			return errorMsg{err}
		}

		networks, err := m.dockerClient.NetworkList(m.ctx, network.ListOptions{})
		if err != nil {
			return errorMsg{err}
		}

		volumeList, err := m.dockerClient.VolumeList(m.ctx, volume.ListOptions{})
		if err != nil {
			return errorMsg{err}
		}

		return dataRefreshedMsg{
			containers: containers,
			images:     images,
			networks:   networks,
			volumes:    volumeList.Volumes,
		}
	}
}

func (m *Model) showLogs() tea.Cmd {
	return func() tea.Msg {
		// TODO
		return statusMsg("Logs view not implemented yet")
	}
}

func (m *Model) showDeleteConfirmation() {
	var message string
	var hasSelection bool

	switch m.currentView {
	case ContainersView:
		if container := m.containerTable.GetSelectedContainer(); container != nil {
			message = fmt.Sprintf("Are you sure you want to delete container '%s'?", container.ID[:12])
			hasSelection = true
			m.pendingAction = func() tea.Cmd {
				return m.deleteContainer(*container)
			}
		}

	case ImagesView:
		if img := m.imageTable.GetSelectedImage(); img != nil {
			repoTag := "<none>:<none>"
			if len(img.RepoTags) > 0 {
				repoTag = img.RepoTags[0]
			}
			message = fmt.Sprintf("Are you sure you want to delete image '%s'?", repoTag)
			hasSelection = true
			m.pendingAction = func() tea.Cmd {
				return m.deleteImage(*img)
			}
		}

	case NetworksView:
		if net := m.networkTable.GetSelectedNetwork(); net != nil {
			message = fmt.Sprintf("Are you sure you want to delete network '%s'?", net.Name)
			hasSelection = true
			m.pendingAction = func() tea.Cmd {
				return m.deleteNetwork(*net)
			}
		}

	case VolumesView:
		if vol := m.volumeTable.GetSelectedVolume(); vol != nil {
			message = fmt.Sprintf("Are you sure you want to delete volume '%s'?", vol.Name)
			hasSelection = true
			m.pendingAction = func() tea.Cmd {
				return m.deleteVolume(vol)
			}
		}
	}

	if hasSelection {
		m.confirmDialog = NewConfirmationDialog(message)
		m.confirmDialog.SetSize(m.width, m.height)
		m.confirmDialog.Show()
	}
}

func (m *Model) deleteContainer(cont container.Summary) tea.Cmd {
	return func() tea.Msg {
		removeOptions := container.RemoveOptions{Force: true}
		err := m.dockerClient.ContainerRemove(m.ctx, cont.ID, removeOptions)
		if err != nil {
			return errorMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s deleted", cont.ID[:12]))
	}
}

func (m *Model) deleteImage(img image.Summary) tea.Cmd {
	return func() tea.Msg {
		_, err := m.dockerClient.ImageRemove(m.ctx, img.ID, image.RemoveOptions{Force: true})
		if err != nil {
			return errorMsg{err}
		}
		return statusMsg(fmt.Sprintf("Image %s deleted", img.ID[:12]))
	}
}

func (m *Model) deleteNetwork(net network.Summary) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.NetworkRemove(m.ctx, net.ID)
		if err != nil {
			return errorMsg{err}
		}
		return statusMsg(fmt.Sprintf("Network %s deleted", net.Name))
	}
}

func (m *Model) deleteVolume(vol *volume.Volume) tea.Cmd {
	return func() tea.Msg {
		err := m.dockerClient.VolumeRemove(m.ctx, vol.Name, true)
		if err != nil {
			return errorMsg{err}
		}
		return statusMsg(fmt.Sprintf("Volume %s deleted", vol.Name))
	}
}

// Group operations for Docker Compose projects
func (m *Model) stopGroup(group ContainerGroup) tea.Cmd {
	return func() tea.Msg {
		projectName := group.Name
		var services []string
		for _, container := range group.Containers {
			if serviceName, exists := container.Labels["com.docker.compose.service"]; exists {
				services = append(services, serviceName)
			}
		}

		err := stopComposeStack(m.dockerCli, projectName, services, nil, 30)
		if err != nil {
			return errorMsg{err}
		}

		return statusMsg(fmt.Sprintf("Compose stack '%s' stopped", projectName))
	}
}

func (m *Model) startGroup(group ContainerGroup) tea.Cmd {
	return func() tea.Msg {
		projectName := group.Name

		var services []string
		for _, container := range group.Containers {
			if serviceName, exists := container.Labels["com.docker.compose.service"]; exists {
				services = append(services, serviceName)
			}
		}

		err := startComposeStack(m.dockerCli, projectName, services, nil)
		if err != nil {
			return errorMsg{err}
		}

		return statusMsg(fmt.Sprintf("Compose stack '%s' started", projectName))
	}
}

func (m *Model) deleteGroup(group ContainerGroup) tea.Cmd {
	return func() tea.Msg {
		projectName := group.Name

		// Delete the entire compose stack (containers, networks, volumes, images)
		err := deleteComposeStack(m.dockerCli, projectName, nil)
		if err != nil {
			return errorMsg{err}
		}

		return statusMsg(fmt.Sprintf("Compose stack '%s' deleted", projectName))
	}
}

type errorMsg struct{ error }
type statusMsg string

type dataRefreshedMsg struct {
	containers []container.Summary
	images     []image.Summary
	networks   []network.Summary
	volumes    []*volume.Volume
}

func (m *Model) handleDataRefresh(msg dataRefreshedMsg) {
	m.containers = msg.containers
	m.images = msg.images
	m.networks = msg.networks
	m.volumes = msg.volumes

	m.containerTable.Update()
	m.imageTable.Update()
	m.networkTable.Update()
	m.volumeTable.Update()

	m.status = fmt.Sprintf("Last updated: %s", time.Now().Format("15:04:05"))
	m.err = nil
}

func (m *Model) showStopConfirmation() {
	if container := m.containerTable.GetSelectedContainer(); container != nil {
		if strings.Contains(strings.ToLower(container.Status), "exited") ||
			strings.Contains(strings.ToLower(container.Status), "created") {
			m.status = fmt.Sprintf("Container %s is already stopped", container.ID[:12])
			return
		}

		message := fmt.Sprintf("Are you sure you want to stop container '%s'?", container.ID[:12])
		m.confirmDialog = NewConfirmationDialog(message)
		m.confirmDialog.SetSize(m.width, m.height)
		m.confirmDialog.Show()
		m.pendingAction = func() tea.Cmd {
			return m.stopContainer(*container)
		}
	}
}

func (m *Model) stopContainer(cont container.Summary) tea.Cmd {
	return func() tea.Msg {
		timeout := 30
		stopOptions := container.StopOptions{
			Timeout: &timeout,
		}

		err := m.dockerClient.ContainerStop(m.ctx, cont.ID, stopOptions)
		if err != nil {
			return errorMsg{err}
		}
		return statusMsg(fmt.Sprintf("Container %s stopped", cont.ID[:12]))
	}
}
