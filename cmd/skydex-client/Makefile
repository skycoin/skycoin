.PHONY: install dev build clean build-go full-build

# Install dependencies
install:
	npm install

# Run development server
dev:
	npm run dev

# Build for production (outputs to commands/static, which is embedded by the Go binary)
build:
	npm run build

# Clean build artifacts (keep commands/static, it is committed and embedded)
clean:
	rm -rf node_modules

# Build the UI then the Go binary with the embedded UI.
# The client runs as a Skywire visor app (`skywire app skydex-client`),
# so it is normally launched by the visor rather than run standalone.
build-go:
	@echo "Building UI..."
	@npm run build
	@echo "Building skydex-client binary..."
	@cd ../../.. && go build -o bin/skydex-client ./cmd/apps/skydex-client

# Full build (install + build)
full-build: install build
