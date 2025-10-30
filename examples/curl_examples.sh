#!/bin/bash

# Exemples de requêtes cURL pour l'API Bureau MLM
# Remplacez les URLs et tokens par vos valeurs

BASE_URL="http://localhost:4000"
ACCESS_TOKEN="your-access-token-here"

# Fonction pour faire des requêtes GraphQL
graphql_request() {
    local query="$1"
    local variables="$2"
    
    curl -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "{\"query\": \"$query\", \"variables\": $variables}" \
        "$BASE_URL/query"
}

# 1. Login Admin
echo "=== Login Admin ==="
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "query": "mutation { adminLogin(input: { email: \"admin@example.com\", password: \"admin123\" }) { accessToken refreshToken admin { id name email role } } }"
    }' \
    "$BASE_URL/query"

echo -e "\n\n"

# 2. Créer un client racine
echo "=== Créer un client racine ==="
graphql_request 'mutation { clientCreate(input: { name: "John Doe", email: "john@example.com" }) { id name email sponsorId position joinDate totalEarnings walletBalance networkVolumeLeft networkVolumeRight binaryPairs } }'

echo -e "\n\n"

# 3. Créer un client avec sponsor
echo "=== Créer un client avec sponsor ==="
graphql_request 'mutation { clientCreate(input: { name: "Jane Smith", email: "jane@example.com", sponsorId: "CLIENT_ID_HERE" }) { id name email sponsorId position joinDate totalEarnings walletBalance networkVolumeLeft networkVolumeRight binaryPairs sponsor { id name email } } }'

echo -e "\n\n"

# 4. Lister les clients
echo "=== Lister les clients ==="
graphql_request 'query { clients(paging: { page: 1, limit: 10 }) { id name email sponsorId position joinDate totalEarnings walletBalance networkVolumeLeft networkVolumeRight binaryPairs sponsor { id name } leftChild { id name } rightChild { id name } } }'

echo -e "\n\n"

# 5. Créer un produit
echo "=== Créer un produit ==="
graphql_request 'mutation { productCreate(input: { name: "Produit Premium", description: "Description du produit", price: 100.0, stock: 50, imageUrl: "https://example.com/image.jpg" }) { id name description price stock imageUrl createdAt updatedAt } }'

echo -e "\n\n"

# 6. Lister les produits
echo "=== Lister les produits ==="
graphql_request 'query { products(paging: { page: 1, limit: 10 }) { id name description price stock imageUrl createdAt updatedAt } }'

echo -e "\n\n"

# 7. Créer une vente
echo "=== Créer une vente ==="
graphql_request 'mutation { saleCreate(input: { clientId: "CLIENT_ID_HERE", productId: "PRODUCT_ID_HERE", amount: 100.0, note: "Vente manuelle" }) { id clientId sponsorId productId amount side date status note client { id name email } sponsor { id name email } product { id name price } } }'

echo -e "\n\n"

# 8. Lister les ventes
echo "=== Lister les ventes ==="
graphql_request 'query { sales(paging: { page: 1, limit: 10 }) { id clientId sponsorId productId amount side date status note client { id name email } sponsor { id name email } product { id name price } } }'

echo -e "\n\n"

# 9. Créer un paiement
echo "=== Créer un paiement ==="
graphql_request 'mutation { paymentCreate(input: { clientId: "CLIENT_ID_HERE", amount: 100.0, method: "mobile-money", description: "Paiement mobile money" }) { id clientId amount date method status description client { id name email } } }'

echo -e "\n\n"

# 10. Lister les paiements
echo "=== Lister les paiements ==="
graphql_request 'query { payments(paging: { page: 1, limit: 10 }) { id clientId amount date method status description client { id name email } } }'

echo -e "\n\n"

# 11. Lister les commissions
echo "=== Lister les commissions ==="
graphql_request 'query { commissions(paging: { page: 1, limit: 10 }) { id clientId sourceClientId amount level type date client { id name email } sourceClient { id name email } } }'

echo -e "\n\n"

