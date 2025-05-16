#!/bin/sh
. .devcontainer/commands/_helper.sh
#
# The "initializeCommand" command is run on the host before the container is created.
# This script is used to set up the environment for the container.
#

drawLogo

drawHeader "📦 Checking local dependencies..."
# Check if OS is macOS and homebrew is installed
checkLocalDependencies

# Create certificate, if not already existing or not valid
drawHeader "🔐 Generating SSL certificates..."
if ! isCertValid app.localhost; then
    installCertificate app.localhost
else
    echo "⚠️ Skipping: Certificate already exists and still valid."
    echo "If you want to update the certificate, please delete it and rebuild the dev container."
    echo "You can also run the command manually: mkcert -install && mkcert app.localhost"
fi

if ! isCertValid hydra.localhost; then
    installCertificate hydra.localhost
else
    echo "⚠️ Skipping: Certificate already exists and still valid."
    echo "If you want to update the certificate, please delete it and rebuild the dev container."
    echo "You can also run the command manually: mkcert -install && mkcert hydra.localhost"
fi

drawHeader "🔗 Setup local hostnames..."
addHostname app.localhost
addHostname hydra.localhost
