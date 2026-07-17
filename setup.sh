#!/bin/bash
set -e

echo "=========================================="
echo "    Droprcam Build & Setup Script         "
echo "=========================================="

GO_VERSION="1.20.14"
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"

# 1. Download and extract Go 1.20 (if not already present)
if [ ! -d ".go1.20" ]; then
    echo "[+] Downloading Go ${GO_VERSION} (Required for Linux 2.6.38 kernels)..."
    wget -q "https://go.dev/dl/${GO_TAR}"
    
    echo "[+] Extracting Go..."
    tar -xzf "${GO_TAR}"
    
    echo "[+] Setting up local Go environment..."
    mv go .go1.20
    rm "${GO_TAR}"
else
    echo "[+] Go ${GO_VERSION} toolchain already exists in .go1.20"
fi

# 2. Build the Daemon
echo "[+] Compiling droprcam binary (GOARM=5 for Ambarella A5s CPU)..."

# Ensure we are using the isolated Go 1.20 binary
GO_BIN="./.go1.20/bin/go"

# Clean modules
${GO_BIN} mod tidy

# Build for ARMv6 with Software Floating Point
GOOS=linux GOARCH=arm GOARM=5 ${GO_BIN} build -ldflags="-s -w" -o droprcam

echo "[+] Build complete!"
echo "[+] Binary size:"
ls -lh droprcam
echo "=========================================="
echo "Run ./install_to_camera.sh to deploy to the camera!"
echo "=========================================="
