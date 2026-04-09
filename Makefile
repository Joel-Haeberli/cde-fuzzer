# Makefile for cde-extractor

# Build targets
.PHONY: all clean build build-linux build-mac build-windows test

# Default target
all: build

# Build for all platforms
build: build-linux build-mac build-windows

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/cde-extractor-linux ./cmd/cli/
	GOOS=linux GOARCH=amd64 go build -o bin/cde-extractor-server-linux ./cmd/server/
	GOOS=linux GOARCH=amd64 go build -o bin/generate-report ./cmd/generate_report/
	GOOS=linux GOARCH=amd64 go build -o bin/generate-diverse-report ./cmd/generate_diverse_report/
	GOOS=linux GOARCH=amd64 go build -o bin/derive-rules ./cmd/derive_rules/
	GOOS=linux GOARCH=amd64 go build -o bin/generate-synthetic ./cmd/generate_synthetic/
	@echo "Built Linux binaries"

# Build for macOS
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/cde-extractor-mac ./cmd/cli/
	GOOS=darwin GOARCH=amd64 go build -o bin/cde-extractor-server-mac ./cmd/server/
	GOOS=darwin GOARCH=amd64 go build -o bin/generate-report-mac ./cmd/generate_report/
	GOOS=darwin GOARCH=amd64 go build -o bin/generate-diverse-report-mac ./cmd/generate_diverse_report/
	GOOS=darwin GOARCH=amd64 go build -o bin/derive-rules-mac ./cmd/derive_rules/
	GOOS=darwin GOARCH=amd64 go build -o bin/generate-synthetic-mac ./cmd/generate_synthetic/
	@echo "Built macOS binaries"

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/cde-extractor-windows.exe ./cmd/cli/
	GOOS=windows GOARCH=amd64 go build -o bin/cde-extractor-server-windows.exe ./cmd/server/
	GOOS=windows GOARCH=amd64 go build -o bin/generate-report-windows.exe ./cmd/generate_report/
	GOOS=windows GOARCH=amd64 go build -o bin/generate-diverse-report-windows.exe ./cmd/generate_diverse_report/
	GOOS=windows GOARCH=amd64 go build -o bin/derive-rules-windows.exe ./cmd/derive_rules/
	GOOS=windows GOARCH=amd64 go build -o bin/generate-synthetic-windows.exe ./cmd/generate_synthetic/
	@echo "Built Windows binaries"

# Clean build artifacts
clean:
	rm -rf bin/
	@echo "Cleaned build artifacts"

# Test the application
test:
	go test ./...
	@echo "Tests completed"

# Derive rules from data
derive-rules:
	./bin/derive-rules -data data/ -output derived_rules/
	@echo "Rules derived"

# Generate synthetic reports
generate-synthetic:
	./bin/generate-synthetic -count 5 -variability 0.8
	@echo "Synthetic reports generated"

# Full pipeline: derive rules then generate synthetic reports
full-pipeline:
	./bin/derive-rules -data data/ -output derived_rules/ -recursive
	./bin/generate-synthetic -count 100 -variability 1.0
	@echo "Full pipeline completed"

release:
	@echo "Enter version number (e.g., v1.0.0):"
	@read -p "Version: " VERSION && \
	git tag $$VERSION && \
	git push origin $$VERSION && \
	echo "Released version $$VERSION"
