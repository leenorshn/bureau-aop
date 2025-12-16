#!/bin/bash

# Script pour configurer un nouveau Droplet DigitalOcean
# Usage: ./scripts/setup-digitalocean-droplet.sh

set -e

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
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

# Vérifier que le script est exécuté en root ou avec sudo
if [ "$EUID" -ne 0 ]; then 
    error "Ce script doit être exécuté en tant que root ou avec sudo"
    exit 1
fi

info "Configuration du Droplet DigitalOcean pour le projet Bureau..."

# Mettre à jour le système
info "Mise à jour du système..."
apt update && apt upgrade -y

# Installer les outils de base
info "Installation des outils de base..."
apt install -y curl wget git ufw fail2ban

# Installer Docker
info "Installation de Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
else
    info "Docker est déjà installé"
fi

# Installer Docker Compose
info "Installation de Docker Compose..."
if ! docker compose version &> /dev/null; then
    apt install -y docker-compose-plugin
else
    info "Docker Compose est déjà installé"
fi

# Créer un utilisateur deploy
info "Création de l'utilisateur deploy..."
if ! id "deploy" &>/dev/null; then
    adduser --disabled-password --gecos "" deploy
    usermod -aG docker deploy
    usermod -aG sudo deploy
    info "Utilisateur 'deploy' créé"
else
    info "Utilisateur 'deploy' existe déjà"
fi

# Configurer les clés SSH pour deploy
info "Configuration des clés SSH pour deploy..."
mkdir -p /home/deploy/.ssh
chmod 700 /home/deploy/.ssh

# Copier les clés depuis root si elles existent
if [ -f /root/.ssh/authorized_keys ]; then
    cp /root/.ssh/authorized_keys /home/deploy/.ssh/authorized_keys
    chown -R deploy:deploy /home/deploy/.ssh
    chmod 600 /home/deploy/.ssh/authorized_keys
    info "✅ Clés SSH configurées pour deploy"
else
    warn "Aucune clé SSH trouvée dans /root/.ssh/authorized_keys"
    touch /home/deploy/.ssh/authorized_keys
    chown -R deploy:deploy /home/deploy/.ssh
    chmod 600 /home/deploy/.ssh/authorized_keys
    warn "Fichier authorized_keys créé (vide) - ajoutez votre clé SSH manuellement"
fi

# Configurer le firewall
info "Configuration du firewall (UFW)..."
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# Installer Nginx
info "Installation de Nginx..."
if ! command -v nginx &> /dev/null; then
    apt install -y nginx
    systemctl enable nginx
    systemctl start nginx
else
    info "Nginx est déjà installé"
fi

# Installer Certbot pour SSL
info "Installation de Certbot..."
if ! command -v certbot &> /dev/null; then
    apt install -y certbot python3-certbot-nginx
else
    info "Certbot est déjà installé"
fi

# Configurer fail2ban
info "Configuration de fail2ban..."
systemctl enable fail2ban
systemctl start fail2ban

# Créer le dossier pour l'application
info "Création du dossier de l'application..."
mkdir -p /home/deploy/bureau
chown deploy:deploy /home/deploy/bureau

# Configurer le swap (optionnel, pour les petits droplets)
info "Configuration du swap..."
if [ ! -f /swapfile ]; then
    fallocate -l 2G /swapfile
    chmod 600 /swapfile
    mkswap /swapfile
    swapon /swapfile
    echo '/swapfile none swap sw 0 0' | tee -a /etc/fstab
    info "Swap de 2GB créé"
else
    info "Swap existe déjà"
fi

# Optimiser les limites système
info "Optimisation des limites système..."
cat >> /etc/security/limits.conf << EOF
* soft nofile 65535
* hard nofile 65535
EOF

# Configurer sysctl pour de meilleures performances
info "Optimisation des paramètres réseau..."
cat >> /etc/sysctl.conf << EOF
# Optimisations réseau
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 2048
net.ipv4.ip_local_port_range = 10000 65535
EOF

sysctl -p

info ""
info "✅ Configuration terminée!"
info ""
info "Prochaines étapes:"
info "1. Connectez-vous en tant que 'deploy':"
info "   ssh deploy@your-droplet-ip"
info ""
info "2. Clonez votre repository:"
info "   git clone https://github.com/your-username/your-repo.git /home/deploy/bureau"
info ""
info "3. Configurez le fichier .env:"
info "   cd /home/deploy/bureau"
info "   cp env.microservices.example .env"
info "   nano .env"
info ""
info "4. Déployez avec Docker Compose:"
info "   docker compose -f docker-compose.production.yml up -d --build"
info ""
info "5. Configurez Nginx et SSL (voir DIGITALOCEAN_DEPLOYMENT.md)"

