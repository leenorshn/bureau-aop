#!/bin/bash

# Script de déploiement avec DigitalOcean CLI (doctl)
# Usage: ./scripts/deploy-doctl.sh [app-platform|droplet]

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

# Vérifier que doctl est installé
if ! command -v doctl &> /dev/null; then
    error "doctl n'est pas installé"
    echo ""
    info "Installation:"
    echo "  macOS:   brew install doctl"
    echo "  Linux:   wget https://github.com/digitalocean/doctl/releases/download/v1.104.0/doctl-1.104.0-linux-amd64.tar.gz"
    echo "  Windows: choco install doctl"
    exit 1
fi

# Vérifier l'authentification
if ! doctl account get &> /dev/null; then
    error "doctl n'est pas authentifié"
    info "Exécutez: doctl auth init"
    exit 1
fi

DEPLOYMENT_TYPE="${1:-droplet}"

if [ "$DEPLOYMENT_TYPE" = "app-platform" ]; then
    info "🚀 Déploiement sur App Platform..."
    
    # Vérifier si app.yaml existe
    if [ ! -f "app.yaml" ]; then
        error "app.yaml n'existe pas"
        info "Créez app.yaml ou utilisez docker-compose.production.yml"
        exit 1
    fi
    
    # Créer ou mettre à jour l'app
    if doctl apps list --format ID,Spec.Name | grep -q "bureau-mlm"; then
        APP_ID=$(doctl apps list --format ID,Spec.Name | grep "bureau-mlm" | awk '{print $1}')
        info "Mise à jour de l'application existante..."
        doctl apps update "$APP_ID" --spec app.yaml
    else
        info "Création d'une nouvelle application..."
        doctl apps create --spec app.yaml
    fi
    
    info "✅ Déploiement sur App Platform terminé!"
    
elif [ "$DEPLOYMENT_TYPE" = "droplet" ]; then
    info "🚀 Déploiement sur Droplet avec doctl..."
    
    # Configuration
    DROPLET_NAME="${DROPLET_NAME:-bureau-droplet}"
    DROPLET_SIZE="${DROPLET_SIZE:-s-2vcpu-2gb}"
    DROPLET_REGION="${DROPLET_REGION:-nyc1}"
    DROPLET_IMAGE="ubuntu-22-04-x64"
    
    # Obtenir ou créer la clé SSH
    SSH_KEYS=$(doctl compute ssh-key list --format ID,Name --no-header 2>/dev/null | head -1 | awk '{print $1}')
    if [ -z "$SSH_KEYS" ]; then
        warn "Aucune clé SSH trouvée sur DigitalOcean"
        
        # Vérifier si une clé SSH locale existe
        if [ -f "$HOME/.ssh/id_rsa.pub" ]; then
            question "Voulez-vous importer votre clé SSH locale automatiquement? (y/n)"
            read -r -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                SSH_KEY_NAME="bureau-ssh-key-$(date +%s)"
                info "Import de la clé SSH: $SSH_KEY_NAME"
                if doctl compute ssh-key import "$SSH_KEY_NAME" --public-key-file "$HOME/.ssh/id_rsa.pub" 2>/dev/null; then
                    SSH_KEYS=$(doctl compute ssh-key list --format ID,Name --no-header | grep "$SSH_KEY_NAME" | awk '{print $1}')
                    info "✅ Clé SSH importée: $SSH_KEYS"
                else
                    error "Échec de l'import de la clé SSH"
                    info "Essayez manuellement: ./scripts/setup-ssh-key-doctl.sh"
                    exit 1
                fi
            else
                error "Une clé SSH est nécessaire pour créer le Droplet"
                info ""
                info "Options:"
                info "  1. Utiliser le script d'aide: ./scripts/setup-ssh-key-doctl.sh"
                info "  2. Ou manuellement: doctl compute ssh-key import <name> --public-key-file ~/.ssh/id_rsa.pub"
                exit 1
            fi
        else
            error "Aucune clé SSH locale trouvée"
            info ""
            info "Solutions:"
            info "  1. Utiliser le script d'aide: ./scripts/setup-ssh-key-doctl.sh"
            info "  2. Ou générer manuellement:"
            echo "     ssh-keygen -t rsa -b 4096 -C 'your-email@example.com'"
            echo "     doctl compute ssh-key import bureau-key --public-key-file ~/.ssh/id_rsa.pub"
            exit 1
        fi
    else
        SSH_KEY_NAME=$(doctl compute ssh-key list --format ID,Name --no-header | head -1 | awk '{print $2}')
        info "✅ Clé SSH trouvée: $SSH_KEY_NAME ($SSH_KEYS)"
    fi
    
    # Vérifier si le Droplet existe déjà
    if doctl compute droplet list --format Name | grep -q "^$DROPLET_NAME$"; then
        warn "Le Droplet '$DROPLET_NAME' existe déjà"
        question "Voulez-vous le supprimer et en créer un nouveau? (y/n)"
        read -r -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            info "Suppression de l'ancien Droplet..."
            doctl compute droplet delete "$DROPLET_NAME" --force
            sleep 5
        else
            info "Utilisation du Droplet existant..."
            DROPLET_IP=$(doctl compute droplet get "$DROPLET_NAME" --format IPAddress --no-header)
            info "IP du Droplet: $DROPLET_IP"
        fi
    fi
    
    # Créer le Droplet si nécessaire
    if [ -z "$DROPLET_IP" ]; then
        info "Création du Droplet..."
        doctl compute droplet create "$DROPLET_NAME" \
            --image "$DROPLET_IMAGE" \
            --size "$DROPLET_SIZE" \
            --region "$DROPLET_REGION" \
            --ssh-keys "$SSH_KEYS" \
            --wait
        
        DROPLET_IP=$(doctl compute droplet get "$DROPLET_NAME" --format IPAddress --no-header)
        info "✅ Droplet créé: $DROPLET_IP"
        
        # Attendre que le Droplet soit prêt
        info "Attente que le Droplet soit prêt..."
        sleep 30
    fi
    
    # Configurer le Droplet
    export DROPLET_HOST="$DROPLET_IP"
    export DROPLET_USER="root"
    
    info "Configuration du Droplet..."
    if [ -f "scripts/setup-digitalocean-droplet.sh" ]; then
        scp scripts/setup-digitalocean-droplet.sh root@"$DROPLET_IP":/tmp/
        ssh root@"$DROPLET_IP" "bash /tmp/setup-digitalocean-droplet.sh"
    else
        warn "Script de setup non trouvé, configuration manuelle nécessaire"
    fi
    
    # Déployer
    info "Déploiement de l'application..."
    if [ -f "scripts/deploy-digitalocean.sh" ]; then
        export DROPLET_USER="deploy"
        ./scripts/deploy-digitalocean.sh
    else
        warn "Script de déploiement non trouvé"
        info "Déployez manuellement avec:"
        echo "  export DROPLET_HOST=$DROPLET_IP"
        echo "  export DROPLET_USER=deploy"
        echo "  ./scripts/deploy-digitalocean.sh"
    fi
    
    info "✅ Déploiement terminé!"
    info ""
    info "Droplet IP: $DROPLET_IP"
    info "Gateway: http://$DROPLET_IP:8080"
    
else
    error "Type de déploiement invalide: $DEPLOYMENT_TYPE"
    echo ""
    info "Usage:"
    echo "  ./scripts/deploy-doctl.sh app-platform  # Déployer sur App Platform"
    echo "  ./scripts/deploy-doctl.sh droplet       # Déployer sur Droplet"
    exit 1
fi

