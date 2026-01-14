#!/bin/bash

# Script de test pour la query clientTree
# Usage: ./test_client_tree.sh <client_id>

CLIENT_ID="${1:-6906e2ca634b66b9c3fb7a07}"
ENDPOINT="${GRAPHQL_ENDPOINT:-http://localhost:8080/query}"

echo "Testing clientTree query with ID: $CLIENT_ID"
echo "Endpoint: $ENDPOINT"
echo ""

# Query simplifiée pour test
QUERY='{
  "query": "query clientTreeQ($id: ID!) { clientTree(id: $id) { root { id name clientId phone } nodes { id name clientId position } totalNodes maxLevel } }",
  "variables": {
    "id": "'"$CLIENT_ID"'"
  }
}'

echo "Sending query..."
echo ""

# Tester avec curl
curl -X POST \
  -H "Content-Type: application/json" \
  -d "$QUERY" \
  "$ENDPOINT" \
  | jq '.'

echo ""
echo "Test completed."







