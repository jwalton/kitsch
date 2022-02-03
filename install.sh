#!/usr/bin/env sh

GITHUB_USER='jwalton'
GITHUB_PROJECT='kitsch'
ARCHIVE_PREFIX=$GITHUB_PROJECT
EXECUTABLE=$GITHUB_PROJECT

error() {
    echo "Error: $*" >&2
    exit 1
}

INSTALL_DIR="/usr/local/bin"

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

# Print usage instructions.
printUsage() {
  cat <<END
usage: $(basename "$0") [-h] [--dir DIR]

Installs "${EXECUTABLE}" from https://github.com/${GITHUB_USER}/${GITHUB_PROJECT}.

Optional Arguments:
    -h        - Show this help.
    --dir DIR - Select the directory to install the executable to.

END
}

while :; do
    case $1 in
        -h|-\?|--help)   # Call a "show_help" function to display a synopsis, then exit.
            printUsage
            exit
            ;;
        -d|--dir)       # Takes an option argument, ensuring it has been specified.
            if [ -n "$2" ]; then
                INSTALL_DIR=$2
                shift
            else
                printf 'ERROR: "--dir" requires a non-empty option argument.\n' >&2
                exit 1
            fi
            ;;
        -?*)
            printf 'Unknown option: "%s"\n\n' "$1" >&2
            printUsage >&2
            exit 1
            ;;
        *)               # Default case: If no more options then break out of the loop.
            break
    esac

    shift
done

if [ $# -ne 0 ]; then
    printUsage >&2
    exit 1
fi

if ! command -v curl > /dev/null; then
    error "curl must be installed to run this script."
fi

if [ ! -d "${INSTALL_DIR}" ]; then
    printf 'Installation directory %s does not exist' "${INSTALL_DIR}" >&2
    exit 1
fi

# Download the archive and unpack it to /tmp.
ARCHIVE="https://github.com/${GITHUB_USER}/${GITHUB_PROJECT}/releases/latest/download/${ARCHIVE_PREFIX}_${OS}_${ARCH}.tar.gz"
echo "Downloading archive ${ARCHIVE}..."
curl --silent --location "${ARCHIVE}" | tar xz -C /tmp

if [ ! -e "/tmp/${EXECUTABLE}" ]; then
    error "Download failed."
fi

DEST="${INSTALL_DIR}/${EXECUTABLE}"
echo "Copying ${EXECUTABLE} to '${INSTALL_DIR}'."
if ! mv "/tmp/${EXECUTABLE}" "${DEST}" >  /dev/null 2>&1 ; then
    sudo -p "Enter your password to install ${EXECUTABLE}:" mv "/tmp/${EXECUTABLE}" "${DEST}" || {
        rm "/tmp/${EXECUTABLE}"
        error "Failed to copy ${EXECUTABLE} to '${INSTALL_DIR}'." >&2
    }
fi

echo "Installed ${EXECUTABLE} to '${INSTALL_DIR}'."

# Show setup instructions
"${INSTALL_DIR}/${EXECUTABLE}" setup