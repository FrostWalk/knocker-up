#!/bin/bash

set -e

# Variables
BINARY_NAME="knocker-up"
VERSION="0.1.0"
SOURCE_PATH="./src"
DEB_BUILD_DIR="./deb-build"
DEB_PACKAGE_NAME="${BINARY_NAME}_${VERSION}_arm64.deb"
MAINTAINER="Me"

# Create temporary directory for Debian packaging
mkdir -p "${DEB_BUILD_DIR}/DEBIAN"
mkdir -p "${DEB_BUILD_DIR}/usr/local/bin"
mkdir -p "${DEB_BUILD_DIR}/etc/systemd/system"

# Compile the Go program
env GOOS=linux GOARCH=arm64 go build -ldflags="-w -s -buildid=" -trimpath -o "${DEB_BUILD_DIR}/usr/local/bin/${BINARY_NAME}" "${SOURCE_PATH}/main.go"

# Copy systemd unit file
cp "./systemd/${BINARY_NAME}.service" "${DEB_BUILD_DIR}/etc/systemd/system/"

# Create control file
cat > "${DEB_BUILD_DIR}/DEBIAN/control" <<EOF
Package: ${BINARY_NAME}
Version: ${VERSION}
Architecture: arm64
Maintainer: ${MAINTAINER}
Description: Knocker up
EOF

# Set permissions
chmod 755 "${DEB_BUILD_DIR}/DEBIAN/control"
chmod 755 "${DEB_BUILD_DIR}/usr/local/bin/${BINARY_NAME}"
chmod 644 "${DEB_BUILD_DIR}/etc/systemd/system/${BINARY_NAME}.service"

# Build the Debian package
dpkg-deb --build "${DEB_BUILD_DIR}" "${DEB_PACKAGE_NAME}"

echo "Debian package created: ${DEB_PACKAGE_NAME}"

rm -rf ${DEB_BUILD_DIR}
