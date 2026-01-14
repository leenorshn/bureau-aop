# ⚡ Déploiement Ultra-Rapide

Déployez Bureau MLM sur Cloud Run en **3 commandes**.

## Prérequis

- ✅ MongoDB Atlas configuré (whitelist 0.0.0.0/0)
- ✅ gcloud CLI installé (`curl https://sdk.cloud.google.com | bash`)
- ✅ Connecté à Google Cloud (`gcloud auth login`)

## 🚀 Déploiement en 3 Commandes

### 1. Configurer MongoDB URI

```bash
# Éditer .env.cloudrun
nano .env.cloudrun

# Remplacer cette ligne avec votre vraie URI MongoDB:
# export MONGO_URI="mongodb+srv://user:password@cluster.mongodb.net/bureau"
```

### 2. Charger les Variables

```bash
source .env.cloudrun
```

### 3. Déployer

```bash
./scripts/deploy-cloudrun.sh
```

**C'est tout! 🎉**

## 🧪 Tester

```bash
# Récupérer l'URL
GATEWAY_URL=$(gcloud run services describe gateway --region us-central1 --format 'value(status.url)')

# Ouvrir dans le navigateur
open $GATEWAY_URL
```

## 📊 Voir les Services

```bash
gcloud run services list --project bureaumlmg
```

## 🔄 Mettre à Jour

```bash
source .env.cloudrun
./scripts/deploy-cloudrun.sh
```

---

**Projet: bureaumlmg | Région: us-central1**




