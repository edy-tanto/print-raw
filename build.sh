#!/bin/bash

# Build script for Print Web Service
# This script builds the Go application as 32-bit executable and packages it into installer-printer.zip

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Print Web Service Build Script ===${NC}"

# Get the script directory (project root)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "Working directory: $SCRIPT_DIR"

# Baca version dari file version.txt (default jika user tidak isi)
CURRENT_VERSION=""
if [ -f "bin/version.txt" ]; then
    CURRENT_VERSION=$(grep -E "^version=" "bin/version.txt" | sed 's/^version=//' | tr -d '\r\n')
fi
# Jika tidak ada file atau kosong, default v0.0
[[ -z "$CURRENT_VERSION" ]] && CURRENT_VERSION="v0.0"

# Prompt untuk version build; kosong = pakai version yang ada
echo ""
read -p "Masukkan version hasil build [$CURRENT_VERSION] (kosong = tidak ubah): " BUILD_VERSION
if [[ -z "$BUILD_VERSION" ]]; then
    BUILD_VERSION="$CURRENT_VERSION"
else
    if [[ "$BUILD_VERSION" != v* ]]; then
        BUILD_VERSION="v${BUILD_VERSION}"
    fi
fi
echo -e "${GREEN}Version build: $BUILD_VERSION${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    exit 1
fi

echo -e "${GREEN}Go version:$(go version)${NC}"

# Create bin directory if it doesn't exist
if [ ! -d "bin" ]; then
    echo "Creating bin directory..."
    mkdir -p bin
fi

# Build the application with 32-bit architecture
echo -e "${YELLOW}Building Go application (32-bit)...${NC}"
export GOARCH=386
export GOOS=windows

# Build command - from project root, build cmd/print_web_service
go build -o bin/print_web_service.exe ./cmd/print_web_service

if [ ! -f "bin/print_web_service.exe" ]; then
    echo -e "${RED}Error: Build failed - executable not found${NC}"
    exit 1
fi

echo -e "${GREEN}Build successful!${NC}"

# Create temporary directory for packaging
TEMP_DIR=$(mktemp -d)
echo "Temporary directory: $TEMP_DIR"

# Copy all required files to temp directory
echo -e "${YELLOW}Packaging files...${NC}"

# Copy executable
cp bin/print_web_service.exe "$TEMP_DIR/"

# Copy bitmap files
if [ -f "bin/captain-order-receipt-header.bmp" ]; then
    cp bin/captain-order-receipt-header.bmp "$TEMP_DIR/"
fi
if [ -f "bin/paradis-q.bmp" ]; then
    cp bin/paradis-q.bmp "$TEMP_DIR/"
fi
if [ -f "bin/qubu-resort-icon.bmp" ]; then
    cp bin/qubu-resort-icon.bmp "$TEMP_DIR/"
fi

# Copy installation scripts
if [ -f "bin/install.bat" ]; then
    cp bin/install.bat "$TEMP_DIR/"
fi
if [ -f "bin/start.bat" ]; then
    cp bin/start.bat "$TEMP_DIR/"
fi
if [ -f "bin/stop.bat" ]; then
    cp bin/stop.bat "$TEMP_DIR/"
fi
if [ -f "bin/uninstall.bat" ]; then
    cp bin/uninstall.bat "$TEMP_DIR/"
fi

# Update version.txt dengan version build lalu copy ke temp
echo "version=$BUILD_VERSION" > bin/version.txt
cp bin/version.txt "$TEMP_DIR/"

# Create ZIP file (nama file mengandung version)
ZIP_NAME="installer-printer-${BUILD_VERSION}.zip"

# Target directory (Windows path converted for Git Bash)
TARGET_DIR="/c/Users/Dream/Downloads"

# Create target directory if it doesn't exist
if [ ! -d "$TARGET_DIR" ]; then
    echo "Creating target directory: $TARGET_DIR"
    mkdir -p "$TARGET_DIR"
fi

echo -e "${YELLOW}Creating ZIP archive...${NC}"
# Convert Unix paths to Windows format for PowerShell
# Handle /c/, /d/ style and /tmp/ paths
convert_to_win_path() {
    local unix_path="$1"
    # Try cygpath first (if available in Git Bash)
    if command -v cygpath &> /dev/null; then
        cygpath -w "$unix_path"
    else
        # Manual conversion for /c/, /d/ style paths
        if [[ "$unix_path" =~ ^/([a-z])/ ]]; then
            echo "$unix_path" | sed "s|^/\\([a-z]\\)/|\\1:/|" | sed 's|/|\\|g'
        # For /tmp/, use Windows temp
        elif [[ "$unix_path" =~ ^/tmp/ ]]; then
            local win_temp=$(powershell.exe -Command "[System.IO.Path]::GetTempPath()" | tr -d '\r\n')
            local rel_path=$(echo "$unix_path" | sed 's|^/tmp/||')
            echo "$win_temp$rel_path" | sed 's|/|\\|g'
        else
            echo "$unix_path" | sed 's|/|\\|g'
        fi
    fi
}

TEMP_DIR_WIN=$(convert_to_win_path "$TEMP_DIR")
TARGET_DIR_WIN=$(convert_to_win_path "$TARGET_DIR")
TARGET_PATH_WIN="$TARGET_DIR_WIN\\$ZIP_NAME"

# Use PowerShell Compress-Archive to create ZIP file
# Escape backslashes properly for PowerShell
TEMP_DIR_PS=$(echo "$TEMP_DIR_WIN" | sed 's|\\|\\\\|g')
TARGET_PATH_PS=$(echo "$TARGET_PATH_WIN" | sed 's|\\|\\\\|g')
powershell.exe -Command "Compress-Archive -Path '$TEMP_DIR_PS\*' -DestinationPath '$TARGET_PATH_PS' -Force"

TARGET_PATH="$TARGET_DIR/$ZIP_NAME"

echo -e "${GREEN}ZIP file created successfully!${NC}"
echo -e "${GREEN}Location: $TARGET_PATH${NC}"

# Clean up temporary directory
rm -rf "$TEMP_DIR"
echo "Temporary files cleaned up."

echo -e "${GREEN}=== Build Complete ===${NC}"
