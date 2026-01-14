#!/usr/bin/env bash
set -euo pipefail

# Script pour réinitialiser le mot de passe de l'admin admin@mlm.com
# Usage: ./scripts/reset-admin-password.sh [new_password] [api_url]
# - new_password (optionnel): nouveau mot de passe (doit respecter les règles de sécurité)
# - api_url (optionnel): URL de l'API (défaut: http://localhost:8080/query)

NEW_PASSWORD="${1:-Admin123@}"
API_URL="${2:-http://localhost:8080/query}"

echo "Réinitialisation du mot de passe pour admin@mlm.com"
echo "API: ${API_URL}"
echo ""

# Vérifier que le mot de passe respecte les règles
if [ ${#NEW_PASSWORD} -lt 8 ]; then
  echo "❌ Erreur: Le mot de passe doit contenir au moins 8 caractères"
  exit 1
fi

# Note: Pour utiliser cette mutation, vous devez être authentifié comme admin
# Si vous n'avez pas de token admin, vous devrez d'abord vous connecter avec un autre compte admin
# ou modifier directement la base de données MongoDB

echo "⚠️  Note: Cette mutation nécessite une authentification admin."
echo "Si vous n'avez pas de token admin, vous pouvez:"
echo "1. Vous connecter avec un autre compte admin"
echo "2. Ou modifier directement dans MongoDB (voir instructions ci-dessous)"
echo ""

# Si un token admin est fourni via variable d'environnement
if [[ -n "${ADMIN_TOKEN:-}" ]]; then
  echo "Utilisation du token admin fourni..."
  
  payload=$(cat <<EOF
{
  "query": "mutation { resetAdminPasswordByEmail(input: { email: \"admin@mlm.com\", newPassword: \"${NEW_PASSWORD}\" }) }"
}
EOF
)

  response=$(curl -sS -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    -d "${payload}" \
    "${API_URL}")

  echo "Réponse: ${response}"
  
  if echo "${response}" | grep -q "true"; then
    echo "✅ Mot de passe réinitialisé avec succès!"
    echo "Nouveau mot de passe: ${NEW_PASSWORD}"
  else
    echo "❌ Erreur lors de la réinitialisation"
    echo "${response}"
    exit 1
  fi
else
  echo "❌ Variable ADMIN_TOKEN non définie"
  echo ""
  echo "Pour réinitialiser via GraphQL, vous devez:"
  echo "1. Obtenir un token admin (se connecter avec un autre compte admin)"
  echo "2. Exporter le token: export ADMIN_TOKEN=votre_token"
  echo "3. Relancer ce script"
  echo ""
  echo "OU modifier directement dans MongoDB:"
  echo ""
  echo "1. Se connecter à MongoDB"
  echo "2. Exécuter:"
  echo "   use mlm_db"
  echo "   db.admins.updateOne("
  echo "     { email: \"admin@mlm.com\" },"
  echo "     { \$set: { passwordHash: \"<hash_bcrypt_du_nouveau_mot_de_passe>\" } }"
  echo "   )"
  echo ""
  echo "Pour générer le hash bcrypt, vous pouvez utiliser:"
  echo "  - Un script Go (voir scripts/generate-password-hash.go)"
  echo "  - Ou un outil en ligne"
  exit 1
fi



