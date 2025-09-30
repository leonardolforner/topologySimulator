# Go parameters
APP_NAME = sim2
BIN_DIR = ./bin
MAIN_PKG = ./cmd/sim

# Default target
.PHONY: all
all: build

# Build the binary into ./bin/
.PHONY: build
build:
	@echo ">> Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PKG)

# Run directly (build + execute)
.PHONY: run
run: build
	@$(BIN_DIR)/$(APP_NAME) -config model.yaml -max-rands 100000

# Clean built binaries
.PHONY: clean
clean:
	@echo ">> Cleaning..."
	@rm -rf $(BIN_DIR)
