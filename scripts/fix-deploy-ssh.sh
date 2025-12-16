#!/bin/bash

# Script pour corriger la configuration SSH de l'utilisateur deploy
# Usage depuis votre machine locale: ./scripts/fix-deploy-ssh.sh

set -e

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
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

question() {
    echo -e "${BLUE}[?]${NC} $1"
}

# Configuration
DROPLET_HOST="${DROPLET_HOST:-}"
DROPLET_USER="${DROPLET_USER:-root}"

if [ -z "$DROPLET_HOST" ]; then
    question "Quelle est l'IP de votre Droplet?"
    read -r DROPLET_HOST
    if [ -z "$DROPLET_HOST" ]; then
        error "IP du Droplet requise"
        exit 1
    fi
fi

info "🔧 Correction de la configuration SSH pour l'utilisateur deploy..."
info "Serveur: $DROPLET_HOST"
info "Utilisateur: $DROPLET_USER"
echo ""

# Vérifier la connexion avec root
info "Vérification de la connexion SSH avec $DROPLET_USER..."
if ! ssh -o ConnectTimeout=5 "$DROPLET_USER@$DROPLET_HOST" "echo 'OK'" > /dev/null 2>&1; then
    error "Impossible de se connecter avec $DROPLET_USER@$DROPLET_HOST"
    info "Vérifiez votre connexion SSH"
    exit 1
fi

info "✅ Connexion SSH réussie"
echo ""

# Copier le script de setup sur le serveur
info "Copie du script de configuration sur le serveur..."
scp scripts/setup-deploy-user.sh "$DROPLET_USER@$DROPLET_HOST:/tmp/setup-deploy-user.sh"

# Exécuter le script sur le serveur
info "Exécution du script de configuration..."
ssh "$DROPLET_USER@$DROPLET_HOST" "bash /tmp/setup-deploy-user.sh"

# Vérifier que la connexion avec deploy fonctionne maintenant
echo ""
info "Vérification de la connexion avec l'utilisateur deploy..."
sleep 2

if ssh -o ConnectTimeout=5 "deploy@$DROPLET_HOST" "echo 'OK'" > /dev/null 2>&1; then
    info "✅ Connexion SSH avec deploy réussie!"
    echo ""
    info "Vous pouvez maintenant utiliser:"
    info "  ssh deploy@$DROPLET_HOST"
    info "  export DROPLET_USER=deploy"
    info "  ./scripts/deploy-digitalocean.sh"
else
    warn "⚠️  La connexion avec deploy ne fonctionne pas encore"
    echo ""
    info "Vérifiez manuellement:"
    info "  1. Se connecter en root: ssh root@$DROPLET_HOST"
    info "  2. Vérifier: cat /home/deploy/.ssh/authorized_keys"
    info "  3. Comparer avec votre clé: cat ~/.ssh/id_rsa.pub"
    echo ""
    info "Si les clés ne correspondent pas, ajoutez votre clé:"
    info "  ssh root@$DROPLET_HOST"
    info "  echo 'VOTRE_CLE_PUBLIQUE' >> /home/deploy/.ssh/authorized_keys"
    info "  chown deploy:deploy /home/deploy/.ssh/authorized_keys"
    info "  chmod 600 /home/deploy/.ssh/authorized_keys"
fi


