# 🚀 Déploiement avec DigitalOcean CLI (doctl)

DigitalOcean propose plusieurs options pour déployer avec Docker Compose via leur CLI.

## 📦 Installation de doctl

### macOS
```bash
brew install doctl
```

### Linux
```bash
cd ~
wget https://github.com/digitalocean/doctl/releases/download/v1.104.0/doctl-1.104.0-linux-amd64.tar.gz
tar xf doctl-1.104.0-linux-amd64.tar.gz
sudo mv doctl /usr/local/bin
```

### Windows
```powershell
# Via Chocolatey
choco install doctl

# Ou télécharger depuis GitHub
# https://github.com/digitalocean/doctl/releases
```

## 🔐 Authentification

```bash
# Authentifier avec votre token DigitalOcean
doctl auth init

# Vérifier l'authentification
doctl account get
```

**Obtenir un token:**
1. Allez sur DigitalOcean > API > Tokens/Keys
2. Générer un nouveau token
3. Copiez-le et utilisez-le avec `doctl auth init`

## 🎯 Option 1: App Platform (Recommandé)

**Note importante:** App Platform ne supporte pas directement `docker-compose.yml`, mais peut déployer des services Docker individuels via un fichier de configuration YAML.

### Méthode A: Déploiement via fichier de configuration (app.yaml)

Créez `app.yaml`:

```yaml
name: bureau-mlm
region: nyc

services:
  - name: gateway
    github:
      repo: your-username/your-repo
      branch: main
      deploy_on_push: true
    dockerfile_path: gateway/Dockerfile
    http_port: 8080
    instance_count: 1
    instance_size_slug: basic-xxs
    envs:
      - key: TREE_SERVICE_URL
        value: ${tree-service.PUBLIC_URL}
      - key: GATEWAY_PORT
        value: "8080"
    routes:
      - path: /
    health_check:
      http_path: /query

  - name: tree-service
    github:
      repo: your-username/your-repo
      branch: main
      deploy_on_push: true
    dockerfile_path: services/tree-service/Dockerfile
    http_port: 8082
    instance_count: 1
    instance_size_slug: basic-xxs
    envs:
      - key: MONGO_URI
        value: ${db.DATABASE_URL}
        type: SECRET
      - key: MONGO_DB_NAME
        value: "bureau"
      - key: TREE_SERVICE_PORT
        value: "8082"
    health_check:
      http_path: /health

databases:
  - name: db
    engine: MONGODB
    version: "7"
    production: false
    cluster_name: bureau-db
    db_name: bureau
    db_user: bureau
```

Déployer:

```bash
doctl apps create --spec app.yaml
```

### Méthode B: Utiliser Defang (Alternative pour docker-compose)

**Defang** est un outil tiers qui permet de déployer docker-compose directement sur DigitalOcean:

```bash
# Installer Defang CLI
curl -fsSL https://raw.githubusercontent.com/DefangLabs/defang/main/install.sh | sh

# Déployer
defang compose up --provider=digitalocean
```

**Avantages:**
- ✅ Support complet de docker-compose
- ✅ Déploiement en une commande
- ✅ Gestion automatique des ressources

**Documentation:** https://docs.defang.io/docs/tutorials/deploy-to-digitalocean

## 🎯 Option 2: Droplet avec doctl + Docker Compose

Cette méthode crée un Droplet et y déploie docker-compose.

### Étape 1: Créer le Droplet

```bash
# Créer un Droplet
doctl compute droplet create bureau-droplet \
  --image ubuntu-22-04-x64 \
  --size s-2vcpu-2gb \
  --region nyc1 \
  --ssh-keys YOUR_SSH_KEY_ID \
  --wait

# Obtenir l'IP
DROPLET_IP=$(doctl compute droplet get bureau-droplet --format IPAddress --no-header)
echo "Droplet IP: $DROPLET_IP"
```

### Étape 2: Configurer et déployer

Utilisez notre script existant:

```bash
export DROPLET_HOST=$DROPLET_IP
export DROPLET_USER="root"
./scripts/setup-digitalocean-droplet.sh
./scripts/deploy-digitalocean.sh
```

## 🎯 Option 3: Script Automatique avec doctl

Nous avons créé un script qui automatise tout le processus.

## 📋 Comparaison des Options

| Option | Avantages | Inconvénients | Coût |
|--------|-----------|---------------|------|
| **App Platform** | ✅ Gestion automatique<br>✅ Scaling automatique<br>✅ HTTPS intégré<br>✅ CI/CD intégré | ❌ Moins de contrôle<br>❌ Limitations docker-compose | ~$10-20/mois |
| **Droplet + doctl** | ✅ Contrôle total<br>✅ Docker-compose complet<br>✅ Flexible | ❌ Gestion manuelle<br>❌ Pas de scaling auto | ~$12/mois |

## 🔧 Scripts Disponibles

1. **`scripts/setup-doctl.sh`** - Installer et configurer doctl
2. **`scripts/setup-ssh-key-doctl.sh`** - Configurer les clés SSH sur DigitalOcean
3. **`scripts/deploy-doctl.sh`** - Déployer avec doctl (App Platform ou Droplet)
4. **`scripts/deploy-defang.sh`** - Déployer docker-compose avec Defang

### Configuration Initiale

Avant de déployer, configurez vos clés SSH:

```bash
# Configurer les clés SSH sur DigitalOcean
./scripts/setup-ssh-key-doctl.sh
```

Ou manuellement:

```bash
# Importer votre clé SSH
doctl compute ssh-key import bureau-key --public-key-file ~/.ssh/id_rsa.pub

# Vérifier les clés
doctl compute ssh-key list
```

## 🎯 Option 4: Defang (Pour docker-compose natif)

Defang permet de déployer directement votre `docker-compose.yml` sur DigitalOcean:

```bash
# Installer Defang
./scripts/deploy-defang.sh  # Installe et déploie automatiquement

# Ou manuellement
defang compose up --provider=digitalocean
```

**Avantages:**
- ✅ Support complet de docker-compose
- ✅ Pas besoin de convertir en app.yaml
- ✅ Déploiement en une commande
- ✅ Gestion automatique des ressources

**Documentation:** https://docs.defang.io/docs/tutorials/deploy-to-digitalocean

## 📚 Ressources

- [doctl Documentation](https://docs.digitalocean.com/reference/doctl/)
- [App Platform Documentation](https://docs.digitalocean.com/products/app-platform/)
- [Docker Compose on App Platform](https://docs.digitalocean.com/products/app-platform/how-to/use-docker-compose/)

