#!/usr/bin/env sh
set -eu

APP_EXECUTABLE=${APP_EXECUTABLE:-a7}
BUILD_DIR=${BUILD_DIR:-build}

os=$(uname -s)
case "$os" in
  Darwin) goos=darwin ;;
  Linux) goos=linux ;;
  *)
    echo "Unsupported OS: $os" >&2
    exit 1
    ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64) goarch=amd64 ;;
  arm64|aarch64) goarch=arm64 ;;
  *)
    echo "Unsupported CPU architecture: $arch" >&2
    exit 1
    ;;
esac

binary_path=${BINARY_PATH:-${BUILD_DIR}/${APP_EXECUTABLE}-${goos}-${goarch}}

if [ ! -f "$binary_path" ]; then
  echo "Binary not found at $binary_path. Run 'make build' first." >&2
  exit 1
fi

if [ -n "${INSTALL_DIR:-}" ]; then
  install_dir=$INSTALL_DIR
elif [ "$goos" = "darwin" ]; then
  if [ -d /opt/homebrew/bin ]; then
    install_dir=/opt/homebrew/bin
  else
    install_dir=/usr/local/bin
  fi
else
  install_dir="${HOME}/.local/bin"
fi

if [ ! -d "$install_dir" ]; then
  if [ "$install_dir" = "${HOME}/.local/bin" ]; then
    mkdir -p "$install_dir"
  elif command -v sudo >/dev/null 2>&1; then
    sudo mkdir -p "$install_dir"
  else
    echo "Install dir does not exist and sudo not available: $install_dir" >&2
    exit 1
  fi
fi

if [ -w "$install_dir" ]; then
  install -m 0755 "$binary_path" "$install_dir/$APP_EXECUTABLE"
elif command -v sudo >/dev/null 2>&1; then
  sudo install -m 0755 "$binary_path" "$install_dir/$APP_EXECUTABLE"
else
  echo "Install dir not writable and sudo not available: $install_dir" >&2
  exit 1
fi

echo "Installed $APP_EXECUTABLE to $install_dir/$APP_EXECUTABLE"
