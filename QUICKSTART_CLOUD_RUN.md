# 🚀 Démarrage Rapide - Google Cloud Run

Guide ultra-rapide pour déployer Bureau MLM sur Google Cloud Run en **moins de 10 minutes**.

## ⚡ Configuration Express (5 minutes)

### 1. Prérequis

- ✅ Compte Google Cloud ([créer un compte](https://cloud.google.com))
- ✅ MongoDB Atlas configuré ([guide rapide](https://www.mongodb.com/cloud/atlas/register))

### 2. Configuration Automatique

```bash
# Cloner le projet (si pas déjà fait)
cd /path/to/bureau

# Lancer la configuration
./scripts/setup-cloudrun.sh
```

Le script va vous demander:
1. **Project ID**: Nom de votre projet (ex: `bureau-mlm-prod`)
2. **Region**: Choisir `us-central1` (ou autre région)
3. **MongoDB URI**: Votre URI MongoDB Atlas
4. **Database Name**: Nom de la DB (défaut: `bureau`)

### 3. Déploiement (3 minutes)

```bash
# Charger les variables
source .env.cloudrun

# Déployer
./scripts/deploy-cloudrun.sh
```

**C'est tout! 🎉**

## 🧪 Tester Votre Déploiement

```bash
# Récupérer l'URL du Gateway
GATEWAY_URL=$(gcloud run services describe gateway --region us-central1 --format 'value(status.url)')

# Ouvrir dans le navigateur
open $GATEWAY_URL

# Ou tester avec curl
curl -X POST $GATEWAY_URL/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}'
```

## 📋 Commandes Utiles

```bash
# Voir les services déployés
gcloud run services list

# Voir les logs
gcloud run logs read gateway --region us-central1 --limit 50

# Mettre à jour
source .env.cloudrun
./scripts/deploy-cloudrun.sh

# Supprimer les services
gcloud run services delete gateway --region us-central1
gcloud run services delete tree-service --region us-central1
```

## 💰 Coûts

- **Tier gratuit**: 2M requêtes/mois
- **Coût estimé**: $5-15/mois après tier gratuit

## 🆘 Besoin d'Aide?

Consultez le guide complet: [CLOUD_RUN_DEPLOYMENT.md](./CLOUD_RUN_DEPLOYMENT.md)

## 📝 Checklist

- [ ] MongoDB Atlas configuré
- [ ] `./scripts/setup-cloudrun.sh` exécuté
- [ ] `source .env.cloudrun` exécuté
- [ ] `./scripts/deploy-cloudrun.sh` exécuté
- [ ] Tests GraphQL fonctionnels

---

**Configuration en 5 minutes, déploiement en 3 minutes. Total: 8 minutes! ⚡**




