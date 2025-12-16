# Architecture Microservices - Bureau MLM

## Vue d'ensemble

Cette architecture sépare l'application monolithique en microservices indépendants pour améliorer les performances, la scalabilité et la maintenabilité.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│              GraphQL Gateway (Port 8080)                │
│  - Point d'entrée unique                               │
│  - Routage vers les microservices                       │
│  - Agrégation des réponses                              │
└─────────────────────────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
        ▼               ▼               ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ Client       │ │ Tree         │ │ Binary       │
│ Service      │ │ Service      │ │ Commission   │
│ (Port 8081)  │ │ (Port 8082)  │ │ Service      │
│              │ │              │ │ (Port 8083)  │
│ - CRUD       │ │ - clientTree │ │ - Calculs    │
│ - Auth       │ │ - Optimisé   │ │ - Cycles     │
│              │ │ - Cache      │ │ - Capping    │
└──────────────┘ └──────────────┘ └──────────────┘
        │               │               │
        └───────────────┼───────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │   MongoDB        │
              │   (Shared DB)    │
              └─────────────────┘
```

## Microservices

### 1. GraphQL Gateway (gateway/)
- **Port**: 8080
- **Responsabilité**: 
  - Point d'entrée GraphQL unique
  - Routage des queries vers les microservices appropriés
  - Agrégation des réponses
  - Gestion de l'authentification

### 2. Client Service (services/client-service/)
- **Port**: 8081
- **Responsabilité**:
  - CRUD des clients
  - Authentification client (clientLogin)
  - Gestion des profils clients
  - API REST/GraphQL interne

### 3. Tree Service (services/tree-service/)
- **Port**: 8082
- **Responsabilité**:
  - **Gestion de l'arbre client (clientTree)**
  - Calcul optimisé de l'arbre
  - Cache Redis pour les arbres
  - Calcul des actifs par jambe
  - API REST/GraphQL interne

### 4. Binary Commission Service (services/binary-commission-service/)
- **Port**: 8083
- **Responsabilité**:
  - Calcul des commissions binaires
  - Gestion des cycles
  - Capping journalier/hebdomadaire
  - Qualification des membres
  - API REST/GraphQL interne

### 5. Product Service (services/product-service/)
- **Port**: 8084
- **Responsabilité**:
  - CRUD des produits
  - Gestion du stock

### 6. Sale Service (services/sale-service/)
- **Port**: 8085
- **Responsabilité**:
  - Gestion des ventes
  - Historique des ventes

### 7. Payment Service (services/payment-service/)
- **Port**: 8086
- **Responsabilité**:
  - Gestion des paiements
  - Transactions

## Communication entre services

### Option 1: gRPC (Recommandé)
- Performance optimale
- Type-safe
- Streaming support

### Option 2: HTTP/REST
- Plus simple à implémenter
- Compatible avec GraphQL
- Facile à déboguer

### Option 3: GraphQL Federation
- Chaque service expose son propre GraphQL
- Le gateway agrège les schémas

## Cache Strategy

### Tree Service
- **Redis** pour cache des arbres
- Cache key: `tree:{clientId}`
- TTL: 5 minutes
- Invalidation lors des modifications

### Binary Commission Service
- Cache des calculs de qualification
- Cache des actifs par jambe

## Base de données

- **MongoDB partagé** (pour commencer)
- Chaque service peut avoir ses propres collections
- Évolution future: base de données par service

## Déploiement

### Docker Compose
- Chaque service dans son propre container
- Réseau Docker pour communication
- Volume partagé pour MongoDB

### Kubernetes (Future)
- Déploiement indépendant
- Auto-scaling par service
- Service mesh (Istio/Linkerd)

## Avantages

1. **Performance**: Tree Service peut être optimisé indépendamment
2. **Scalabilité**: Scale uniquement le Tree Service si nécessaire
3. **Maintenabilité**: Code isolé par domaine
4. **Déploiement**: Déploiement indépendant des services
5. **Résilience**: Un service en panne n'affecte pas les autres

## Migration Strategy

1. **Phase 1**: Créer le Tree Service et le Gateway
2. **Phase 2**: Migrer clientTree vers Tree Service
3. **Phase 3**: Séparer les autres services progressivement
4. **Phase 4**: Optimiser et ajouter le cache


