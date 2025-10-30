#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/seed-yetu-clients.sh [password] [api_url]
# - password (optional): default "yetu@2025"
# - api_url (optional): default "http://localhost:8080/query"
#
# If your API requires admin auth, export ADMIN_TOKEN beforehand:
#   export ADMIN_TOKEN=eyJhbGciOi...

PASSWORD="${1:-yetu@2025}"
API_URL="${2:-http://localhost:8080/query}"

AUTH_HEADER_ARG=
if [[ -n "${ADMIN_TOKEN:-}" ]]; then
  AUTH_HEADER_ARG="-H Authorization: Bearer ${ADMIN_TOKEN}"
fi

echo "Seeding clients yetu1..yetu7 to ${API_URL}"

for i in {1..7}; do
  name="yetu${i}"
  echo "Creating client: ${name}"

  payload=$(cat <<EOF
{ "query": "mutation { clientCreate(input: { name: \"${name}\", password: \"${PASSWORD}\" }) { id name clientId } }" }
EOF
)

  if [[ -n "${AUTH_HEADER_ARG}" ]]; then
    response=$(curl -sS -X POST \
      -H "Content-Type: application/json" \
      ${AUTH_HEADER_ARG} \
      -d "${payload}" \
      "${API_URL}")
  else
    response=$(curl -sS -X POST \
      -H "Content-Type: application/json" \
      -d "${payload}" \
      "${API_URL}")
  fi

  echo "Response: ${response}"
  echo
done

echo "Done."


