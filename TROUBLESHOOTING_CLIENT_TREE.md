# Guide de Diagnostic - Erreur "Service Unavailable" avec clientTree

## Problème
L'erreur `"Service Unavailable"` signifie que le serveur GraphQL ne répond pas correctement ou qu'il y a un problème de connexion.

## Vérifications à faire

### 1. Vérifier que le serveur démarre correctement

```bash
# Démarrer le serveur et vérifier les logs
go run server.go
```

**Vérifiez :**
- ✅ Le serveur démarre sans erreur
- ✅ La connexion MongoDB est réussie
- ✅ Le message "Starting server" apparaît
- ✅ Aucune erreur de panique (panic)

### 2. Vérifier l'URL du endpoint GraphQL

Assurez-vous que vous utilisez la bonne URL :
- **Local** : `http://localhost:8080/query`
- **Production** : Vérifiez votre variable d'environnement

### 3. Tester avec une query simple

Testez d'abord avec une query plus simple pour isoler le problème :

```graphql
query {
  client(id: "6906e2ca634b66b9c3fb7a07") {
    id
    name
    clientId
  }
}
```

Si cette query fonctionne, le problème est spécifique à `clientTree`.

### 4. Vérifier les logs du serveur

Lorsque vous exécutez la query `clientTree`, regardez les logs du serveur pour voir :
- Des erreurs de connexion MongoDB
- Des timeouts
- Des erreurs dans `enrichClientTreeNode`
- Des messages "Erreur lors de l'enrichissement du nœud"

### 5. Vérifier la connexion MongoDB

Le problème peut venir d'une connexion MongoDB qui timeout :

```bash
# Vérifier que MongoDB est accessible
# Vérifier les variables d'environnement
echo $MONGO_URI
```

### 6. Tester avec une version simplifiée

Si le problème persiste, testez avec une query qui ne demande pas les champs calculés :

```graphql
query clientTreeSimple($id: ID!) {
  clientTree(id: $id) {
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
    }
    totalNodes
    maxLevel
  }
}
```

## Solutions possibles

### Solution 1 : Timeout MongoDB

Si MongoDB est lent ou inaccessible, ajoutez un timeout plus long dans `server.go` :

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Augmenter le timeout
```

### Solution 2 : Désactiver temporairement les calculs coûteux

Si le problème vient des calculs d'actifs, vous pouvez temporairement désactiver l'enrichissement :

Dans `graph/schema.resolvers.go`, commentez l'appel à `enrichClientTreeNodeOptimized` :

```go
// Temporairement désactivé pour debug
// if err := r.enrichClientTreeNodeOptimized(ctx, rootNode, client, 0, 3, activeCache); err != nil {
//     fmt.Printf("Erreur lors de l'enrichissement du nœud racine: %v\n", err)
// }
```

### Solution 3 : Vérifier les variables d'environnement

Assurez-vous que toutes les variables d'environnement sont correctement définies :

```bash
# Vérifier le fichier .env ou les variables d'environnement
cat .env
```

### Solution 4 : Redémarrer le serveur

Parfois, un simple redémarrage résout le problème :

```bash
# Arrêter le serveur (Ctrl+C)
# Redémarrer
go run server.go
```

## Query de test recommandée

Utilisez cette query pour tester progressivement :

```graphql
# Étape 1 : Test basique
query test1($id: ID!) {
  clientTree(id: $id) {
    totalNodes
    maxLevel
  }
}

# Étape 2 : Ajouter root
query test2($id: ID!) {
  clientTree(id: $id) {
    root {
      id
      name
      clientId
    }
    totalNodes
  }
}

# Étape 3 : Ajouter nodes (sans champs calculés)
query test3($id: ID!) {
  clientTree(id: $id) {
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
    }
    totalNodes
  }
}

# Étape 4 : Ajouter les champs calculés progressivement
query test4($id: ID!) {
  clientTree(id: $id) {
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
    }
    totalNodes
  }
}
```

## Debug avancé

Si le problème persiste, activez les logs détaillés :

Dans `graph/schema.resolvers.go`, ajoutez des logs :

```go
func (r *queryResolver) ClientTree(ctx context.Context, id string) (*model.ClientTree, error) {
    fmt.Printf("DEBUG: ClientTree appelé avec id=%s\n", id)
    
    client, err := r.Resolver.clientService.GetByID(ctx, id)
    if err != nil {
        fmt.Printf("DEBUG: Erreur GetByID: %v\n", err)
        return nil, fmt.Errorf("client introuvable: %w", err)
    }
    
    fmt.Printf("DEBUG: Client trouvé: %s\n", client.Name)
    // ... reste du code
}
```

## Contact

Si le problème persiste après avoir suivi ces étapes, vérifiez :
1. Les logs complets du serveur
2. Les logs MongoDB
3. La configuration réseau/firewall




