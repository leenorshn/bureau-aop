# Guide de Configuration - Architecture Microservices

## Structure du Projet

```
bureau/
├── gateway/                    # GraphQL Gateway
│   ├── main.go
│   ├── graph/                  # Schema GraphQL et resolvers
│   └── internal/
│       ├── client/             # Clients pour les microservices
│       ├── config/
│       └── models/
│
├── services/
│   └── tree-service/           # Service dédié à l'arbre client
│       ├── main.go
│       └── internal/
│           ├── service/        # Logique métier
│           ├── handler/        # Handlers HTTP
│           ├── cache/          # Cache (Memory/Redis)
│           ├── store/          # Repositories MongoDB
│           ├── config/
│           └── models/
│
└── docker-compose.microservices.yml
```

## Installation

### 1. Prérequis

- Docker et Docker Compose
- Go 1.21+ (pour développement local)

### 2. Configuration

Copiez et modifiez les variables d'environnement :

```bash
# Gateway
export TREE_SERVICE_URL=http://localhost:8082
export GATEWAY_PORT=8080

# Tree Service
export MONGO_URI=mongodb://localhost:27017
export MONGO_DB_NAME=bureau
export TREE_SERVICE_PORT=8082
export REDIS_URL=redis://localhost:6379  # Optionnel
```

### 3. Démarrage avec Docker Compose

```bash
# Démarrer tous les services
docker-compose -f docker-compose.microservices.yml up -d

# Vérifier les logs
docker-compose -f docker-compose.microservices.yml logs -f

# Arrêter les services
docker-compose -f docker-compose.microservices.yml down
```

### 4. Développement Local

#### Tree Service

```bash
cd services/tree-service
go mod tidy
go run main.go
```

#### Gateway

```bash
cd gateway
go mod tidy
go generate ./graph
go run main.go
```

## Tests

### Tester le Tree Service directement

```bash
curl http://localhost:8082/api/v1/tree/6906e2ca634b66b9c3fb7a07
```

### Tester via le Gateway GraphQL

```graphql
query {
  clientTree(id: "6906e2ca634b66b9c3fb7a07") {
    root {
      id
      name
      clientId
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
    }
    totalNodes
    maxLevel
  }
}
```

## Avantages de cette Architecture

1. **Performance**: Tree Service optimisé indépendamment
2. **Scalabilité**: Scale uniquement le Tree Service si nécessaire
3. **Cache**: Cache Redis dédié pour les arbres
4. **Isolation**: Un problème dans un service n'affecte pas les autres
5. **Déploiement**: Déploiement indépendant

## Prochaines Étapes

1. Ajouter Redis pour le cache distribué
2. Créer les autres microservices (Client, Binary Commission, etc.)
3. Ajouter la gestion d'erreurs et retry logic
4. Ajouter le monitoring (Prometheus, Grafana)
5. Ajouter les health checks



