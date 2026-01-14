# Guide de Réinitialisation de Mot de Passe

Ce document explique comment réinitialiser un mot de passe oublié dans le système Bureau MLM.

## Options Disponibles

### 1. Changer son propre mot de passe (utilisateur authentifié)

Si vous êtes connecté (admin ou client), vous pouvez changer votre propre mot de passe en utilisant la mutation `changePassword`.

**Mutation GraphQL :**
```graphql
mutation {
  changePassword(input: {
    currentPassword: "ancien_mot_de_passe"
    newPassword: "NouveauMotDePasse123@"
  })
}
```

**Headers requis :**
```
Authorization: Bearer <votre_token_jwt>
```

**Exigences pour le nouveau mot de passe :**
- Minimum 8 caractères
- Au moins une lettre majuscule
- Au moins une lettre minuscule
- Au moins un chiffre
- Au moins un caractère spécial parmi : `@$!%*?&`

### 2. Réinitialiser le mot de passe d'un admin (admin uniquement)

Un administrateur peut réinitialiser le mot de passe d'un autre administrateur.

**Option A : Par ID**
```graphql
mutation {
  resetAdminPassword(input: {
    id: "507f1f77bcf86cd799439011"
    newPassword: "NouveauMotDePasse123@"
  })
}
```

**Option B : Par email (plus pratique)**
```graphql
mutation {
  resetAdminPasswordByEmail(input: {
    email: "admin@mlm.com"
    newPassword: "NouveauMotDePasse123@"
  })
}
```

**Headers requis :**
```
Authorization: Bearer <token_admin>
```

### 3. Réinitialiser le mot de passe d'un client (admin uniquement)

Un administrateur peut réinitialiser le mot de passe d'un client en utilisant son `clientId`.

**Mutation GraphQL :**
```graphql
mutation {
  resetClientPassword(input: {
    clientId: "12345678"
    newPassword: "NouveauMotDePasse123@"
  })
}
```

**Headers requis :**
```
Authorization: Bearer <token_admin>
```

## Exemples avec cURL

### Changer son propre mot de passe (admin)
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "query": "mutation { changePassword(input: { currentPassword: \"ancien\", newPassword: \"Nouveau123@\" }) }"
  }'
```

### Réinitialiser le mot de passe d'un admin par email
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "query": "mutation { resetAdminPasswordByEmail(input: { email: \"admin@mlm.com\", newPassword: \"Nouveau123@\" }) }"
  }'
```

### Réinitialiser le mot de passe d'un client (admin)
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "query": "mutation { resetClientPassword(input: { clientId: \"12345678\", newPassword: \"Nouveau123@\" }) }"
  }'
```

## Réinitialisation du compte admin@mlm.com

Si vous avez oublié le mot de passe du compte `admin@mlm.com`, voici les options :

### Option 1 : Via GraphQL (si vous avez un autre compte admin)

1. Connectez-vous avec un autre compte admin
2. Utilisez la mutation `resetAdminPasswordByEmail` :

```bash
# 1. Se connecter avec un autre admin
TOKEN=$(curl -s -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { userLogin(input: { email: \"autre_admin@example.com\", password: \"mot_de_passe\" }) { accessToken } }"
  }' | jq -r '.data.userLogin.accessToken')

# 2. Réinitialiser le mot de passe de admin@mlm.com
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "query": "mutation { resetAdminPasswordByEmail(input: { email: \"admin@mlm.com\", newPassword: \"NouveauMotDePasse123@\" }) }"
  }'
```

### Option 2 : Via MongoDB (si vous n'avez pas d'autre compte admin)

1. Connectez-vous à MongoDB
2. Générez le hash du nouveau mot de passe :
```bash
cd scripts
go run generate-password-hash.go "NouveauMotDePasse123@"
```

3. Mettez à jour dans MongoDB :
```javascript
use mlm_db
db.admins.updateOne(
  { email: "admin@mlm.com" },
  { $set: { passwordHash: "<hash_généré>" } }
)
```

### Option 3 : Utiliser le script fourni

```bash
# Avec un token admin
export ADMIN_TOKEN="votre_token_admin"
./scripts/reset-admin-password.sh "NouveauMotDePasse123@"
```

## Notes Importantes

1. **Sécurité** : Toutes les mutations de réinitialisation nécessitent une authentification admin, sauf `changePassword` qui peut être utilisée par n'importe quel utilisateur authentifié pour changer son propre mot de passe.

2. **Validation** : Le nouveau mot de passe doit respecter les règles de sécurité (minimum 8 caractères, majuscule, minuscule, chiffre, caractère spécial).

3. **Hachage** : Les mots de passe sont automatiquement hashés avec bcrypt avant d'être stockés en base de données.

4. **En cas d'oubli complet** : Si vous avez oublié votre mot de passe et que vous n'êtes pas connecté, vous devez contacter un administrateur pour qu'il réinitialise votre mot de passe.

## Dépannage

- **Erreur "authentification requise"** : Vérifiez que vous avez inclus le header `Authorization: Bearer <token>` dans votre requête.
- **Erreur "mot de passe actuel incorrect"** : Vérifiez que vous avez saisi correctement votre mot de passe actuel.
- **Erreur de validation du mot de passe** : Assurez-vous que le nouveau mot de passe respecte toutes les exigences de sécurité.

