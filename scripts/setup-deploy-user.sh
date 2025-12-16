#!/bin/bash

# Script pour créer et configurer l'utilisateur deploy sur le serveur
# Usage: ./scripts/setup-deploy-user.sh
# Ou sur le serveur: bash setup-deploy-user.sh

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

# Vérifier que le script est exécuté en root
if [ "$EUID" -ne 0 ]; then 
    error "Ce script doit être exécuté en tant que root"
    info "Utilisez: sudo bash $0"
    info "Ou connectez-vous en root: ssh root@your-server"
    exit 1
fi

info "🔧 Configuration de l'utilisateur deploy..."
echo ""

# Vérifier si l'utilisateur deploy existe
if id "deploy" &>/dev/null; then
    info "✅ L'utilisateur 'deploy' existe déjà"
else
    info "Création de l'utilisateur 'deploy'..."
    adduser --disabled-password --gecos "" deploy
    usermod -aG docker deploy
    usermod -aG sudo deploy
    info "✅ Utilisateur 'deploy' créé"
fi

# Créer le dossier .ssh
info "Configuration du dossier .ssh..."
mkdir -p /home/deploy/.ssh
chmod 700 /home/deploy/.ssh

# Copier les clés autorisées depuis root
if [ -f /root/.ssh/authorized_keys ]; then
    info "Copie des clés SSH depuis root..."
    cp /root/.ssh/authorized_keys /home/deploy/.ssh/authorized_keys
    info "✅ Clés SSH copiées"
else
    warn "Aucune clé SSH trouvée dans /root/.ssh/authorized_keys"
    question "Voulez-vous créer un fichier authorized_keys vide? (y/n)"
    read -r -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        touch /home/deploy/.ssh/authorized_keys
        info "Fichier authorized_keys créé (vide)"
        warn "Vous devrez ajouter votre clé SSH manuellement"
    fi
fi

# Définir les permissions correctes
info "Configuration des permissions..."
chown -R deploy:deploy /home/deploy/.ssh
chmod 700 /home/deploy/.ssh
chmod 600 /home/deploy/.ssh/authorized_keys

# Vérifier que docker est installé et ajouter deploy au groupe
if command -v docker &> /dev/null; then
    if ! groups deploy | grep -q docker; then
        usermod -aG docker deploy
        info "✅ Utilisateur 'deploy' ajouté au groupe docker"
    else
        info "✅ Utilisateur 'deploy' est déjà dans le groupe docker"
    fi
else
    warn "Docker n'est pas installé"
fi

# Vérifier que sudo est configuré
if groups deploy | grep -q sudo; then
    info "✅ Utilisateur 'deploy' a les privilèges sudo"
else
    usermod -aG sudo deploy
    info "✅ Utilisateur 'deploy' ajouté au groupe sudo"
fi

# Afficher les clés configurées
echo ""
info "Clés SSH configurées pour deploy:"
if [ -f /home/deploy/.ssh/authorized_keys ] && [ -s /home/deploy/.ssh/authorized_keys ]; then
    cat /home/deploy/.ssh/authorized_keys
else
    warn "⚠️  Aucune clé SSH configurée"
    info "Ajoutez votre clé avec:"
    echo "  echo 'VOTRE_CLE_PUBLIQUE' >> /home/deploy/.ssh/authorized_keys"
fi

echo ""
info "✅ Configuration terminée!"
info ""
info "Vous pouvez maintenant vous connecter avec:"
info "  ssh deploy@$(hostname -I | awk '{print $1}')"
info ""
info "Ou depuis votre machine locale:"
info "  ssh deploy@VOTRE_IP_SERVEUR"


