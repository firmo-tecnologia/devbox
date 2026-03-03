#!/bin/bash
set -e

DOTBINS_ARCH="linux/$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"

if [ -f "/root/.config/dotbins/config.yaml" ]; then
    echo "[devbox] Syncing tools via dotbins..."
    dotbins sync --quiet || true
    export PATH="/root/.dotbins/$DOTBINS_ARCH/bin:$PATH"
fi

exec "$@"
