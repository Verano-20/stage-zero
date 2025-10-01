#!/bin/bash

# run-e2e-tests.sh - Script to run E2E tests with proper setup and cleanup

set -e  # Exit on any error

echo "🚀 Starting E2E Test Suite for Stage Zero API"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to cleanup
cleanup() {
    print_status "Cleaning up test environment..."
    docker-compose -f docker-compose.test.yml down -v --remove-orphans 2>/dev/null || true
    print_success "Cleanup completed"
}

# Set trap to cleanup on script exit
trap cleanup EXIT

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js and try again."
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    print_error "npm is not installed. Please install npm and try again."
    exit 1
fi

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    print_status "Installing Node.js dependencies..."
    npm install
    print_success "Dependencies installed"
fi

# Install Playwright browsers if needed
print_status "Installing Playwright browsers..."
npx playwright install --with-deps
print_success "Playwright browsers installed"

# Clean up any existing test containers
print_status "Cleaning up any existing test containers..."
docker-compose -f docker-compose.test.yml down -v --remove-orphans 2>/dev/null || true

# Start test environment
print_status "Starting test environment..."
docker-compose -f docker-compose.test.yml up -d --build

# Wait for services to be ready
print_status "Waiting for services to be ready..."
timeout=120
counter=0

while [ $counter -lt $timeout ]; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        print_success "API is ready!"
        break
    fi
    
    if [ $counter -eq 0 ]; then
        echo -n "Waiting for API"
    fi
    echo -n "."
    sleep 2
    counter=$((counter + 2))
done

echo ""

if [ $counter -ge $timeout ]; then
    print_error "API failed to start within ${timeout} seconds"
    print_status "Checking container logs..."
    docker-compose -f docker-compose.test.yml logs app-test
    exit 1
fi

# Run the tests
print_status "Running E2E tests..."
echo ""

# Parse command line arguments
TEST_ARGS=""
HEADED=false
DEBUG=false
UI=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --headed)
            HEADED=true
            shift
            ;;
        --debug)
            DEBUG=true
            shift
            ;;
        --ui)
            UI=true
            shift
            ;;
        --grep)
            TEST_ARGS="$TEST_ARGS --grep $2"
            shift 2
            ;;
        --project)
            TEST_ARGS="$TEST_ARGS --project $2"
            shift 2
            ;;
        *)
            TEST_ARGS="$TEST_ARGS $1"
            shift
            ;;
    esac
done

# Run tests based on flags
if [ "$UI" = true ]; then
    print_status "Running tests in UI mode..."
    npx playwright test --ui $TEST_ARGS
elif [ "$DEBUG" = true ]; then
    print_status "Running tests in debug mode..."
    npx playwright test --debug $TEST_ARGS
elif [ "$HEADED" = true ]; then
    print_status "Running tests in headed mode..."
    npx playwright test --headed $TEST_ARGS
else
    print_status "Running tests in headless mode..."
    npx playwright test $TEST_ARGS
fi

TEST_EXIT_CODE=$?

# Show test results
echo ""
if [ $TEST_EXIT_CODE -eq 0 ]; then
    print_success "All E2E tests passed! 🎉"
else
    print_error "Some E2E tests failed. Check the output above for details."
fi

# Show test report
if [ -f "playwright-report/index.html" ]; then
    print_status "Test report available at: playwright-report/index.html"
    print_status "To view the report, run: npx playwright show-report"
fi

echo ""
print_status "E2E test run completed"
exit $TEST_EXIT_CODE
