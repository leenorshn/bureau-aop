# ⚡ Cloud Run Cheatsheet

Aide-mémoire rapide pour Google Cloud Run.

## 🚀 Déploiement

```bash
# Configuration initiale (une seule fois)
./scripts/setup-cloudrun.sh

# Charger les variables
source .env.cloudrun

# Déployer
./scripts/deploy-cloudrun.sh
```

## 📊 Gestion des Services

```bash
# Lister les services
gcloud run services list

# Détails d'un service
gcloud run services describe gateway --region us-central1

# Mettre à jour une variable
gcloud run services update gateway \
  --region us-central1 \
  --set-env-vars="VAR=value"

# Supprimer un service
gcloud run services delete gateway --region us-central1
```

## 📝 Logs

```bash
# Logs temps réel
gcloud run logs read gateway --region us-central1 --follow

# Derniers 100 logs
gcloud run logs read gateway --region us-central1 --limit 100

# Logs d'erreur seulement
gcloud run logs read gateway --region us-central1 --log-filter="severity>=ERROR"
```

## 🔄 Scaling

```bash
# Changer min/max instances
gcloud run services update gateway \
  --region us-central1 \
  --min-instances 1 \
  --max-instances 10

# Changer mémoire/CPU
gcloud run services update gateway \
  --region us-central1 \
  --memory 1Gi \
  --cpu 2
```

## 🧪 Tests

```bash
# Health check
curl $(gcloud run services describe tree-service --region us-central1 --format 'value(status.url)')/health

# GraphQL
GATEWAY_URL=$(gcloud run services describe gateway --region us-central1 --format 'value(status.url)')
curl -X POST $GATEWAY_URL/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}'

# Ouvrir GraphQL Playground
open $GATEWAY_URL
```

## 🔧 Debugging

```bash
# Voir les révisions
gcloud run revisions list --service gateway --region us-central1

# Rollback
gcloud run services update-traffic gateway \
  --region us-central1 \
  --to-revisions REVISION_NAME=100

# Exécuter un shell dans le container
gcloud run services proxy gateway --region us-central1
```

## 💰 Coûts

```bash
# Voir l'utilisation
gcloud run services describe gateway --region us-central1 --format="value(status.traffic)"

# Calculer les coûts
# https://cloud.google.com/products/calculator
```

## 🗑️ Nettoyage

```bash
# Supprimer tous les services
./scripts/destroy-cloudrun.sh

# Ou manuellement
gcloud run services delete gateway --region us-central1
gcloud run services delete tree-service --region us-central1
```

## 🔒 Sécurité

```bash
# Nécessiter l'authentification
gcloud run services update gateway \
  --region us-central1 \
  --no-allow-unauthenticated

# Appeler avec token
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
  $GATEWAY_URL/query
```

## 📚 Variables d'Environnement Utiles

```bash
# Charger depuis .env.cloudrun
source .env.cloudrun

# Ou définir manuellement
export GCP_PROJECT_ID="your-project-id"
export GCP_REGION="us-central1"
export MONGO_URI="mongodb+srv://..."
```

## 🔗 Liens Rapides

```bash
# Console Cloud Run
echo "https://console.cloud.google.com/run?project=$GCP_PROJECT_ID"

# Logs
echo "https://console.cloud.google.com/logs?project=$GCP_PROJECT_ID"

# Facturation
echo "https://console.cloud.google.com/billing?project=$GCP_PROJECT_ID"
```

## 🆘 Commandes de Dépannage

```bash
# Service ne démarre pas
gcloud run logs read SERVICE_NAME --region us-central1 --limit 200

# Build échoue
gcloud builds list --limit 5
gcloud builds log BUILD_ID

# Connection MongoDB échoue
# Vérifier whitelist 0.0.0.0/0 dans MongoDB Atlas

# Timeout
gcloud run services update SERVICE_NAME \
  --region us-central1 \
  --timeout 600s
```

---

**Gardez ce fichier à portée de main pour une référence rapide! 📖**




