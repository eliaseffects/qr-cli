#!/usr/bin/env bash
set -euo pipefail

REPO="eliaseffects/qr-cli"
BIN_NAME="qr"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

case "$OS" in
  darwin|linux) ;; 
  *)
    echo "Unsupported OS: $OS" >&2
    exit 1
    ;;
esac

TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | \
  grep -m1 '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
  echo "Failed to determine latest release tag" >&2
  exit 1
fi

ASSET="qr-cli_${TAG#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"

BIN_DIR=${BIN_DIR:-/usr/local/bin}
if [ ! -w "$BIN_DIR" ]; then
  BIN_DIR="$HOME/.local/bin"
  mkdir -p "$BIN_DIR"
fi

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -sSL "$URL" -o "$TMP_DIR/${ASSET}"
tar -xzf "$TMP_DIR/${ASSET}" -C "$TMP_DIR"

if [ ! -f "$TMP_DIR/$BIN_NAME" ]; then
  echo "Binary not found in archive" >&2
  exit 1
fi

install -m 0755 "$TMP_DIR/$BIN_NAME" "$BIN_DIR/$BIN_NAME"

printf "Installed %s to %s\n" "$BIN_NAME" "$BIN_DIR"
