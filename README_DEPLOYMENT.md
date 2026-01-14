# 📚 Guide de Déploiement - Bureau MLM

Ce document centralise toutes les options de déploiement pour l'application Bureau MLM.

## 🎯 Options de Déploiement

### 1. Google Cloud Run (Recommandé)

**Avantages:**
- ✅ Déploiement ultra-rapide (< 10 minutes)
- ✅ Scalabilité automatique
- ✅ Tier gratuit généreux (2M requêtes/mois)
- ✅ Coût faible ($5-15/mois après tier gratuit)
- ✅ Pas de gestion d'infrastructure

**Documentation:**
- [Guide Complet](./CLOUD_RUN_DEPLOYMENT.md) - Tout savoir sur le déploiement
- [Démarrage Rapide](./QUICKSTART_CLOUD_RUN.md) - Déployer en 10 minutes

**Commandes:**
```bash
# Configuration initiale
./scripts/setup-cloudrun.sh

# Déploiement
source .env.cloudrun
./scripts/deploy-cloudrun.sh

# Supprimer les services
./scripts/destroy-cloudrun.sh
```

### 2. Docker Compose Local

**Pour le développement et les tests locaux**

**Commandes:**
```bash
# Démarrer
./start.sh

# Arrêter
./stop.sh

# Redémarrer
./restart.sh
```

**Fichiers:**
- `docker-compose.microservices.yml` - Configuration de développement
- `docker-compose.production.yml` - Configuration de production

### 3. Déploiement Manuel

Pour d'autres plateformes (AWS, Azure, etc.), utilisez les Dockerfiles fournis:

**Dockerfiles disponibles:**
- `services/tree-service/Dockerfile` - Version standard
- `services/tree-service/Dockerfile.cloudrun` - Optimisé Cloud Run
- `gateway/Dockerfile` - Version standard
- `gateway/Dockerfile.cloudrun` - Optimisé Cloud Run

## 🏗️ Architecture

```
┌─────────────────┐
│   Client Web    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐      ┌──────────────┐
│   Gateway       │─────▶│ Tree Service │
│   (GraphQL)     │      │              │
└────────┬────────┘      └──────┬───────┘
         │                      │
         └──────────┬───────────┘
                    ▼
            ┌───────────────┐
            │ MongoDB Atlas │
            └───────────────┘
```

## 📋 Prérequis Communs

### MongoDB

**Option 1: MongoDB Atlas (Recommandé)**
1. Créer un compte sur [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Créer un cluster M0 (gratuit)
3. Whitelist les IPs appropriées
4. Copier l'URI de connexion

**Option 2: MongoDB Self-Hosted**
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

### Variables d'Environnement

Tous les déploiements nécessitent ces variables:

```bash
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/bureau
MONGO_DB_NAME=bureau
REDIS_URL=  # Optionnel
```

## 🚀 Déploiement Rapide par Plateforme

### Google Cloud Run

```bash
./scripts/setup-cloudrun.sh
source .env.cloudrun
./scripts/deploy-cloudrun.sh
```

**Temps**: 10 minutes  
**Coût**: $5-15/mois  
**Difficulté**: ⭐ Facile

### Docker Compose (Local)

```bash
cp env.microservices.example .env
# Éditer .env avec vos valeurs
./start.sh
```

**Temps**: 5 minutes  
**Coût**: Gratuit  
**Difficulté**: ⭐ Facile

### AWS (ECS/Fargate)

```bash
# Build et push vers ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com

docker build -t tree-service -f services/tree-service/Dockerfile services/tree-service
docker tag tree-service:latest YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/tree-service:latest
docker push YOUR_ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/tree-service:latest

# Créer services ECS via console ou CLI
```

**Temps**: 30 minutes  
**Coût**: $15-30/mois  
**Difficulté**: ⭐⭐⭐ Moyen

### Azure (Container Instances)

```bash
# Login
az login

# Créer resource group
az group create --name bureau-rg --location eastus

# Créer container registry
az acr create --resource-group bureau-rg --name bureauacr --sku Basic

# Build et push
az acr build --registry bureauacr --image tree-service:latest services/tree-service

# Deploy container
az container create --resource-group bureau-rg --name tree-service \
  --image bureauacr.azurecr.io/tree-service:latest \
  --dns-name-label bureau-tree --ports 8080
```

**Temps**: 30 minutes  
**Coût**: $20-40/mois  
**Difficulté**: ⭐⭐⭐ Moyen

## 📊 Comparaison des Plateformes

| Plateforme | Coût/mois | Setup | Scalabilité | Maintenance |
|------------|-----------|-------|-------------|-------------|
| **Cloud Run** | $5-15 | 10 min | Auto | Faible |
| Local Docker | $0 | 5 min | Manuelle | Moyenne |
| AWS ECS | $15-30 | 30 min | Auto | Moyenne |
| Azure ACI | $20-40 | 30 min | Auto | Moyenne |
| Kubernetes | $50+ | 2h+ | Auto | Élevée |

## 🔧 Scripts Disponibles

### Développement Local
- `start.sh` - Démarrer les services
- `stop.sh` - Arrêter les services
- `restart.sh` - Redémarrer les services

### Cloud Run
- `scripts/setup-cloudrun.sh` - Configuration initiale
- `scripts/deploy-cloudrun.sh` - Déploiement
- `scripts/destroy-cloudrun.sh` - Suppression

## 📚 Documentation Complète

- [CLOUD_RUN_DEPLOYMENT.md](./CLOUD_RUN_DEPLOYMENT.md) - Guide complet Cloud Run
- [QUICKSTART_CLOUD_RUN.md](./QUICKSTART_CLOUD_RUN.md) - Démarrage rapide
- [README.md](./README.md) - Documentation générale du projet

## 🆘 Support

Pour obtenir de l'aide:
1. Consultez la documentation spécifique à votre plateforme
2. Vérifiez les logs avec les commandes appropriées
3. Consultez les sections Troubleshooting

## 🎯 Recommandations

**Pour débuter:**  
→ Utilisez **Google Cloud Run** (simple, rapide, économique)

**Pour le développement:**  
→ Utilisez **Docker Compose local**

**Pour la production à grande échelle:**  
→ Envisagez **Kubernetes** (GKE, EKS, AKS)

---

**Choisissez votre plateforme et déployez en quelques minutes! 🚀**




