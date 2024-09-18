# Define variables
BUILD_DIR=cmd/plant-report-lambda/
LAMBDA_DIR=infrastructure/bin
OUTPUT_DIR=../../dist/plant-report-lambda
DIST_DIR := dist
ZIP_FILE := $(DIST_DIR)/plant-report-lambda.zip
BINARY_NAME := plant-report-lambda

# Default target
all: build deploy

# Build the Go application
build:
	mkdir -p $(OUTPUT_DIR)
	cd $(BUILD_DIR) && GOOS=linux GOARCH=amd64 go build -o ../../dist/bootstrap
	# Package the binary into a zip file
	cd $(DIST_DIR) && zip -r9 plant-report-lambda.zip bootstrap

# Deploy using CDK
deploy:
	cd infrastructure && npx cdk deploy

# Clean up build artifacts
clean:
	rm -rf $(OUTPUT_DIR)

# Run all targets
run: build deploy
