# ⚡ Solution Rapide - Problème SSH

Vous avez l'erreur: **"Impossible de se connecter au droplet"** ?

## 🎯 Solution en 3 Étapes

### Étape 1: Utiliser l'Assistant Interactif

```bash
./scripts/fix-ssh-connection.sh
```

Cet assistant va:
- ✅ Vérifier la connectivité réseau
- ✅ Vérifier les clés SSH
- ✅ Configurer ssh-agent
- ✅ Tester la connexion
- ✅ Vous guider étape par étape

### Étape 2: Si l'Assistant Échoue

**Option A: Utiliser root au lieu de deploy**

```bash
export DROPLET_HOST="64.227.180.21"
export DROPLET_USER="root"  # Utilisez root d'abord
./scripts/test-ssh-connection.sh
```

**Option B: Configurer les clés SSH**

```bash
# 1. Configurer les clés SSH
./scripts/setup-ssh-keys.sh

# 2. Ajouter la clé au Droplet:
#    - DigitalOcean > Droplets > Settings > Security > Add SSH Key
#    - Collez la clé affichée

# 3. Tester
export DROPLET_HOST="64.227.180.21"
export DROPLET_USER="root"
./scripts/test-ssh-connection.sh
```

### Étape 3: Déployer

Une fois la connexion SSH fonctionnelle:

```bash
# Si vous utilisez root, setup le serveur d'abord
export DROPLET_USER="root"
scp scripts/setup-digitalocean-droplet.sh root@64.227.180.21:/tmp/
ssh root@64.227.180.21 "bash /tmp/setup-digitalocean-droplet.sh"

# Puis déployez avec deploy
export DROPLET_USER="deploy"
./scripts/deploy-digitalocean.sh
```

## 🔍 Diagnostic Rapide

```bash
# Test 1: Ping
ping -c 3 64.227.180.21

# Test 2: Port SSH
nc -z -w 2 64.227.180.21 22

# Test 3: Connexion SSH
ssh root@64.227.180.21
```

## 📚 Documentation Complète

Pour plus de détails, consultez:
- `TROUBLESHOOTING_SSH.md` - Guide complet de dépannage
- `DIGITALOCEAN_QUICKSTART.md` - Guide de démarrage rapide

