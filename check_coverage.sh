#!/bin/bash

# Run tests and generate coverage profile
go test -coverprofile=coverage.out ./...

# Check if the tests passed
if [ $? -ne 0 ]; then
  echo "Tests failed. Coverage report not generated."
  exit 1
fi

# Display coverage report in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

echo "HTML coverage report generated: coverage.html"
