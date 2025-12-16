#!/bin/bash

# Script pour arrêter les services du projet
# Usage: ./stop.sh

set -e

# Couleurs pour les messages
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"
PID_FILE="$PROJECT_ROOT/.services.pid"

# Fonction pour afficher les messages
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Vérifier si le fichier PID existe
if [ ! -f "$PID_FILE" ]; then
    warn "Aucun fichier PID trouvé. Les services ne semblent pas être en cours d'exécution."
    exit 0
fi

info "Arrêt des services..."

# Lire les PIDs et arrêter les processus
PIDS_STOPPED=0
PIDS_NOT_FOUND=0

while read pid; do
    if [ -n "$pid" ] && ps -p "$pid" > /dev/null 2>&1; then
        info "Arrêt du processus $pid..."
        kill "$pid" 2>/dev/null || true
        # Attendre un peu pour que le processus se termine
        sleep 1
        # Si le processus existe encore, forcer l'arrêt
        if ps -p "$pid" > /dev/null 2>&1; then
            warn "Le processus $pid ne s'est pas arrêté, arrêt forcé..."
            kill -9 "$pid" 2>/dev/null || true
        fi
        PIDS_STOPPED=$((PIDS_STOPPED + 1))
    else
        PIDS_NOT_FOUND=$((PIDS_NOT_FOUND + 1))
    fi
done < "$PID_FILE"

# Supprimer le fichier PID
rm -f "$PID_FILE"

if [ $PIDS_STOPPED -gt 0 ]; then
    info "✅ $PIDS_STOPPED processus arrêté(s)."
fi

if [ $PIDS_NOT_FOUND -gt 0 ]; then
    warn "$PIDS_NOT_FOUND processus n'ont pas été trouvés (peut-être déjà arrêtés)."
fi

info "Services arrêtés."

