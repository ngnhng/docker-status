package tui

import (
	"context"
	"time"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

func stopComposeStack(
	dockerCli command.Cli,
	projectName string,
	services []string,
	project *types.Project,
	timeoutSeconds int) error {
	svc := compose.NewComposeService(dockerCli)
	timeout := time.Duration(timeoutSeconds) * time.Second
	opts := api.StopOptions{
		Timeout:  &timeout,
		Services: services, // or a []string{"service1", "service2"} to stop specific services
		Project:  project,  // optional: pass a parsed *types.Project if you have one
	}
	ctx := context.Background()
	return svc.Stop(ctx, projectName, opts)
}

func startComposeStack(
	dockerCli command.Cli,
	projectName string,
	services []string,
	project *types.Project) error {
	svc := compose.NewComposeService(dockerCli)
	opts := api.StartOptions{
		Services: services,
		Project:  project,
	}
	ctx := context.Background()
	return svc.Start(ctx, projectName, opts)
}

func deleteComposeStack(
	dockerCli command.Cli,
	projectName string,
	project *types.Project) error {
	svc := compose.NewComposeService(dockerCli)
	opts := api.DownOptions{
		RemoveOrphans: true,
		Project:       project,
		Volumes:       true,  // Remove volumes
		Images:        "all", // Remove images
	}
	ctx := context.Background()
	return svc.Down(ctx, projectName, opts)
}
