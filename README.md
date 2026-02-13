# SOCKS5 UDP Checker

A command-line tool for testing UDP support in SOCKS5 proxy servers using NTP.

[![Go Report Card](https://goreportcard.com/badge/github.com/iproxy-online/socks5-udp-checker)](https://goreportcard.com/report/github.com/iproxy-online/socks5-udp-checker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

Many SOCKS5 proxies only support TCP and don't relay UDP traffic. This tool checks UDP support by sending an NTP request through the proxy. If a time response comes back, UDP works.

How it works:

1. Connects to your SOCKS5 proxy
2. Establishes a UDP relay through the proxy
3. Sends an NTP time request to a public time server
4. Shows the result — success with NTP response details, or an error explaining what went wrong

## Installation

Download a pre-built binary from the [Releases page](https://github.com/iproxy-online/socks5-udp-checker/releases):

- **Linux (amd64)**: `socks5-udp-checker-linux-amd64`
- **Linux (arm64)**: `socks5-udp-checker-linux-arm64`
- **macOS (Intel)**: `socks5-udp-checker-darwin-amd64`
- **macOS (Apple Silicon)**: `socks5-udp-checker-darwin-arm64`
- **Windows**: `socks5-udp-checker-windows-amd64.exe`

### Linux and macOS

```bash
# Download (replace with your platform)
wget https://github.com/iproxy-online/socks5-udp-checker/releases/latest/download/socks5-udp-checker-linux-amd64

# Make executable
chmod +x socks5-udp-checker-linux-amd64

# Optional: move to PATH
sudo mv socks5-udp-checker-linux-amd64 /usr/local/bin/socks5-udp-checker
```

### Windows

1. Download `socks5-udp-checker-windows-amd64.exe` from the releases page
2. Place the executable in a directory of your choice
3. Optionally add the directory to your system PATH

## Usage

```bash
socks5-udp-checker
```

The interactive TUI will prompt you for your SOCKS5 proxy URL and an NTP server (default: `time.google.com:123`).

## SOCKS5 URL Formats

**With authentication:**
```
socks5://username:password@host:port
socks5://host:port:username:password
socks5:host:port:username:password
```

**Without authentication:**
```
socks5://host:port
socks5:host:port
```

## Controls

| Key | Action |
|-----|--------|
| `Enter` | Submit form / retry test |
| `Esc` | Cancel test / change settings |
| `Tab` | Navigate form fields |
| `Ctrl+C` | Exit |

## Troubleshooting

If the test fails, the error message will indicate the cause. Common errors:

- **`Invalid Username or Password for Auth`** — wrong credentials. Check your username and password.

- **`connect: connection refused`** — can't reach the proxy. Check that:
  - The proxy server is running
  - The hostname and port are correct
  - No firewall is blocking the connection

- **`no such host`** — the proxy hostname can't be resolved. Check for typos or DNS issues.

- **`Host unreachable`** — the proxy can't reach the NTP server. Possible causes:
  - The proxy doesn't support UDP relay
  - Firewall is blocking UDP traffic
  - Network routing issues

## License

MIT License. See [LICENSE](LICENSE) for details.
