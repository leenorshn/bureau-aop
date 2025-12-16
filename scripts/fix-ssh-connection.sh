#!/bin/bash

# Script interactif pour résoudre les problèmes de connexion SSH
# Usage: ./scripts/fix-ssh-connection.sh

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

echo ""
info "🔧 Assistant de Résolution des Problèmes SSH"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Demander l'IP du Droplet
if [ -z "$DROPLET_HOST" ]; then
    question "Quelle est l'IP de votre Droplet DigitalOcean?"
    read -r DROPLET_HOST
    export DROPLET_HOST
else
    info "Droplet IP: $DROPLET_HOST"
fi

# Demander l'utilisateur
if [ -z "$DROPLET_USER" ]; then
    question "Quel utilisateur voulez-vous utiliser? (root/deploy) [root]"
    read -r DROPLET_USER
    DROPLET_USER="${DROPLET_USER:-root}"
    export DROPLET_USER
else
    info "Utilisateur: $DROPLET_USER"
fi

echo ""
info "Étape 1: Vérification de la connectivité réseau..."
if ping -c 1 -W 2 "$DROPLET_HOST" > /dev/null 2>&1; then
    info "✅ Le serveur répond au ping"
else
    error "❌ Le serveur ne répond pas au ping"
    warn "Vérifiez que:"
    echo "  - Le Droplet est démarré sur DigitalOcean"
    echo "  - Votre connexion internet fonctionne"
    exit 1
fi

echo ""
info "Étape 2: Vérification du port SSH..."
if nc -z -w 2 "$DROPLET_HOST" 22 2>/dev/null; then
    info "✅ Le port 22 est accessible"
else
    error "❌ Le port 22 n'est pas accessible"
    warn "Vérifiez le firewall DigitalOcean:"
    echo "  - DigitalOcean > Droplets > Settings > Networking > Firewalls"
    echo "  - Assurez-vous que le port 22 est ouvert"
    exit 1
fi

echo ""
info "Étape 3: Vérification des clés SSH..."
if [ ! -f "$HOME/.ssh/id_rsa" ]; then
    warn "⚠️  Aucune clé SSH trouvée"
    question "Voulez-vous créer une nouvelle clé SSH? (y/n)"
    read -r -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        ./scripts/setup-ssh-keys.sh
    else
        error "Une clé SSH est nécessaire. Exécutez: ./scripts/setup-ssh-keys.sh"
        exit 1
    fi
else
    info "✅ Clé SSH trouvée: $HOME/.ssh/id_rsa"
fi

echo ""
info "Étape 4: Vérification de ssh-agent..."
if [ -z "$SSH_AUTH_SOCK" ]; then
    warn "⚠️  ssh-agent n'est pas démarré"
    eval "$(ssh-agent -s)" > /dev/null
    info "✅ ssh-agent démarré"
fi

if ssh-add -l | grep -q "$HOME/.ssh/id_rsa" 2>/dev/null; then
    info "✅ La clé est dans ssh-agent"
else
    warn "⚠️  La clé n'est pas dans ssh-agent"
    ssh-add "$HOME/.ssh/id_rsa" 2>/dev/null || {
        error "Impossible d'ajouter la clé. Essayez manuellement: ssh-add ~/.ssh/id_rsa"
        exit 1
    }
    info "✅ Clé ajoutée à ssh-agent"
fi

echo ""
info "Étape 5: Affichage de la clé publique..."
echo ""
warn "IMPORTANT: Assurez-vous que cette clé est ajoutée au Droplet:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat "$HOME/.ssh/id_rsa.pub"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
question "La clé est-elle ajoutée au Droplet? (DigitalOcean > Settings > Security) (y/n)"
read -r -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    warn "Ajoutez la clé au Droplet avant de continuer:"
    echo "  1. Allez sur DigitalOcean > Droplets > Votre Droplet"
    echo "  2. Settings > Security"
    echo "  3. Add SSH Key"
    echo "  4. Collez la clé ci-dessus"
    exit 1
fi

echo ""
info "Étape 6: Test de connexion SSH..."
if ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$DROPLET_USER@$DROPLET_HOST" "echo 'OK'" > /dev/null 2>&1; then
    info "✅ Connexion SSH réussie!"
    echo ""
    info "🎉 Tout fonctionne! Vous pouvez maintenant déployer:"
    echo ""
    echo "  export DROPLET_HOST=\"$DROPLET_HOST\""
    echo "  export DROPLET_USER=\"$DROPLET_USER\""
    echo "  ./scripts/deploy-digitalocean.sh"
    echo ""
else
    error "❌ Échec de la connexion SSH"
    echo ""
    warn "Causes possibles:"
    echo "  1. La clé SSH n'est pas correctement ajoutée au Droplet"
    echo "  2. L'utilisateur '$DROPLET_USER' n'existe pas sur le serveur"
    echo "  3. Les permissions SSH sont incorrectes sur le serveur"
    echo ""
    info "Solutions:"
    echo "  1. Si l'utilisateur 'deploy' n'existe pas, utilisez 'root' d'abord:"
    echo "     export DROPLET_USER=root"
    echo "     ./scripts/setup-digitalocean-droplet.sh"
    echo ""
    echo "  2. Testez manuellement:"
    echo "     ssh $DROPLET_USER@$DROPLET_HOST"
    echo ""
    echo "  3. Consultez TROUBLESHOOTING_SSH.md pour plus de détails"
    exit 1
fi

