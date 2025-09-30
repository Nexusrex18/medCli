#!/bin/bash

set -e

# Configuration
REPO_OWNER="Nexusrex18"
REPO_NAME="medCli"
BINARY_NAME="medCli"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and Architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            BINARY_NAME="${BINARY_NAME}.exe"
            ;;
        *)
            echo -e "${RED}Unsupported operating system: $OS${NC}"
            exit 1
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        i386|i686)
            ARCH="386"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${GREEN}Detected platform: ${OS}/${ARCH}${NC}"
}

# Check if required commands exist
check_dependencies() {
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        echo -e "${RED}Error: curl or wget is required but not installed.${NC}"
        exit 1
    fi
    
    if ! command -v tar >/dev/null 2>&1 && [ "$OS" != "windows" ]; then
        echo -e "${RED}Error: tar is required but not installed.${NC}"
        exit 1
    fi
}

# Get latest release version
get_latest_version() {
    echo -e "${YELLOW}Fetching latest release...${NC}"
    
    if command -v curl >/dev/null 2>&1; then
        LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        LATEST_VERSION=$(wget -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    
    if [ -z "$LATEST_VERSION" ]; then
        echo -e "${RED}Failed to fetch latest version${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Latest version: ${LATEST_VERSION}${NC}"
}

# Download binary
download_binary() {
    DOWNLOAD_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${LATEST_VERSION}/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
    
    echo -e "${YELLOW}Downloading from: ${DOWNLOAD_URL}${NC}"
    
    TMP_DIR=$(mktemp -d)
    TMP_FILE="${TMP_DIR}/${BINARY_NAME}.tar.gz"
    
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"
    else
        wget -O "$TMP_FILE" "$DOWNLOAD_URL"
    fi
    
    if [ ! -f "$TMP_FILE" ]; then
        echo -e "${RED}Download failed${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Download complete${NC}"
}

# Install binary
install_binary() {
    echo -e "${YELLOW}Installing ${BINARY_NAME}...${NC}"
    
    cd "$TMP_DIR"
    tar -xzf "${BINARY_NAME}.tar.gz"
    
    if [ "$OS" = "windows" ]; then
        # For Windows, install to user's local bin or current directory
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
        mv "$BINARY_NAME" "$INSTALL_DIR/"
        echo -e "${YELLOW}Add ${INSTALL_DIR} to your PATH if not already added${NC}"
    else
        # For Unix-like systems
        if [ -w "$INSTALL_DIR" ]; then
            mv "$BINARY_NAME" "$INSTALL_DIR/"
        else
            echo -e "${YELLOW}Requesting sudo access to install to ${INSTALL_DIR}${NC}"
            sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        fi
        
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    # Cleanup
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}${BINARY_NAME} installed successfully!${NC}"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        VERSION=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        echo -e "${GREEN}✓ Installation verified${NC}"
        echo -e "${GREEN}Run '${BINARY_NAME}' to get started${NC}"
    else
        echo -e "${YELLOW}Warning: ${BINARY_NAME} not found in PATH${NC}"
        echo -e "${YELLOW}You may need to add ${INSTALL_DIR} to your PATH${NC}"
    fi
}

# Build from source (fallback)
build_from_source() {
    echo -e "${YELLOW}Building from source...${NC}"
    
    if ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}Go is not installed. Please install Go first.${NC}"
        exit 1
    fi
    
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    echo -e "${YELLOW}Cloning repository...${NC}"
    git clone "https://github.com/${REPO_OWNER}/${REPO_NAME}.git"
    cd "$REPO_NAME"
    
    echo -e "${YELLOW}Building binary...${NC}"
    go build -o "$BINARY_NAME" ./cmd/main.go
    
    if [ "$OS" = "windows" ]; then
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    else
        if [ -w "$INSTALL_DIR" ]; then
            mv "$BINARY_NAME" "$INSTALL_DIR/"
        else
            sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        fi
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    cd /
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}Build and installation complete!${NC}"
}

# Main installation flow
main() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  ${BINARY_NAME} Installation Script${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    
    detect_platform
    check_dependencies
    
    # Try to download pre-built binary, fallback to building from source
    if get_latest_version && download_binary && install_binary; then
        verify_installation
    else
        echo -e "${YELLOW}Failed to download pre-built binary${NC}"
        echo -e "${YELLOW}Attempting to build from source...${NC}"
        build_from_source
        verify_installation
    fi
    
    echo ""
    echo -e "${GREEN}Installation complete! 🎉${NC}"
}

# Run main function
main