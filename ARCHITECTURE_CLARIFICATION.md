# Clarification de l'Architecture

## Structure des Dossiers

### 📁 `graph/` - Ancien Code Monolithique (À CONSERVER pour compatibilité)
- **Usage**: Code GraphQL de l'ancien serveur monolithique (`server.go`)
- **Contenu**: Schéma GraphQL complet avec tous les types (Product, Client, Sale, Payment, etc.)
- **Status**: ⚠️ **Legacy** - Utilisé uniquement par `server.go` (monolithique)
- **Ne pas utiliser** pour les nouveaux microservices

### 📁 `gateway/` - GraphQL Gateway (Microservices)
- **Usage**: Point d'entrée GraphQL pour l'architecture microservices
- **Contenu**: 
  - `graph/` - Schéma GraphQL et resolvers du gateway
  - `internal/` - Logique interne (clients, config, models)
  - `main.go` - Point d'entrée du gateway
- **Status**: ✅ **Actif** - Utilisé pour les microservices
- **Indépendant** de `graph/` et `internal/`

### 📁 `services/` - Microservices
- **Usage**: Services backend indépendants
- **Contenu**: 
  - `tree-service/` - Service dédié à l'arbre client
  - (autres services à venir)

### 📁 `internal/` - Code Partagé (Monolithique)
- **Usage**: Code partagé pour le monolithique (`server.go`)
- **Status**: ⚠️ **Legacy** - Utilisé uniquement par le monolithique
- Les microservices ont leur propre code dans `services/{service}/internal/`

## Recommandation

Pour éviter la confusion :
1. **Utiliser `gateway/`** pour tout le code du GraphQL Gateway
2. **Ignorer `graph/`** si vous utilisez les microservices
3. **Utiliser `graph/`** uniquement si vous utilisez le monolithique (`server.go`)

## Migration

Si vous voulez migrer complètement vers les microservices :
- Le code dans `graph/` peut être supprimé une fois que tous les services sont migrés
- Pour l'instant, il est conservé pour compatibilité avec `server.go`

