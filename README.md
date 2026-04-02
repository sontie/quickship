# QuickShip

[English](README.md) | [中文](README.zh-CN.md)

Lightweight Go deployment tool with parallel execution via SSH Agent Forwarding.

## Features

- **SSH Agent Forwarding** - Forward local private keys to remote servers without storing keys on servers
- **Parallel Deployment** - Concurrent execution across multiple hosts for improved efficiency
- **Real-time Logs** - Colored output with multi-host differentiation
- **Idempotent Operations** - Automatic Git clone or pull detection for repeatable execution

## Installation

### Option 1: One-Click Install (Recommended)

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/sontie/quickship/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/sontie/quickship/main/install.ps1 | iex
```

### Uninstall

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/sontie/quickship/main/uninstall.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/sontie/quickship/main/uninstall.ps1 | iex
```

Or manually remove:
- Linux/macOS: `sudo rm /usr/local/bin/qship && rm -rf ~/.quickship`
- Windows: Delete `qship.exe` and `%USERPROFILE%\.quickship` directory

### Upgrade

**If you have v0.1.1+ installed:**
```bash
qship upgrade
```

**If you have an older version (no upgrade command):**

Re-run the install script to automatically overwrite:

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/sontie/quickship/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/sontie/quickship/main/install.ps1 | iex
```

Or manually download and replace the binary:
```bash
# Get the latest version from GitHub Releases
curl -L -o qship https://github.com/sontie/quickship/releases/latest/download/qship-darwin-arm64
chmod +x qship
sudo mv qship /usr/local/bin/qship
```

### Option 2: Download Binary

Download the binary for your platform from [Releases](https://github.com/sontie/quickship/releases):

**Linux/macOS:**
```bash
chmod +x qship-*
sudo mv qship-* /usr/local/bin/qship
```

**Windows:**
Download `qship-windows-amd64.exe`, rename to `qship.exe` and add to PATH.

### Option 3: Build from Source

Requires Go 1.19+:

```bash
git clone https://github.com/sontie/quickship.git
cd quickship
go build -o qship
```

## Quick Start

### 1. Initialize Configuration

```bash
qship init
```

This generates a `qship.yaml` configuration template.

### 2. Edit Configuration

Edit `qship.yaml` to configure hosts and projects.

#### Configuration File Structure

```yaml
# Configuration version
version: "1.0"

# Host list
hosts:
  ali-dev:                    # Host alias
    addr: "192.168.1.100:22"  # SSH address and port
    user: "deploy"            # SSH username
  ali-prod:
    addr: "192.168.1.101:22"
    user: "deploy"

# Project list
projects:
  - name: "api-server"                        # Project name
    repo: "git@github.com:user/api-server.git" # Git repository URL
    path: "/opt/api-server"                   # Remote deployment path
    envs:                                     # Environment configuration
      dev: ["ali-dev"]                        # dev environment deploys to ali-dev
      prod: ["ali-prod"]                      # prod environment deploys to ali-prod
    scripts:                                  # Deployment scripts
      deploy: |                               # deploy script content
        docker-compose down
        docker-compose up -d --build
```

### 3. Environment Setup

Run the following command to check your SSH environment:

```bash
qship check
```

If you're familiar with SSH, configure it yourself. Otherwise, follow the prompts.

### 4. View Configuration

```bash
qship list
# or use shorthand
qship ls
```

Displays all configured hosts and projects.

### 5. Deploy

```bash
qship go dev
```

Deploy to the specified environment (dev/prod).

## Commands

### qship version / -v / --version
Display version information.

```bash
qship version
qship -v
qship --version
```

### qship init
Generate default configuration file `qship.yaml`.

### qship check
Check SSH Agent status and loaded keys.

### qship auth <host>
Copy local public key to remote host (equivalent to ssh-copy-id).

```bash
qship auth deploy@192.168.1.100
```

### qship list / ls
List all configured hosts and projects.

### qship go <env>
Deploy to specified environment.

```bash
qship go dev    # Deploy to development environment
qship go prod   # Deploy to production environment
```

### qship upgrade
Automatically check and upgrade to the latest version.

```bash
qship upgrade
```

If installed in system directory, requires permissions:
```bash
sudo qship upgrade
```

### qship exec "<command>" [hosts]
Execute command on specified hosts.

```bash
qship exec "uptime"                    # Execute on all hosts
qship exec "df -h" ali-dev,ali-prod   # Execute on specified hosts
```

## How It Works

1. Get ssh-agent connection from local `SSH_AUTH_SOCK`
2. Establish SSH connection using Agent Forwarding
3. Execute Git operations (clone/pull) on remote server
4. Run configured deployment scripts
5. Stream remote execution logs in real-time

## License

MIT
