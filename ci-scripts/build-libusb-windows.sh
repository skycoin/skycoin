#!/bin/bash
set -e

# Install hidapi and libusb for Windows builds using vcpkg
# This script should be run on Windows (Git Bash/MSYS2)

VCPKG_ROOT="${VCPKG_ROOT:-C:/vcpkg}"

echo "Installing hidapi and libusb via vcpkg..."

# Install vcpkg if not already installed
if [ ! -d "$VCPKG_ROOT" ]; then
    echo "Installing vcpkg to $VCPKG_ROOT"
    git clone https://github.com/Microsoft/vcpkg.git "$VCPKG_ROOT"
    cd "$VCPKG_ROOT"
    ./bootstrap-vcpkg.bat
else
    echo "vcpkg already installed at $VCPKG_ROOT"
    cd "$VCPKG_ROOT"
fi

# Install hidapi and libusb for x64, x86, and arm64
echo "Installing hidapi:x64-windows-static..."
./vcpkg install hidapi:x64-windows-static

echo "Installing hidapi:x86-windows-static..."
./vcpkg install hidapi:x86-windows-static

echo "Installing hidapi:arm64-windows-static..."
./vcpkg install hidapi:arm64-windows-static || echo "Warning: arm64 may not be available"

echo "✓ hidapi installed for Windows"
echo "VCPKG_ROOT=$VCPKG_ROOT"
echo ""
echo "Set the following environment variables for builds:"
echo "For amd64:"
echo "  PKG_CONFIG_PATH=$VCPKG_ROOT/installed/x64-windows-static/lib/pkgconfig"
echo "  CGO_CFLAGS=-I$VCPKG_ROOT/installed/x64-windows-static/include"
echo "  CGO_LDFLAGS=-L$VCPKG_ROOT/installed/x64-windows-static/lib"
echo ""
echo "For 386:"
echo "  PKG_CONFIG_PATH=$VCPKG_ROOT/installed/x86-windows-static/lib/pkgconfig"
echo "  CGO_CFLAGS=-I$VCPKG_ROOT/installed/x86-windows-static/include"
echo "  CGO_LDFLAGS=-L$VCPKG_ROOT/installed/x86-windows-static/lib"
