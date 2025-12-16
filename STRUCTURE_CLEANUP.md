# Structure du Projet - Nettoyage Effectué

## Structure Actuelle

### ✅ Gateway (Microservices)
```
gateway/
├── main.go              # Point d'entrée du Gateway
├── graph/               # GraphQL Schema et Resolvers (TOUT EST ICI)
│   ├── schema.graphqls
│   ├── resolver.go
│   └── schema.resolvers.go
└── internal/            # Code interne (clients, config, models)
```

**Tout le code GraphQL du Gateway est maintenant dans `gateway/graph/`**

### ✅ Tree Service (Microservices)
```
services/tree-service/
├── main.go
└── internal/            # Service, handlers, cache, store
```

### ⚠️ Ancien Serveur Monolithique (Optionnel)
```
server.go                # Ancien serveur monolithique
graph/                   # GraphQL pour le monolithique (utilisé par server.go)
internal/                # Services et repositories partagés
```

**Note**: Le dossier `graph/` à la racine est utilisé par `server.go` (monolithique). Si vous n'utilisez plus le monolithique, vous pouvez le supprimer.

## Fichiers Supprimés

✅ `gateway/api/` - Dossier dupliqué supprimé
✅ `graph/schema.resolvers.go.backup` - Fichier backup supprimé
✅ `graph/schema.resolvers.go.tmp` - Fichier temporaire supprimé
✅ `graph/schemam.mgraphqls` - Fichier erroné supprimé
✅ `bureau` - Binaire compilé supprimé

## Fichiers Créés

✅ `.gitignore` - Pour éviter de commiter les fichiers inutiles
✅ `gateway/README.md` - Documentation du Gateway

## Recommandations

### Si vous utilisez uniquement les microservices:

Vous pouvez supprimer:
- `server.go` (remplacé par `gateway/main.go`)
- `graph/` à la racine (remplacé par `gateway/graph/`)
- `gqlgen.yml` à la racine (remplacé par `gateway/gqlgen.yml`)

### Si vous gardez les deux (monolithique + microservices):

Gardez tout tel quel. Les deux architectures peuvent coexister.

## Prochaines Étapes

1. Générer le code GraphQL du Gateway:
   ```bash
   cd gateway
   go generate ./graph
   ```

2. Tester le Gateway:
   ```bash
   cd gateway
   go run main.go
   ```


