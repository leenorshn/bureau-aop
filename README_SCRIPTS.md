# Scripts de Lancement du Projet

Ce projet inclut des scripts pour faciliter le lancement et la gestion des services microservices.

## 📋 Scripts Disponibles

### 1. `start.sh` - Lancer les services

Lance le Tree Service et le Gateway.

**Usage:**
```bash
# Mode interactif (affiche les logs en temps réel)
./start.sh

# Mode arrière-plan (services en background)
./start.sh --background
```

**Fonctionnalités:**
- ✅ Vérifie que MongoDB est accessible
- ✅ Charge les variables d'environnement depuis `.env` (si présent)
- ✅ Lance le Tree Service sur le port 8082
- ✅ Lance le Gateway sur le port 8080
- ✅ Vérifie que les services sont prêts
- ✅ Affiche les URLs et informations utiles
- ✅ Gère proprement l'arrêt avec Ctrl+C

**Variables d'environnement:**
- `MONGO_URI` - URI MongoDB (défaut: `mongodb://localhost:27017`)
- `MONGO_DB_NAME` - Nom de la base de données (défaut: `bureau`)
- `TREE_SERVICE_PORT` - Port du Tree Service (défaut: `8082`)
- `TREE_SERVICE_URL` - URL du Tree Service (défaut: `http://localhost:8082`)
- `GATEWAY_PORT` - Port du Gateway (défaut: `8080`)

**Logs:**
Les logs sont sauvegardés dans le dossier `logs/`:
- `logs/tree-service.log` - Logs du Tree Service
- `logs/gateway.log` - Logs du Gateway

### 2. `stop.sh` - Arrêter les services

Arrête tous les services en cours d'exécution.

**Usage:**
```bash
./stop.sh
```

**Fonctionnalités:**
- ✅ Lit les PIDs depuis `.services.pid`
- ✅ Arrête proprement tous les processus
- ✅ Force l'arrêt si nécessaire
- ✅ Nettoie le fichier PID

### 3. `restart.sh` - Redémarrer les services

Redémarre tous les services.

**Usage:**
```bash
# Mode interactif
./restart.sh

# Mode arrière-plan
./restart.sh --background
```

## 🚀 Démarrage Rapide

### 1. Créer un fichier `.env` (optionnel)

```bash
cp env.example .env
```

Puis éditez `.env` avec vos configurations:
```env
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=bureau
TREE_SERVICE_PORT=8082
TREE_SERVICE_URL=http://localhost:8082
GATEWAY_PORT=8080
```

### 2. Lancer les services

```bash
./start.sh
```

### 3. Accéder aux services

- **GraphQL Playground**: http://localhost:8080/
- **GraphQL Endpoint**: http://localhost:8080/query
- **Tree Service API**: http://localhost:8082/api/v1/tree/{clientId}

### 4. Arrêter les services

```bash
./stop.sh
```

Ou appuyez sur `Ctrl+C` si vous avez lancé en mode interactif.

## 📝 Exemple de Query GraphQL

Une fois les services lancés, vous pouvez tester avec cette query:

```graphql
query {
  clientTree(id: "6906e2ca634b66b9c3fb7a07") {
    root {
      id
      name
      clientId
      totalEarnings
      walletBalance
    }
    nodes {
      id
      name
      clientId
      position
      level
      totalEarnings
      walletBalance
      leftActives
      rightActives
      isActive
      isQualified
    }
    totalNodes
    maxLevel
  }
}
```

## 🔧 Dépannage

### Les services ne démarrent pas

1. Vérifiez que MongoDB est en cours d'exécution:
   ```bash
   mongosh --eval "db.adminCommand('ping')"
   ```

2. Vérifiez les logs:
   ```bash
   tail -f logs/tree-service.log
   tail -f logs/gateway.log
   ```

3. Vérifiez que les ports ne sont pas déjà utilisés:
   ```bash
   lsof -i :8080
   lsof -i :8082
   ```

### Les services ne répondent pas

1. Vérifiez que les services sont bien démarrés:
   ```bash
   ps aux | grep "go run main.go"
   ```

2. Testez les endpoints directement:
   ```bash
   curl http://localhost:8082/api/v1/tree/test
   curl http://localhost:8080/query
   ```

### Arrêt forcé

Si les services ne s'arrêtent pas proprement:
```bash
# Trouver les processus
ps aux | grep "go run main.go"

# Arrêter manuellement
kill -9 <PID>
```

## 📦 Structure des Fichiers

```
.
├── start.sh          # Script de lancement
├── stop.sh           # Script d'arrêt
├── restart.sh        # Script de redémarrage
├── .env              # Variables d'environnement (optionnel)
├── .services.pid     # Fichier PID (créé automatiquement)
└── logs/             # Dossier des logs
    ├── tree-service.log
    └── gateway.log
```

## 🎯 Notes

- Les scripts sont compatibles avec bash et zsh
- Les logs sont sauvegardés dans `logs/` pour faciliter le débogage
- Le fichier `.services.pid` est créé automatiquement et contient les PIDs des services
- Les scripts gèrent proprement l'arrêt avec Ctrl+C

