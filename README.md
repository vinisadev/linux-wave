# LinuxWave - Linux Facial Recognition Authentication System

A Linux facial recognition authentication system that integrates with PAM (Pluggable Authentication Modules) for seamless user authentication using facial biometrics.

## Overview

LinuxWave provides a secure, privacy-focused facial authentication solution for Linux systems. The system runs as a background service and integrates with the PAM authentication stack to enable passwordless login via facial recognition.

### Components

- **linuxwave-service** - Main systemd service that manages authentication
- **linuxwave-pam** - PAM helper module for authentication integration
- **linuxwave-cli** - Command-line management tool for configuration and user management
- **linuxwave-enroll** - GUI application for enrolling user facial data

## Prerequisites

### Build Requirements

- **Go 1.21+** - Required for compilation
- **Make** - Build automation

### Runtime Dependencies (Future Stories)

The following dependencies are not required for the current version but will be needed in future releases:

- OpenCV 4.x development headers (for GoCV)
- TensorFlow C library
- V4L2 development headers (camera interface)
- BlueZ (Bluetooth libraries)
- GTK or Qt development libraries (for enrollment GUI)

## Building

### Quick Build

```bash
make build
```

This will compile all four binaries into the `bin/` directory.

### Available Make Targets

- `make build` - Compile all binaries
- `make test` - Run all tests
- `make install` - Install binaries and systemd unit files (requires sudo)
- `make clean` - Remove build artifacts
- `make lint` - Run linting (golangci-lint or go vet)
- `make help` - Show available targets

## Development Environment Setup

1. **Install Go 1.21+**

   ```bash
   # Arch Linux
   sudo pacman -S go

   # Ubuntu/Debian
   sudo apt install golang-1.21
   ```

2. **Clone Repository**

   ```bash
   git clone https://github.com/vinisadev/linux-wave.git
   cd linux-wave
   ```

3. **Install Dependencies**

   ```bash
   go mod download
   ```

4. **Build Project**

   ```bash
   make build
   ```

5. **Run Tests**

   ```bash
   make test
   ```

## Testing

The project uses Go's standard testing framework with the `testify` library for assertions.

```bash
# Run all tests
make test

# Run tests with coverage
go test ./... -cover

# Run tests verbosely
go test ./... -v
```

## Project Structure

```
linux-wave/
├── cmd/                        # Binary entry points
│   ├── linuxwave-service/       # Main systemd service binary
│   ├── linuxwave-pam/           # PAM helper module binary
│   ├── linuxwave-cli/           # CLI management tool
│   └── linuxwave-enroll/        # Enrollment GUI application
├── internal/                   # Private packages
│   ├── camera/                 # V4L2 camera interface (future)
│   ├── recognition/            # Face detection + recognition engine (future)
│   ├── liveness/               # Liveness detection (future)
│   ├── bluetooth/              # BLE phone proximity detection (future)
│   ├── storage/                # Encrypted embedding storage (future)
│   └── pam/                    # PAM integration logic (future)
├── pkg/                        # Exportable packages
├── models/                     # Pre-trained ML models (future)
├── config/                     # Configuration file templates (future)
├── systemd/                    # Systemd service unit files (future)
├── Makefile                    # Build automation
└── go.mod                      # Go module definition
```

## Quick Start

### Running Binaries

After building, binaries are available in the `bin/` directory:

```bash
# Run service binary
./bin/linuxwave-service

# Run CLI tool
./bin/linuxwave-cli

# Run enrollment GUI
./bin/linuxwave-enroll

# Run PAM helper
./bin/linuxwave-pam
```

**Note:** These are currently stub implementations. Full functionality will be implemented in future stories.

## Installation

To install binaries system-wide:

```bash
sudo make install
```

This installs binaries to `/usr/local/bin/` and systemd unit files to `/etc/systemd/system/`.

## Contributing

### Code Standards

- Follow Go standard formatting (use `gofmt`)
- Run `make lint` before committing
- Ensure all tests pass with `make test`
- Maintain test coverage above 70% for `internal/` packages

### Development Workflow

1. Create feature branch
2. Implement changes
3. Add tests
4. Run `make test` and `make lint`
5. Submit pull request

## License

TBD

## Current Status

**Version:** 0.1.0 (Early Development)

This is the initial project structure setup. Core authentication features will be implemented in subsequent development iterations.
