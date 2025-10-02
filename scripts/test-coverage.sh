#!/bin/bash
# Test Coverage Helper Script for Sloth Runner

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Command functions
run_all_tests() {
    print_header "Running All Tests"
    go test ./...
    if [ $? -eq 0 ]; then
        print_success "All tests passed!"
    else
        print_error "Some tests failed!"
        exit 1
    fi
}

run_tests_with_coverage() {
    print_header "Running Tests with Coverage"
    go test ./... -coverprofile=coverage.out -coverpkg=./...
    
    if [ $? -eq 0 ]; then
        print_success "Tests completed!"
        
        # Calculate total coverage
        total_coverage=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}')
        print_info "Total Coverage: $total_coverage"
        
        # Check if coverage meets threshold
        coverage_value=$(echo $total_coverage | sed 's/%//')
        if (( $(echo "$coverage_value >= 80" | bc -l) )); then
            print_success "Coverage threshold met! (>= 80%)"
        else
            print_warning "Coverage below 80% threshold"
        fi
    else
        print_error "Tests failed!"
        exit 1
    fi
}

show_coverage_report() {
    print_header "Detailed Coverage Report"
    
    if [ ! -f "coverage.out" ]; then
        print_error "No coverage file found. Run tests with coverage first."
        exit 1
    fi
    
    go tool cover -func=coverage.out | grep -v "100.0%" | sort -k3 -n
}

show_coverage_by_package() {
    print_header "Coverage by Package"
    
    if [ ! -f "coverage.out" ]; then
        print_error "No coverage file found. Run tests with coverage first."
        exit 1
    fi
    
    go tool cover -func=coverage.out | \
        awk '{print $1}' | \
        grep -v "^total:" | \
        sed 's/:.*$//' | \
        sort -u | \
        while read package; do
            coverage=$(go tool cover -func=coverage.out | grep "$package" | \
                awk '{sum+=$3; count++} END {if(count>0) printf "%.1f", sum/count; else print "0"}')
            printf "%-60s %6s%%\n" "$package" "$coverage"
        done | sort -k2 -n
}

open_html_report() {
    print_header "Generating HTML Coverage Report"
    
    if [ ! -f "coverage.out" ]; then
        print_error "No coverage file found. Run tests with coverage first."
        exit 1
    fi
    
    go tool cover -html=coverage.out -o coverage.html
    print_success "HTML report generated: coverage.html"
    
    # Try to open in browser
    if command -v xdg-open > /dev/null; then
        xdg-open coverage.html
    elif command -v open > /dev/null; then
        open coverage.html
    else
        print_info "Open coverage.html in your browser manually"
    fi
}

run_specific_package() {
    if [ -z "$1" ]; then
        print_error "Please specify a package"
        echo "Usage: $0 package <package-path>"
        echo "Example: $0 package ./internal/core"
        exit 1
    fi
    
    print_header "Testing Package: $1"
    go test -v -cover "$1"
}

run_with_race_detector() {
    print_header "Running Tests with Race Detector"
    go test -race ./...
    if [ $? -eq 0 ]; then
        print_success "No race conditions detected!"
    else
        print_error "Race conditions found!"
        exit 1
    fi
}

run_benchmarks() {
    print_header "Running Benchmarks"
    go test -bench=. -benchmem ./... | grep -E "(Benchmark|PASS|FAIL)"
}

check_coverage_threshold() {
    threshold=${1:-80}
    
    if [ ! -f "coverage.out" ]; then
        run_tests_with_coverage
    fi
    
    total_coverage=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
    
    print_info "Current coverage: ${total_coverage}%"
    print_info "Required threshold: ${threshold}%"
    
    if (( $(echo "$total_coverage >= $threshold" | bc -l) )); then
        print_success "Coverage threshold met!"
        exit 0
    else
        print_error "Coverage below threshold!"
        print_warning "Missing: $(echo "$threshold - $total_coverage" | bc)% points"
        exit 1
    fi
}

find_untested_files() {
    print_header "Finding Files Without Tests"
    
    find ./internal ./cmd -name "*.go" -not -name "*_test.go" -not -name "*.pb.go" 2>/dev/null | \
        while read file; do
            test_file=$(echo $file | sed 's/.go$/_test.go/')
            if [ ! -f "$test_file" ]; then
                echo "$file"
            fi
        done
}

generate_coverage_badge() {
    if [ ! -f "coverage.out" ]; then
        run_tests_with_coverage
    fi
    
    coverage=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
    
    # Determine color based on coverage
    if (( $(echo "$coverage >= 80" | bc -l) )); then
        color="brightgreen"
    elif (( $(echo "$coverage >= 60" | bc -l) )); then
        color="yellow"
    elif (( $(echo "$coverage >= 40" | bc -l) )); then
        color="orange"
    else
        color="red"
    fi
    
    badge_url="https://img.shields.io/badge/coverage-${coverage}%25-${color}"
    print_success "Coverage Badge URL:"
    echo "$badge_url"
}

show_help() {
    cat << EOF
${BLUE}Sloth Runner - Test Coverage Helper${NC}

${GREEN}Usage:${NC}
    $0 [command] [options]

${GREEN}Commands:${NC}
    test                Run all tests
    coverage            Run tests with coverage report
    report              Show detailed coverage report
    package             Show coverage by package
    html                Generate and open HTML coverage report
    race                Run tests with race detector
    bench               Run benchmarks
    threshold [num]     Check if coverage meets threshold (default: 80)
    untested            Find files without tests
    badge               Generate coverage badge URL
    specific <pkg>      Test specific package
    help                Show this help message

${GREEN}Examples:${NC}
    $0 test                     # Run all tests
    $0 coverage                 # Run tests with coverage
    $0 specific ./internal/core # Test specific package
    $0 threshold 70             # Check if coverage >= 70%
    $0 html                     # Open HTML coverage report

${BLUE}Shortcuts:${NC}
    ./scripts/test-coverage.sh              # Run tests with coverage (default)
    ./scripts/test-coverage.sh -t           # Run all tests
    ./scripts/test-coverage.sh -r           # Show coverage report
    ./scripts/test-coverage.sh -h           # Open HTML report

EOF
}

# Main script
case "${1:-coverage}" in
    test|-t)
        run_all_tests
        ;;
    coverage|-c)
        run_tests_with_coverage
        ;;
    report|-r)
        show_coverage_report
        ;;
    package|-p)
        show_coverage_by_package
        ;;
    html|-h)
        open_html_report
        ;;
    race)
        run_with_race_detector
        ;;
    bench|-b)
        run_benchmarks
        ;;
    threshold)
        check_coverage_threshold "$2"
        ;;
    untested|-u)
        find_untested_files
        ;;
    badge)
        generate_coverage_badge
        ;;
    specific)
        run_specific_package "$2"
        ;;
    help|--help|-?)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
