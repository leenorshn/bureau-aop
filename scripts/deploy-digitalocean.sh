#!/bin/bash

# Script de déploiement sur DigitalOcean Droplet
# Usage: ./scripts/deploy-digitalocean.sh [--production|--staging]

set -e

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
ENVIRONMENT="${1:---staging}"
DROPLET_HOST="${DROPLET_HOST:-}"
DROPLET_USER="${DROPLET_USER:-deploy}"
REMOTE_DIR="${REMOTE_DIR:-/home/deploy/bureau}"
COMPOSE_FILE="docker-compose.production.yml"

# Fonctions
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Vérifications
if [ -z "$DROPLET_HOST" ]; then
    error "DROPLET_HOST n'est pas défini. Exportez-le ou définissez-le dans .env"
    exit 1
fi

if [ ! -f ".env" ]; then
    error "Le fichier .env n'existe pas. Créez-le à partir de env.microservices.example"
    exit 1
fi

info "Déploiement sur DigitalOcean Droplet..."
info "Host: $DROPLET_HOST"
info "User: $DROPLET_USER"
info "Environment: $ENVIRONMENT"

# Vérifier la connexion SSH
info "Vérification de la connexion SSH..."
if ! ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$DROPLET_USER@$DROPLET_HOST" "echo 'Connection OK'" > /dev/null 2>&1; then
    error "Impossible de se connecter au droplet. Vérifiez votre configuration SSH."
    echo ""
    warn "Causes possibles:"
    echo "  1. La clé SSH n'est pas configurée ou ajoutée à ssh-agent"
    echo "  2. L'utilisateur '$DROPLET_USER' n'existe pas sur le serveur"
    echo "  3. Le firewall bloque la connexion SSH"
    echo ""
    info "Solutions:"
    echo "  1. Testez la connexion manuellement:"
    echo "     ssh $DROPLET_USER@$DROPLET_HOST"
    echo ""
    echo "  2. Utilisez le script de test:"
    echo "     ./scripts/test-ssh-connection.sh"
    echo ""
    echo "  3. Si l'utilisateur 'deploy' n'existe pas, utilisez 'root' d'abord:"
    echo "     export DROPLET_USER=root"
    echo "     ./scripts/setup-digitalocean-droplet.sh"
    echo ""
    echo "  4. Vérifiez que la clé SSH est ajoutée:"
    echo "     ssh-add ~/.ssh/id_rsa"
    echo ""
    exit 1
fi

# Créer le dossier distant si nécessaire
info "Création du dossier distant..."
ssh "$DROPLET_USER@$DROPLET_HOST" "mkdir -p $REMOTE_DIR"

# Copier les fichiers nécessaires
info "Copie des fichiers..."
rsync -avz --exclude '.git' \
    --exclude 'node_modules' \
    --exclude '.env' \
    --exclude 'logs' \
    --exclude '*.log' \
    ./ "$DROPLET_USER@$DROPLET_HOST:$REMOTE_DIR/"

# Copier le fichier .env séparément (sécurisé)
info "Copie du fichier .env..."
scp .env "$DROPLET_USER@$DROPLET_HOST:$REMOTE_DIR/.env"

# Déployer sur le droplet
info "Déploiement des services..."
ssh "$DROPLET_USER@$DROPLET_HOST" << EOF
    set -e
    cd $REMOTE_DIR
    
    # Vérifier que Docker est installé
    if ! command -v docker &> /dev/null; then
        echo "Installation de Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sh get-docker.sh
    fi
    
    # Vérifier que Docker Compose est installé
    if ! docker compose version &> /dev/null; then
        echo "Installation de Docker Compose..."
        apt-get update
        apt-get install -y docker-compose-plugin
    fi
    
    # Arrêter les anciens conteneurs
    echo "Arrêt des anciens conteneurs..."
    docker compose -f $COMPOSE_FILE down || true
    
    # Construire et lancer les nouveaux conteneurs
    echo "Construction et démarrage des services..."
    docker compose -f $COMPOSE_FILE pull || true
    docker compose -f $COMPOSE_FILE up -d --build
    
    # Nettoyer les images inutilisées
    echo "Nettoyage des images Docker..."
    docker image prune -f
    
    # Afficher les logs
    echo "Services démarrés. Logs:"
    docker compose -f $COMPOSE_FILE ps
EOF

info "✅ Déploiement terminé!"
info ""
info "Services disponibles:"
info "  - Gateway: http://$DROPLET_HOST:8080"
info ""
info "Pour voir les logs:"
info "  ssh $DROPLET_USER@$DROPLET_HOST 'cd $REMOTE_DIR && docker compose -f $COMPOSE_FILE logs -f'"

