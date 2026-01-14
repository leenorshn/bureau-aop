# ☁️ Bureau MLM - Cloud Deployment

API GraphQL pour MLM binaire déployé sur **Google Cloud Run** (Projet: `bureaumlmg`).

## ⚡ Déploiement Rapide

### 🎯 Commencez ici: [START_HERE.md](./START_HERE.md)

Pour déployer en production sur Cloud Run, suivez simplement ces étapes:

```bash
# 1. Initialiser la configuration
./scripts/init-env.sh

# 2. Charger les variables
source .env.cloudrun

# 3. Déployer
./scripts/deploy-cloudrun.sh
```

**Temps total: 5-10 minutes ⚡**

## 📚 Documentation

### Pour Commencer
- **[START_HERE.md](./START_HERE.md)** - Guide de démarrage (LISEZ EN PREMIER)
- **[QUICK_DEPLOY.md](./QUICK_DEPLOY.md)** - Déploiement en 3 commandes
- **[DEPLOY_BUREAUMLMG.md](./DEPLOY_BUREAUMLMG.md)** - Guide pour le projet bureaumlmg

### Documentation Complète
- **[CLOUD_RUN_DEPLOYMENT.md](./CLOUD_RUN_DEPLOYMENT.md)** - Guide détaillé Cloud Run
- **[CLOUD_RUN_CHEATSHEET.md](./CLOUD_RUN_CHEATSHEET.md)** - Commandes utiles
- **[README_DEPLOYMENT.md](./README_DEPLOYMENT.md)** - Comparaison des options

### Code Source
- **[README.md](./README.md)** - Documentation du code et développement local

## 🏗️ Architecture

```
Client → Gateway (GraphQL) → Tree Service → MongoDB Atlas
         Cloud Run              Cloud Run
```

## 🛠️ Scripts Disponibles

```bash
# Configuration
./scripts/init-env.sh          # Initialiser .env.cloudrun
./scripts/setup-cloudrun.sh    # Configuration interactive complète

# Déploiement
./scripts/deploy-cloudrun.sh   # Déployer sur Cloud Run

# Nettoyage
./scripts/destroy-cloudrun.sh  # Supprimer les services
```

## 📦 Services Déployés

- **Gateway** - API GraphQL principale
- **Tree Service** - Gestion de l'arbre MLM binaire

## 💰 Coûts

- **Tier gratuit**: 2M requêtes/mois
- **Estimation**: $5-15/mois après tier gratuit

## 🧪 URLs de Production

Après déploiement:
- GraphQL API: `https://gateway-xxx-uc.a.run.app/query`
- Playground: `https://gateway-xxx-uc.a.run.app`

## 📊 Commandes Courantes

```bash
# Voir les services
gcloud run services list --project bureaumlmg

# Logs en temps réel
gcloud run logs read gateway --region us-central1 --follow

# Mettre à jour
source .env.cloudrun && ./scripts/deploy-cloudrun.sh
```

## 🎯 Configuration Requise

1. **Google Cloud**
   - Projet: `bureaumlmg`
   - gcloud CLI installé
   - Authentifié: `gcloud auth login`

2. **MongoDB Atlas**
   - Cluster M0 (gratuit)
   - IP `0.0.0.0/0` whitelisted
   - URI de connexion

## 🆘 Support

**Problèmes?** Consultez:
1. [DEPLOY_BUREAUMLMG.md](./DEPLOY_BUREAUMLMG.md#troubleshooting)
2. [CLOUD_RUN_DEPLOYMENT.md](./CLOUD_RUN_DEPLOYMENT.md#troubleshooting)

## ✅ Checklist

- [ ] gcloud CLI installé
- [ ] Connecté à Google Cloud
- [ ] MongoDB Atlas configuré
- [ ] `./scripts/init-env.sh` exécuté
- [ ] `source .env.cloudrun` exécuté
- [ ] `./scripts/deploy-cloudrun.sh` exécuté

---

**Projet: bureaumlmg | Région: us-central1**

**Prêt à déployer?** → [START_HERE.md](./START_HERE.md)




