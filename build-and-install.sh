#!/bin/bash

if [ "$(id -u)" -ne 0 ]; then
    echo "Error: This script must be executed as root."
    exit 1
fi

# Define the paths
SOURCE_PATH="./src"
BINARY_NAME="knocker-up"
DESTINATION_PATH="/usr/local/bin/"
SYSTEMD_PATH="./systemd/"
SYSTEMD_DESTINATION_PATH="/etc/systemd/system/"

export GOPATH="/tmp/gopath"

# Compile the Go program trimming the path and removing the debug information
cd ${SOURCE_PATH} && go build -ldflags="-w -s -buildid=" -trimpath -o "/tmp/${BINARY_NAME}" "main.go" && cd ..

# Moving the binary to the destination path
mv "/tmp/${BINARY_NAME}" ${DESTINATION_PATH}

# Setting execution permission
chmod +x "${DESTINATION_PATH}/${BINARY_NAME}"

# Copy the systemd service file to the systemd directory
cp "${SYSTEMD_PATH}knocker-up.service" "${SYSTEMD_DESTINATION_PATH}"

# Enabling the unprivileged ping
echo 'net.ipv4.ping_group_range=0 2147483647' | tee /etc/sysctl.d/99-allow-unprivileged-ping.conf >/dev/null

# Load the new sysctl configuration
sysctl -p /etc/sysctl.d/99-allow-unprivileged-ping.conf

# Reload systemd services
systemctl daemon-reload

# Start and enable the service
systemctl start knocker-up
systemctl enable knocker-up