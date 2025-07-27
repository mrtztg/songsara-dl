.PHONY: build clean test install release snapshot

# Build the binary
build:
	go build -o songsara-dl

# Clean build artifacts
clean:
	rm -f songsara-dl
	rm -rf downloads/
	rm -rf dist/

# Install dependencies
deps:
	go mod tidy

# Run tests (if any)
test:
	go test ./...

# Install the binary to /usr/local/bin (requires sudo)
install: build
	sudo cp songsara-dl /usr/local/bin/

# Development: build and run with example URL
dev: build
	./songsara-dl --help

# Show current version and build info
version:
	@echo "SongSara Downloader"
	@echo "Go version: $(shell go version)"
	@echo "Build time: $(shell date)"

# Build release binaries for all platforms
release:
	goreleaser release --clean

# Build snapshot binaries for all platforms
snapshot:
	goreleaser build --snapshot --clean

# Build and test locally
local-build: snapshot
	@echo "Built binaries:"
	@ls -la dist/ 