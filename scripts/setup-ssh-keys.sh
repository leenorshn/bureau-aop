#!/bin/bash

# Script pour configurer les clés SSH pour DigitalOcean
# Usage: ./scripts/setup-ssh-keys.sh

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

info "Configuration des clés SSH pour DigitalOcean..."
echo ""

# Vérifier si une clé SSH existe
SSH_KEY_PATH="$HOME/.ssh/id_rsa"
SSH_PUB_KEY_PATH="$HOME/.ssh/id_rsa.pub"

if [ ! -f "$SSH_KEY_PATH" ]; then
    warn "Aucune clé SSH trouvée à $SSH_KEY_PATH"
    read -p "Voulez-vous créer une nouvelle clé SSH? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        info "Génération d'une nouvelle clé SSH..."
        ssh-keygen -t rsa -b 4096 -C "your-email@example.com" -f "$SSH_KEY_PATH"
        info "✅ Clé SSH créée"
    else
        error "Une clé SSH est nécessaire pour se connecter au droplet"
        exit 1
    fi
else
    info "✅ Clé SSH trouvée: $SSH_KEY_PATH"
fi

# Afficher la clé publique
if [ -f "$SSH_PUB_KEY_PATH" ]; then
    echo ""
    info "Votre clé publique SSH:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    cat "$SSH_PUB_KEY_PATH"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    warn "IMPORTANT: Ajoutez cette clé à votre Droplet DigitalOcean:"
    echo "  1. Allez sur DigitalOcean > Droplets > Votre Droplet > Settings > Security"
    echo "  2. Cliquez sur 'Add SSH Key'"
    echo "  3. Collez la clé ci-dessus"
    echo ""
fi

# Ajouter la clé à ssh-agent
info "Ajout de la clé à ssh-agent..."
if [ -z "$SSH_AUTH_SOCK" ]; then
    eval "$(ssh-agent -s)" > /dev/null
fi

if ssh-add -l | grep -q "$SSH_KEY_PATH" 2>/dev/null; then
    info "✅ La clé est déjà dans ssh-agent"
else
    ssh-add "$SSH_KEY_PATH" 2>/dev/null || {
        warn "Impossible d'ajouter automatiquement la clé à ssh-agent"
        info "Ajoutez-la manuellement: ssh-add $SSH_KEY_PATH"
    }
fi

# Tester la connexion si DROPLET_HOST est défini
if [ -n "$DROPLET_HOST" ]; then
    echo ""
    info "Test de connexion au droplet..."
    if ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "root@$DROPLET_HOST" "echo 'OK'" > /dev/null 2>&1; then
        info "✅ Connexion réussie!"
    else
        warn "⚠️  Impossible de se connecter. Vérifiez que:"
        echo "  1. La clé SSH est ajoutée au Droplet"
        echo "  2. Le Droplet est démarré"
        echo "  3. Le firewall DigitalOcean autorise le port 22"
    fi
fi

echo ""
info "✅ Configuration terminée!"


