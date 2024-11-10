#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🚀 Running local CI tests..."

# Function to check step status
check_status() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $1 passed${NC}"
    else
        echo -e "${RED}✗ $1 failed${NC}"
        exit 1
    fi
}

# 1. Check Go installation
echo "📝 Checking Go installation..."
go version
check_status "Go version check"

# 2. Verify dependencies
echo "📝 Verifying dependencies..."
go mod verify
check_status "Dependency verification"

# 3. Check formatting
echo "📝 Checking code formatting..."
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo -e "${RED}✗ Code formatting check failed${NC}"
    gofmt -s -l .
    exit 1
fi
check_status "Code formatting"

# 4. Run go vet
echo "📝 Running go vet..."
go vet ./...
check_status "Go vet"

# 5. Run unit tests
echo "📝 Running unit tests..."
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
check_status "Unit tests"

# 6. Build
echo "📝 Building project..."
go build -v ./...
check_status "Build"

# 7. Show coverage (if tests passed)
if [ -f coverage.txt ]; then
    echo "📝 Coverage report:"
    go tool cover -func=coverage.txt
fi

# 8. Check if we should run integration tests
if [ "$1" == "--with-integration" ]; then
    echo "📝 Running integration tests..."
    
    # Check if kind is installed
    if ! command -v kind &> /dev/null; then
        echo -e "${RED}✗ kind is not installed. Please install it first.${NC}"
        exit 1
    fi
    
    # Check if kubectl is installed
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}✗ kubectl is not installed. Please install it first.${NC}"
        exit 1
    fi
    
    # Create test cluster if it doesn't exist
    if ! kind get clusters | grep -q "microkube-test"; then
        kind create cluster --name microkube-test
    fi
    
    # Run integration tests
    echo "Running ingress tests..."
    ./scripts/ingress-test.sh
    check_status "Ingress tests"
    
    echo "Running canary tests..."
    ./scripts/canary-test.sh
    check_status "Canary tests"
    
    echo "Running traffic tests..."
    ./scripts/traffic-test.sh
    check_status "Traffic tests"
    
    echo "Running load tests..."
    ./scripts/advanced-load-test.sh
    check_status "Load tests"
    
    # Optional: Delete test cluster
    if [ "$2" == "--cleanup" ]; then
        kind delete cluster --name microkube-test
    fi
fi

echo -e "${GREEN}✓ All CI checks completed successfully${NC}"