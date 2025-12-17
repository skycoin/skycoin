#!/bin/bash
set -e

# Build libusb-1.0 static library for a musl target
# Usage: ./build-libusb-musl.sh <arch> <toolchain-prefix> <sysroot>
# Example: ./build-libusb-musl.sh amd64 x86_64-linux-musl ./musl-data/x86_64-linux-musl-cross

ARCH=$1
TOOLCHAIN_PREFIX=$2
SYSROOT=$3

if [ -z "$ARCH" ] || [ -z "$TOOLCHAIN_PREFIX" ] || [ -z "$SYSROOT" ]; then
    echo "Usage: $0 <arch> <toolchain-prefix> <sysroot>"
    echo "Example: $0 amd64 x86_64-linux-musl ./musl-data/x86_64-linux-musl-cross"
    exit 1
fi

LIBUSB_VERSION="v1.0.27"
BUILD_DIR="/tmp/libusb-build-${ARCH}"

echo "Building libusb ${LIBUSB_VERSION} for ${ARCH} (${TOOLCHAIN_PREFIX})"

# Clean previous build
rm -rf "${BUILD_DIR}"
mkdir -p "${BUILD_DIR}"
cd "${BUILD_DIR}"

# Download libusb
wget -q "https://github.com/libusb/libusb/archive/refs/tags/${LIBUSB_VERSION}.tar.gz"
tar -xzf "${LIBUSB_VERSION}.tar.gz"
cd "libusb-${LIBUSB_VERSION#v}"

# Configure and build
./autogen.sh --no-configure
./configure \
    --host="${TOOLCHAIN_PREFIX}" \
    CC="${SYSROOT}/bin/${TOOLCHAIN_PREFIX}-gcc" \
    --prefix="${SYSROOT}/${TOOLCHAIN_PREFIX}" \
    --enable-static \
    --disable-shared \
    --disable-udev

make -j"$(nproc)"

# Install to sysroot
make install

echo "✓ Installed libusb-1.0.a to ${SYSROOT}/${TOOLCHAIN_PREFIX}/lib/"
ls -lh "${SYSROOT}/${TOOLCHAIN_PREFIX}/lib/libusb-1.0.a"
