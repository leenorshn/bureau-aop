#!/bin/bash

# Script de test de connexion SSH
# Usage: ./scripts/test-ssh-connection.sh

set -e

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Configuration
DROPLET_HOST="${DROPLET_HOST:-}"
DROPLET_USER="${DROPLET_USER:-root}"

if [ -z "$DROPLET_HOST" ]; then
    error "DROPLET_HOST n'est pas défini"
    echo ""
    echo "Usage:"
    echo "  export DROPLET_HOST=\"your-droplet-ip\""
    echo "  export DROPLET_USER=\"root\"  # ou 'deploy'"
    echo "  ./scripts/test-ssh-connection.sh"
    exit 1
fi

info "Test de connexion SSH..."
info "Host: $DROPLET_HOST"
info "User: $DROPLET_USER"
echo ""

# Test 1: Ping
info "Test 1: Ping du serveur..."
if ping -c 1 -W 2 "$DROPLET_HOST" > /dev/null 2>&1; then
    info "✅ Le serveur répond au ping"
else
    error "❌ Le serveur ne répond pas au ping"
    exit 1
fi

# Test 2: Port SSH
info "Test 2: Vérification du port SSH (22)..."
if nc -z -w 2 "$DROPLET_HOST" 22 2>/dev/null; then
    info "✅ Le port 22 est ouvert"
else
    error "❌ Le port 22 n'est pas accessible"
    warn "Vérifiez le firewall DigitalOcean et UFW sur le serveur"
    exit 1
fi

# Test 3: Connexion SSH
info "Test 3: Connexion SSH..."
if ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$DROPLET_USER@$DROPLET_HOST" "echo 'Connexion SSH réussie'" 2>&1; then
    info "✅ Connexion SSH réussie!"
else
    error "❌ Échec de la connexion SSH"
    echo ""
    warn "Causes possibles:"
    echo "  1. La clé SSH n'est pas ajoutée à ssh-agent"
    echo "  2. L'utilisateur '$DROPLET_USER' n'existe pas sur le serveur"
    echo "  3. Le serveur SSH n'est pas configuré correctement"
    echo ""
    info "Solutions:"
    echo "  1. Ajouter la clé SSH: ssh-add ~/.ssh/id_rsa"
    echo "  2. Tester avec root: export DROPLET_USER=root"
    echo "  3. Se connecter manuellement: ssh $DROPLET_USER@$DROPLET_HOST"
    exit 1
fi

# Test 4: Docker
info "Test 4: Vérification de Docker..."
if ssh "$DROPLET_USER@$DROPLET_HOST" "command -v docker" > /dev/null 2>&1; then
    DOCKER_VERSION=$(ssh "$DROPLET_USER@$DROPLET_HOST" "docker --version" 2>/dev/null)
    info "✅ Docker installé: $DOCKER_VERSION"
else
    warn "⚠️  Docker n'est pas installé"
    warn "Exécutez: ./scripts/setup-digitalocean-droplet.sh"
fi

# Test 5: Docker Compose
info "Test 5: Vérification de Docker Compose..."
if ssh "$DROPLET_USER@$DROPLET_HOST" "docker compose version" > /dev/null 2>&1; then
    COMPOSE_VERSION=$(ssh "$DROPLET_USER@$DROPLET_HOST" "docker compose version" 2>/dev/null)
    info "✅ Docker Compose installé: $COMPOSE_VERSION"
else
    warn "⚠️  Docker Compose n'est pas installé"
    warn "Exécutez: ./scripts/setup-digitalocean-droplet.sh"
fi

echo ""
info "✅ Tous les tests sont passés!"
info "Vous pouvez maintenant déployer avec: ./scripts/deploy-digitalocean.sh"

