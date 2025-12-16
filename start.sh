#!/bin/bash

# Script pour lancer le projet (Gateway + Tree Service)
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
TREE_SERVICE_DIR="$PROJECT_ROOT/services/tree-service"
GATEWAY_DIR="$PROJECT_ROOT/gateway"
PID_FILE="$PROJECT_ROOT/.services.pid"
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
    info "Arrêt des services..."
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

# Vérifier si les services sont déjà en cours d'exécution
if [ -f "$PID_FILE" ]; then
    warn "Des services semblent déjà être en cours d'exécution."
    read -p "Voulez-vous les arrêter et redémarrer? (y/N) " -n 1 -r
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
export MONGO_DB_NAME="${MONGO_DB_NAME:-bureau}"
export TREE_SERVICE_PORT="${TREE_SERVICE_PORT:-8082}"
export TREE_SERVICE_URL="${TREE_SERVICE_URL:-http://localhost:8082}"
export GATEWAY_PORT="${GATEWAY_PORT:-8080}"

info "Configuration:"
info "  MongoDB URI: $MONGO_URI"
info "  MongoDB DB: $MONGO_DB_NAME"
info "  Tree Service Port: $TREE_SERVICE_PORT"
info "  Gateway Port: $GATEWAY_PORT"
info "  Tree Service URL: $TREE_SERVICE_URL"

# Vérifier que les dossiers existent
if [ ! -d "$TREE_SERVICE_DIR" ]; then
    error "Le dossier Tree Service n'existe pas: $TREE_SERVICE_DIR"
    exit 1
fi

if [ ! -d "$GATEWAY_DIR" ]; then
    error "Le dossier Gateway n'existe pas: $GATEWAY_DIR"
    exit 1
fi

# Vider le fichier PID
> "$PID_FILE"

# Lancer le Tree Service
info "Démarrage du Tree Service..."
cd "$TREE_SERVICE_DIR"
if [ "$1" == "--background" ]; then
    go run main.go > "$LOG_DIR/tree-service.log" 2>&1 &
    TREE_PID=$!
    echo "$TREE_PID" >> "$PID_FILE"
    info "Tree Service démarré (PID: $TREE_PID, logs: $LOG_DIR/tree-service.log)"
else
    go run main.go > "$LOG_DIR/tree-service.log" 2>&1 &
    TREE_PID=$!
    echo "$TREE_PID" >> "$PID_FILE"
    info "Tree Service démarré (PID: $TREE_PID)"
fi

# Attendre que le Tree Service soit prêt
info "Attente du Tree Service..."
sleep 3
for i in {1..10}; do
    if curl -s "http://localhost:$TREE_SERVICE_PORT/health" > /dev/null 2>&1 || \
       curl -s "http://localhost:$TREE_SERVICE_PORT/api/v1/tree" > /dev/null 2>&1; then
        info "Tree Service est prêt!"
        break
    fi
    if [ $i -eq 10 ]; then
        warn "Tree Service ne répond pas encore, mais on continue..."
    else
        sleep 1
    fi
done

# Lancer le Gateway
info "Démarrage du Gateway..."
cd "$GATEWAY_DIR"
if [ "$1" == "--background" ]; then
    go run main.go > "$LOG_DIR/gateway.log" 2>&1 &
    GATEWAY_PID=$!
    echo "$GATEWAY_PID" >> "$PID_FILE"
    info "Gateway démarré (PID: $GATEWAY_PID, logs: $LOG_DIR/gateway.log)"
else
    go run main.go > "$LOG_DIR/gateway.log" 2>&1 &
    GATEWAY_PID=$!
    echo "$GATEWAY_PID" >> "$PID_FILE"
    info "Gateway démarré (PID: $GATEWAY_PID)"
fi

# Attendre que le Gateway soit prêt
info "Attente du Gateway..."
sleep 3
for i in {1..10}; do
    if curl -s "http://localhost:$GATEWAY_PORT/query" > /dev/null 2>&1 || \
       curl -s "http://localhost:$GATEWAY_PORT/" > /dev/null 2>&1; then
        info "Gateway est prêt!"
        break
    fi
    if [ $i -eq 10 ]; then
        warn "Gateway ne répond pas encore, mais on continue..."
    else
        sleep 1
    fi
done

info ""
info "=========================================="
info "✅ Services démarrés avec succès!"
info "=========================================="
info ""
info "Tree Service: http://localhost:$TREE_SERVICE_PORT"
info "Gateway GraphQL: http://localhost:$GATEWAY_PORT"
info "GraphQL Playground: http://localhost:$GATEWAY_PORT/"
info ""
info "Logs:"
info "  Tree Service: $LOG_DIR/tree-service.log"
info "  Gateway: $LOG_DIR/gateway.log"
info ""
info "Pour arrêter les services, utilisez: ./stop.sh"
info ""

# Si en mode background, on sort
if [ "$1" == "--background" ]; then
    info "Services lancés en arrière-plan."
    exit 0
fi

# Sinon, on attend et affiche les logs
info "Appuyez sur Ctrl+C pour arrêter les services..."
wait


