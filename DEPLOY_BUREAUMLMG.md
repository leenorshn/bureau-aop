# 🚀 Déploiement Rapide - Projet bureaumlmg

Guide de déploiement pour votre projet Google Cloud `bureaumlmg`.

## ⚡ Déploiement Express (5 minutes)

### Étape 1: Configurer MongoDB

1. Connectez-vous à [MongoDB Atlas](https://cloud.mongodb.com)
2. Créez un cluster M0 (gratuit) si pas déjà fait
3. **Important**: Whitelist l'IP `0.0.0.0/0` dans Network Access
4. Copiez votre URI de connexion

### Étape 2: Configurer les Variables

```bash
# Éditer le fichier .env.cloudrun
nano .env.cloudrun
```

Remplacez la ligne `MONGO_URI` avec votre vraie URI MongoDB :

```bash
export MONGO_URI="mongodb+srv://VOTRE_USER:VOTRE_PASSWORD@cluster.mongodb.net/bureau?retryWrites=true&w=majority"
```

### Étape 3: Installer gcloud CLI (si pas déjà fait)

```bash
# macOS
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# Linux
curl https://sdk.cloud.google.com | bash
exec -l $SHELL
```

### Étape 4: Se Connecter à Google Cloud

```bash
# Login
gcloud auth login

# Vérifier que bureaumlmg est accessible
gcloud projects list | grep bureaumlmg
```

### Étape 5: Charger les Variables

```bash
source .env.cloudrun
```

### Étape 6: Déployer

```bash
./scripts/deploy-cloudrun.sh
```

**Le script va automatiquement:**
- ✅ Utiliser le projet `bureaumlmg`
- ✅ Activer les APIs nécessaires
- ✅ Builder les images Docker
- ✅ Déployer Tree Service et Gateway
- ✅ Afficher les URLs de vos services

## 🧪 Tester le Déploiement

```bash
# Récupérer l'URL du Gateway
GATEWAY_URL=$(gcloud run services describe gateway --region us-central1 --format 'value(status.url)')

# Ouvrir GraphQL Playground
open $GATEWAY_URL

# Ou tester avec curl
curl -X POST $GATEWAY_URL/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}'
```

## 📊 Voir les Services

```bash
# Lister les services Cloud Run
gcloud run services list --project bureaumlmg

# Voir les logs du Gateway
gcloud run logs read gateway --region us-central1 --limit 50

# Voir les logs du Tree Service
gcloud run logs read tree-service --region us-central1 --limit 50
```

## 🔄 Mettre à Jour l'Application

```bash
# Recharger les variables
source .env.cloudrun

# Redéployer
./scripts/deploy-cloudrun.sh
```

## 🌐 URLs de Production

Une fois déployé, vos services seront accessibles à:

- **Gateway (GraphQL)**: `https://gateway-xxx-uc.a.run.app`
- **Tree Service**: `https://tree-service-xxx-uc.a.run.app`

## 📋 Configuration Recommandée

Pour éviter les cold starts sur le Gateway (recommandé):

```bash
gcloud run services update gateway \
  --region us-central1 \
  --min-instances 1 \
  --project bureaumlmg
```

## 💰 Coûts Estimés

Avec votre configuration actuelle:

- **Tree Service** (min 0): $0-5/mois
- **Gateway** (min 1): $5-10/mois
- **Total**: ~$5-15/mois

### Optimiser les Coûts

Si vous voulez réduire les coûts (accepter cold starts):

```bash
gcloud run services update gateway \
  --region us-central1 \
  --min-instances 0 \
  --project bureaumlmg
```

## 🗑️ Supprimer les Services

```bash
./scripts/destroy-cloudrun.sh
```

## 🆘 Troubleshooting

### Erreur "Project not found"

```bash
# Vérifier que vous avez accès au projet
gcloud projects list | grep bureaumlmg

# Si nécessaire, se reconnecter
gcloud auth login
```

### Erreur "Permission denied"

```bash
# Vérifier vos permissions
gcloud projects get-iam-policy bureaumlmg

# Vous devez avoir au minimum le rôle "Editor" ou "Owner"
```

### Service ne démarre pas

```bash
# Voir les logs détaillés
gcloud run logs read tree-service --region us-central1 --limit 200 --project bureaumlmg

# Vérifier la configuration
gcloud run services describe tree-service --region us-central1 --project bureaumlmg
```

### MongoDB Connection Failed

1. Vérifier que `0.0.0.0/0` est whitelisted dans MongoDB Atlas
2. Vérifier que l'URI dans `.env.cloudrun` est correct
3. Tester la connexion depuis Cloud Shell:

```bash
gcloud cloud-shell ssh --project bureaumlmg
mongosh "$MONGO_URI"
```

## 📚 Documentation Complète

Pour plus de détails, consultez:
- [CLOUD_RUN_DEPLOYMENT.md](./CLOUD_RUN_DEPLOYMENT.md) - Guide complet
- [CLOUD_RUN_CHEATSHEET.md](./CLOUD_RUN_CHEATSHEET.md) - Commandes utiles

## ✅ Checklist de Déploiement

- [ ] MongoDB Atlas configuré avec IP 0.0.0.0/0
- [ ] gcloud CLI installé et connecté
- [ ] `.env.cloudrun` configuré avec le bon MONGO_URI
- [ ] Variables chargées: `source .env.cloudrun`
- [ ] Déploiement exécuté: `./scripts/deploy-cloudrun.sh`
- [ ] Tests GraphQL fonctionnels
- [ ] Logs vérifiés

---

**Projet: bureaumlmg | Temps total: ~5 minutes ⚡**




