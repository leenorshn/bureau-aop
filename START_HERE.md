# 🚀 COMMENCEZ ICI - Déploiement Cloud Run

Guide de démarrage pour déployer Bureau MLM sur Google Cloud Run (Projet: **bureaumlmg**).

## ⚡ Déploiement en 4 Étapes

### Étape 1: Installer gcloud CLI (si pas déjà fait)

```bash
# macOS ou Linux
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# Se connecter à Google Cloud
gcloud auth login
```

### Étape 2: Configurer MongoDB Atlas

1. Allez sur [MongoDB Atlas](https://cloud.mongodb.com)
2. Créez un cluster gratuit M0 (si pas déjà fait)
3. **Important**: Dans "Network Access", ajoutez l'IP `0.0.0.0/0`
4. Dans "Database Access", créez un utilisateur
5. Copiez l'URI de connexion (ex: `mongodb+srv://user:password@cluster.mongodb.net/bureau`)

### Étape 3: Initialiser la Configuration

```bash
./scripts/init-env.sh
```

Ce script va vous demander:
- URI MongoDB (collez l'URI copié à l'étape 2)
- Nom de la base de données (appuyez sur Entrée pour `bureau`)
- Région Cloud Run (appuyez sur Entrée pour `us-central1`)

### Étape 4: Déployer

```bash
# Charger les variables
source .env.cloudrun

# Déployer
./scripts/deploy-cloudrun.sh
```

**Attendez 3-5 minutes... ☕**

## ✅ Tester le Déploiement

```bash
# Récupérer l'URL du Gateway
GATEWAY_URL=$(gcloud run services describe gateway --region us-central1 --format 'value(status.url)')

# Ouvrir GraphQL Playground dans le navigateur
echo $GATEWAY_URL
open $GATEWAY_URL

# Ou tester avec curl
curl -X POST $GATEWAY_URL/query \
  -H "Content-Type: application/json" \
  -d '{"query":"{ __typename }"}'
```

## 📊 Commandes Utiles

```bash
# Voir les services déployés
gcloud run services list --project bureaumlmg

# Voir les logs en temps réel
gcloud run logs read gateway --region us-central1 --follow

# Mettre à jour l'application
source .env.cloudrun
./scripts/deploy-cloudrun.sh

# Supprimer les services
./scripts/destroy-cloudrun.sh
```

## 🆘 Problèmes Courants

### "Permission denied"
```bash
# Vérifiez que vous avez accès au projet
gcloud projects list | grep bureaumlmg

# Reconnectez-vous
gcloud auth login
```

### "MongoDB connection failed"
1. Vérifiez que `0.0.0.0/0` est dans Network Access
2. Vérifiez l'URI dans `.env.cloudrun`
3. Testez: `mongosh "YOUR_MONGO_URI"`

### "Build failed"
```bash
# Voir les logs de build
gcloud builds list --limit 5 --project bureaumlmg
```

## 💰 Coûts

- **Tier gratuit**: 2M requêtes/mois
- **Estimation**: $5-15/mois après tier gratuit
- **MongoDB**: Gratuit (M0)

## 📚 Documentation Complète

Pour plus de détails:
- [DEPLOY_BUREAUMLMG.md](./DEPLOY_BUREAUMLMG.md) - Guide complet
- [CLOUD_RUN_DEPLOYMENT.md](./CLOUD_RUN_DEPLOYMENT.md) - Documentation détaillée
- [CLOUD_RUN_CHEATSHEET.md](./CLOUD_RUN_CHEATSHEET.md) - Commandes utiles

## 🎯 Checklist

- [ ] gcloud CLI installé et connecté
- [ ] MongoDB Atlas configuré (0.0.0.0/0 whitelisted)
- [ ] `./scripts/init-env.sh` exécuté
- [ ] `source .env.cloudrun` exécuté
- [ ] `./scripts/deploy-cloudrun.sh` exécuté
- [ ] Tests GraphQL fonctionnels

---

**Projet: bureaumlmg | Temps total: 5-10 minutes ⚡**

**Besoin d'aide?** Consultez [DEPLOY_BUREAUMLMG.md](./DEPLOY_BUREAUMLMG.md)




