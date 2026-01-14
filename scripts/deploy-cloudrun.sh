#!/bin/bash

# Script de déploiement sur Google Cloud Run
# Usage: ./scripts/deploy-cloudrun.sh [--prod|--staging]

set -e

# Couleurs pour les logs
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
ENVIRONMENT="${1:---staging}"
PROJECT_ID="${GCP_PROJECT_ID:-bureaumlmg}"
REGION="${GCP_REGION:-us-central1}"

# Vérifier que gcloud est installé
if ! command -v gcloud &> /dev/null; then
    error "gcloud CLI n'est pas installé"
    info "Installez-le: curl https://sdk.cloud.google.com | bash"
    exit 1
fi

# Vérifier les variables d'environnement requises
if [ -z "$PROJECT_ID" ]; then
    error "GCP_PROJECT_ID n'est pas défini"
    info "Exportez votre project ID: export GCP_PROJECT_ID=your-project-id"
    exit 1
fi

if [ -z "$MONGO_URI" ]; then
    error "MONGO_URI n'est pas défini"
    info "Exportez votre URI MongoDB: export MONGO_URI=mongodb+srv://..."
    exit 1
fi

info "🚀 Déploiement sur Google Cloud Run"
info "Project: $PROJECT_ID"
info "Region: $REGION"
info "Environment: $ENVIRONMENT"
echo ""

# Configurer le projet gcloud
gcloud config set project "$PROJECT_ID"

# Vérifier que les APIs nécessaires sont activées
info "Vérification des APIs..."
gcloud services enable run.googleapis.com cloudbuild.googleapis.com containerregistry.googleapis.com --quiet

# Étape 1: Build et déploiement du Tree Service
info "📦 Building Tree Service..."
cd services/tree-service

gcloud builds submit \
    --tag "gcr.io/$PROJECT_ID/tree-service:latest" \
    --dockerfile Dockerfile.cloudrun \
    --quiet

info "🚀 Déploiement du Tree Service..."
gcloud run deploy tree-service \
    --image "gcr.io/$PROJECT_ID/tree-service:latest" \
    --platform managed \
    --region "$REGION" \
    --allow-unauthenticated \
    --set-env-vars="MONGO_URI=$MONGO_URI,MONGO_DB_NAME=${MONGO_DB_NAME:-bureau},REDIS_URL=${REDIS_URL:-}" \
    --memory 512Mi \
    --cpu 1 \
    --min-instances 0 \
    --max-instances 10 \
    --timeout 300s \
    --port 8080 \
    --quiet

# Récupérer l'URL du Tree Service
TREE_SERVICE_URL=$(gcloud run services describe tree-service \
    --region "$REGION" \
    --format 'value(status.url)')

info "✅ Tree Service déployé: $TREE_SERVICE_URL"
cd ../..

# Étape 2: Build et déploiement du Gateway
info "📦 Building Gateway..."
cd gateway

gcloud builds submit \
    --tag "gcr.io/$PROJECT_ID/gateway:latest" \
    --dockerfile Dockerfile.cloudrun \
    --quiet

info "🚀 Déploiement du Gateway..."
gcloud run deploy gateway \
    --image "gcr.io/$PROJECT_ID/gateway:latest" \
    --platform managed \
    --region "$REGION" \
    --allow-unauthenticated \
    --set-env-vars="TREE_SERVICE_URL=$TREE_SERVICE_URL" \
    --memory 512Mi \
    --cpu 1 \
    --min-instances 1 \
    --max-instances 10 \
    --timeout 300s \
    --port 8080 \
    --quiet

# Récupérer l'URL du Gateway
GATEWAY_URL=$(gcloud run services describe gateway \
    --region "$REGION" \
    --format 'value(status.url)')

info "✅ Gateway déployé: $GATEWAY_URL"
cd ..

echo ""
info "🎉 Déploiement terminé avec succès!"
echo ""
info "📝 URLs des services:"
info "  Tree Service: $TREE_SERVICE_URL"
info "  Gateway (GraphQL): $GATEWAY_URL"
echo ""
info "🧪 Tester le déploiement:"
info "  GraphQL Playground: $GATEWAY_URL"
info "  Health Check: $TREE_SERVICE_URL/health"
echo ""
info "📊 Commandes utiles:"
info "  Voir les logs du Gateway:"
info "    gcloud run logs read gateway --region $REGION --limit 50"
info ""
info "  Voir les logs du Tree Service:"
info "    gcloud run logs read tree-service --region $REGION --limit 50"
echo ""

# Tester l'endpoint
info "🧪 Test de l'endpoint GraphQL..."
if curl -s -X POST "$GATEWAY_URL/query" \
    -H "Content-Type: application/json" \
    -d '{"query":"{ __typename }"}' | grep -q "Query"; then
    info "✅ GraphQL endpoint fonctionne!"
else
    warn "⚠️  GraphQL endpoint ne répond pas comme attendu"
fi

echo ""
info "💰 Estimation des coûts:"
info "  - Tree Service (min 0): ~\$0-5/mois"
info "  - Gateway (min 1): ~\$5-10/mois"
info "  - Total estimé: ~\$5-15/mois"
echo ""

