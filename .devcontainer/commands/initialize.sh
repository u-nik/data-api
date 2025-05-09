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
if ! isCertValid localhost; then
    installCertificate localhost
else
    echo "⚠️ Skipping: Certificate already exists and still valid."
    echo "If you want to update the certificate, please delete it and rebuild the dev container."
    echo "You can also run the command manually: mkcert -install && mkcert localhost"
fi

HYDRA_HOST="hydra.localhost"
HOSTS_FILE="/etc/hosts"
ENTRY="127.0.0.1 $HYDRA_HOST"
drawHeader "🔗 Setup ${HYDRA_HOST} domain..."

# Prüfen, ob der Eintrag schon vorhanden ist
if grep -qE "^\s*127\.0\.0\.1\s+.*\bhydra\.localhost\b" "$HOSTS_FILE"; then
    printf "✅ '%s' is already in the hosts file.\n" "${HYDRA_HOST}"
else
    echo "➕ Füge 'hydra.localhost' zur hosts-Datei hinzu..."
    # Temporäre Datei erzeugen
    TMP_FILE=$(mktemp)
    cp "$HOSTS_FILE" "$TMP_FILE"
    echo "$ENTRY" >>"$TMP_FILE"

    # Mit sudo zurückkopieren
    sudo cp "$TMP_FILE" "$HOSTS_FILE"
    rm "$TMP_FILE"

    printf "✅ '%s' wurde erfolgreich hinzugefügt.\n" "${HYDRA_HOST}"
fi
