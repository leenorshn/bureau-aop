#!/bin/bash

# Script pour installer et configurer doctl
# Usage: ./scripts/setup-doctl.sh

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

info "🔧 Installation et configuration de doctl (DigitalOcean CLI)"
echo ""

# Détecter l'OS
OS="$(uname -s)"
ARCH="$(uname -m)"

# Vérifier si doctl est déjà installé
if command -v doctl &> /dev/null; then
    VERSION=$(doctl version --format Version --no-header)
    info "✅ doctl est déjà installé: $VERSION"
    question "Voulez-vous le réinstaller? (y/n)"
    read -r -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        info "Installation annulée"
        exit 0
    fi
fi

# Installation selon l'OS
case "$OS" in
    Darwin)
        info "Installation sur macOS..."
        if command -v brew &> /dev/null; then
            brew install doctl
        else
            error "Homebrew n'est pas installé"
            info "Installez Homebrew: https://brew.sh"
            exit 1
        fi
        ;;
    Linux)
        info "Installation sur Linux..."
        DOCTL_VERSION="1.104.0"
        cd /tmp
        wget "https://github.com/digitalocean/doctl/releases/download/v${DOCTL_VERSION}/doctl-${DOCTL_VERSION}-linux-amd64.tar.gz"
        tar xf "doctl-${DOCTL_VERSION}-linux-amd64.tar.gz"
        sudo mv doctl /usr/local/bin
        rm "doctl-${DOCTL_VERSION}-linux-amd64.tar.gz"
        ;;
    *)
        error "OS non supporté: $OS"
        info "Installez doctl manuellement: https://docs.digitalocean.com/reference/doctl/how-to/install/"
        exit 1
        ;;
esac

# Vérifier l'installation
if command -v doctl &> /dev/null; then
    VERSION=$(doctl version --format Version --no-header)
    info "✅ doctl installé avec succès: $VERSION"
else
    error "❌ Échec de l'installation"
    exit 1
fi

echo ""
info "🔐 Configuration de l'authentification..."
echo ""

# Vérifier si déjà authentifié
if doctl account get &> /dev/null; then
    info "✅ doctl est déjà authentifié"
    ACCOUNT=$(doctl account get --format Email --no-header)
    info "Compte: $ACCOUNT"
    question "Voulez-vous vous ré-authentifier? (y/n)"
    read -r -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        info "Configuration terminée!"
        exit 0
    fi
fi

# Authentification
info "Pour vous authentifier, vous avez besoin d'un token DigitalOcean:"
echo ""
warn "1. Allez sur: https://cloud.digitalocean.com/account/api/tokens"
warn "2. Cliquez sur 'Generate New Token'"
warn "3. Donnez-lui un nom (ex: 'doctl-cli')"
warn "4. Copiez le token généré"
echo ""
question "Avez-vous votre token? (y/n)"
read -r -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    warn "Ouvrez le lien ci-dessus et générez un token, puis relancez ce script"
    exit 0
fi

question "Collez votre token DigitalOcean:"
read -r -s TOKEN
echo

if [ -z "$TOKEN" ]; then
    error "Token vide"
    exit 1
fi

# Authentifier
doctl auth init --access-token "$TOKEN"

# Vérifier l'authentification
if doctl account get &> /dev/null; then
    ACCOUNT=$(doctl account get --format Email --no-header)
    info "✅ Authentification réussie!"
    info "Compte: $ACCOUNT"
else
    error "❌ Échec de l'authentification"
    exit 1
fi

echo ""
info "✅ Configuration terminée!"
echo ""
info "Commandes utiles:"
echo "  doctl account get              # Voir les infos du compte"
echo "  doctl compute droplet list      # Lister les Droplets"
echo "  doctl apps list                 # Lister les apps"
echo "  ./scripts/deploy-doctl.sh      # Déployer l'application"


