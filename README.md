- [SOCKS5 UDP Checker](#socks5-udp-checker)
  - [Overview](#overview)
    - [What This Tool Does](#what-this-tool-does)
  - [Installation](#installation)
    - [Linux and macOS Installation](#linux-and-macos-installation)
    - [Windows Installation](#windows-installation)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
    - [SOCKS5 Proxy Formats](#socks5-proxy-formats)
    - [Examples](#examples)
  - [Understanding Results](#understanding-results)
    - [Success](#success)
    - [Common Failures](#common-failures)
  - [Controls](#controls)
  - [Troubleshooting](#troubleshooting)
  - [License](#license)


# SOCKS5 UDP Checker

A simple command-line tool for testing UDP support in SOCKS5 proxy servers using the Network Time Protocol (NTP).

[![Go Report Card](https://goreportcard.com/badge/github.com/iproxy-online/socks5-udp-checker)](https://goreportcard.com/report/github.com/iproxy-online/socks5-udp-checker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

SOCKS5 UDP Checker is a cross-platform command-line utility that verifies whether SOCKS5 proxy servers properly support UDP traffic forwarding. The tool performs this verification by attempting to execute Network Time Protocol (NTP) queries through the specified proxy connection.

### What This Tool Does

Many SOCKS5 proxy implementations support only TCP traffic and lack proper UDP relay functionality. This tool helps administrators and users verify UDP support by:

1. Connecting to your SOCKS5 proxy server
2. Establishing a UDP relay connection through the proxy
3. Sending an NTP time request to a public time server
4. Analyzing the response to confirm UDP traffic flows correctly

## Installation

Download a pre-built binary for your operating system:

1. Visit the [Releases page](https://github.com/iproxy-online/socks5-udp-checker/releases)
2. Download the appropriate binary for your system:
   - **Linux**: `socks5-udp-checker-linux-amd64`
   - **macOS**: `socks5-udp-checker-darwin-amd64`
   - **Windows**: `socks5-udp-checker-windows-amd64.exe`

### Linux and macOS Installation
```bash
# Download the binary (replace URL with actual release URL)
wget https://github.com/iproxy-online/socks5-udp-checker/releases/latest/download/socks5-udp-checker-linux-amd64

# Make it executable
chmod +x socks5-udp-checker-linux-amd64

# Move to system PATH (optional, requires administrator privileges)
sudo mv socks5-udp-checker-linux-amd64 /usr/local/bin/socks5-udp-checker
```

### Windows Installation
1. Download `socks5-udp-checker-windows-amd64.exe` from the releases page
2. Place the executable in a directory of your choice
3. Optionally, add the directory to your system PATH for global access

## Getting Started

Run the application from your terminal:

```bash
socks5-udp-checker
```

The interactive interface will guide you through:
1. Enter your SOCKS5 proxy details
2. Specify an NTP server (default: `time.google.com:123`)
3. Press Enter to test UDP connectivity
4. Review the results

## Configuration

### SOCKS5 Proxy Formats

**With Authentication:**
```
socks5://username:password@hostname:port
socks5://hostname:port:username:password
```

**Without Authentication:**
```
socks5://hostname:port
```

### Examples

```
# No authentication
socks5://proxy.example.com:1080

# With authentication
socks5://user:pass@proxy.example.com:1080
```

## Understanding Results

### Success
```
✅ Test Successful
Server Time: 2025-07-16 14:30:25
Round Trip Time: 45.2ms
```

### Common Failures

**Connection refused**: Proxy server not accessible
**Authentication failed**: Wrong username/password
**Timeout**: Proxy doesn't support UDP or blocks NTP

## Controls

- **Enter**: Start test or run another test
- **Tab**: Navigate form fields
- **v**: Show version
- **Ctrl+C**: Exit

## Troubleshooting

**Command not found**: Add binary to PATH or use `./socks5-udp-checker`
**Permission denied**: Run `chmod +x socks5-udp-checker`
**Connection timeout**: Check proxy address and network connectivity
**Authentication failed**: Verify username and password

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for complete license terms.
