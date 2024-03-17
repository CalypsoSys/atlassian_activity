# Makefile for a simple Go project

# The binary to build (just the basename).
BIN := atlassian_collector

# This will output the binary in the current directory under 'bin'
BIN_DIR := ./bin

# Where to find the source code. Assuming it's in the current directory
SRC_DIR := ./cmd

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GORUN := $(GOCMD) run

# Build the project
build:
	$(GOBUILD) -o $(BIN_DIR)/$(BIN) $(SRC_DIR)/atlassian_collector/

# Clean up the project
clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/$(BIN)

# Run the project
run:
	$(GORUN) $(SRC_DIR)/main.go

# Make directory if it doesn't exist
$(shell mkdir -p $(BIN_DIR))

.PHONY: build clean run