# 12. Obtenir les statistiques du dashboard
echo "=== Statistiques du dashboard ==="
graphql_request 'query { dashboardStats(range: "30d") { totalClients totalSales totalCommissions totalProducts activeClients } }'

echo -e "\n\n"

# 13. Vérifier les commissions binaires
echo "=== Vérifier les commissions binaires ==="
graphql_request 'mutation { runBinaryCommissionCheck(clientId: "CLIENT_ID_HERE") { commissionsCreated totalAmount message } }'

echo -e "\n\n"

# 14. Refresh token
echo "=== Refresh token ==="
curl -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "query": "mutation { refreshToken(input: { token: \"REFRESH_TOKEN_HERE\" }) { accessToken refreshToken admin { id name email } } }"
    }' \
    "$BASE_URL/query"

echo -e "\n\n"

# 15. Recherche avec filtres
echo "=== Recherche avec filtres ==="
graphql_request 'query { clients(filter: { search: "john" }, paging: { page: 1, limit: 5 }) { id name email totalEarnings walletBalance } }'

echo -e "\n\n"

# 16. Filtrage par date
echo "=== Filtrage par date ==="
graphql_request 'query { sales(filter: { dateFrom: "2024-01-01T00:00:00Z", dateTo: "2024-12-31T23:59:59Z" }, paging: { page: 1, limit: 10 }) { id amount date status client { name email } } }'

echo -e "\n\n"

# 17. Filtrage par statut
echo "=== Filtrage par statut ==="
graphql_request 'query { payments(filter: { status: "completed" }, paging: { page: 1, limit: 10 }) { id amount status date client { name email } } }'

echo -e "\n\n"

# 18. Obtenir un client par ID
echo "=== Obtenir un client par ID ==="
graphql_request 'query { client(id: "CLIENT_ID_HERE") { id name email sponsorId position joinDate totalEarnings walletBalance networkVolumeLeft networkVolumeRight binaryPairs sponsor { id name email } leftChild { id name email } rightChild { id name email } } }'

echo -e "\n\n"

# 19. Obtenir un produit par ID
echo "=== Obtenir un produit par ID ==="
graphql_request 'query { product(id: "PRODUCT_ID_HERE") { id name description price stock imageUrl createdAt updatedAt } }'

echo -e "\n\n"

# 20. Obtenir une vente par ID
echo "=== Obtenir une vente par ID ==="
graphql_request 'query { sale(id: "SALE_ID_HERE") { id clientId sponsorId productId amount side date status note client { id name email } sponsor { id name email } product { id name price } } }'

echo -e "\n\n"

# 21. Mettre à jour un client
echo "=== Mettre à jour un client ==="
graphql_request 'mutation { clientUpdate(id: "CLIENT_ID_HERE", input: { name: "John Updated", email: "john.updated@example.com" }) { id name email } }'

echo -e "\n\n"

# 22. Mettre à jour un produit
echo "=== Mettre à jour un produit ==="
graphql_request 'mutation { productUpdate(id: "PRODUCT_ID_HERE", input: { name: "Produit Updated", description: "Description mise à jour", price: 120.0, stock: 75, imageUrl: "https://example.com/new-image.jpg" }) { id name description price stock imageUrl updatedAt } }'

echo -e "\n\n"

# 23. Supprimer un client
echo "=== Supprimer un client ==="
graphql_request 'mutation { clientDelete(id: "CLIENT_ID_HERE") }'

echo -e "\n\n"

# 24. Supprimer un produit
echo "=== Supprimer un produit ==="
graphql_request 'mutation { productDelete(id: "PRODUCT_ID_HERE") }'

echo -e "\n\n"

# 25. Créer une commission manuelle
echo "=== Créer une commission manuelle ==="
graphql_request 'mutation { commissionManualCreate(input: { clientId: "CLIENT_ID_HERE", sourceClientId: "SOURCE_CLIENT_ID_HERE", amount: 50.0, level: 1, type: "override" }) { id clientId sourceClientId amount level type date client { id name email } sourceClient { id name email } } }'

echo -e "\n\n"

echo "=== Toutes les requêtes terminées ==="
