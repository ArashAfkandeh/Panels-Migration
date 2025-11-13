#!/bin/bash

# Panels Migration Installation Script
# This script downloads and extracts the Panels Migration tool

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Error: This script must be run as root (use sudo).${NC}"
    exit 1
fi

# Download URL
DOWNLOAD_URL="https://github.com/ArashAfkandeh/Panels-Migration/releases/download/Panels_Migration/Panels_Migration_v0.0.1.tar.gz"

# Installation directory
INSTALL_DIR="/root/Panels_Migration"
TEMP_DIR=$(mktemp -d)

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Panels Migration Installer${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Update package manager
echo -e "${YELLOW}Updating package manager...${NC}"
apt-get update -y

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to update package manager.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Package manager updated${NC}"
echo ""

# Install prerequisites
echo -e "${YELLOW}Installing prerequisites...${NC}"
PACKAGES="curl wget tar gzip"

for package in $PACKAGES; do
    if ! command -v $package &> /dev/null; then
        echo -e "${YELLOW}Installing $package...${NC}"
        apt-get install -y $package
        if [ $? -ne 0 ]; then
            echo -e "${RED}Error: Failed to install $package.${NC}"
            exit 1
        fi
    else
        echo -e "${GREEN}✓ $package is already installed${NC}"
    fi
done

echo -e "${GREEN}✓ All prerequisites installed${NC}"
echo ""

# Check if curl or wget is available
if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
    echo -e "${RED}Error: curl or wget is required but not installed.${NC}"
    exit 1
fi

# Download the file
echo -e "${YELLOW}Downloading Panels Migration...${NC}"
if command -v curl &> /dev/null; then
    curl -L -o "/root/Panels_Migration.tar.gz" "${DOWNLOAD_URL}"
else
    wget -O "/root/Panels_Migration.tar.gz" "${DOWNLOAD_URL}"
fi

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to download the file.${NC}"
    rm -rf "${TEMP_DIR}"
    exit 1
fi

echo -e "${GREEN}✓ Download completed${NC}"
echo ""

# Check if tar is available
if ! command -v tar &> /dev/null; then
    echo -e "${RED}Error: tar is required but not installed.${NC}"
    rm -rf "${TEMP_DIR}"
    exit 1
fi

# Extract the file in /root directory
echo -e "${YELLOW}Extracting Panels Migration in /root...${NC}"
tar -xzf "/root/Panels_Migration.tar.gz" -C "/root"

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to extract the file.${NC}"
    rm -rf "${TEMP_DIR}"
    exit 1
fi

echo -e "${GREEN}✓ Extraction completed${NC}"
echo ""

# Make the binary executable
echo -e "${YELLOW}Setting executable permissions...${NC}"
if [ -f "${INSTALL_DIR}" ]; then
    chmod +x "${INSTALL_DIR}"
    echo -e "${GREEN}✓ Executable permissions set${NC}"
else
    echo -e "${RED}Error: Panels_Migration binary not found at ${INSTALL_DIR}${NC}"
    rm -rf "${TEMP_DIR}"
    exit 1
fi

# Clean up temporary files and archive
rm -rf "${TEMP_DIR}"
rm -f "/root/Panels_Migration.tar.gz"

# Run the application
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Starting Panels Migration...${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

"${INSTALL_DIR}"
