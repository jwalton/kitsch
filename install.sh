#!/usr/bin/env sh

GITHUB_USER='jwalton'
GITHUB_PROJECT='kitsch-prompt'
ARCHIVE_PREFIX=$GITHUB_PROJECT
EXECUTABLE=$GITHUB_PROJECT

error() {
    echo "Error: $*" >&2
    exit 1
}

OS=$(uname -s)

ARCH=$(uname -m)
case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="armv7";;
    aarch64) ARCH="arm64";;
    x86) ARCH="i386";;
    x86_64) ARCH="x86_64";;
    i686) ARCH="i386";;
    i386) ARCH="i386";;
esac

if ! command -v curl > /dev/null; then
    error "curl must be installed to run this script."
fi

# Download the archive and unpack it to /tmp.
ARCHIVE="https://github.com/${GITHUB_USER}/${GITHUB_PROJECT}/releases/latest/download/${ARCHIVE_PREFIX}_${OS}_${ARCH}.tar.gz"
echo "Downloading archive ${ARCHIVE}..."
curl --silent --location "${ARCHIVE}" | tar xz -C /tmp

if [ ! -e "/tmp/${EXECUTABLE}" ]; then
    error "Download failed."
fi

echo "Copying ${EXECUTABLE} to /usr/local/bin."
sudo -p "Enter your password to install ${EXECUTABLE}:" mv "/tmp/${EXECUTABLE}" /usr/local/bin || {
    rm "/tmp/${EXECUTABLE}"
    error "Failed to copy ${EXECUTABLE} to /usr/local/bin." >&2
}

echo "Installed ${EXECUTABLE} to /usr/local/bin"

# Show setup instructions
/usr/local/bin/${EXECUTABLE} setup