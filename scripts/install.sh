#!/usr/bin/env bash

set -euo pipefail

SB_INSTALL_DIR=${SB_INSTALL_DIR:-${HOME}/.local/bin}

get_latest_release() {
  curl --silent "https://api.github.com/repos/$1/releases/latest" | jq --raw-output .tag_name
}

get_architecture() {
  case "$(uname -m)" in
    aarch64|arm64)
      echo "arm64"
      ;;
    x86_64|amd64)
      echo "amd64"
      ;;
    *)
      echo ""
      ;;
  esac
}

get_kernel() {
  case "$(uname -s)" in
    Darwin)
      echo "darwin"
      ;;
    Linux)
      echo "linux"
      ;;
    *)
      echo ""
      ;;
  esac
}

ARCH=$(get_architecture)
KERNEL=$(get_kernel)

if [ -z "$ARCH" ]; then
    echo 1>&2 "ERROR: Architecture $(uname -m) is not supported."
    exit 1
fi

if [ -z "$KERNEL" ]; then
    echo 1>&2 "ERROR: Kernel $(uname -s) is not supported."
    exit 1
fi

RELEASE_TAG=$(get_latest_release aadam-ali/second-brain-cli)
TARBALL="sb-${RELEASE_TAG}-${KERNEL}-${ARCH}.tar.gz"
TARBALL_PATH="/tmp/${TARBALL}"
TARBALL_MD5="${TARBALL}.md5"
TARBALL_MD5_PATH="/tmp/${TARBALL_MD5}"

curl -Lo "${TARBALL_PATH}" "https://github.com/aadam-ali/second-brain-cli/releases/download/${RELEASE_TAG}/${TARBALL}"
curl -Lo "${TARBALL_PATH}.md5" "https://github.com/aadam-ali/second-brain-cli/releases/download/${RELEASE_TAG}/${TARBALL_MD5}"

echo "$(cat "${TARBALL_MD5_PATH}")  ${TARBALL_PATH}" | md5sum -c -

tar -C "${SB_INSTALL_DIR}" -xzf "${TARBALL_PATH}"
