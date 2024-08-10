# Variables
BINARY_NAME=app
FRAMES_DIR=frames

# Build the Go binary
build:
	@echo "Building the binary..."
	go build -o $(BINARY_NAME) .

# Run the application
run: build
	@echo "Running the application..."
	mkdir -p $(FRAMES_DIR)
	rm -f $(FRAMES_DIR)/*
	./$(BINARY_NAME)

# Format source code
fmt:
	@echo "Formatting the code"
	go fmt ./...

# Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(FRAMES_DIR)

.PHONY: build fmt run clean

