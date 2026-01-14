#!/bin/bash

# Script de configuration initiale pour Google Cloud Run
# Usage: ./scripts/setup-cloudrun.sh

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

info "🔧 Configuration de Google Cloud Run"
echo ""

# Vérifier si gcloud est installé
if ! command -v gcloud &> /dev/null; then
    error "gcloud CLI n'est pas installé"
    info "Installation en cours..."
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        info "Téléchargement du SDK Google Cloud..."
        curl https://sdk.cloud.google.com | bash
        exec -l $SHELL
    else
        # Linux
        info "Téléchargement du SDK Google Cloud..."
        curl https://sdk.cloud.google.com | bash
        exec -l $SHELL
    fi
fi

info "✅ gcloud CLI est installé"
echo ""

# Se connecter à Google Cloud
question "Vous connecter à Google Cloud? (y/n)"
read -r -n 1 REPLY
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    gcloud auth login
fi

# Utiliser le projet existant
PROJECT_ID="bureaumlmg"

info "Projets Google Cloud disponibles:"
gcloud projects list --format="table(projectId,name,projectNumber)"
echo ""

info "✅ Utilisation du projet: $PROJECT_ID"
question "Voulez-vous utiliser un autre projet? (y/n)"
read -r -n 1 USE_OTHER
echo

if [[ $USE_OTHER =~ ^[Yy]$ ]]; then
    question "ID du projet à utiliser:"
    read -r PROJECT_ID
    info "✅ Projet sélectionné: $PROJECT_ID"
fi

# Configurer le projet
gcloud config set project "$PROJECT_ID"
info "✅ Projet configuré: $PROJECT_ID"
echo ""

# Vérifier la facturation
warn "⚠️  Assurez-vous que la facturation est activée pour ce projet"
info "Ouvrez: https://console.cloud.google.com/billing/linkedaccount?project=$PROJECT_ID"
question "La facturation est-elle activée? (y/n)"
read -r -n 1 BILLING
echo

if [[ ! $BILLING =~ ^[Yy]$ ]]; then
    error "Activez la facturation avant de continuer"
    exit 1
fi

# Activer les APIs nécessaires
info "Activation des APIs Google Cloud..."
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable containerregistry.googleapis.com
info "✅ APIs activées"
echo ""

# Configurer la région
info "Régions Cloud Run populaires:"
info "  us-central1 (Iowa)"
info "  us-east1 (Caroline du Sud)"
info "  europe-west1 (Belgique)"
info "  asia-northeast1 (Tokyo)"
echo ""
question "Région à utiliser (défaut: us-central1):"
read -r REGION
REGION=${REGION:-us-central1}

# Créer le fichier .env.cloudrun
info "Configuration des variables d'environnement..."
echo ""

question "URI MongoDB (ex: mongodb+srv://user:pass@cluster.mongodb.net/bureau):"
read -r MONGO_URI

question "Nom de la base de données (défaut: bureau):"
read -r MONGO_DB_NAME
MONGO_DB_NAME=${MONGO_DB_NAME:-bureau}

question "URL Redis (optionnel, appuyez sur Entrée pour ignorer):"
read -r REDIS_URL

# Créer le fichier de configuration
cat > .env.cloudrun << EOF
# Configuration Google Cloud Run
export GCP_PROJECT_ID="$PROJECT_ID"
export GCP_REGION="$REGION"

# MongoDB Configuration
export MONGO_URI="$MONGO_URI"
export MONGO_DB_NAME="$MONGO_DB_NAME"
export REDIS_URL="$REDIS_URL"
EOF

info "✅ Configuration sauvegardée dans .env.cloudrun"
echo ""

info "📝 Configuration terminée!"
echo ""
info "Prochaines étapes:"
info "  1. Chargez les variables: source .env.cloudrun"
info "  2. Déployez l'application: ./scripts/deploy-cloudrun.sh"
echo ""
info "Commandes utiles:"
info "  Voir les services: gcloud run services list"
info "  Voir les logs: gcloud run logs read gateway --limit 50"
info "  Ouvrir la console: https://console.cloud.google.com/run?project=$PROJECT_ID"
echo ""

