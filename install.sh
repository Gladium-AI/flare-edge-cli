#!/bin/sh
set -eu

REPO="${REPO:-Gladium-AI/flare-edge-cli}"
BINARY="${BINARY:-flare-edge-cli}"

resolve_install_dir() {
	if [ -n "${INSTALL_DIR:-}" ]; then
		printf '%s\n' "${INSTALL_DIR}"
		return
	fi

	if [ -n "${XDG_BIN_HOME:-}" ]; then
		printf '%s\n' "${XDG_BIN_HOME}"
		return
	fi

	if [ -d "${HOME}/.local/bin" ]; then
		printf '%s\n' "${HOME}/.local/bin"
		return
	fi

	if [ -d "${HOME}/bin" ]; then
		printf '%s\n' "${HOME}/bin"
		return
	fi

	printf '%s\n' "${HOME}/.local/bin"
}

detect_platform() {
	case "$(uname -s)" in
		Linux*)
			GOOS="linux"
			;;
		Darwin*)
			GOOS="darwin"
			;;
		*)
			echo "Unsupported OS: $(uname -s)" >&2
			exit 1
			;;
	esac

	case "$(uname -m)" in
		x86_64|amd64)
			GOARCH="amd64"
			;;
		arm64|aarch64)
			GOARCH="arm64"
			;;
		*)
			echo "Unsupported architecture: $(uname -m)" >&2
			exit 1
			;;
	esac
}

fetch_latest_tag() {
	api_url="https://api.github.com/repos/${REPO}/releases/latest"

	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "${api_url}" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1
		return
	fi

	if command -v wget >/dev/null 2>&1; then
		wget -qO- "${api_url}" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -n 1
		return
	fi

	echo "Error: curl or wget is required" >&2
	exit 1
}

download_archive() {
	url="$1"
	out="$2"

	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "${url}" -o "${out}"
		return
	fi

	wget -q "${url}" -O "${out}"
}

INSTALL_DIR="$(resolve_install_dir)"
detect_platform

echo "Detected platform: ${GOOS}/${GOARCH}"

TAG="${TAG:-$(fetch_latest_tag)}"
if [ -z "${TAG}" ]; then
	echo "Error: could not determine latest release tag" >&2
	exit 1
fi

ARCHIVE="${BINARY}_${TAG}_${GOOS}_${GOARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${TAG}/${ARCHIVE}"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "${TMPDIR}"' EXIT INT TERM

echo "Latest release: ${TAG}"
echo "Downloading ${URL}"
download_archive "${URL}" "${TMPDIR}/${ARCHIVE}"

tar -xzf "${TMPDIR}/${ARCHIVE}" -C "${TMPDIR}"
EXTRACTED_DIR="${TMPDIR}/${BINARY}_${TAG}_${GOOS}_${GOARCH}"

mkdir -p "${INSTALL_DIR}"
install -m 0755 "${EXTRACTED_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"

echo "Installed ${BINARY} ${TAG} to ${INSTALL_DIR}/${BINARY}"

case ":${PATH}:" in
	*":${INSTALL_DIR}:"*)
		;;
	*)
		echo
		echo "Add ${INSTALL_DIR} to your PATH:"
		echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
		;;
esac
