# Améliorations Nécessaires pour les Tests

## 🔴 Problèmes Critiques à Corriger

### 1. **Configuration d'Environnement de Test**
**Problème**: Le fichier `env.test` n'est pas chargé automatiquement.

**Solution**: Modifier `SetupTestEnvironment` pour charger explicitement `env.test`:
```go
// Dans test_helpers.go
func SetupTestEnvironment(t *testing.T) *TestConfig {
    // Charger env.test explicitement
    if err := godotenv.Load("env.test"); err != nil {
        t.Logf("Warning: Could not load env.test: %v", err)
    }
    // ... reste du code
}
```

### 2. **Middleware d'Authentification Manquant**
**Problème**: Le serveur de test n'a pas le middleware d'authentification comme le serveur principal.

**Solution**: Ajouter le middleware dans `SetupTestEnvironment`:
```go
// Créer un wrapper HTTP avec middleware d'authentification
authMiddleware := func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
            token := authHeader[7:]
            claims, err := jwtService.ValidateAccessToken(token)
            if err == nil && claims != nil {
                // Ajouter au contexte GraphQL
                ctx := context.WithValue(r.Context(), "user", claims)
                r = r.WithContext(ctx)
            }
        }
        next.ServeHTTP(w, r)
    })
}

testServer := httptest.NewServer(authMiddleware(srv))
```

### 3. **Fermeture Propre de MongoDB**
**Problème**: Le client MongoDB n'est pas fermé proprement dans `TeardownTestEnvironment`.

**Solution**: Stocker le client MongoDB et le fermer:
```go
type TestConfig struct {
    MongoDB      *mongo.Database
    MongoClient  *mongo.Client  // Ajouter ce champ
    // ... autres champs
}

// Dans SetupTestEnvironment
return &TestConfig{
    MongoDB:     db,
    MongoClient: mongoClient,  // Stocker le client
    // ...
}

// Dans TeardownTestEnvironment
if tc.MongoClient != nil {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := tc.MongoClient.Disconnect(ctx); err != nil {
        t.Logf("Warning: Failed to disconnect MongoDB: %v", err)
    }
}
```

### 4. **Gestion des Variables GraphQL**
**Problème**: Les requêtes GraphQL utilisent des variables mais elles ne sont pas correctement formatées dans certains tests.

**Solution**: Utiliser des requêtes paramétrées correctement ou utiliser `fmt.Sprintf` de manière sécurisée:
```go
// Au lieu de variables GraphQL complexes, utiliser des requêtes directes
query := fmt.Sprintf(`
    mutation {
        productCreate(input: {
            name: "%s"
            description: "%s"
            price: %.2f
            stock: %d
            points: %.2f
            imageUrl: "%s"
        }) {
            id
        }
    }
`, name, description, price, stock, points, imageUrl)
```

## 🟡 Améliorations Importantes

### 5. **Tests avec MongoDB Local ou Docker**
**Problème**: Les tests nécessitent MongoDB mais il n'y a pas de fallback si MongoDB n'est pas disponible.

**Solution**: Ajouter un check et utiliser Docker Compose pour les tests:
```yaml
# docker-compose.test.yml
version: '3.8'
services:
  mongodb-test:
    image: mongo:latest
    ports:
      - "27018:27017"
    environment:
      MONGO_INITDB_DATABASE: mlm_test_db
```

### 6. **Fixtures de Test Réutilisables**
**Problème**: Chaque test crée ses propres données, ce qui est lent et répétitif.

**Solution**: Créer des fixtures réutilisables:
```go
// Dans test_helpers.go
type TestFixtures struct {
    RootClientID   string
    LeftClientID   string
    RightClientID  string
    TestProductID  string
    TestAdminID    string
}

func CreateTestFixtures(t *testing.T, tc *TestConfig) *TestFixtures {
    rootID := CreateTestClient(t, tc, "Root", nil)
    leftID := CreateTestClient(t, tc, "Left", &rootID)
    rightID := CreateTestClient(t, tc, "Right", &rootID)
    productID := CreateTestProduct(t, tc, "Test Product")
    
    return &TestFixtures{
        RootClientID:  rootID,
        LeftClientID:  leftID,
        RightClientID: rightID,
        TestProductID: productID,
        TestAdminID:   tc.TestAdminID,
    }
}
```

