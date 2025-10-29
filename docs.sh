#!/bin/bash

# Documentation helper script for PipeOps Go SDK

set -e

DOCS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DOCS_DIR"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_usage() {
    cat << EOF
Usage: $0 [command]

Commands:
    serve       Start local documentation server
    build       Build static documentation
    install     Install documentation dependencies
    deploy      Deploy documentation to GitHub Pages
    clean       Clean built documentation
    help        Show this help message

Examples:
    $0 serve              # Start local server at http://127.0.0.1:8000
    $0 build              # Build static HTML to site/ directory
    $0 install            # Install MkDocs and dependencies

EOF
}

check_dependencies() {
    if ! command -v python3 &> /dev/null; then
        echo -e "${RED}Error: python3 is required but not installed${NC}"
        exit 1
    fi

    if ! command -v pip3 &> /dev/null; then
        echo -e "${RED}Error: pip3 is required but not installed${NC}"
        exit 1
    fi
}

install_deps() {
    echo -e "${GREEN}Installing documentation dependencies...${NC}"
    pip3 install -r docs/requirements.txt
    echo -e "${GREEN}Dependencies installed successfully!${NC}"
}

serve_docs() {
    echo -e "${GREEN}Starting documentation server...${NC}"
    echo -e "${YELLOW}Documentation will be available at: http://127.0.0.1:8000${NC}"
    echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"
    mkdocs serve
}

build_docs() {
    echo -e "${GREEN}Building documentation...${NC}"
    mkdocs build
    echo -e "${GREEN}Documentation built successfully!${NC}"
    echo -e "${YELLOW}View the documentation by opening: site/index.html${NC}"
}

deploy_docs() {
    echo -e "${GREEN}Deploying documentation to GitHub Pages...${NC}"
    mkdocs gh-deploy
    echo -e "${GREEN}Documentation deployed successfully!${NC}"
}

clean_docs() {
    echo -e "${GREEN}Cleaning built documentation...${NC}"
    rm -rf site/
    echo -e "${GREEN}Documentation cleaned!${NC}"
}

# Main script
check_dependencies

case "${1:-}" in
    serve)
        serve_docs
        ;;
    build)
        build_docs
        ;;
    install)
        install_deps
        ;;
    deploy)
        deploy_docs
        ;;
    clean)
        clean_docs
        ;;
    help|--help|-h)
        print_usage
        ;;
    *)
        if [ -z "$1" ]; then
            print_usage
        else
            echo -e "${RED}Error: Unknown command '$1'${NC}"
            echo ""
            print_usage
            exit 1
        fi
        ;;
esac
