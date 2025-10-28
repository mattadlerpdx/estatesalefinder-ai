#!/bin/bash

# Script to install Go in WSL
# Run with: bash install_go_wsl.sh

set -e

echo "======================================"
echo "Installing Go in WSL"
echo "======================================"
echo ""

# Check if Go is already installed
if command -v go &> /dev/null; then
    echo "✓ Go is already installed: $(go version)"
    echo ""
    read -p "Do you want to reinstall/upgrade? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Installation cancelled."
        exit 0
    fi
fi

# Determine architecture
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    GO_ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    GO_ARCH="arm64"
else
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
fi

echo "✓ Detected architecture: $ARCH (Go: $GO_ARCH)"
echo ""

# Set Go version (check https://go.dev/dl/ for latest)
GO_VERSION="1.21.6"
GO_TARBALL="go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
GO_URL="https://go.dev/dl/${GO_TARBALL}"

echo "Will install Go version: $GO_VERSION"
echo "Download URL: $GO_URL"
echo ""

# Download Go
echo "Step 1: Downloading Go..."
cd /tmp
if [ -f "$GO_TARBALL" ]; then
    echo "✓ Tarball already exists in /tmp"
else
    if wget "$GO_URL"; then
        echo "✓ Download complete"
    else
        echo "❌ Download failed"
        exit 1
    fi
fi
echo ""

# Remove old Go installation if it exists
echo "Step 2: Removing old Go installation (if any)..."
if [ -d "/usr/local/go" ]; then
    echo "Removing /usr/local/go (requires sudo)"
    sudo rm -rf /usr/local/go
    echo "✓ Removed old installation"
else
    echo "✓ No previous installation found"
fi
echo ""

# Extract new Go installation
echo "Step 3: Installing Go to /usr/local/go..."
echo "Extracting (requires sudo)..."
sudo tar -C /usr/local -xzf "$GO_TARBALL"
echo "✓ Go installed to /usr/local/go"
echo ""

# Setup PATH
echo "Step 4: Setting up PATH..."

# Detect shell
if [ -n "$BASH_VERSION" ]; then
    SHELL_RC="$HOME/.bashrc"
elif [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="$HOME/.zshrc"
else
    SHELL_RC="$HOME/.profile"
fi

echo "Detected shell config: $SHELL_RC"

# Check if PATH is already set
if grep -q "/usr/local/go/bin" "$SHELL_RC"; then
    echo "✓ PATH already configured in $SHELL_RC"
else
    echo "Adding Go to PATH in $SHELL_RC..."
    echo '' >> "$SHELL_RC"
    echo '# Go installation' >> "$SHELL_RC"
    echo 'export PATH=$PATH:/usr/local/go/bin' >> "$SHELL_RC"
    echo 'export PATH=$PATH:$HOME/go/bin' >> "$SHELL_RC"
    echo "✓ PATH configured"
fi
echo ""

# Set PATH for current session
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$HOME/go/bin

# Verify installation
echo "Step 5: Verifying installation..."
if /usr/local/go/bin/go version; then
    echo "✓ Go installed successfully!"
else
    echo "❌ Installation verification failed"
    exit 1
fi
echo ""

# Setup GOPATH
echo "Step 6: Setting up GOPATH..."
if [ ! -d "$HOME/go" ]; then
    mkdir -p "$HOME/go"
    echo "✓ Created $HOME/go directory"
else
    echo "✓ $HOME/go already exists"
fi
echo ""

# Cleanup
echo "Step 7: Cleaning up..."
rm -f "/tmp/$GO_TARBALL"
echo "✓ Removed temporary files"
echo ""

echo "======================================"
echo "✅ Go Installation Complete!"
echo "======================================"
echo ""
echo "Go version: $(/usr/local/go/bin/go version)"
echo "Go location: /usr/local/go"
echo "GOPATH: $HOME/go"
echo ""
echo "⚠️  IMPORTANT: Restart your terminal or run:"
echo "    source $SHELL_RC"
echo ""
echo "Then verify with:"
echo "    go version"
echo ""
echo "Next steps:"
echo "1. Restart terminal or: source $SHELL_RC"
echo "2. Verify: go version"
echo "3. Run backend tests: cd backend && ./run_tests.sh"
echo ""
