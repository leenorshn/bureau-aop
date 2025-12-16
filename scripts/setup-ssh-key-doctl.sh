#!/bin/bash

# Script pour configurer une clé SSH sur DigitalOcean avec doctl
# Usage: ./scripts/setup-ssh-key-doctl.sh

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

info "🔑 Configuration de la clé SSH sur DigitalOcean"
echo ""

# Vérifier que doctl est installé
if ! command -v doctl &> /dev/null; then
    error "doctl n'est pas installé"
    info "Installez-le avec: ./scripts/setup-doctl.sh"
    exit 1
fi

# Vérifier l'authentification
if ! doctl account get &> /dev/null; then
    error "doctl n'est pas authentifié"
    info "Authentifiez-vous avec: doctl auth init"
    exit 1
fi

# Vérifier les clés existantes
EXISTING_KEYS=$(doctl compute ssh-key list --format ID,Name --no-header)
if [ -n "$EXISTING_KEYS" ]; then
    info "Clés SSH existantes sur DigitalOcean:"
    echo "$EXISTING_KEYS"
    echo ""
    question "Voulez-vous en ajouter une nouvelle? (y/n)"
    read -r -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        info "Configuration terminée. Utilisation des clés existantes."
        exit 0
    fi
fi

# Vérifier si une clé SSH locale existe
SSH_KEY_PUB="$HOME/.ssh/id_rsa.pub"
SSH_KEY_PRIV="$HOME/.ssh/id_rsa"

if [ ! -f "$SSH_KEY_PUB" ]; then
    warn "Aucune clé SSH publique trouvée à $SSH_KEY_PUB"
    question "Voulez-vous créer une nouvelle clé SSH? (y/n)"
    read -r -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        question "Entrez votre email pour la clé SSH:"
        read -r EMAIL
        if [ -z "$EMAIL" ]; then
            EMAIL="bureau@example.com"
        fi
        
        info "Génération de la clé SSH..."
        ssh-keygen -t rsa -b 4096 -C "$EMAIL" -f "$SSH_KEY_PRIV" -N ""
        info "✅ Clé SSH créée"
    else
        error "Une clé SSH est nécessaire"
        exit 1
    fi
fi

# Afficher la clé publique
info "Clé SSH publique:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat "$SSH_KEY_PUB"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Demander le nom de la clé
question "Quel nom voulez-vous donner à cette clé SSH? [bureau-ssh-key]"
read -r KEY_NAME
KEY_NAME="${KEY_NAME:-bureau-ssh-key}"

# Importer la clé
info "Import de la clé SSH sur DigitalOcean..."
if doctl compute ssh-key import "$KEY_NAME" --public-key-file "$SSH_KEY_PUB"; then
    info "✅ Clé SSH importée avec succès!"
    
    # Afficher les clés
    echo ""
    info "Clés SSH sur DigitalOcean:"
    doctl compute ssh-key list --format ID,Name,FingerPrint
    
    echo ""
    info "✅ Configuration terminée!"
    info "Vous pouvez maintenant utiliser: ./scripts/deploy-doctl.sh"
else
    error "❌ Échec de l'import de la clé SSH"
    exit 1
fi

