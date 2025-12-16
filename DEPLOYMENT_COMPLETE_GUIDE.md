# 🚀 Guide Complet de Déploiement - DigitalOcean

Guide étape par étape pour déployer votre application sur DigitalOcean.

## ✅ Prérequis Complétés

Vous avez déjà exécuté:
- ✅ `./scripts/setup-ssh-key-doctl.sh` - Clés SSH configurées
- ✅ `./scripts/deploy-doctl.sh droplet` - Droplet créé

## 🔧 Étape 1: Corriger la Configuration SSH pour deploy

Si vous avez l'erreur "Permission denied (publickey)" avec deploy:

```bash
# Option A: Utiliser le script automatique (recommandé)
export DROPLET_HOST="165.227.84.113"
./scripts/fix-deploy-ssh.sh

# Option B: Manuellement
ssh root@165.227.84.113
bash /tmp/setup-deploy-user.sh  # Le script sera copié automatiquement
```

**Ou depuis votre machine locale:**

```bash
export DROPLET_HOST="165.227.84.113"
scp scripts/setup-deploy-user.sh root@$DROPLET_HOST:/tmp/
ssh root@$DROPLET_HOST "bash /tmp/setup-deploy-user.sh"
```

## 📋 Étape 2: Vérifier la Connexion

```bash
# Tester la connexion avec deploy
ssh deploy@165.227.84.113

# Si ça fonctionne, vous êtes prêt pour la suite!
```

## 🐳 Étape 3: Vérifier les Services Docker

```bash
# Se connecter au serveur
ssh deploy@165.227.84.113

# Vérifier que les services tournent
cd /home/deploy/bureau
docker compose -f docker-compose.production.yml ps

# Si les services ne tournent pas, les démarrer
docker compose -f docker-compose.production.yml up -d

# Vérifier les logs
docker compose -f docker-compose.production.yml logs -f gateway
```

## ⚙️ Étape 4: Configurer les Variables d'Environnement

```bash
# Sur votre machine locale
cp env.microservices.example .env
nano .env  # Configurez avec vos valeurs MongoDB
```

**Configuration minimale:**

```env
MONGO_URI=mongodb+srv://user:password@cluster.mongodb.net/bureau?retryWrites=true&w=majority
MONGO_DB_NAME=bureau
TREE_SERVICE_PORT=8082
TREE_SERVICE_URL=http://localhost:8082
GATEWAY_PORT=8080
REDIS_URL=
```

**Copier sur le serveur:**

```bash
# Depuis votre machine locale
scp .env deploy@165.227.84.113:/home/deploy/bureau/.env

# Redémarrer les services
ssh deploy@165.227.84.113 "cd /home/deploy/bureau && docker compose -f docker-compose.production.yml restart"
```

## 🌐 Étape 5: Configurer Nginx

```bash
# Se connecter au serveur
ssh deploy@165.227.84.113
sudo su

# Copier la configuration Nginx
cp /home/deploy/bureau/nginx/bureau.conf /etc/nginx/sites-available/bureau

# Éditer avec votre domaine
nano /etc/nginx/sites-available/bureau
# Remplacez "your-domain.com" par votre vrai domaine (2 fois)

# Activer la configuration
ln -s /etc/nginx/sites-available/bureau /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Tester
nginx -t

# Redémarrer
systemctl restart nginx
systemctl status nginx

exit  # Quitter root
```

## 🔒 Étape 6: Configurer SSL avec Let's Encrypt

```bash
# Sur le serveur
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Suivre les instructions:
# - Email
# - Accepter les termes
# - Rediriger HTTP vers HTTPS (option 2)

# Vérifier le renouvellement
sudo certbot renew --dry-run
```

## 🌍 Étape 7: Configurer le DNS

Sur votre fournisseur DNS:

```
Type: A
Name: @ (ou votre-domaine.com)
Value: 165.227.84.113
TTL: Auto

Type: A
Name: www
Value: 165.227.84.113
TTL: Auto
```

**Vérifier la propagation:**

```bash
dig your-domain.com +short
nslookup your-domain.com
```

## 🔥 Étape 8: Configurer le Firewall

```bash
# Sur le serveur
sudo ufw status
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
```

## ✅ Étape 9: Tester l'Application

```bash
# Depuis votre machine locale

# Test HTTP (redirection)
curl -I http://your-domain.com

# Test HTTPS
curl https://your-domain.com/query

# Test GraphQL
curl -X POST https://your-domain.com/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}'
```

## 📊 Étape 10: Monitoring et Logs

```bash
# Se connecter au serveur
ssh deploy@165.227.84.113

# Logs des services
cd /home/deploy/bureau
docker compose -f docker-compose.production.yml logs -f

# Logs Nginx
sudo tail -f /var/log/nginx/bureau-access.log
sudo tail -f /var/log/nginx/bureau-error.log

# Statut
docker compose -f docker-compose.production.yml ps
```

## 🔄 Commandes Utiles

```bash
# Redémarrer les services
ssh deploy@165.227.84.113 "cd /home/deploy/bureau && docker compose -f docker-compose.production.yml restart"

# Voir les logs
ssh deploy@165.227.84.113 "cd /home/deploy/bureau && docker compose -f docker-compose.production.yml logs -f gateway"

# Mettre à jour l'application
export DROPLET_HOST="165.227.84.113"
export DROPLET_USER="deploy"
./scripts/deploy-digitalocean.sh

# Vérifier le statut
ssh deploy@165.227.84.113 "cd /home/deploy/bureau && docker compose -f docker-compose.production.yml ps"
```

## 🎯 URLs Finales

- **GraphQL Endpoint**: `https://your-domain.com/query`
- **GraphQL Playground**: `https://your-domain.com/playground`
- **Health Check**: `https://your-domain.com/health`

## 📝 Checklist Finale

- [ ] Connexion SSH avec deploy fonctionne
- [ ] Services Docker en cours d'exécution
- [ ] Variables d'environnement configurées
- [ ] Nginx configuré et actif
- [ ] SSL/HTTPS configuré
- [ ] DNS configuré et propagé
- [ ] Firewall configuré
- [ ] Application accessible via HTTPS
- [ ] Tests GraphQL fonctionnels

## 🆘 Dépannage

### Problème: "Permission denied (publickey)" avec deploy

```bash
./scripts/fix-deploy-ssh.sh
```

### Problème: Services Docker ne démarrent pas

```bash
ssh deploy@165.227.84.113
cd /home/deploy/bureau
docker compose -f docker-compose.production.yml logs
docker compose -f docker-compose.production.yml up -d
```

### Problème: Nginx ne démarre pas

```bash
sudo nginx -t  # Vérifier la configuration
sudo systemctl status nginx
sudo journalctl -u nginx -f  # Voir les logs
```

---

**Votre application est maintenant déployée! 🎉**


