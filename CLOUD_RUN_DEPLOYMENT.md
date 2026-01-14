# 🚀 Déploiement sur Google Cloud Run

Guide complet pour déployer l'application Bureau MLM sur Google Cloud Run.

## 📋 Table des Matières

- [Prérequis](#prérequis)
- [Architecture](#architecture)
- [Configuration Initiale](#configuration-initiale)
- [Déploiement](#déploiement)
- [Gestion et Monitoring](#gestion-et-monitoring)
- [Coûts](#coûts)
- [Troubleshooting](#troubleshooting)

## 🎯 Prérequis

### 1. Compte Google Cloud

- Créer un compte sur [Google Cloud](https://cloud.google.com)
- Activer la facturation (carte bancaire requise, mais tier gratuit disponible)
- **Tier gratuit**: 2 millions de requêtes/mois

### 2. Outils Locaux

```bash
# Vérifier que Git est installé
git --version

# Vérifier que Docker est installé (optionnel)
docker --version
```

### 3. MongoDB

Utilisez **MongoDB Atlas** (gratuit jusqu'à 512MB):

1. Créer un compte sur [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Créer un cluster gratuit (M0)
3. Whitelist l'IP `0.0.0.0/0` (Cloud Run utilise des IPs dynamiques)
4. Créer un utilisateur database
5. Copier l'URI de connexion

## 🏗️ Architecture

```
┌─────────────────────┐
│   Client (Web)      │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Cloud Load Balancer│
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐      ┌──────────────────┐
│  Gateway Service    │─────▶│  Tree Service    │
│  (Cloud Run)        │      │  (Cloud Run)     │
│  Port: 8080         │      │  Port: 8080      │
└──────────┬──────────┘      └────────┬─────────┘
           │                          │
           └──────────┬───────────────┘
                      ▼
              ┌───────────────┐
              │ MongoDB Atlas │
              │   (External)  │
              └───────────────┘
```

### Services Déployés

1. **Tree Service** - Service de gestion de l'arbre MLM
   - URL: `https://tree-service-xxx-uc.a.run.app`
   - Scaling: 0-10 instances
   - Memory: 512Mi
   - CPU: 1

2. **Gateway** - API GraphQL Gateway
   - URL: `https://gateway-xxx-uc.a.run.app`
   - Scaling: 1-10 instances (min 1 pour éviter cold starts)
   - Memory: 512Mi
   - CPU: 1

## ⚙️ Configuration Initiale

### Étape 1: Configuration Automatique

```bash
# Lancer le script de configuration
./scripts/setup-cloudrun.sh
```

Ce script va:
- ✅ Installer/vérifier gcloud CLI
- ✅ Vous connecter à Google Cloud
- ✅ Créer ou sélectionner un projet
- ✅ Activer les APIs nécessaires
- ✅ Configurer la région
- ✅ Créer le fichier `.env.cloudrun`

### Étape 2: Configuration Manuelle (Alternative)

```bash
# 1. Installer gcloud CLI
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# 2. Se connecter
gcloud auth login

# 3. Créer un projet
gcloud projects create bureau-mlm-prod --name="Bureau MLM"

# 4. Configurer le projet
gcloud config set project bureau-mlm-prod

# 5. Activer les APIs
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable containerregistry.googleapis.com

# 6. Créer .env.cloudrun
cat > .env.cloudrun << EOF
export GCP_PROJECT_ID="bureau-mlm-prod"
export GCP_REGION="us-central1"
export MONGO_URI="mongodb+srv://user:pass@cluster.mongodb.net/bureau"
export MONGO_DB_NAME="bureau"
export REDIS_URL=""
EOF
```

## 🚀 Déploiement

### Déploiement Automatique

```bash
# 1. Charger les variables d'environnement
source .env.cloudrun

# 2. Déployer sur Cloud Run
./scripts/deploy-cloudrun.sh
```

### Déploiement Manuel

#### Tree Service

```bash
# Build et push l'image
cd services/tree-service
gcloud builds submit \
  --tag gcr.io/bureau-mlm-prod/tree-service:latest \
  --dockerfile Dockerfile.cloudrun

# Déployer
gcloud run deploy tree-service \
  --image gcr.io/bureau-mlm-prod/tree-service:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="MONGO_URI=$MONGO_URI,MONGO_DB_NAME=bureau" \
  --memory 512Mi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 10 \
  --port 8080

# Récupérer l'URL
TREE_SERVICE_URL=$(gcloud run services describe tree-service \
  --region us-central1 \
  --format 'value(status.url)')
echo "Tree Service: $TREE_SERVICE_URL"

cd ../..
```

#### Gateway

```bash
# Build et push l'image
cd gateway
gcloud builds submit \
  --tag gcr.io/bureau-mlm-prod/gateway:latest \
  --dockerfile Dockerfile.cloudrun

# Déployer
gcloud run deploy gateway \
  --image gcr.io/bureau-mlm-prod/gateway:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="TREE_SERVICE_URL=$TREE_SERVICE_URL" \
  --memory 512Mi \
  --cpu 1 \
  --min-instances 1 \
  --max-instances 10 \
  --port 8080

# Récupérer l'URL
GATEWAY_URL=$(gcloud run services describe gateway \
  --region us-central1 \
  --format 'value(status.url)')
echo "Gateway: $GATEWAY_URL"

cd ..
```

## 📊 Gestion et Monitoring

### Voir les Services

```bash
# Lister tous les services
gcloud run services list

# Détails d'un service
gcloud run services describe gateway --region us-central1
```

### Logs

```bash
# Logs du Gateway (temps réel)
gcloud run logs read gateway --region us-central1 --limit 50 --follow

# Logs du Tree Service
gcloud run logs read tree-service --region us-central1 --limit 50

# Filtrer par niveau
gcloud run logs read gateway --region us-central1 --log-filter="severity>=ERROR"
```

### Métriques

```bash
# Ouvrir la console Cloud Run
gcloud run services list --uri

# Ou directement dans la console
# https://console.cloud.google.com/run
```

### Mise à Jour

```bash
# Recharger les variables
source .env.cloudrun

# Redéployer
./scripts/deploy-cloudrun.sh
```

### Variables d'Environnement

```bash
# Mettre à jour une variable
gcloud run services update gateway \
  --region us-central1 \
  --set-env-vars="NEW_VAR=value"

# Supprimer une variable
gcloud run services update gateway \
  --region us-central1 \
  --remove-env-vars="VAR_NAME"
```

### Rollback

```bash
# Voir les révisions
gcloud run revisions list --service gateway --region us-central1

# Revenir à une révision précédente
gcloud run services update-traffic gateway \
  --region us-central1 \
  --to-revisions REVISION_NAME=100
```

## 💰 Coûts Estimés

### Tier Gratuit

- **Requêtes**: 2 millions/mois
- **CPU**: 180,000 vCPU-secondes/mois
- **Mémoire**: 360,000 GiB-secondes/mois
- **Réseau sortant**: 1 GB/mois

### Estimation Mensuelle (après tier gratuit)

**Scénario 1: Faible Trafic (< 100k requêtes/mois)**
- Tree Service (min 0): **$0-2/mois**
- Gateway (min 1): **$5-8/mois**
- **Total: $5-10/mois**

**Scénario 2: Trafic Moyen (500k requêtes/mois)**
- Tree Service: **$5-10/mois**
- Gateway: **$10-15/mois**
- **Total: $15-25/mois**

**Scénario 3: Fort Trafic (2M requêtes/mois)**
- Tree Service: **$15-20/mois**
- Gateway: **$20-30/mois**
- **Total: $35-50/mois**

### Optimisation des Coûts

```bash
# Réduire min instances à 0 pour Gateway (augmente cold starts)
gcloud run services update gateway \
  --region us-central1 \
  --min-instances 0

# Réduire la mémoire
gcloud run services update tree-service \
  --region us-central1 \
  --memory 256Mi

# Réduire max instances
gcloud run services update gateway \
  --region us-central1 \
  --max-instances 5
```

## 🧪 Tests

### Test Health Check

```bash
# Récupérer les URLs
GATEWAY_URL=$(gcloud run services describe gateway --region us-central1 --format 'value(status.url)')
TREE_URL=$(gcloud run services describe tree-service --region us-central1 --format 'value(status.url)')

# Tester Tree Service
curl $TREE_URL/health

# Tester Gateway
curl $GATEWAY_URL
```

### Test GraphQL

```bash
# Query simple
curl -X POST $GATEWAY_URL/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}'

# Client Tree Query
curl -X POST $GATEWAY_URL/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query GetClientTree($id: ID!) { clientTree(id: $id) { root { id name } } }",
    "variables": {"id": "YOUR_CLIENT_ID"}
  }'
```

### GraphQL Playground

Ouvrez simplement l'URL du Gateway dans votre navigateur:

```
https://gateway-xxx-uc.a.run.app
```

## 🔒 Sécurité

### Authentification (Optionnel)

Par défaut, les services sont publics (`--allow-unauthenticated`). Pour les sécuriser:

```bash
# Nécessiter l'authentification
gcloud run services update gateway \
  --region us-central1 \
  --no-allow-unauthenticated

# Appeler avec authentification
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  $GATEWAY_URL/query
```

### Variables Secrètes

Utilisez **Secret Manager** pour les données sensibles:

```bash
# Créer un secret
echo -n "mongodb+srv://..." | gcloud secrets create mongo-uri --data-file=-

# Utiliser dans Cloud Run
gcloud run services update tree-service \
  --region us-central1 \
  --update-secrets=MONGO_URI=mongo-uri:latest
```

## 🆘 Troubleshooting

### Problème: Build échoue

```bash
# Vérifier les logs de build
gcloud builds list --limit 5

# Voir les détails d'un build
gcloud builds log BUILD_ID
```

### Problème: Service ne démarre pas

```bash
# Voir les logs
gcloud run logs read tree-service --region us-central1 --limit 100

# Vérifier la configuration
gcloud run services describe tree-service --region us-central1
```

### Problème: Cold Starts

```bash
# Augmenter min instances
gcloud run services update gateway \
  --region us-central1 \
  --min-instances 1

# Ou activer le CPU boost
gcloud run services update gateway \
  --region us-central1 \
  --cpu-boost
```

### Problème: MongoDB Connection

```bash
# Vérifier que l'IP 0.0.0.0/0 est whitelisted dans MongoDB Atlas
# Tester la connexion depuis Cloud Shell
gcloud cloud-shell ssh
mongosh "$MONGO_URI"
```

### Problème: Service Timeout

```bash
# Augmenter le timeout (max 3600s)
gcloud run services update gateway \
  --region us-central1 \
  --timeout 600s
```

## 📚 Ressources

- [Documentation Cloud Run](https://cloud.google.com/run/docs)
- [Pricing Calculator](https://cloud.google.com/products/calculator)
- [Best Practices](https://cloud.google.com/run/docs/best-practices)
- [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)

## 🎯 Checklist de Déploiement

- [ ] Compte Google Cloud créé et facturation activée
- [ ] MongoDB Atlas configuré avec IP 0.0.0.0/0 whitelisté
- [ ] gcloud CLI installé et configuré
- [ ] Variables d'environnement configurées dans `.env.cloudrun`
- [ ] Tree Service déployé avec succès
- [ ] Gateway déployé avec succès
- [ ] Tests GraphQL fonctionnels
- [ ] Logs et monitoring configurés
- [ ] Plan de backup MongoDB en place

---

**Votre application est maintenant déployée sur Google Cloud Run! 🎉**




