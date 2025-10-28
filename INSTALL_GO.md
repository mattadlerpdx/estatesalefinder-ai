# Installing Go in WSL

You have two options to install Go in your WSL environment:

## Option 1: Automated Installation (Recommended)

Run the provided installation script:

```bash
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai
bash install_go_wsl.sh
```

This script will:
1. Download Go 1.21.6 for your architecture
2. Install it to `/usr/local/go`
3. Configure your PATH automatically
4. Set up GOPATH
5. Verify the installation

After the script completes:
```bash
# Restart your terminal OR run:
source ~/.bashrc

# Verify installation:
go version

# Run backend tests:
cd backend
./run_tests.sh
```

## Option 2: Manual Installation

If you prefer to install manually:

### Step 1: Download Go

```bash
cd /tmp
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
```

**Note**: For ARM64 (like Apple M1/M2 via WSL), use:
```bash
wget https://go.dev/dl/go1.21.6.linux-arm64.tar.gz
```

### Step 2: Remove Old Installation (if exists)

```bash
sudo rm -rf /usr/local/go
```

### Step 3: Extract Go

```bash
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
```

### Step 4: Configure PATH

Add to your `~/.bashrc` (or `~/.zshrc` if using zsh):

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
```

### Step 5: Apply Changes

```bash
source ~/.bashrc
```

### Step 6: Verify Installation

```bash
go version
```

You should see something like:
```
go version go1.21.6 linux/amd64
```

### Step 7: Test It

```bash
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai/backend
go version
go build ./...
```

## Verifying the Installation

Once Go is installed, run these commands to verify everything works:

```bash
# Check Go version
go version

# Check Go environment
go env

# Navigate to project
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai/backend

# Try to build the project
go build ./...

# If build succeeds, run tests
./run_tests.sh
```

## Troubleshooting

### "go: command not found"

**Problem**: PATH not set correctly

**Solution**:
```bash
# Check if Go is installed
ls -la /usr/local/go/bin/go

# If it exists, add to PATH manually for current session:
export PATH=$PATH:/usr/local/go/bin

# Then add permanently:
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Permission Denied

**Problem**: Need sudo for installation

**Solution**: The automated script and manual steps use `sudo` - you may be prompted for your password.

### Wrong Architecture

**Problem**: Downloaded wrong Go version for your CPU

**Solution**: Check your architecture:
```bash
uname -m
```

- If `x86_64`: Use `linux-amd64`
- If `aarch64` or `arm64`: Use `linux-arm64`

### WSL 1 vs WSL 2

This installation works for both WSL 1 and WSL 2. If you're unsure which you have:

```bash
wsl -l -v
```

## After Installation - Running Tests

Once Go is installed:

```bash
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai/backend

# Set up database URL for integration tests (optional)
export DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"

# Run all tests
./run_tests.sh

# Or run specific tests
go test -v ./internal/domain/listing -run TestListing
go test -v ./internal/domain/listing -run TestListingIntegrationSuite
```

## Quick Start Summary

```bash
# 1. Install Go
cd /mnt/c/Users/matt/Desktop/Stuff/estatesalefinder-ai
bash install_go_wsl.sh

# 2. Restart terminal or reload config
source ~/.bashrc

# 3. Verify
go version

# 4. Run validation
./validate_refactor.sh

# 5. Run backend tests
cd backend
./run_tests.sh

# 6. Run frontend tests
cd ../frontend
npm install
./run_tests.sh
```

## Alternative: Use Windows Go Instead

If you prefer not to install Go in WSL, you can use Go installed on Windows:

```powershell
# In PowerShell (not WSL)
cd C:\Users\matt\Desktop\Stuff\estatesalefinder-ai\backend
go test -v .\...
```

You can download Go for Windows from: https://go.dev/dl/

## Need Help?

If you encounter issues:
1. Check the [official Go installation guide](https://go.dev/doc/install)
2. Verify your WSL version: `wsl -l -v`
3. Check WSL2 docs if using WSL2
4. Make sure you have sudo access in WSL

## Summary

- **Recommended**: Use `bash install_go_wsl.sh` for automated installation
- **Manual**: Follow the step-by-step guide above
- **Verification**: Run `go version` and `./run_tests.sh`
- **Alternative**: Install Go on Windows and run tests from PowerShell
