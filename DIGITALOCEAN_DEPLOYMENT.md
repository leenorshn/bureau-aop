# Guide de Déploiement sur DigitalOcean

Ce guide vous aidera à déployer votre application MLM sur DigitalOcean.

## 🏗️ Options d'Architecture

### Option 1: Droplets avec Docker Compose (Recommandé pour débuter)

**Avantages:**
- ✅ Contrôle total
- ✅ Coût prévisible
- ✅ Facile à gérer
- ✅ Idéal pour MVP et petites/moyennes charges

**Recommandations:**
- **Droplet**: 2GB RAM / 1 vCPU minimum (4GB RAM recommandé)
- **OS**: Ubuntu 22.04 LTS
- **Stockage**: 25GB SSD minimum

### Option 2: App Platform (Recommandé pour production)

**Avantages:**
- ✅ Gestion automatique
- ✅ Scaling automatique
- ✅ HTTPS intégré
- ✅ CI/CD intégré
- ✅ Monitoring intégré

**Recommandations:**
- **Plan**: Basic ($5/mois) pour commencer
- **Scaling**: Auto-scaling selon la charge

### Option 3: Kubernetes (Pour charges élevées)

**Avantages:**
- ✅ Haute disponibilité
- ✅ Scaling avancé
- ✅ Gestion de multiples environnements

**Recommandations:**
- **Cluster**: 3 nodes minimum
- **Node Size**: 2GB RAM / 1 vCPU minimum

## 📋 Prérequis

1. **Compte DigitalOcean**
2. **Droplet créé** (si Option 1)
3. **MongoDB Atlas** ou **DigitalOcean Managed MongoDB**
4. **Domain name** (optionnel mais recommandé)
5. **SSH Key** configurée sur DigitalOcean

## 🚀 Option 1: Déploiement sur Droplet

### Étape 1: Préparer le Droplet

```bash
# Se connecter au droplet
#64.227.180.21
ssh root@your-droplet-ip

# Mettre à jour le système
apt update && apt upgrade -y

# Installer Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Installer Docker Compose
apt install docker-compose-plugin -y

# Créer un utilisateur non-root
adduser deploy
usermod -aG docker deploy
```

### Étape 2: Cloner le projet

```bash
# Se connecter en tant que deploy
su - deploy

# Cloner le repository
git clone https://github.com/your-username/your-repo.git
cd your-repo
```

### Étape 3: Configurer les variables d'environnement

```bash
# Créer le fichier .env
cp env.microservices.example .env
nano .env
```

Configurer avec vos valeurs de production:
```env
MONGO_URI=mongodb+srv://user:password@cluster.mongodb.net/bureau?retryWrites=true&w=majority
MONGO_DB_NAME=bureau
TREE_SERVICE_PORT=8082
TREE_SERVICE_URL=http://localhost:8082
GATEWAY_PORT=8080
REDIS_URL=redis://localhost:6379
```

### Étape 4: Déployer avec Docker Compose

```bash
# Construire et lancer les services
docker compose -f docker-compose.production.yml up -d --build

# Vérifier les logs
docker compose -f docker-compose.production.yml logs -f
```

### Étape 5: Configurer Nginx (Reverse Proxy)

```bash
# Installer Nginx
sudo apt install nginx -y

# Configurer le reverse proxy
sudo nano /etc/nginx/sites-available/bureau
```

Configuration Nginx:
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # Gateway GraphQL
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Tree Service (optionnel, si vous voulez l'exposer)
    location /api/v1/tree/ {
        proxy_pass http://localhost:8082;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

```bash
# Activer le site
sudo ln -s /etc/nginx/sites-available/bureau /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Étape 6: Configurer SSL avec Let's Encrypt

```bash
# Installer Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtenir le certificat SSL
sudo certbot --nginx -d your-domain.com

# Vérifier le renouvellement automatique
sudo certbot renew --dry-run
```

## 🚀 Option 2: Déploiement sur App Platform

### Étape 1: Préparer le repository

Assurez-vous que votre code est sur GitHub/GitLab.

### Étape 2: Créer l'application

1. Allez sur DigitalOcean App Platform
2. Cliquez sur "Create App"
3. Connectez votre repository
4. Sélectionnez le dossier racine

### Étape 3: Configurer les services

#### Service 1: Gateway

- **Type**: Web Service
- **Build Command**: `cd gateway && go build -o gateway ./main.go`
- **Run Command**: `./gateway`
- **Port**: 8080
- **Environment Variables**:
  - `TREE_SERVICE_URL`: URL du Tree Service
  - `GATEWAY_PORT`: 8080

#### Service 2: Tree Service

- **Type**: Web Service
- **Build Command**: `cd services/tree-service && go build -o tree-service ./main.go`
- **Run Command**: `./tree-service`
- **Port**: 8082
- **Environment Variables**:
  - `MONGO_URI`: Votre URI MongoDB
  - `MONGO_DB_NAME`: bureau
  - `TREE_SERVICE_PORT`: 8082

### Étape 4: Configurer la base de données

Utilisez **DigitalOcean Managed MongoDB** ou **MongoDB Atlas**.

## 🔒 Sécurité

### 1. Firewall

```bash
# Configurer UFW
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 2. Secrets Management

Utilisez les **Secrets** d'App Platform ou **DigitalOcean Secrets** pour stocker:
- MongoDB URI
- JWT Secrets
- API Keys

### 3. Variables d'environnement sensibles

Ne jamais commiter:
- `.env` (ajouté à `.gitignore`)
- Secrets
- Clés privées

## 📊 Monitoring

### Option 1: DigitalOcean Monitoring

Activez le monitoring sur votre Droplet ou App Platform.

### Option 2: Logs

```bash
# Voir les logs Docker
docker compose logs -f gateway
docker compose logs -f tree-service

# Logs système
journalctl -u nginx -f
```

## 🔄 CI/CD

### GitHub Actions

Créez `.github/workflows/deploy.yml`:

```yaml
name: Deploy to DigitalOcean

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to Droplet
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.DROPLET_HOST }}
          username: ${{ secrets.DROPLET_USER }}
          key: ${{ secrets.DROPLET_SSH_KEY }}
          script: |
            cd /home/deploy/your-repo
            git pull
            docker compose -f docker-compose.production.yml up -d --build
```

## 💰 Estimation des Coûts

### Option 1: Droplet
- **Droplet 2GB**: $12/mois
- **MongoDB Atlas M0**: Gratuit (512MB)
- **Total**: ~$12/mois

### Option 2: App Platform
- **Basic Plan**: $5/mois par service
- **2 services**: $10/mois
- **MongoDB Atlas M0**: Gratuit
- **Total**: ~$10/mois

### Option 3: Managed MongoDB (DigitalOcean)
- **Droplet 2GB**: $12/mois
- **Managed MongoDB**: $15/mois
- **Total**: ~$27/mois

## 🎯 Checklist de Déploiement

- [ ] Droplet/App Platform créé
- [ ] MongoDB configuré (Atlas ou Managed)
- [ ] Variables d'environnement configurées
- [ ] Docker Compose configuré
- [ ] Nginx configuré (si Droplet)
- [ ] SSL configuré (Let's Encrypt)
- [ ] Firewall configuré
- [ ] Monitoring activé
- [ ] Backups configurés
- [ ] Domain name configuré
- [ ] Tests de charge effectués

## 📚 Ressources

- [DigitalOcean Documentation](https://docs.digitalocean.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)


