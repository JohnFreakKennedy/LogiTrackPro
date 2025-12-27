#!/bin/bash
# Test runner script for Go backend

set -e

echo "Running Go backend tests..."

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

echo "Coverage report generated: coverage.html"
echo "Test coverage:"
go tool cover -func=coverage.out | tail -1
