# Bureau MLM Backend

Un serveur backend complet en Go pour une application d'administration MLM (marketing de réseau) avec API GraphQL, MongoDB et logique binaire.

## 🚀 Fonctionnalités

- **API GraphQL** complète avec gqlgen
- **Base de données MongoDB** avec collections optimisées
- **Authentification JWT** (access + refresh tokens)
- **Logique MLM binaire** avec placement automatique
- **Génération automatique de ventes** lors de l'ajout de clients
- **Calcul des commissions binaires** automatique
- **Gestion des paiements** et statistiques
- **Tests unitaires** et d'intégration
- **Docker** et docker-compose pour le déploiement

## 🏗️ Architecture

```
bureau/
├── cmd/server/           # Point d'entrée de l'application
├── internal/
│   ├── config/          # Configuration
│   ├── models/          # Modèles de données
│   ├── graphql/         # Schéma et resolvers GraphQL
│   ├── store/           # Repositories MongoDB
│   ├── service/         # Logique métier
│   └── auth/            # Authentification JWT
├── scripts/             # Scripts utilitaires
├── docker/              # Configuration Docker
└── tests/               # Tests
```

## 🛠️ Installation

### Prérequis

- Go 1.21+
- MongoDB (local ou Atlas)
- Docker (optionnel)

### Installation locale

1. **Cloner le repository**
```bash
git clone <repository-url>
cd bureau
```

2. **Installer les dépendances**
```bash
make deps
```

3. **Configurer l'environnement**
```bash
cp env.example .env
# Éditer .env avec vos paramètres
```

4. **Générer le code GraphQL**
```bash
make generate-gql
```

5. **Créer l'utilisateur admin**
```bash
make seed-admin
```

6. **Lancer l'application**
```bash
make run
```

### Installation avec Docker

1. **Lancer avec Docker Compose**
```bash
make docker-run
```

2. **Créer l'utilisateur admin**
```bash
make seed-admin
```

## 🔧 Configuration

### Variables d'environnement

```env
# MongoDB
MONGO_URI=mongodb+srv://<user>:<pass>@cluster0.mongodb.net/mlm?retryWrites=true&w=majority
MONGO_DB_NAME=mlm_db

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_REFRESH_SECRET=your-super-secret-refresh-key
JWT_ACCESS_EXP=15m
JWT_REFRESH_EXP=7d

# Admin
ADMIN_SEED_EMAIL=admin@example.com
ADMIN_SEED_PASSWORD=admin123

# Server
APP_PORT=4000
APP_ENV=development

# MLM Configuration
BINARY_THRESHOLD=100.0
BINARY_COMMISSION_RATE=0.1
DEFAULT_PRODUCT_PRICE=50.0
```

## 📊 Modèles de données

### Client
- Informations personnelles
- Structure binaire (sponsor, enfants gauche/droite)
- Volumes de réseau et commissions
- Portefeuille et gains

### Sale
- Ventes automatiques et manuelles
- Association client-sponsor
- Statuts et montants

### Commission
- Commissions binaires
- Niveaux et types
- Historique des gains

## 🔐 Authentification

### Login Admin
```graphql
mutation {
  adminLogin(input: {
    email: "admin@example.com"
    password: "admin123"
  }) {
    accessToken
    refreshToken
    admin {
      id
      name
      email
      role
    }
  }
}
```

### Refresh Token
```graphql
mutation {
  refreshToken(input: {
    token: "your-refresh-token"
  }) {
    accessToken
    refreshToken
    admin {
      id
      name
      email
    }
  }
}
```

## 🌐 API GraphQL

### Endpoints

- **GraphQL Playground**: http://localhost:4000
- **GraphQL Endpoint**: http://localhost:4000/query

### Exemples de requêtes

#### Créer un client
```graphql
mutation {
  clientCreate(input: {
    name: "John Doe"
    email: "john@example.com"
    sponsorId: "sponsor-id"
  }) {
    id
    name
    email
    sponsorId
    position
    networkVolumeLeft
    networkVolumeRight
  }
}
```

#### Obtenir les statistiques
```graphql
query {
  dashboardStats(range: "30d") {
    totalClients
    totalSales
    totalCommissions
    totalProducts
    activeClients
  }
}
```

#### Lister les clients
```graphql
query {
  clients(filter: { search: "john" }, paging: { page: 1, limit: 10 }) {
    id
    name
    email
    totalEarnings
    walletBalance
    binaryPairs
  }
}
```

## 🧪 Tests

### Lancer les tests
```bash
make test
```

### Tests unitaires
- Logique de placement binaire
- Calcul des commissions
- Authentification JWT

### Tests d'intégration
- Création de clients
- Génération de ventes
- Mise à jour des volumes

## 🐳 Docker

### Build
```bash
make docker-build
```

### Run
```bash
make docker-run
```

### Stop
```bash
make docker-stop
```

## 📈 Logique MLM Binaire

### Placement automatique
1. Nouveau client ajouté
2. Recherche de position dans l'arbre binaire
3. Placement en position gauche ou droite
4. Mise à jour des volumes de réseau

### Commissions binaires
1. Vérification des seuils (gauche et droite)
2. Calcul du montant de commission
3. Création de l'enregistrement de commission
4. Mise à jour des gains du client

### Génération de ventes
- Vente automatique lors de l'ajout d'un client
- Association avec le sponsor
- Mise à jour des volumes de réseau

## 🔍 Monitoring et logs

- Logs structurés avec Zap
- Métriques de performance
- Surveillance des erreurs

## 🚀 Déploiement

### Production
1. Configurer MongoDB Atlas
2. Définir les variables d'environnement
3. Build et déployer avec Docker
4. Configurer le reverse proxy (nginx)

### Variables de production
- `APP_ENV=production`
- `JWT_SECRET` sécurisé
- `MONGO_URI` Atlas
- Configuration SSL

## 📚 Documentation API

### GraphQL Schema
Le schéma GraphQL est défini dans `internal/graphql/schema.graphql`

### Types principaux
- `Product`: Produits
- `Client`: Clients avec structure binaire
- `Sale`: Ventes
- `Payment`: Paiements
- `Commission`: Commissions
- `Admin`: Administrateurs

## 🤝 Contribution

1. Fork le projet
2. Créer une branche feature
3. Commiter les changements
4. Push vers la branche
5. Ouvrir une Pull Request

## 📄 Licence

Ce projet est sous licence MIT.

## 🆘 Support

Pour toute question ou problème:
1. Vérifier les logs
2. Consulter la documentation
3. Ouvrir une issue GitHub

## 🔄 Changelog

### v1.0.0
- API GraphQL complète
- Logique MLM binaire
- Authentification JWT
- Tests unitaires
- Docker support


gcloud builds submit --tag gcr.io/bureaumlmg/bureau

