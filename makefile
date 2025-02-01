# Makefile for building and running a Go project.
help:
	@echo "Available make commands:"
	@cat makefile | grep '^[a-z][^:]*:' | grep -v 'install_' | cut -d: -f1 | sort | sed 's/^/  /'

# Set the output binary name.
BINARY_NAME=count_locs

# Set the source directory.
SRC_DIR=./src

# Default make command should build the project.
all: build

# Build the project.
build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME) $(SRC_DIR)

# Run the project.
run:
	@echo "Running..."
	go run $(SRC_DIR)
# Run the tests.
test:
	@echo "Testing..."
	go test -v $(SRC_DIR)

# Clean up the binary.
clean:
	@echo "Cleaning..."
	go clean
	rm -f bin/$(BINARY_NAME)

# "make build" will compile the app and produce a binary.
# "make run" will run the app using 'go run'.
# "make test" will run the tests.
# "make clean" will remove the binary and clean up the project.