### 7. **Tests de Subscription WebSocket**
**Problème**: Les tests de subscription sont des placeholders.

**Solution**: Implémenter avec un client WebSocket:
```go
import "github.com/gorilla/websocket"

func TestSubscription_OnNewSale_Real(t *testing.T) {
    // Convertir httptest.Server en WebSocket
    wsURL := strings.Replace(tc.Server.URL, "http://", "ws://", 1) + "/query"
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    
    // Envoyer subscription
    subscription := `{"query": "subscription { onNewSale { id amount } }"}`
    conn.WriteMessage(websocket.TextMessage, []byte(subscription))
    
    // Créer une vente
    CreateTestSale(t, tc, clientID, productID, 100.0, "paid")
    
    // Lire message
    _, message, err := conn.ReadMessage()
    // Vérifier le message
}
```

### 8. **Tests de Performance avec Benchmarks**
**Problème**: Les tests de performance sont basiques.

**Solution**: Utiliser les benchmarks Go:
```go
func BenchmarkProductCreation(b *testing.B) {
    tc := SetupTestEnvironment(&testing.T{})
    defer TeardownTestEnvironment(&testing.T{}, tc)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        CreateTestProduct(&testing.T{}, tc, fmt.Sprintf("Product %d", i))
    }
}
```

### 9. **Couverture de Code**
**Problème**: Pas de suivi de la couverture de code.

**Solution**: Ajouter dans le Makefile:
```makefile
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
```

### 10. **Tests Parallèles**
**Problème**: Les tests ne sont pas parallélisables car ils partagent la même DB.

**Solution**: Utiliser des bases de données séparées par test:
```go
func SetupTestEnvironment(t *testing.T) *TestConfig {
    // Utiliser un nom de DB unique par test
    testDBName := fmt.Sprintf("mlm_test_%d_%d", time.Now().Unix(), t.Name())
    os.Setenv("MONGO_DB_NAME", testDBName)
    // ...
}
```

## 🟢 Améliorations Optionnelles

### 11. **Mocks pour Services Complexes**
Créer des mocks pour les services qui font des appels externes ou sont complexes.

### 12. **Tests de Charge avec Artillery ou k6**
Ajouter des tests de charge réels avec des outils dédiés.

### 13. **CI/CD Integration**
Ajouter les tests dans GitHub Actions ou GitLab CI.

### 14. **Documentation des Tests**
Créer un guide d'utilisation des tests pour les développeurs.

### 15. **Tests de Regression**
Ajouter des tests qui vérifient que les bugs passés ne réapparaissent pas.

## 📋 Checklist d'Implémentation

- [ ] Corriger le chargement de `env.test`
- [ ] Ajouter le middleware d'authentification dans les tests
- [ ] Fermer proprement MongoDB dans TeardownTestEnvironment
- [ ] Créer docker-compose.test.yml pour MongoDB de test
- [ ] Implémenter les fixtures réutilisables
- [ ] Implémenter les tests WebSocket pour subscriptions
- [ ] Ajouter des benchmarks de performance
- [ ] Configurer la couverture de code
- [ ] Rendre les tests parallélisables
- [ ] Ajouter des tests dans CI/CD

## 🚀 Ordre de Priorité

1. **Urgent**: Points 1, 2, 3 (les tests ne fonctionneront pas sans ces corrections)
2. **Important**: Points 4, 5, 6 (améliorent la stabilité et la vitesse)
3. **Souhaitable**: Points 7-15 (améliorent la qualité globale)


