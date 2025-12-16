# 🚀 Déploiement DigitalOcean - Guide Complet

Ce guide vous accompagne pour déployer votre application MLM sur DigitalOcean.

## 📚 Documentation Disponible

1. **`DIGITALOCEAN_QUICKSTART.md`** - Guide rapide (10 minutes)
2. **`DIGITALOCEAN_DEPLOYMENT.md`** - Guide complet avec toutes les options
3. **`DIGITALOCEAN_SUMMARY.md`** - Résumé des fichiers et checklist

## 🎯 Options d'Architecture

### Option 1: Droplet avec Docker Compose ⭐ (Recommandé)
- **Coût**: ~$12/mois
- **Contrôle**: Total
- **Idéal pour**: MVP, petites/moyennes charges

### Option 2: App Platform
- **Coût**: ~$10/mois
- **Avantages**: Gestion automatique, scaling
- **Idéal pour**: Production avec scaling automatique

### Option 3: Kubernetes
- **Coût**: Variable
- **Avantages**: Haute disponibilité, scaling avancé
- **Idéal pour**: Charges élevées, multiples environnements

## 🚀 Options de Déploiement

### Option A: Avec Defang ⭐ Pour docker-compose natif

**Le plus simple pour docker-compose:**

```bash
# Déployer directement votre docker-compose.yml
./scripts/deploy-defang.sh
```

Defang supporte nativement docker-compose et déploie sur DigitalOcean.

### Option B: Avec DigitalOcean CLI (doctl)

**Pour App Platform (nécessite conversion en app.yaml):**

```bash
# 1. Installer et configurer doctl
./scripts/setup-doctl.sh

# 2. Déployer sur App Platform
./scripts/deploy-doctl.sh app-platform

# Ou déployer sur Droplet
./scripts/deploy-doctl.sh droplet
```

Voir `DIGITALOCEAN_CLI_DEPLOYMENT.md` pour plus de détails.

### Option B: Déploiement Manuel (Scripts SSH)

## 🚀 Démarrage Rapide (5 étapes)

### 1. Créer le Droplet
Sur DigitalOcean, créez un Droplet Ubuntu 22.04, 2GB RAM.

### 2. Configurer le Droplet
```bash
export DROPLET_HOST="your-droplet-ip"
scp scripts/setup-digitalocean-droplet.sh root@$DROPLET_HOST:/tmp/
ssh root@$DROPLET_HOST "bash /tmp/setup-digitalocean-droplet.sh"
```

### 3. Configurer les Variables
```bash
cp env.microservices.example .env
nano .env  # Configurez vos valeurs
```

### 4. Déployer
```bash
export DROPLET_HOST="your-droplet-ip"
export DROPLET_USER="deploy"
./scripts/deploy-digitalocean.sh
```

### 5. Configurer Nginx et SSL
```bash
ssh deploy@your-droplet-ip
sudo cp nginx/bureau.conf /etc/nginx/sites-available/bureau
sudo nano /etc/nginx/sites-available/bureau  # Modifiez your-domain.com
sudo ln -s /etc/nginx/sites-available/bureau /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl restart nginx
sudo certbot --nginx -d your-domain.com
```

## 📦 Fichiers Créés

### Scripts
- `scripts/setup-digitalocean-droplet.sh` - Configuration initiale
- `scripts/deploy-digitalocean.sh` - Déploiement automatique

### Configurations
- `docker-compose.production.yml` - Docker Compose production
- `nginx/bureau.conf` - Configuration Nginx avec SSL
- `env.microservices.example` - Template variables d'environnement

### CI/CD
- `.github/workflows/deploy-digitalocean.yml` - GitHub Actions

### Dockerfiles Optimisés
- `gateway/Dockerfile` - Optimisé pour production
- `services/tree-service/Dockerfile` - Optimisé pour production

## 🔒 Sécurité

- ✅ Firewall UFW configuré
- ✅ Fail2ban installé
- ✅ SSL/TLS avec Let's Encrypt
- ✅ Headers de sécurité Nginx
- ✅ Utilisateur non-root dans Docker
- ✅ Health checks configurés

## 🔄 CI/CD

Le workflow GitHub Actions déploie automatiquement:
- Sur push vers `main`
- Build automatique
- Déploiement sur Droplet
- Health check automatique

**Secrets GitHub requis:**
- `DROPLET_HOST`
- `DROPLET_USER`
- `DROPLET_SSH_KEY`

## 📊 Monitoring

### Logs
```bash
ssh deploy@your-droplet-ip
cd /home/deploy/bureau
docker compose -f docker-compose.production.yml logs -f
```

### Health Checks
- Gateway: `http://localhost:8080/query`
- Tree Service: `http://localhost:8082/health`

## 💰 Coûts Estimés

| Service | Coût/Mois |
|---------|-----------|
| Droplet 2GB | $12 |
| MongoDB Atlas M0 | Gratuit |
| Domain (optionnel) | $0-15 |
| **Total** | **~$12-27** |

## 🎯 Checklist

- [ ] Droplet créé
- [ ] Script de setup exécuté
- [ ] Variables d'environnement configurées
- [ ] Application déployée
- [ ] Nginx configuré
- [ ] SSL configuré
- [ ] Domain configuré
- [ ] Monitoring activé
- [ ] Backups configurés
- [ ] Tests effectués

## 🆘 Support

### Problèmes de Connexion SSH

Si vous obtenez une erreur de connexion SSH:

```bash
# Utiliser l'assistant interactif
./scripts/fix-ssh-connection.sh

# Ou tester manuellement
./scripts/test-ssh-connection.sh

# Consulter le guide de dépannage
cat TROUBLESHOOTING_SSH.md
```

### Autres Problèmes

1. Vérifiez les logs: `docker compose logs -f`
2. Vérifiez Nginx: `sudo nginx -t`
3. Vérifiez les services: `docker compose ps`
4. Consultez `DIGITALOCEAN_DEPLOYMENT.md` section Dépannage

## 📚 Ressources

- [DigitalOcean Documentation](https://docs.digitalocean.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)

---

**Prêt à déployer ?** Commencez par `DIGITALOCEAN_QUICKSTART.md` ! 🚀

