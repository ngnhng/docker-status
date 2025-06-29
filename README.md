# Docker Status TUI

A K9s-inspired terminal user interface (TUI) for Docker management.

## Installation

### As a Docker CLI Plugin

1. Build the plugin:

   ```bash
   go build -o docker-status
   ```

2. Install as a Docker CLI plugin:

   ```bash
   mkdir -p ~/.docker/cli-plugins
   cp docker-status ~/.docker/cli-plugins/
   chmod +x ~/.docker/cli-plugins/docker-status
   ```

3. Verify installation:
   ```bash
   docker status --help
   ```

### Standalone Usage

You can also run it directly:

```bash
./docker-status
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Acknowledgments

- Inspired by [K9s](https://k9scli.io/) for Kubernetes
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- Uses [Docker Engine API](https://docs.docker.com/engine/api/) for Docker operations
