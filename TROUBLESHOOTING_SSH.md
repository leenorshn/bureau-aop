# 🔧 Dépannage SSH - DigitalOcean

Guide pour résoudre les problèmes de connexion SSH lors du déploiement.

## ❌ Erreur: "Impossible de se connecter au droplet"

### Étape 1: Vérifier la connexion de base

```bash
# Testez la connexion manuellement
ssh root@64.227.180.21
```

**Si ça ne fonctionne pas**, continuez avec les étapes ci-dessous.

### Étape 2: Vérifier que le Droplet est accessible

```bash
# Ping du serveur
ping -c 3 64.227.180.21

# Vérifier le port SSH
nc -z -w 2 64.227.180.21 22
```

**Si le ping échoue:**
- Vérifiez que le Droplet est démarré sur DigitalOcean
- Vérifiez votre connexion internet

**Si le port 22 n'est pas accessible:**
- Vérifiez le firewall DigitalOcean (Settings > Networking > Firewalls)
- Vérifiez que le port 22 est ouvert

### Étape 3: Vérifier les clés SSH

```bash
# Utiliser le script de configuration
./scripts/setup-ssh-keys.sh
```

**Ou manuellement:**

```bash
# Vérifier si vous avez une clé SSH
ls -la ~/.ssh/id_rsa

# Si pas de clé, créez-en une
ssh-keygen -t rsa -b 4096 -C "your-email@example.com"

# Afficher votre clé publique
cat ~/.ssh/id_rsa.pub
```

**Ajouter la clé au Droplet:**
1. Allez sur DigitalOcean > Droplets > Votre Droplet
2. Settings > Security
3. Cliquez sur "Add SSH Key"
4. Collez votre clé publique (`cat ~/.ssh/id_rsa.pub`)

### Étape 4: Ajouter la clé à ssh-agent

```bash
# Démarrer ssh-agent
eval "$(ssh-agent -s)"

# Ajouter votre clé
ssh-add ~/.ssh/id_rsa

# Vérifier que la clé est ajoutée
ssh-add -l
```

### Étape 5: Tester la connexion avec le script

```bash
# Utiliser le script de test
export DROPLET_HOST="64.227.180.21"
export DROPLET_USER="root"  # Commencez avec root
./scripts/test-ssh-connection.sh
```

### Étape 6: Si l'utilisateur 'deploy' n'existe pas

Si vous essayez de vous connecter avec `deploy` mais que cet utilisateur n'existe pas encore:

```bash
# 1. Connectez-vous d'abord avec root
export DROPLET_USER="root"
ssh root@64.227.180.21

# 2. Sur le serveur, exécutez le script de setup
# (depuis votre machine locale)
export DROPLET_HOST="64.227.180.21"
export DROPLET_USER="root"
scp scripts/setup-digitalocean-droplet.sh root@64.227.180.21:/tmp/
ssh root@64.227.180.21 "bash /tmp/setup-digitalocean-droplet.sh"

# 3. Maintenant vous pouvez utiliser 'deploy'
export DROPLET_USER="deploy"
./scripts/test-ssh-connection.sh
```

## 🔍 Diagnostic Détaillé

### Vérifier la configuration SSH

```bash
# Voir la configuration SSH actuelle
ssh -v root@64.227.180.21

# Cela affichera des informations détaillées sur la connexion
```

### Vérifier les clés autorisées sur le serveur

```bash
# Se connecter au serveur
ssh root@64.227.180.21

# Vérifier les clés autorisées
cat ~/.ssh/authorized_keys
```

### Vérifier les permissions SSH

Sur le serveur, les permissions doivent être correctes:

```bash
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
chmod 644 ~/.ssh/authorized_keys  # Si vous avez plusieurs clés
```

## 🚨 Problèmes Courants

### Problème 1: "Permission denied (publickey)"

**Solution:**
1. Vérifiez que votre clé SSH est ajoutée au Droplet
2. Vérifiez que la clé est dans `~/.ssh/authorized_keys` sur le serveur
3. Vérifiez les permissions (voir ci-dessus)

### Problème 2: "Connection timed out"

**Solution:**
1. Vérifiez que le Droplet est démarré
2. Vérifiez le firewall DigitalOcean
3. Vérifiez que le port 22 est ouvert

### Problème 3: "Host key verification failed"

**Solution:**
```bash
# Supprimer l'ancienne clé du known_hosts
ssh-keygen -R 64.227.180.21

# Réessayer la connexion
ssh root@64.227.180.21
```

### Problème 4: L'utilisateur 'deploy' n'existe pas

**Solution:**
Exécutez d'abord le script de setup avec `root`:
```bash
export DROPLET_USER="root"
./scripts/setup-digitalocean-droplet.sh
```

## ✅ Checklist de Vérification

Avant de déployer, vérifiez:

- [ ] Le Droplet est démarré sur DigitalOcean
- [ ] Le ping fonctionne: `ping 64.227.180.21`
- [ ] Le port 22 est accessible: `nc -z 64.227.180.21 22`
- [ ] Vous avez une clé SSH: `ls ~/.ssh/id_rsa`
- [ ] La clé est ajoutée au Droplet (DigitalOcean > Settings > Security)
- [ ] La clé est dans ssh-agent: `ssh-add -l`
- [ ] La connexion SSH fonctionne: `ssh root@64.227.180.21`
- [ ] L'utilisateur 'deploy' existe (si vous l'utilisez)

## 🎯 Solution Rapide

Si vous voulez une solution rapide:

```bash
# 1. Configurer les clés SSH
./scripts/setup-ssh-keys.sh

# 2. Tester la connexion
export DROPLET_HOST="64.227.180.21"
export DROPLET_USER="root"
./scripts/test-ssh-connection.sh

# 3. Si root fonctionne, setup le serveur
export DROPLET_USER="root"
scp scripts/setup-digitalocean-droplet.sh root@$DROPLET_HOST:/tmp/
ssh root@$DROPLET_HOST "bash /tmp/setup-digitalocean-droplet.sh"

# 4. Maintenant utilisez deploy
export DROPLET_USER="deploy"
./scripts/deploy-digitalocean.sh
```

## 📞 Support

Si le problème persiste:

1. Vérifiez les logs SSH: `ssh -v root@64.227.180.21`
2. Vérifiez les logs du serveur: `journalctl -u ssh`
3. Contactez le support DigitalOcean si le problème vient de leur infrastructure

