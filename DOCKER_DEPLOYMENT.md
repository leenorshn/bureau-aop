# 🐳 Déploiement Docker Local - Bureau MLM

Ce guide vous explique comment déployer l'application Bureau MLM en local avec Docker.

## 📋 Prérequis

- Docker (version 20.10+)
- Docker Compose (version 2.0+)
- Git

## 🚀 Déploiement rapide

### 1. Cloner le projet
```bash
git clone <votre-repo>
cd bureau
```

### 2. Déployer avec Docker Compose
```bash
# Déploiement automatique
./scripts/deploy-local.sh
```

### 3. Initialiser l'admin
```bash
# Créer l'utilisateur admin
./scripts/seed-admin.sh
```

## 🔧 Déploiement manuel

### 1. Démarrer MongoDB
```bash
docker-compose -f docker-compose.local.yml up mongodb -d
```

### 2. Construire et démarrer l'application
```bash
docker-compose -f docker-compose.local.yml up --build -d
```

### 3. Vérifier les logs
```bash
docker-compose -f docker-compose.local.yml logs -f bureau-backend
```

## 🌐 Services disponibles

Une fois déployé, les services suivants sont disponibles :

- **GraphQL Playground**: http://localhost:8080
- **API GraphQL**: http://localhost:8080/query
- **MongoDB**: localhost:27017

## 🔑 Informations de connexion

- **Admin Email**: admin@mlm.com
- **Admin Password**: admin123
- **MongoDB**: admin/password123

## 📝 Commandes utiles

### Gestion des conteneurs
```bash
# Voir le statut
docker-compose -f docker-compose.local.yml ps

# Voir les logs
docker-compose -f docker-compose.local.yml logs -f

# Redémarrer un service
docker-compose -f docker-compose.local.yml restart bureau-backend

# Arrêter tous les services
docker-compose -f docker-compose.local.yml down

# Arrêter et supprimer les volumes
docker-compose -f docker-compose.local.yml down -v
```

### Accès aux conteneurs
```bash
# Accéder au conteneur backend
docker exec -it bureau-backend sh

# Accéder à MongoDB
docker exec -it bureau-mongodb mongosh
```

### Nettoyage
```bash
# Supprimer les images
docker rmi bureau-mlm-backend:latest

# Nettoyer tout
docker system prune -a
```

## 🐛 Dépannage

### Problèmes courants

1. **Port déjà utilisé**
   ```bash
   # Vérifier les ports utilisés
   lsof -i :8080
   lsof -i :27017
   ```

2. **Erreur de connexion MongoDB**
   ```bash
   # Vérifier les logs MongoDB
   docker-compose -f docker-compose.local.yml logs mongodb
   ```

3. **Erreur de build**
   ```bash
   # Nettoyer et reconstruire
   docker-compose -f docker-compose.local.yml down
   docker-compose -f docker-compose.local.yml up --build --force-recreate
   ```

### Logs détaillés
```bash
# Logs du backend
docker-compose -f docker-compose.local.yml logs bureau-backend

# Logs de MongoDB
docker-compose -f docker-compose.local.yml logs mongodb

# Tous les logs
docker-compose -f docker-compose.local.yml logs
```

## 🔄 Mise à jour

Pour mettre à jour l'application :

1. Arrêter les services
2. Puller les dernières modifications
3. Reconstruire et redémarrer

```bash
docker-compose -f docker-compose.local.yml down
git pull
docker-compose -f docker-compose.local.yml up --build -d
```

## 📊 Monitoring

### Vérifier la santé des services
```bash
# Statut des conteneurs
docker ps

# Utilisation des ressources
docker stats

# Espace disque
docker system df
```

## 🗄️ Base de données

### Sauvegarde
```bash
# Créer une sauvegarde
docker exec bureau-mongodb mongodump --out /backup --db mlm_db

# Copier la sauvegarde
docker cp bureau-mongodb:/backup ./backup
```

### Restauration
```bash
# Copier la sauvegarde
docker cp ./backup bureau-mongodb:/backup

# Restaurer
docker exec bureau-mongodb mongorestore /backup
```

## 🔒 Sécurité

⚠️ **Important**: Cette configuration est pour le développement local uniquement. Pour la production :

- Changez tous les mots de passe par défaut
- Utilisez des secrets Docker
- Configurez un réseau privé
- Activez l'authentification MongoDB
- Utilisez HTTPS

## 📞 Support

En cas de problème, vérifiez :

1. Les logs des conteneurs
2. La configuration des ports
3. La connectivité réseau
4. Les variables d'environnement

Pour plus d'aide, consultez la documentation du projet ou créez une issue.






