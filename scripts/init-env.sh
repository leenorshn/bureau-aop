#!/bin/bash

# Script pour initialiser .env.cloudrun avec le projet bureaumlmg
# Usage: ./scripts/init-env.sh

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

info "🔧 Initialisation de .env.cloudrun pour bureaumlmg"
echo ""

# Vérifier si .env.cloudrun existe déjà
if [ -f .env.cloudrun ]; then
    warn "⚠️  .env.cloudrun existe déjà"
    question "Voulez-vous le recréer? (y/n)"
    read -r -n 1 RECREATE
    echo
    if [[ ! $RECREATE =~ ^[Yy]$ ]]; then
        info "Annulation..."
        exit 0
    fi
fi

# Demander l'URI MongoDB
echo ""
info "Configuration MongoDB Atlas"
info "Assurez-vous que 0.0.0.0/0 est whitelisted dans Network Access"
echo ""
question "URI MongoDB (format: mongodb+srv://user:pass@cluster.mongodb.net/bureau):"
read -r MONGO_URI

if [ -z "$MONGO_URI" ]; then
    error "URI MongoDB requise"
    exit 1
fi

# Demander le nom de la base de données
question "Nom de la base de données (défaut: bureau):"
read -r MONGO_DB_NAME
MONGO_DB_NAME=${MONGO_DB_NAME:-bureau}

# Demander la région
info ""
info "Régions Cloud Run disponibles:"
info "  us-central1 (Iowa) - Recommandé"
info "  us-east1 (Caroline du Sud)"
info "  europe-west1 (Belgique)"
info "  asia-northeast1 (Tokyo)"
echo ""
question "Région (défaut: us-central1):"
read -r REGION
REGION=${REGION:-us-central1}

# Redis (optionnel)
question "URL Redis (optionnel, appuyez sur Entrée pour ignorer):"
read -r REDIS_URL

# Créer le fichier .env.cloudrun
cat > .env.cloudrun << EOF
# Configuration Google Cloud Run pour Bureau MLM
# Projet: bureaumlmg

# Google Cloud Project Configuration
export GCP_PROJECT_ID="bureaumlmg"
export GCP_REGION="$REGION"

# MongoDB Configuration
export MONGO_URI="$MONGO_URI"
export MONGO_DB_NAME="$MONGO_DB_NAME"

# Redis Configuration (Optional)
export REDIS_URL="$REDIS_URL"

# Instructions:
# 1. Chargez ce fichier: source .env.cloudrun
# 2. Déployez: ./scripts/deploy-cloudrun.sh
EOF

info "✅ Fichier .env.cloudrun créé avec succès!"
echo ""
info "Configuration:"
info "  Project: bureaumlmg"
info "  Region: $REGION"
info "  Database: $MONGO_DB_NAME"
echo ""
info "Prochaines étapes:"
info "  1. Chargez les variables: source .env.cloudrun"
info "  2. Déployez: ./scripts/deploy-cloudrun.sh"
echo ""
warn "⚠️  Important: Ne commitez pas .env.cloudrun (contient des secrets)"
echo ""




