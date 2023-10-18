#!/bin/sh

# Migrash install script
# install example commands
#   curl -sSL https://raw.githubusercontent.com/psykka/migrash/main/install.sh | sh
#   wget -qO- https://raw.githubusercontent.com/psykka/migrash/main/install.sh | sh

# check if has wget or curl installed
if [ -x "$(command -v wget)" ]; then
    DOWNLOADER="wget"
elif [ -x "$(command -v curl)" ]; then
    DOWNLOADER="curl"
else
    echo "Error: wget or curl not found"
    exit 1
fi

run_root_installer() {
    if [ "$DOWNLOADER" = "wget" ]; then
        wget -qO- https://raw.githubusercontent.com/psykka/migrash/main/install.sh | sh
    elif [ "$DOWNLOADER" = "curl" ]; then
        curl -sSL https://raw.githubusercontent.com/psykka/migrash/main/install.sh | sh
    fi
}

# Upgrade to root if not root
if [ "$(id -u)" != "0" ]; then
    echo "Installing migrash"
    sudo sh -c "$(declare -f run_root_installer); run_root_installer"
    exit 0
fi

# Download migrash.sh to /usr/local/bin
donwload_file

# Make migrash executable
chmod 555 /usr/local/bin/migrash

# Check if migrash was installed successfully
if [ ! -x "$(command -v migrash)" ]; then
    echo "Error: migrash not installed correctly"
    exit 1
fi

echo "Migrash installed successfully"