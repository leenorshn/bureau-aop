#!/bin/bash

# Script pour lancer le projet (Monolithe)
# Usage: ./start.sh [--background]

set -e

# Couleurs pour les messages
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"
PID_FILE="$PROJECT_ROOT/.server.pid"
LOG_DIR="$PROJECT_ROOT/logs"

# Créer le dossier de logs
mkdir -p "$LOG_DIR"

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

# Fonction pour nettoyer les processus en cas d'interruption
cleanup() {
    info "Arrêt du serveur..."
    if [ -f "$PID_FILE" ]; then
        while read pid; do
            if ps -p "$pid" > /dev/null 2>&1; then
                kill "$pid" 2>/dev/null || true
            fi
        done < "$PID_FILE"
        rm -f "$PID_FILE"
    fi
    exit 0
}

# Capturer les signaux pour nettoyer proprement
trap cleanup SIGINT SIGTERM

# Vérifier si le serveur est déjà en cours d'exécution
if [ -f "$PID_FILE" ]; then
    warn "Le serveur semble déjà être en cours d'exécution."
    read -p "Voulez-vous l'arrêter et redémarrer? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        ./stop.sh
    else
        exit 1
    fi
fi

# Vérifier MongoDB
info "Vérification de MongoDB..."
if ! command -v mongosh &> /dev/null && ! command -v mongo &> /dev/null; then
    warn "MongoDB client non trouvé. Vérifiez que MongoDB est installé et accessible."
fi

# Charger les variables d'environnement si .env existe
if [ -f "$PROJECT_ROOT/.env" ]; then
    info "Chargement des variables d'environnement depuis .env"
    export $(cat "$PROJECT_ROOT/.env" | grep -v '^#' | xargs)
fi

# Variables d'environnement par défaut
export MONGO_URI="${MONGO_URI:-mongodb://localhost:27017}"
export MONGO_DB_NAME="${MONGO_DB_NAME:-mlm_db}"
export APP_PORT="${APP_PORT:-4000}"

info "Configuration:"
info "  MongoDB URI: $MONGO_URI"
info "  MongoDB DB: $MONGO_DB_NAME"
info "  Server Port: $APP_PORT"

# Vérifier que Go est installé
if ! command -v go &> /dev/null; then
    error "Go n'est pas installé. Veuillez installer Go 1.21+"
    exit 1
fi

# Vider le fichier PID
> "$PID_FILE"

# Lancer le serveur monolithique
info "Démarrage du serveur..."
cd "$PROJECT_ROOT"
if [ "$1" == "--background" ]; then
    go run server.go > "$LOG_DIR/server.log" 2>&1 &
    SERVER_PID=$!
    echo "$SERVER_PID" >> "$PID_FILE"
    info "Serveur démarré (PID: $SERVER_PID, logs: $LOG_DIR/server.log)"
else
    go run server.go > "$LOG_DIR/server.log" 2>&1 &
    SERVER_PID=$!
    echo "$SERVER_PID" >> "$PID_FILE"
    info "Serveur démarré (PID: $SERVER_PID)"
fi

# Attendre que le serveur soit prêt
info "Attente du serveur..."
sleep 3
for i in {1..10}; do
    if curl -s "http://localhost:$APP_PORT/query" > /dev/null 2>&1 || \
       curl -s "http://localhost:$APP_PORT/" > /dev/null 2>&1; then
        info "Serveur est prêt!"
        break
    fi
    if [ $i -eq 10 ]; then
        warn "Serveur ne répond pas encore, mais on continue..."
    else
        sleep 1
    fi
done

info ""
info "=========================================="
info "✅ Serveur démarré avec succès!"
info "=========================================="
info ""
info "GraphQL Endpoint: http://localhost:$APP_PORT/query"
info "GraphQL Playground: http://localhost:$APP_PORT/"
info ""
info "Logs: $LOG_DIR/server.log"
info ""
info "Pour arrêter le serveur, utilisez: ./stop.sh"
info ""

# Si en mode background, on sort
if [ "$1" == "--background" ]; then
    info "Serveur lancé en arrière-plan."
    exit 0
fi

# Sinon, on attend et affiche les logs
info "Appuyez sur Ctrl+C pour arrêter le serveur..."
wait





