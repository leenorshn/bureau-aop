# Quick Start - Déploiement DigitalOcean

Guide rapide pour déployer votre application sur DigitalOcean en 10 minutes.

## 🚀 Option Rapide: Script Automatique

### 1. Créer un Droplet

1. Allez sur [DigitalOcean](https://www.digitalocean.com)
2. Créez un nouveau Droplet:
   - **Image**: Ubuntu 22.04 LTS
   - **Plan**: 2GB RAM / 1 vCPU ($12/mois)
   - **Region**: Choisissez la plus proche
   - **Authentication**: SSH Key

### 1.5. Configurer les Clés SSH (IMPORTANT)

**Avant de continuer**, configurez vos clés SSH:

```bash
# Configurer les clés SSH
./scripts/setup-ssh-keys.sh

# Suivez les instructions pour ajouter la clé au Droplet
# DigitalOcean > Droplets > Settings > Security > Add SSH Key
```

**Puis testez la connexion:**

```bash
export DROPLET_HOST="your-droplet-ip"
export DROPLET_USER="root"
./scripts/test-ssh-connection.sh
```

### 2. Configurer le Droplet

```bash
# Sur votre machine locale
export DROPLET_HOST="your-droplet-ip"
export DROPLET_USER="root"

# Tester la connexion d'abord
./scripts/test-ssh-connection.sh

# Si la connexion fonctionne, copier le script de setup
scp scripts/setup-digitalocean-droplet.sh root@$DROPLET_HOST:/tmp/

# Exécuter le script de setup
ssh root@$DROPLET_HOST "bash /tmp/setup-digitalocean-droplet.sh"
```

**Note:** Si vous avez des problèmes de connexion SSH, consultez `TROUBLESHOOTING_SSH.md`

### 3. Configurer les Variables d'Environnement

```bash
# Sur votre machine locale
cp env.microservices.example .env
nano .env  # Configurez avec vos valeurs
```

### 4. Déployer

```bash
# Sur votre machine locale
export DROPLET_HOST="your-droplet-ip"
export DROPLET_USER="deploy"
./scripts/deploy-digitalocean.sh
```

### 5. Configurer Nginx et SSL

```bash
# Sur le droplet
ssh deploy@your-droplet-ip

# Copier la configuration Nginx
sudo cp nginx/bureau.conf /etc/nginx/sites-available/bureau
sudo nano /etc/nginx/sites-available/bureau  # Modifiez your-domain.com

# Activer le site
sudo ln -s /etc/nginx/sites-available/bureau /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Obtenir le certificat SSL
sudo certbot --nginx -d your-domain.com
```

## 📋 Checklist Rapide

- [ ] Droplet créé (Ubuntu 22.04, 2GB RAM)
- [ ] Script de setup exécuté
- [ ] Fichier `.env` configuré
- [ ] Déploiement effectué
- [ ] Nginx configuré
- [ ] SSL configuré
- [ ] Test de l'application

## 🔧 Commandes Utiles

### Voir les logs
```bash
ssh deploy@your-droplet-ip
cd /home/deploy/bureau
docker compose -f docker-compose.production.yml logs -f
```

### Redémarrer les services
```bash
ssh deploy@your-droplet-ip
cd /home/deploy/bureau
docker compose -f docker-compose.production.yml restart
```

### Mettre à jour
```bash
./scripts/deploy-digitalocean.sh
```

## 🎯 URLs

- **GraphQL Playground**: https://your-domain.com/
- **GraphQL Endpoint**: https://your-domain.com/query
- **Health Check**: https://your-domain.com/health

## 💡 Astuces

1. **MongoDB**: Utilisez MongoDB Atlas (gratuit jusqu'à 512MB)
2. **Monitoring**: Activez DigitalOcean Monitoring
3. **Backups**: Configurez des snapshots automatiques
4. **Domain**: Utilisez Cloudflare pour DNS gratuit

## 📚 Documentation Complète

Voir `DIGITALOCEAN_DEPLOYMENT.md` pour plus de détails.

