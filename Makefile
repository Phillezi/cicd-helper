# Variables
BINARY_NAME=cicd-helper
BUILD_DIR=bin
CMD_DIR=cmd/cicd-helper
MAIN_FILE=$(CMD_DIR)/main.go

# Targets
.PHONY: all clean build run

all: build

build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@CGOENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."
