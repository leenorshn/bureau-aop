#!/bin/bash

# Script pour redémarrer les services du projet
# Usage: ./restart.sh [--background]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🔄 Redémarrage des services..."
"$SCRIPT_DIR/stop.sh"
sleep 2
"$SCRIPT_DIR/start.sh" "$@"

