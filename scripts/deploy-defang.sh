#!/bin/bash

# Script de déploiement avec Defang (docker-compose sur DigitalOcean)
# Usage: ./scripts/deploy-defang.sh

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

info "🚀 Déploiement avec Defang (docker-compose sur DigitalOcean)"
echo ""

# Vérifier que Defang est installé
if ! command -v defang &> /dev/null; then
    warn "Defang n'est pas installé"
    info "Installation de Defang..."
    
    # Détecter l'OS
    OS="$(uname -s)"
    case "$OS" in
        Darwin)
            if command -v brew &> /dev/null; then
                brew install defang-io/defang/defang
            else
                error "Homebrew n'est pas installé"
                info "Installez Defang manuellement: https://docs.defang.io/docs/getting-started/install"
                exit 1
            fi
            ;;
        Linux)
            curl -fsSL https://raw.githubusercontent.com/DefangLabs/defang/main/install.sh | sh
            ;;
        *)
            error "OS non supporté: $OS"
            info "Installez Defang manuellement: https://docs.defang.io/docs/getting-started/install"
            exit 1
            ;;
    esac
fi

# Vérifier l'authentification Defang
if ! defang whoami &> /dev/null; then
    warn "Defang n'est pas authentifié"
    info "Authentification Defang..."
    defang login
fi

# Vérifier que docker-compose.production.yml existe
COMPOSE_FILE="docker-compose.production.yml"
if [ ! -f "$COMPOSE_FILE" ]; then
    error "Fichier $COMPOSE_FILE non trouvé"
    exit 1
fi

info "Fichier docker-compose: $COMPOSE_FILE"
echo ""

# Demander confirmation
question "Voulez-vous déployer sur DigitalOcean? (y/n)"
read -r -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    info "Déploiement annulé"
    exit 0
fi

# Déployer
info "Déploiement en cours..."
defang compose up --provider=digitalocean --file "$COMPOSE_FILE"

info "✅ Déploiement terminé!"
info ""
info "Pour voir les logs:"
info "  defang compose logs"
info ""
info "Pour arrêter:"
info "  defang compose down"

