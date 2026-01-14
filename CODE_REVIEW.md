# 🔍 Code Review - Bureau MLM API

**Date:** $(date)  
**Projet:** Bureau MLM Backend  
**Langage:** Go + GraphQL

---

## 📋 Table des Matières

1. [Résumé Exécutif](#résumé-exécutif)
2. [Problèmes Critiques](#problèmes-critiques)
3. [Problèmes Majeurs](#problèmes-majeurs)
4. [Améliorations Recommandées](#améliorations-recommandées)
5. [Points Positifs](#points-positifs)
6. [TODOs Identifiés](#todos-identifiés)

---

## 📊 Résumé Exécutif

**Statut Global:** ⚠️ **Nécessite des améliorations**

- ✅ Architecture bien structurée (microservices, GraphQL)
- ⚠️ Problèmes de sécurité et de validation
- ⚠️ Incohérences dans le schéma GraphQL
- ⚠️ Gestion d'erreurs à améliorer
- ⚠️ Performance potentielle dans le service de commission binaire

---

## 🚨 Problèmes Critiques

### 1. **Sécurité - Validation des mots de passe manquante**

**Fichier:** `graph/schema.resolvers.go:163`

```go
PasswordHash: input.Password, // service will hash
```

**Problème:** Le commentaire indique que le service hash le mot de passe, mais il n'y a pas de validation de la force du mot de passe avant le hashage.

**Recommandation:**
- Ajouter une validation de la force du mot de passe (min 8 caractères, complexité)
- Vérifier que le service hash bien le mot de passe avant stockage

### 2. **Sécurité - Gestion des erreurs expose des informations**

**Fichier:** `internal/service/binary_commission_service.go`

**Problème:** Les messages d'erreur peuvent exposer des détails internes de l'application.

**Exemple:**
```go
Reason: fmt.Sprintf("Erreur lors de la vérification de qualification: %v", err)
```

**Recommandation:**
- Ne pas exposer les erreurs brutes aux clients
- Utiliser des codes d'erreur personnalisés
- Logger les erreurs détaillées côté serveur uniquement

### 3. **Race Condition dans BinaryCommissionService**

**Fichier:** `internal/service/binary_commission_service.go:154-176`

**Problème:** Double vérification de la limite journalière, mais pas de transaction atomique.

```go
s.mu.Lock()
defer s.mu.Unlock()

// Double vérification après verrouillage
cyclesToPayFinal, err := s.applyDailyLimit(ctx, client.ID, cyclesAvailable)
```

**Reblème:** Le mutex protège seulement la fonction, mais `applyDailyLimit` fait un appel DB qui peut avoir des conditions de course avec d'autres instances du service.

**Recommandation:**
- Utiliser des transactions MongoDB ou des opérations atomiques
- Implémenter un verrouillage distribué si plusieurs instances

---

## ⚠️ Problèmes Majeurs

### 4. **Incohérences dans le schéma GraphQL**

**Fichier:** `graph/schema.graphqls`

#### 4.1 Espacement manquant
```graphql
phone:String  # Ligne 19 - manque un espace
nn:String     # Ligne 20 - manque un espace
```

**Recommandation:** Ajouter des espaces pour la cohérence:
```graphql
phone: String
nn: String
```

#### 4.2 Champs optionnels non marqués comme nullable

**Ligne 19-22:** Les champs `phone`, `nn`, `address`, `avatar` sont définis comme `String` mais devraient être `String` (nullable) car ils sont optionnels dans `ClientInput`.

**Recommandation:**
```graphql
phone: String    # Devrait être nullable
nn: String       # Devrait être nullable
address: String  # Devrait être nullable
avatar: String   # Devrait être nullable
```

#### 4.3 Incohérence entre `Client` et `ClientInput`

Dans `Client` (ligne 19-22), les champs sont `String` (non-nullable), mais dans `ClientInput` (ligne 220-223), ils sont optionnels. Cela crée une incohérence.

### 5. **Gestion d'erreurs inconsistante**

**Fichier:** `graph/schema.resolvers.go`

**Problème:** Certaines fonctions retournent directement les erreurs sans contexte.

**Exemple:**
```go
func (r *mutationResolver) ProductDelete(ctx context.Context, id string) (bool, error) {
	return r.Resolver.productService.Delete(ctx, id)
}
```

**Recommandation:**
- Ajouter un contexte d'erreur avec `fmt.Errorf` et `%w`
- Logger les erreurs avant de les retourner
- Utiliser des erreurs typées pour un meilleur handling

### 6. **Performance - Comptage récursif des actifs**

**Fichier:** `internal/service/binary_commission_service.go:286-327`

**Problème:** La fonction `countActivesInLeg` fait des appels DB récursifs qui peuvent être très coûteux pour de grands arbres.

```go
func (s *BinaryCommissionService) countActivesInLeg(ctx context.Context, rootID *primitive.ObjectID, side string) (int, error) {
	// ... boucle avec appels DB pour chaque nœud
	client, err := s.clientRepo.GetByID(ctx, currentID.Hex())
	isActive, err := s.isClientActive(ctx, currentID.Hex()) // Appel DB supplémentaire
}
```

**Recommandation:**
- Implémenter un cache (Redis) comme suggéré dans `services/tree-service/internal/cache/cache.go`
- Utiliser des requêtes batch pour récupérer plusieurs clients en une fois
- Limiter la profondeur de recherche (déjà partiellement implémenté dans `countActivesInLegWithCache`)

### 7. **Self-assignment détecté par le linter**

**Fichier:** `internal/service/caisse_service.go:98`

**Problème:** 
```go
caisse.UpdatedAt = caisse.UpdatedAt
```

**Recommandation:** Corriger cette ligne pour mettre à jour avec `time.Now()`.

---

## 💡 Améliorations Recommandées

### 8. **Validation des entrées GraphQL**

**Problème:** Pas de validation explicite des entrées dans les resolvers.

**Recommandation:**
- Ajouter des validations pour les montants (positifs)
- Valider les formats (email, dates)
- Valider les IDs (format ObjectID)

**Exemple:**
```go
func validateProductInput(input model.ProductInput) error {
	if input.Price < 0 {
		return errors.New("price must be positive")
	}
	if input.Stock < 0 {
		return errors.New("stock must be positive")
	}
	return nil
}
```

### 9. **Documentation du schéma GraphQL**

**Problème:** Le schéma GraphQL manque de descriptions pour les types et champs.

**Recommandation:** Ajouter des descriptions:
```graphql
"""
Représente un client dans le système MLM
"""
type Client {
  """
  Identifiant unique du client
  """
  id: ID!
  
  """
  Solde du portefeuille en FCFA
  """
  walletBalance: Float!
}
```

### 10. **Gestion des contextes**

**Problème:** Pas de timeout explicite dans les contextes pour les opérations longues.

**Recommandation:**
```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
```

### 11. **Tests manquants**

**Problème:** Seul `binary_commission_service_test.go` existe, mais pas de tests pour les resolvers GraphQL.

**Recommandation:**
- Ajouter des tests unitaires pour les resolvers
- Ajouter des tests d'intégration pour les mutations critiques
- Tester les cas limites (volumes négatifs, cycles, etc.)

### 12. **Subscriptions non implémentées**

**Fichier:** `graph/schema.resolvers.go:1437-1446`

**Problème:** Les subscriptions retournent des channels vides.

```go
func (r *subscriptionResolver) OnNewSale(ctx context.Context) (<-chan *model.Sale, error) {
	ch := make(chan *model.Sale, 1)
	return ch, nil
}
```

**Recommandation:**
- Implémenter un système de pub/sub (Redis, NATS, etc.)
- Connecter les subscriptions aux événements réels (création de vente, commission)

### 13. **Conversion de types répétitive**

**Problème:** Beaucoup de code répétitif pour convertir entre modèles internes et GraphQL.

**Recommandation:**
- Créer des fonctions helper de conversion
- Utiliser des mappers automatiques (copier, etc.)

### 14. **Configuration hardcodée**

**Problème:** Certaines valeurs sont hardcodées dans le code.

**Recommandation:**
- Déplacer toutes les configurations vers des variables d'environnement
- Utiliser un fichier de configuration structuré

---

## ✅ Points Positifs

1. **Architecture propre:** Séparation claire entre services, repositories, et resolvers
2. **Microservices:** Bonne séparation avec le Tree Service
3. **Interfaces:** Utilisation d'interfaces pour faciliter les tests (repositories)
4. **Logging:** Utilisation de zap pour le logging structuré
5. **Mutex pour éviter les doubles paiements:** Bonne pratique dans `BinaryCommissionService`
6. **Version avec cache:** `GetLegsVolumesWithCache` montre une bonne réflexion sur la performance
7. **Documentation:** Bonne documentation dans `CHANGES_BUREAUMLMG.md`

---

## 📝 TODOs Identifiés

### Cache Redis non implémenté
**Fichier:** `services/tree-service/internal/cache/cache.go`

```go
// TODO: Implémenter avec go-redis
```

**Impact:** Performance dégradée pour les grands arbres

### BinaryCycle Repository manquant
**Fichier:** `internal/service/binary_commission_service.go:397`

```go
// TODO: Créer un repository pour BinaryCycle si nécessaire
```

**Impact:** Historique des cycles binaires non enregistré

---

## 🎯 Priorités d'Action

### 🔴 Urgent (À faire immédiatement)
1. Corriger le self-assignment dans `caisse_service.go`
2. Ajouter validation des mots de passe
3. Corriger les incohérences du schéma GraphQL (espaces, nullability)

### 🟡 Important (Cette semaine)
4. Implémenter les transactions atomiques pour les commissions
5. Améliorer la gestion d'erreurs (ne pas exposer les détails)
6. Ajouter des validations d'entrée

### 🟢 Amélioration (Ce mois)
7. Implémenter le cache Redis
8. Ajouter des tests pour les resolvers
9. Implémenter les subscriptions
10. Documenter le schéma GraphQL

---

## 📚 Ressources Recommandées

- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [Go Error Handling](https://go.dev/blog/error-handling-and-go)
- [MongoDB Transactions](https://www.mongodb.com/docs/manual/core/transactions/)
- [Redis Caching Patterns](https://redis.io/docs/manual/patterns/)

---

**Review effectué par:** Auto (AI Assistant)  
**Prochaine review recommandée:** Après correction des problèmes critiques



