# GraphQL Gateway

## Structure

```
gateway/
├── main.go                 # Point d'entrée du Gateway
├── go.mod                  # Dépendances Go
├── gqlgen.yml              # Configuration gqlgen
├── Dockerfile              # Configuration Docker
│
├── graph/                  # GraphQL Schema et Resolvers
│   ├── schema.graphqls     # Schéma GraphQL
│   ├── resolver.go         # Resolver principal
│   ├── schema.resolvers.go # Implémentations des resolvers
│   ├── generated.go        # Code généré (généré automatiquement)
│   └── model/
│       └── models_gen.go   # Modèles générés (généré automatiquement)
│
└── internal/               # Code interne
    ├── client/             # Clients pour les microservices
    │   └── tree_client.go
    ├── config/             # Configuration
    │   └── config.go
    └── models/             # Modèles internes
        └── tree.go
```

## Fonctionnalités

- **GraphQL Gateway**: Point d'entrée unique pour toutes les queries GraphQL
- **Routing**: Route les queries vers les microservices appropriés
- **Tree Service Integration**: Appelle le Tree Service pour `clientTree`

## Développement

### Générer le code GraphQL

```bash
go generate ./graph
```

### Lancer le Gateway

```bash
go run main.go
```

### Variables d'environnement

- `TREE_SERVICE_URL`: URL du Tree Service (défaut: http://localhost:8082)
- `GATEWAY_PORT`: Port du Gateway (défaut: 8080)

## Endpoints

- `GET /`: GraphQL Playground
- `POST /query`: GraphQL Endpoint
