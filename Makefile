BINARY_NAME=docker-status

GO_FILES=main.go

# The path to the Docker CLI plugins directory.
PLUGIN_DIR=$(HOME)/.docker/cli-plugins

PLUGIN_COMMAND=$(subst docker-,,$(BINARY_NAME))

LDFLAGS=-ldflags "-s -w"
BUILD_ENV=CGO_ENABLED=0

.PHONY: all build build-small install uninstall clean deps run check-docker fmt help

all: build

build:
	@echo "Building plugin binary..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) $(GO_FILES)

build-small:
	@echo "Building optimized plugin binary (smallest size)..."
	@go build $(LDFLAGS) -trimpath -o $(BINARY_NAME) $(GO_FILES)
	@echo "Binary size: $$(du -h $(BINARY_NAME) | cut -f1)"

install: build
	@echo "Installing plugin to $(PLUGIN_DIR)..."
	@mkdir -p $(PLUGIN_DIR)
	@cp $(BINARY_NAME) $(PLUGIN_DIR)/$(BINARY_NAME)
	@chmod +x $(PLUGIN_DIR)/$(BINARY_NAME)
	@echo "Installation complete."
	@echo "Run 'docker --help' to see if '$(PLUGIN_COMMAND)' appears in the list of commands."
	@echo "You can now run your plugin with: docker $(PLUGIN_COMMAND)"

uninstall:
	@echo "Uninstalling plugin..."
	@rm -f $(PLUGIN_DIR)/$(BINARY_NAME)
	@echo "Uninstallation complete."

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(BINARY_NAME)

deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

run: build
	@echo "Running docker-status TUI..."
	@./$(BINARY_NAME)

check-docker:
	@docker info > /dev/null 2>&1 || (echo "Docker is not running. Please start Docker first." && exit 1)
	@echo "Docker is running ✓"

fmt:
	@echo "Formatting code..."
	@go fmt ./...

help:
	@echo "Docker Status TUI - Available commands:"
	@echo ""
	@echo "  build        Build the binary"
	@echo "  build-small  Build with maximum size optimization"
	@echo "  install      Install as Docker CLI plugin"
	@echo "  uninstall    Remove the plugin"
	@echo "  deps         Install dependencies"
	@echo "  clean        Clean build artifacts"
	@echo "  run          Run the application directly"
	@echo "  check-docker Check if Docker is running"
	@echo "  fmt          Format code"
	@echo "  help         Show this help message"
	@echo ""
	@echo "Quick start:"
	@echo "  make deps && make install"
	@echo "  docker status"