# Réinitialiser le mot de passe de admin@mlm.com

## Méthode 1 : Via GraphQL (recommandé si vous avez un autre compte admin)

### Étape 1 : Se connecter avec un autre compte admin

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { userLogin(input: { email: \"autre_admin@example.com\", password: \"mot_de_passe\" }) { accessToken } }"
  }'
```

Copiez le `accessToken` de la réponse.

### Étape 2 : Réinitialiser le mot de passe de admin@mlm.com

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer VOTRE_TOKEN_ADMIN" \
  -d '{
    "query": "mutation { resetAdminPasswordByEmail(input: { email: \"admin@mlm.com\", newPassword: \"NouveauMotDePasse123@\" }) }"
  }'
```

**Important** : Le nouveau mot de passe doit respecter ces règles :
- Minimum 8 caractères
- Au moins une majuscule (A-Z)
- Au moins une minuscule (a-z)
- Au moins un chiffre (0-9)
- Au moins un caractère spécial parmi : `@$!%*?&`

## Méthode 2 : Via MongoDB (si vous n'avez pas d'autre compte admin)

### Étape 1 : Générer le hash du nouveau mot de passe

```bash
cd /Users/apple/Documents/devs/apis/bureau/scripts
go run generate-password-hash.go "NouveauMotDePasse123@"
```

Cela affichera le hash bcrypt du mot de passe.

### Étape 2 : Se connecter à MongoDB

```bash
mongosh "votre_connection_string"
```

### Étape 3 : Mettre à jour le mot de passe

```javascript
use mlm_db
db.admins.updateOne(
  { email: "admin@mlm.com" },
  { $set: { passwordHash: "LE_HASH_GÉNÉRÉ_À_L_ÉTAPE_1" } }
)
```

### Étape 4 : Vérifier

```javascript
db.admins.findOne({ email: "admin@mlm.com" })
```

Vous devriez voir le nouveau `passwordHash`.

## Méthode 3 : Utiliser le script fourni

```bash
# 1. Obtenir un token admin (voir Méthode 1, Étape 1)
export ADMIN_TOKEN="votre_token_admin"

# 2. Exécuter le script
cd /Users/apple/Documents/devs/apis/bureau
./scripts/reset-admin-password.sh "NouveauMotDePasse123@"
```

## Exemple complet avec GraphQL Playground

1. Ouvrez http://localhost:8080 dans votre navigateur
2. Dans l'onglet "HTTP HEADERS", ajoutez :
```json
{
  "Authorization": "Bearer VOTRE_TOKEN_ADMIN"
}
```
3. Exécutez la mutation :
```graphql
mutation {
  resetAdminPasswordByEmail(input: {
    email: "admin@mlm.com"
    newPassword: "NouveauMotDePasse123@"
  })
}
```

## Vérification

Après la réinitialisation, testez la connexion :

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { userLogin(input: { email: \"admin@mlm.com\", password: \"NouveauMotDePasse123@\" }) { accessToken user { id name email } } }"
  }'
```

Si la connexion réussit, le mot de passe a été correctement réinitialisé !



