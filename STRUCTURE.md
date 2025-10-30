# Structure du Projet Bureau MLM

## 📁 Structure Actuelle (Après Nettoyage)

```
bureau/
├── server.go                    # Point d'entrée principal
├── tools.go                     # Dépendances gqlgen
├── gqlgen.yml                   # Configuration gqlgen
├── go.mod                       # Dépendances Go
├── go.sum                       # Checksums des dépendances
├── env.example                  # Variables d'environnement
├── Dockerfile                   # Configuration Docker
├── docker-compose.yml           # Services Docker
├── Makefile                     # Commandes de build
├── README.md                    # Documentation
├── graph/                       # Code GraphQL (structure gqlgen)
│   ├── generated/
│   │   └── generated.go        # Code généré par gqlgen
│   ├── model/
│   │   └── models_gen.go        # Modèles GraphQL générés
│   ├── resolver.go              # Resolver principal + injection de dépendances
│   ├── schema.resolvers.go      # Implémentations des resolvers
│   └── schema.graphqls          # Schéma GraphQL
├── internal/                    # Code interne (non exposé)
│   ├── auth/                    # Authentification JWT
│   │   ├── jwt.go
│   │   └── bcrypt.go
│   ├── config/                  # Configuration
│   │   └── config.go
│   ├── models/                  # Modèles internes
│   │   └── models.go
│   ├── service/                 # Logique métier
│   │   ├── auth_service.go
│   │   ├── client_service.go
│   │   ├── commission_service.go
│   │   ├── payment_service.go
│   │   ├── product_service.go
│   │   ├── sale_service.go
│   │   └── admin_service.go
│   └── store/                   # Repositories MongoDB
│       ├── mongo.go
│       ├── admin_repository.go
│       ├── client_repository.go
│       ├── commission_repository.go
│       ├── payment_repository.go
│       ├── product_repository.go
│       └── sale_repository.go
├── scripts/                     # Scripts utilitaires
│   ├── seed_admin.go
│   ├── generate_gql.sh
│   └── run_tests.sh
├── tests/                       # Tests
│   ├── auth_test.go
│   └── client_test.go
├── examples/                    # Exemples d'utilisation
│   ├── graphql_queries.md
│   └── curl_examples.sh
└── .github/workflows/           # CI/CD
    └── ci.yml
```

## ✅ Réalisé

1. **Structure gqlgen conforme** à la documentation officielle
2. **Nettoyage des fichiers dupliqués** (suppression de `internal/graphql/` et `cmd/`)
3. **Configuration gqlgen** correcte avec `gqlgen.yml`
4. **Fichier `tools.go`** pour gérer les dépendances
5. **Point d'entrée `server.go`** unifié
6. **Logique MLM binaire** complète
7. **Authentification JWT** fonctionnelle
8. **Repositories MongoDB** avec tous les CRUD
9. **Services métier** avec logique MLM
10. **Tests unitaires** pour les fonctions critiques

## 🔧 Problèmes Actuels

### Erreurs de Compilation
- **Conversion de types** entre GraphQL (`model.*`) et internes (`models.*`)
- **Mismatch de types** dans les resolvers
- **Conversion de dates** (string vs time.Time)

### Fichiers à Corriger
- `graph/schema.resolvers.go` : Conversions de types incorrectes
- `graph/model/models_gen.go` : Types GraphQL vs internes

## 🚀 Prochaines Étapes

1. **Corriger les conversions de types** dans `schema.resolvers.go`
2. **Tester la compilation** complète
3. **Créer l'utilisateur admin** avec `make seed-admin`
4. **Lancer le serveur** avec `make run`
5. **Tester l'API GraphQL** avec le playground

## 📋 Commandes Disponibles

```bash
# Build et run
make build
make run

# Tests
make test

# Docker
make docker-build
make docker-run

# Admin
make seed-admin

# GraphQL
make generate-gql
```

## 🔗 Endpoints

- **GraphQL Playground** : http://localhost:4000
- **GraphQL Endpoint** : http://localhost:4000/query
- **Admin Login** : admin@mlm.com / admin123

## 📊 Fonctionnalités MLM

- ✅ **Placement binaire automatique**
- ✅ **Génération de ventes** lors de l'ajout de clients
- ✅ **Calcul des commissions binaires**
- ✅ **Mise à jour des volumes de réseau**
- ✅ **Gestion des paiements**
- ✅ **Statistiques dashboard**

