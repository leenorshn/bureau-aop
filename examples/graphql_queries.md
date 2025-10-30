# GraphQL Queries Examples

Ce fichier contient des exemples de requêtes GraphQL pour tester l'API Bureau MLM.

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
      createdAt
    }
  }
}
```

### Refresh Token
```graphql
mutation {
  refreshToken(input: {
    token: "your-refresh-token-here"
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

## 👥 Gestion des Clients

### Créer un client (racine)
```graphql
mutation {
  clientCreate(input: {
    name: "John Doe"
    email: "john@example.com"
  }) {
    id
    name
    email
    sponsorId
    position
    joinDate
    totalEarnings
    walletBalance
    networkVolumeLeft
    networkVolumeRight
    binaryPairs
  }
}
```

### Créer un client avec sponsor
```graphql
mutation {
  clientCreate(input: {
    name: "Jane Smith"
    email: "jane@example.com"
    sponsorId: "64f8a1b2c3d4e5f6a7b8c9d0"
  }) {
    id
    name
    email
    sponsorId
    position
    joinDate
    totalEarnings
    walletBalance
    networkVolumeLeft
    networkVolumeRight
    binaryPairs
    sponsor {
      id
      name
      email
    }
  }
}
```

### Lister les clients
```graphql
query {
  clients(filter: { search: "john" }, paging: { page: 1, limit: 10 }) {
    id
    name
    email
    sponsorId
    position
    joinDate
    totalEarnings
    walletBalance
    networkVolumeLeft
    networkVolumeRight
    binaryPairs
    sponsor {
      id
      name
    }
    leftChild {
      id
      name
    }
    rightChild {
      id
      name
    }
  }
}
```

### Obtenir un client par ID
```graphql
query {
  client(id: "64f8a1b2c3d4e5f6a7b8c9d0") {
    id
    name
    email
    sponsorId
    position
    joinDate
    totalEarnings
    walletBalance
    networkVolumeLeft
    networkVolumeRight
    binaryPairs
    sponsor {
      id
      name
      email
    }
    leftChild {
      id
      name
      email
    }
    rightChild {
      id
      name
      email
    }
  }
}
```

### Mettre à jour un client
```graphql
mutation {
  clientUpdate(id: "64f8a1b2c3d4e5f6a7b8c9d0", input: {
    name: "John Updated"
    email: "john.updated@example.com"
  }) {
    id
    name
    email
  }
}
```

### Supprimer un client
```graphql
mutation {
  clientDelete(id: "64f8a1b2c3d4e5f6a7b8c9d0")
}
```

## 🛍️ Gestion des Produits

### Créer un produit
```graphql
mutation {
  productCreate(input: {
    name: "Produit Premium"
    description: "Description du produit premium"
    price: 100.0
    stock: 50
    imageUrl: "https://example.com/image.jpg"
  }) {
    id
    name
    description
    price
    stock
    imageUrl
    createdAt
    updatedAt
  }
}
```

### Lister les produits
```graphql
query {
  products(filter: { search: "premium" }, paging: { page: 1, limit: 10 }) {
    id
    name
    description
    price
    stock
    imageUrl
    createdAt
    updatedAt
  }
}
```

### Obtenir un produit par ID
```graphql
query {
  product(id: "64f8a1b2c3d4e5f6a7b8c9d0") {
    id
    name
    description
    price
    stock
    imageUrl
    createdAt
    updatedAt
  }
}
```

### Mettre à jour un produit
```graphql
mutation {
  productUpdate(id: "64f8a1b2c3d4e5f6a7b8c9d0", input: {
    name: "Produit Premium Updated"
    description: "Description mise à jour"
    price: 120.0
    stock: 75
    imageUrl: "https://example.com/new-image.jpg"
  }) {
    id
    name
    description
    price
    stock
    imageUrl
    updatedAt
  }
}
```

### Supprimer un produit
```graphql
mutation {
  productDelete(id: "64f8a1b2c3d4e5f6a7b8c9d0")
}
```

## 💰 Gestion des Ventes

### Créer une vente
```graphql
mutation {
  saleCreate(input: {
    clientId: "64f8a1b2c3d4e5f6a7b8c9d0"
    productId: "64f8a1b2c3d4e5f6a7b8c9d1"
    amount: 100.0
    note: "Vente manuelle"
  }) {
    id
    clientId
    sponsorId
    productId
    amount
    side
    date
    status
    note
    client {
      id
      name
      email
    }
    sponsor {
      id
      name
      email
    }
    product {
      id
      name
      price
    }
  }
}
```

### Lister les ventes
```graphql
query {
  sales(filter: { 
    status: "paid"
    dateFrom: "2024-01-01T00:00:00Z"
    dateTo: "2024-12-31T23:59:59Z"
  }, paging: { page: 1, limit: 10 }) {
    id
    clientId
    sponsorId
    productId
    amount
    side
    date
    status
    note
    client {
      id
      name
      email
    }
    sponsor {
      id
      name
      email
    }
    product {
      id
      name
      price
    }
  }
}
```

### Obtenir une vente par ID
```graphql
query {
  sale(id: "64f8a1b2c3d4e5f6a7b8c9d0") {
    id
    clientId
    sponsorId
    productId
    amount
    side
    date
    status
    note
    client {
      id
      name
      email
    }
    sponsor {
      id
      name
      email
    }
    product {
      id
      name
      price
    }
  }
}
```

## 💳 Gestion des Paiements

### Créer un paiement
```graphql
mutation {
  paymentCreate(input: {
    clientId: "64f8a1b2c3d4e5f6a7b8c9d0"
    amount: 100.0
    method: "mobile-money"
    description: "Paiement mobile money"
  }) {
    id
    clientId
    amount
    date
    method
    status
    description
    client {
      id
      name
      email
    }
  }
}
```

### Lister les paiements
```graphql
query {
  payments(filter: { 
    status: "completed"
    dateFrom: "2024-01-01T00:00:00Z"
  }, paging: { page: 1, limit: 10 }) {
    id
    clientId
    amount
    date
    method
    status
    description
    client {
      id
      name
      email
    }
  }
}
```

### Obtenir un paiement par ID
```graphql
query {
  payment(id: "64f8a1b2c3d4e5f6a7b8c9d0") {
    id
    clientId
    amount
    date
    method
    status
    description
    client {
      id
      name
      email
    }
  }
}
```

## 🎯 Gestion des Commissions

### Créer une commission manuelle
```graphql
mutation {
  commissionManualCreate(input: {
    clientId: "64f8a1b2c3d4e5f6a7b8c9d0"
    sourceClientId: "64f8a1b2c3d4e5f6a7b8c9d1"
    amount: 50.0
    level: 1
    type: "override"
  }) {
    id
    clientId
    sourceClientId
    amount
    level
    type
    date
    client {
      id
      name
      email
    }
    sourceClient {
      id
      name
      email
    }
  }
}
```

### Lister les commissions
```graphql
query {
  commissions(filter: { 
    dateFrom: "2024-01-01T00:00:00Z"
  }, paging: { page: 1, limit: 10 }) {
    id
    clientId
    sourceClientId
    amount
    level
    type
    date
    client {
      id
      name
      email
    }
    sourceClient {
      id
      name
      email
    }
  }
}
```

### Obtenir une commission par ID
```graphql
query {
  commission(id: "64f8a1b2c3d4e5f6a7b8c9d0") {
    id
    clientId
    sourceClientId
    amount
    level
    type
    date
    client {
      id
      name
      email
    }
    sourceClient {
      id
      name
      email
    }
  }
}
```

## 🔄 Opérations MLM

### Vérifier les commissions binaires
```graphql
mutation {
  runBinaryCommissionCheck(clientId: "64f8a1b2c3d4e5f6a7b8c9d0") {
    commissionsCreated
    totalAmount
    message
  }
}
```

## 📊 Statistiques du Dashboard

### Obtenir les statistiques
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

### Statistiques pour différentes périodes
```graphql
# 7 derniers jours
query {
  dashboardStats(range: "7d") {
    totalClients
    totalSales
    totalCommissions
    totalProducts
    activeClients
  }
}

# 90 derniers jours
query {
  dashboardStats(range: "90d") {
    totalClients
    totalSales
    totalCommissions
    totalProducts
    activeClients
  }
}

# 1 an
query {
  dashboardStats(range: "1y") {
    totalClients
    totalSales
    totalCommissions
    totalProducts
    activeClients
  }
}
```

## 🔍 Requêtes avec Filtres

### Recherche par texte
```graphql
query {
  clients(filter: { search: "john" }) {
    id
    name
    email
  }
}
```

### Filtrage par date
```graphql
query {
  sales(filter: { 
    dateFrom: "2024-01-01T00:00:00Z"
    dateTo: "2024-01-31T23:59:59Z"
  }) {
    id
    amount
    date
    client {
      name
    }
  }
}
```

### Filtrage par statut
```graphql
query {
  payments(filter: { status: "completed" }) {
    id
    amount
    status
    client {
      name
    }
  }
}
```

### Pagination
```graphql
query {
  clients(paging: { page: 2, limit: 5 }) {
    id
    name
    email
  }
}
```

## 🔄 Subscriptions (Temps réel)

### Écouter les nouvelles ventes
```graphql
subscription {
  onNewSale {
    id
    clientId
    amount
    date
    status
    client {
      name
      email
    }
  }
}
```

### Écouter les nouvelles commissions
```graphql
subscription {
  onNewCommission {
    id
    clientId
    amount
    type
    date
    client {
      name
      email
    }
  }
}
```

## 🧪 Tests de Performance

### Requête complexe avec relations
```graphql
query {
  clients(paging: { limit: 100 }) {
    id
    name
    email
    totalEarnings
    walletBalance
    networkVolumeLeft
    networkVolumeRight
    binaryPairs
    sponsor {
      id
      name
      email
    }
    leftChild {
      id
      name
      email
    }
    rightChild {
      id
      name
      email
    }
  }
}
```

### Requête d'agrégation
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

## 📝 Notes d'utilisation

1. **Authentification**: Toutes les mutations nécessitent un token d'accès valide
2. **Pagination**: Utilisez `page` et `limit` pour la pagination
3. **Filtres**: Les filtres sont optionnels et peuvent être combinés
4. **Dates**: Utilisez le format ISO 8601 pour les dates
5. **IDs**: Utilisez les ObjectIDs MongoDB (24 caractères hexadécimaux)
6. **Subscriptions**: Nécessitent une connexion WebSocket
7. **Performance**: Limitez le nombre de relations pour de meilleures performances

