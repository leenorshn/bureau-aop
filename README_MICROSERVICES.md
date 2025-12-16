# Architecture Microservices - Bureau MLM

## 🎯 Objectif

Réorganiser l'application monolithique en architecture microservices pour améliorer les performances, notamment pour le chargement de l'arbre client.

## 📋 Structure Créée

```
bureau/
├── gateway/                          # GraphQL Gateway (Port 8080)
│   ├── main.go
│   ├── graph/                        # Schema GraphQL
│   │   ├── schema.graphqls
│   │   ├── schema.resolvers.go
│   │   └── resolver.go
│   └── internal/
│       ├── client/                   # Clients pour microservices
│       ├── config/
│       └── models/
│
├── services/
│   └── tree-service/                 # Tree Service (Port 8082)
│       ├── main.go
│       └── internal/
│           ├── service/              # Logique métier optimisée
│           ├── handler/              # Handlers HTTP REST
│           ├── cache/                # Cache (Memory/Redis)
│           ├── store/                # Repositories MongoDB
│           ├── config/
│           └── models/
│
└── docker-compose.microservices.yml  # Configuration Docker
```

## 🚀 Démarrage Rapide

### Avec Docker Compose

```bash
# Démarrer tous les services
docker-compose -f docker-compose.microservices.yml up -d

# Voir les logs
docker-compose -f docker-compose.microservices.yml logs -f

# Arrêter
docker-compose -f docker-compose.microservices.yml down
```

### Développement Local

#### 1. Tree Service

```bash
cd services/tree-service
go mod tidy
go run main.go
```

#### 2. Gateway

```bash
cd gateway
go mod tidy
go generate ./graph
go run main.go
```

## 🔧 Configuration

### Variables d'environnement

**Gateway:**
- `TREE_SERVICE_URL`: URL du Tree Service (défaut: http://localhost:8082)
- `GATEWAY_PORT`: Port du Gateway (défaut: 8080)

**Tree Service:**
- `MONGO_URI`: URI MongoDB (défaut: mongodb://localhost:27017)
- `MONGO_DB_NAME`: Nom de la base (défaut: bureau)
- `TREE_SERVICE_PORT`: Port du service (défaut: 8082)
- `REDIS_URL`: URL Redis (optionnel, utilise Memory cache si vide)

## 📊 Performance

### Optimisations du Tree Service

1. **Cache**: Cache des arbres complets (TTL: 5 minutes)
2. **Limite de profondeur**: Calcul des actifs limité aux 3 premiers niveaux
3. **Cache d'activité**: Évite les appels DB répétés pour vérifier si un client est actif
4. **Calculs conditionnels**: Pas de calculs coûteux pour les niveaux profonds

### Résultats Attendus

- **Avant**: ~44+ appels DB pour un arbre de 11 nœuds
- **Après**: ~15-20 appels DB avec cache
- **Avec cache hit**: 0 appels DB (réponse instantanée)

## 🧪 Tests

### Tester le Tree Service directement

```bash
curl http://localhost:8082/api/v1/tree/6906e2ca634b66b9c3fb7a07
```

### Tester via Gateway GraphQL

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

## 📝 Prochaines Étapes

1. ✅ Tree Service créé
2. ✅ Gateway GraphQL créé
3. ⏳ Ajouter Redis pour cache distribué
4. ⏳ Créer Client Service
5. ⏳ Créer Binary Commission Service
6. ⏳ Ajouter monitoring (Prometheus/Grafana)
7. ⏳ Ajouter health checks

## 🔄 Migration depuis Monolithique

L'ancien serveur (`server.go`) continue de fonctionner. Vous pouvez:

1. **Option A**: Utiliser les microservices (recommandé pour production)
2. **Option B**: Garder le monolithique (pour développement)

Les deux peuvent coexister pendant la migration.

## 📚 Documentation

- `MICROSERVICES_ARCHITECTURE.md`: Architecture détaillée
- `MICROSERVICES_SETUP.md`: Guide de configuration
- `MICROSERVICES_MIGRATION.md`: Guide de migration



