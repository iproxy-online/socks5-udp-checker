- [SOCKS5 UDP Checker](#socks5-udp-checker)
  - [Overview](#overview)
    - [What This Tool Does](#what-this-tool-does)
  - [Installation](#installation)
    - [Linux and macOS Installation](#linux-and-macos-installation)
    - [Windows Installation](#windows-installation)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
    - [SOCKS5 Proxy Formats](#socks5-proxy-formats)
  - [Understanding Results](#understanding-results)
    - [Success](#success)
    - [Common Failures](#common-failures)
  - [Controls](#controls)
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
   - **Linux (amd64)**: `socks5-udp-checker-linux-amd64`
   - **Linux (arm64)**: `socks5-udp-checker-linux-arm64`
   - **macOS (Intel)**: `socks5-udp-checker-darwin-amd64`
   - **macOS (Apple Silicon)**: `socks5-udp-checker-darwin-arm64`
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

## Understanding Results

### Success

When the test is successful, you will see output similar to:

```
✅ Test Successful

╭──────────────────────────────────────────────────╮
│                                                  │
│  NTP Response Summary:                           │
│    Server Time:     2025-07-19 10:04:16.500 UTC  │
│    Round Trip Time: 84.405832ms                  │
│    Clock Offset:    164.200067ms                 │
│    Stratum:         1                            │
│                                                  │
╰──────────────────────────────────────────────────╯
```

### Common Failures

When the test fails, you may encounter various error messages. Here is an example of a failure due to connection issues:

```
❌ Test Failed

╭───────────────────────────────────────────────────────────────────────────────────╮
│                                                                                   │
│  Error: NTP request failed: dial tcp 127.0.0.1:1080: connect: connection refused  │
│                                                                                   │
╰───────────────────────────────────────────────────────────────────────────────────╯
```

**Common Error Messages and Solutions:**

- **`Invalid Username or Password for Auth`**: The provided credentials are incorrect or not accepted by the proxy server. Double-check your username and password, and ensure they match the proxy server's authentication requirements.

- **`connect: connection refused`**: The proxy server is not reachable or not accepting connections. Verify that:
  - The proxy server is running and accessible
  - The hostname and port are correct
  - No firewall is blocking the connection

- **`lookup unexistent.proxy: no such host`**: The proxy hostname could not be resolved. This typically indicates:
  - DNS resolution issues
  - Incorrect or misspelled hostname
  - Network connectivity problems

- **`Host unreachable`**: The proxy server cannot reach the specified NTP server. This may occur due to:
  - The proxy server not supporting UDP traffic relay
  - Network routing issues between the proxy and target server
  - Firewall restrictions blocking UDP traffic

## Controls

- **Enter**: Start test or run another test
- **Tab**: Navigate form fields
- **v**: Show version
- **Ctrl+C**: Exit

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for complete license terms.
