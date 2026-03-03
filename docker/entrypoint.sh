#!/bin/bash
set -e

DOTBINS_ARCH="linux/$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"
DOTBINS_BIN="/home/claude/.dotbins/$DOTBINS_ARCH/bin"

if [ -f "/home/claude/.config/dotbins/dotbins.yaml" ]; then
    if [ -d "$DOTBINS_BIN" ] && [ -n "$(ls -A "$DOTBINS_BIN" 2>/dev/null)" ]; then
        echo "[devbox] Tools already cached, skipping sync..."
    else
        echo "[devbox] Syncing tools via dotbins..."
        dotbins sync || true
    fi
    export PATH="$DOTBINS_BIN:$PATH"
fi

exec "$@"
