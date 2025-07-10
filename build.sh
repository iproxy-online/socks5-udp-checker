#!/bin/bash
set -e

echo "🚀 Building SOCKS5 UDP Checker for multiple platforms..."

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf dist/

# Run GoReleaser build
echo "📦 Running GoReleaser build..."
goreleaser build --snapshot --clean

# List generated binaries
echo "✅ Build completed! Generated binaries:"
find dist/ -name "socks5-udp-checker*" -type f -exec ls -lh {} \;

echo ""
echo "📁 Binaries are located in the dist/ directory:"
echo "  • Linux (amd64): dist/socks5-udp-checker_linux_amd64_v1/socks5-udp-checker"
echo "  • Linux (arm64): dist/socks5-udp-checker_linux_arm64_v8.0/socks5-udp-checker"
echo "  • macOS (Intel): dist/socks5-udp-checker_darwin_amd64_v1/socks5-udp-checker"
echo "  • macOS (Apple Silicon): dist/socks5-udp-checker_darwin_arm64_v8.0/socks5-udp-checker"
echo "  • Windows (amd64): dist/socks5-udp-checker_windows_amd64_v1/socks5-udp-checker.exe"
echo ""
echo "🎉 Ready to distribute!"
