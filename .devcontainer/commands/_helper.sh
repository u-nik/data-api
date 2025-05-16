#!/bin/bash

# ----------------------------------------------------------
# Function: drawLine
# Description: Draws a horizontal line in the console for visual separation.
# Usage: drawLine [width]
# Parameters:
#   - width: Optional. The length of the line to draw. If not provided, a default width is used.
# Returns: None. Outputs a line directly to stdout.
# ----------------------------------------------------------
drawLine() {
    printf '\n%*s\n' "$(tput cols)" '' | tr ' ' '-'
}

# Helper functions for devcontainer commands
#
# @description Draws a formatted header with text
# @param $1 The text to display in the header
# @example
#   drawHeader "Starting installation"
drawHeader() {
    drawLine
    printf "$1\n\n"
}

# Function to draw a logo in the terminal.
#
# This function is used to display a visually appealing ASCII art or text-based logo
# in the terminal. The logo is typically used as part of the command line interface
# to provide branding or visual identification.
#
# Usage:
#   drawLogo
#
# Parameters:
#   None
#
# Returns:
#   Outputs the logo to stdout
drawLogo() {
    drawLine
    printf "Open Source Data-API\n"
}

# Checks if Homebrew is installed and properly set up in the environment
# This function verifies the availability of the 'brew' command and
# ensures that Homebrew is correctly configured for the current user
#
# Usage:
#   checkHomebrew
#
# Returns:
#   0 if Homebrew is installed and functioning
#   non-zero otherwise
#
checkHomebrew() {
    # Check if Homebrew is installed
    if ! command -v brew &>/dev/null; then
        return 1
    fi
}

# Check if the current operating system is macOS.
#
# Returns:
#   0 (true) if running on macOS, 1 (false) otherwise.
#
# Example:
#   if isMacOS; then
#     echo "Running on macOS"
#   else
#     echo "Not running on macOS"
#   fi
isMacOS() {
    # Check if the OS is macOS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        return 0
    else
        return 1
    fi
}

# Check if the script is running on a Linux-based system
#
# Returns:
#   0 (true) if the system is Linux
#   1 (false) otherwise
#
# Usage:
#   if isLinux; then
#     echo "Running on Linux"
#   fi
isLinux() {
    # Check if the OS is Linux
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        return 0
    else
        return 1
    fi
}

# Function: checkLocalDependencies
# Description: Validates that required local dependencies are installed on the system
#
# This function checks if necessary tools and applications are available
# in the local environment before proceeding with further operations.
#
# Usage:
#   checkLocalDependencies
#
# Returns:
#   0 if all dependencies are present
#   Non-zero if any dependency is missing
checkLocalDependencies() {
    printf "OS Type: %s\n" $OSTYPE
    if isMacOS; then
        if ! checkHomebrew; then
            echo "ERROR: Homebrew is not installed. Please install it first."
            exit 1
        fi
    fi
    if ! hasAdminPrivileges; then
        printf "⚠️ WARNING: The current user does not have admin privileges.\n"
    fi
}

# Check if the certificate is valid.
#
# This function verifies the validity of a certificate.
#
# Returns:
#   0 if the certificate is valid
#   1 if the certificate is invalid or does not exist
#
isCertValid() {
    if [ -f ${1}.pem ]; then
        # Check if openssl is installed
        if ! command -v openssl &>/dev/null; then
            echo "⚠️ WARNING: openssl is not installed. Certificate validity check skipped."
            echo "If certificate is not valid, delete the certificate manually and rebuild the container."
            # Only a warning, not an error
            return 0
        fi

        # Check if the certificate is valid for more than 30 days
        cert_date=$(openssl x509 -in ${1}.pem -noout -enddate | cut -d= -f2)
        if isMacOS; then
            # macOS compatible date parsing
            cert_date_seconds=$(date -j -f "%b %d %H:%M:%S %Y %Z" "$cert_date" +%s 2>/dev/null)
        else
            # Linux date parsing
            cert_date_seconds=$(date -d "$cert_date" +%s)
        fi
        current_date_seconds=$(date +%s)
        days_left=$(((cert_date_seconds - current_date_seconds) / 86400))
        if [ $days_left -le 30 ]; then
            return 0
        fi
    fi
    return 1
}

# Function: installCertificate
#
# Installs SSL certificates into the trusted certificate store.
# This function is used during the container setup to ensure that custom certificates
# are recognized as trusted by the host operating system and applications.
#
# Usage:
#   installCertificate
#
# Parameters:
#   None
#
# Returns:
#   0 on success, non-zero on failure
installCertificate() {
    # Install mkcert and create certificate, if not already existing
    if ! command -v mkcert &>/dev/null; then
        echo "installing mkcert..."
        if isLinux; then
            sudo apt-get install libnss3-tools mkcert
        elif isMacOS; then
            brew install mkcert nss
        fi
    fi

    # Install the trusted root certificate
    if hasAdminPrivileges; then
        mkcert -install
    else
        echo "⚠️ WARNING: Cannot install trusted root certificate due to missing admin privileges."
    fi

    mkdir -p .certs
    mkcert -cert-file .certs/${1}.pem -key-file .certs/${1}-key.pem ${1}
}

# Checks if the current user has administrative privileges
# This function determines whether the script is running with elevated/root permissions
# Used to ensure certain operations can be performed that require administrative access
#
# Returns:
#   0 if the user has admin privileges (root/sudo)
#   1 otherwise
#
hasAdminPrivileges() {
    # Check if the current user is root
    if [ "$(id -u)" -eq 0 ]; then
        return 0
    fi

    if isLinux; then
        # Check if the current user is in the admin group
        if groups | grep -q '\bsudo\b'; then
            return 0
        fi

    elif isMacOS; then
        # Check if the current user is in the sudo group
        if groups | grep -q '\badmin\b'; then
            return 0
        fi
    fi

    # Check if the current user has sudo privileges without a password
    if sudo -n true 2>/dev/null; then
        return 0
    fi

    return 1
}

# Fügt einen Hostnamen zu /etc/hosts hinzu, falls er noch nicht existiert
# Usage: addHostname app.localhost
addHostname() {
    local hostname="$1"
    if [ -z "$hostname" ]; then
        echo "Hostname fehlt!"
        return 1
    fi

    if grep -q "[[:space:]]$hostname" /etc/hosts; then
        printf "✅ '%s' is already in the hosts file.\n" "${hostname}"
    else
        echo "127.0.0.1 $hostname" | sudo tee -a /etc/hosts >/dev/null
        printf "✅ '%s' wurde erfolgreich hinzugefügt.\n" "${hostname}"
    fi
}
