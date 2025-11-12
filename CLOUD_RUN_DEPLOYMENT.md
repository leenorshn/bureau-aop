# 🚀 Déploiement sur Google Cloud Run

Ce guide explique comment déployer l'application MLM Bureau sur Google Cloud Run.

## 📋 Prérequis

1. **Google Cloud CLI** installé et configuré
2. **Compte Google Cloud** avec facturation activée
3. **Projet Google Cloud** créé
4. **APIs activées** : Cloud Run, Cloud Build, Container Registry

## 🛠️ Configuration

### 1. Installation de Google Cloud CLI

```bash
# macOS
brew install google-cloud-sdk

# Linux
curl https://sdk.cloud.google.com | bash

# Windows
# Télécharger depuis https://cloud.google.com/sdk/docs/install
```

### 2. Authentification

```bash
gcloud auth login
gcloud auth configure-docker
```

### 3. Configuration du projet

```bash
# Remplacer YOUR_PROJECT_ID par votre ID de projet
export PROJECT_ID="your-project-id"
gcloud config set project $PROJECT_ID
```

## 🚀 Déploiement

### Déploiement automatique

```bash
# Utiliser le script de déploiement
./scripts/deploy-cloudrun.sh YOUR_PROJECT_ID us-central1
```

### Déploiement manuel

#### 1. Construire l'image

```bash
# Construire avec Cloud Build
gcloud builds submit --tag gcr.io/$PROJECT_ID/bureau-mlm-backend --file Dockerfile.cloudrun .
```

#### 2. Déployer sur Cloud Run

```bash
gcloud run deploy bureau-mlm-backend \
    --image gcr.io/$PROJECT_ID/bureau-mlm-backend \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated \
    --port 8080 \
    --memory 1Gi \
    --cpu 1 \
    --min-instances 0 \
    --max-instances 10 \
    --timeout 300 \
    --concurrency 100
```

#### 3. Configurer les variables d'environnement

```bash
gcloud run services update bureau-mlm-backend \
    --region us-central1 \
    --set-env-vars "MONGO_URI=mongodb+srv://leenor:avenir23@clusterzone1.b45aacv.mongodb.net/mlm?retryWrites=true&w=majority,MONGO_DB_NAME=mlm_db,JWT_SECRET=your-super-secret-jwt-key-change-this-in-production,JWT_REFRESH_SECRET=your-super-secret-refresh-key-change-this-in-production,JWT_ACCESS_EXP=15m,JWT_REFRESH_EXP=7d,ADMIN_SEED_EMAIL=admin@mlm.com,ADMIN_SEED_PASSWORD=admin123,APP_PORT=8080,APP_ENV=production,BINARY_THRESHOLD=100.0,BINARY_COMMISSION_RATE=0.1,DEFAULT_PRODUCT_PRICE=50.0,PORT=8080"
```

## 🔧 Configuration avancée

### Variables d'environnement sécurisées

Pour la production, utilisez Google Secret Manager :

```bash
# Créer des secrets
gcloud secrets create jwt-secret --data-file=- <<< "your-super-secret-jwt-key"
gcloud secrets create jwt-refresh-secret --data-file=- <<< "your-super-secret-refresh-key"
gcloud secrets create mongo-uri --data-file=- <<< "mongodb+srv://..."

# Accorder les permissions
gcloud secrets add-iam-policy-binding jwt-secret \
    --member="serviceAccount:YOUR_SERVICE_ACCOUNT" \
    --role="roles/secretmanager.secretAccessor"
```

### Configuration avec cloud-run.yaml

```bash
# Déployer avec le fichier de configuration
gcloud run services replace cloud-run.yaml
```

## 📊 Monitoring et logs

### Voir les logs

```bash
# Logs en temps réel
gcloud run logs tail bureau-mlm-backend --region us-central1

# Logs historiques
gcloud run logs read bureau-mlm-backend --region us-central1
```

### Monitoring

- **Cloud Console** : https://console.cloud.google.com/run
- **Métriques** : CPU, mémoire, requêtes, latence
- **Alertes** : Configurer des alertes sur les erreurs

## 🔒 Sécurité

### 1. Authentification

```bash
# Désactiver l'accès public (optionnel)
gcloud run services remove-iam-policy-binding bureau-mlm-backend \
    --member="allUsers" \
    --role="roles/run.invoker" \
    --region us-central1
```

### 2. HTTPS uniquement

```bash
# Forcer HTTPS
gcloud run services update bureau-mlm-backend \
    --region us-central1 \
    --set-env-vars "FORCE_HTTPS=true"
```

### 3. CORS

```bash
# Configurer CORS pour le frontend
gcloud run services update bureau-mlm-backend \
    --region us-central1 \
    --set-env-vars "CORS_ORIGINS=https://your-frontend-domain.com"
```

## 🚀 CI/CD avec GitHub Actions

Créer `.github/workflows/cloud-run.yml` :

```yaml
name: Deploy to Cloud Run

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Google Cloud CLI
      uses: google-github-actions/setup-gcloud@v1
      with:
        service_account_key: ${{ secrets.GCP_SA_KEY }}
        project_id: ${{ secrets.GCP_PROJECT_ID }}
    
    - name: Deploy to Cloud Run
      run: |
        gcloud builds submit --tag gcr.io/${{ secrets.GCP_PROJECT_ID }}/bureau-mlm-backend --file Dockerfile.cloudrun .
        gcloud run deploy bureau-mlm-backend \
          --image gcr.io/${{ secrets.GCP_PROJECT_ID }}/bureau-mlm-backend \
          --platform managed \
          --region us-central1 \
          --allow-unauthenticated
```

## 📈 Optimisations

### 1. Performance

- **Cold start** : Min instances = 1 pour éviter les cold starts
- **Memory** : Ajuster selon l'utilisation (512Mi - 2Gi)
- **CPU** : 1-2 vCPU selon la charge

### 2. Coûts

- **Min instances** : 0 pour économiser
- **Max instances** : Limiter selon le budget
- **Timeout** : 300s max pour éviter les coûts élevés

### 3. Scaling

```bash
# Configuration de scaling
gcloud run services update bureau-mlm-backend \
    --region us-central1 \
    --min-instances 1 \
    --max-instances 20 \
    --concurrency 100
```

## 🧪 Tests

### Test local

```bash
# Tester l'image localement
docker build -f Dockerfile.cloudrun -t bureau-mlm-backend .
docker run -p 8080:8080 -e PORT=8080 bureau-mlm-backend
```

### Test de déploiement

```bash
# Tester l'endpoint
curl https://your-service-url.run.app/

# Tester GraphQL
curl -X POST https://your-service-url.run.app/query \
  -H "Content-Type: application/json" \
  -d '{"query": "query { __typename }"}'
```

## 🆘 Dépannage

### Problèmes courants

1. **Cold start lent** : Augmenter min-instances
2. **Mémoire insuffisante** : Augmenter memory
3. **Timeout** : Augmenter timeout ou optimiser le code
4. **Erreurs de connexion** : Vérifier les variables d'environnement

### Commandes utiles

```bash
# Voir les détails du service
gcloud run services describe bureau-mlm-backend --region us-central1

# Voir les révisions
gcloud run revisions list --service bureau-mlm-backend --region us-central1

# Rollback
gcloud run services update-traffic bureau-mlm-backend \
    --to-revisions REVISION_NAME=100 \
    --region us-central1
```

## 📞 Support

- **Documentation Cloud Run** : https://cloud.google.com/run/docs
- **Pricing** : https://cloud.google.com/run/pricing
- **Quotas** : https://cloud.google.com/run/quotas




















