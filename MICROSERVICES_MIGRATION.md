# Guide de Migration vers Microservices

## Résumé

L'architecture a été réorganisée en microservices pour améliorer les performances du chargement de l'arbre client. Le Tree Service est maintenant séparé et optimisé avec cache.

## Architecture Créée

### 1. GraphQL Gateway (`gateway/`)
- Point d'entrée unique GraphQL
- Route les queries vers les microservices
- Port: 8080

### 2. Tree Service (`services/tree-service/`)
- Service dédié à la gestion de l'arbre client
- Cache intégré (Memory/Redis)
- Optimisations de performance
- Port: 8082

## Avantages

1. **Performance**: Tree Service optimisé indépendamment
2. **Cache**: Cache dédié pour les arbres (5 min TTL)
3. **Scalabilité**: Scale uniquement le Tree Service si nécessaire
4. **Isolation**: Un problème dans un service n'affecte pas les autres

## Migration

### Option 1: Utiliser les microservices (Recommandé)

1. Démarrer les services:
```bash
docker-compose -f docker-compose.microservices.yml up -d
```

2. Le Gateway est disponible sur `http://localhost:8080/query`
3. Utiliser la même query GraphQL qu'avant

### Option 2: Garder l'architecture monolithique

L'ancien serveur (`server.go`) continue de fonctionner normalement.

## Prochaines Étapes

1. ✅ Tree Service créé et fonctionnel
2. ⏳ Ajouter Redis pour cache distribué
3. ⏳ Créer Client Service
4. ⏳ Créer Binary Commission Service
5. ⏳ Migrer progressivement les autres services

## Tests

### Tester le Tree Service directement
```bash
curl http://localhost:8082/api/v1/tree/6906e2ca634b66b9c3fb7a07
```

### Tester via Gateway
```graphql
query {
  clientTree(id: "6906e2ca634b66b9c3fb7a07") {
    root { id name clientId }
    nodes { id name position level }
    totalNodes
  }
}
```


