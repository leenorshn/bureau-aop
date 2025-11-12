#!/bin/bash

echo "=== Test du conteneur Docker ==="

# Construire l'image avec la plateforme spécifique
echo "1. Construction de l'image..."
docker build --platform linux/amd64 -f Dockerfile.local -t bureau-backend .

# Démarrer le conteneur
echo "2. Démarrage du conteneur..."
docker run -d --name bureau-test -p 4000:4000 --env-file env.example bureau-backend

# Attendre que le serveur démarre
echo "3. Attente du démarrage du serveur..."
sleep 10

# Tester l'endpoint
echo "4. Test de l'endpoint GraphQL..."
curl -s http://localhost:4000/ | head -5

# Tester une requête GraphQL simple
echo "5. Test d'une requête GraphQL..."
curl -X POST http://localhost:4000/query \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __schema { types { name } } }"}' | head -5

# Nettoyer
echo "6. Nettoyage..."
docker stop bureau-test
docker rm bureau-test

echo "=== Test terminé ==="



















