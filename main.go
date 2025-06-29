package main

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ngnhng/docker-status/tui"
)

func main() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "status [OPTIONS]",
			Short: "Docker container and image management TUI",
			Long: `A Docker CLI plugin for managing Docker containers and images in a terminal user interface.
Provides an interactive way to view and manage your Docker resources.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				return runPlugin(dockerCli, args)
			},
		}
		return cmd
	},
		manager.Metadata{
			SchemaVersion:    "0.1.0",
			Vendor:           "https://github.com/ngnhng",
			Version:          "v0.1.0",
			ShortDescription: "A Docker CLI plugin for managing Docker containers and images in a TUI",
			URL:              "https://github.com/ngnhng/docker-status",
		})
}

func runPlugin(dockerCli command.Cli, args []string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error creating Docker client: %w", err)
	}

	_, err = cli.Info(ctx)
	if err != nil {
		return fmt.Errorf("error connecting to Docker daemon: %w", err)
	}

	model := tui.NewModel(cli, dockerCli)
	program := tea.NewProgram(model, tea.WithAltScreen())

	_, err = program.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}
