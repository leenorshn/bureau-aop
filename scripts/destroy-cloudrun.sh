#!/bin/bash

# Script pour supprimer les services Cloud Run
# Usage: ./scripts/destroy-cloudrun.sh

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
PROJECT_ID="${GCP_PROJECT_ID:-bureaumlmg}"
REGION="${GCP_REGION:-us-central1}"

if [ -z "$PROJECT_ID" ]; then
    error "GCP_PROJECT_ID n'est pas défini"
    info "Exportez: export GCP_PROJECT_ID=your-project-id"
    info "Ou chargez: source .env.cloudrun"
    exit 1
fi

warn "⚠️  ATTENTION: Cette action va supprimer tous les services Cloud Run!"
echo ""
info "Project: $PROJECT_ID"
info "Region: $REGION"
echo ""

question "Êtes-vous sûr de vouloir continuer? (tapez 'yes' pour confirmer)"
read -r CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    info "Annulation..."
    exit 0
fi

echo ""
info "🗑️  Suppression des services Cloud Run..."

# Configurer le projet
gcloud config set project "$PROJECT_ID"

# Supprimer le Gateway
if gcloud run services describe gateway --region "$REGION" &>/dev/null; then
    info "Suppression du Gateway..."
    gcloud run services delete gateway --region "$REGION" --quiet
    info "✅ Gateway supprimé"
else
    warn "Gateway n'existe pas"
fi

# Supprimer le Tree Service
if gcloud run services describe tree-service --region "$REGION" &>/dev/null; then
    info "Suppression du Tree Service..."
    gcloud run services delete tree-service --region "$REGION" --quiet
    info "✅ Tree Service supprimé"
else
    warn "Tree Service n'existe pas"
fi

echo ""
question "Voulez-vous aussi supprimer les images Docker? (y/n)"
read -r -n 1 DELETE_IMAGES
echo

if [[ $DELETE_IMAGES =~ ^[Yy]$ ]]; then
    info "Suppression des images Docker..."
    
    # Supprimer les images du Gateway
    if gcloud container images list --repository=gcr.io/$PROJECT_ID | grep -q gateway; then
        gcloud container images delete gcr.io/$PROJECT_ID/gateway --quiet --force-delete-tags || true
        info "✅ Images Gateway supprimées"
    fi
    
    # Supprimer les images du Tree Service
    if gcloud container images list --repository=gcr.io/$PROJECT_ID | grep -q tree-service; then
        gcloud container images delete gcr.io/$PROJECT_ID/tree-service --quiet --force-delete-tags || true
        info "✅ Images Tree Service supprimées"
    fi
fi

echo ""
info "✅ Nettoyage terminé!"
echo ""
info "Services restants:"
gcloud run services list --region "$REGION" || info "Aucun service Cloud Run"
echo ""
info "Pour réactiver les services, exécutez:"
info "  source .env.cloudrun"
info "  ./scripts/deploy-cloudrun.sh"
echo ""

