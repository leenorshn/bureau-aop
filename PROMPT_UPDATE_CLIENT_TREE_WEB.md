# Prompt: Mise à jour de l'Arbre de Clients sur le Web

## Contexte

Le système MLM binaire a été récemment mis à jour avec de nouvelles fonctionnalités de commissions binaires. L'interface web de l'arbre de clients doit être mise à jour pour afficher ces nouvelles informations et permettre une meilleure visualisation du réseau binaire.

## Modifications Backend Récentes

### 1. Nouveaux champs dans le modèle Client
- `networkVolumeLeft` : Volume total du réseau à gauche
- `networkVolumeRight` : Volume total du réseau à droite
- `binaryPairs` : Nombre de paires binaires formées
- `totalEarnings` : Gains totaux du membre
- `walletBalance` : Solde du portefeuille

### 2. Nouveau système de commissions binaires
- Calcul automatique des cycles basé sur les actifs (membres ayant fait au moins 1 vente)
- Qualification : nécessite 1 direct actif à gauche ET 1 direct actif à droite
- Limites journalières : maximum 4 cycles par jour (configurable)
- Valeur du cycle : 20$ par cycle (configurable)

### 3. Query GraphQL existante
La query `clientTree(id: ID!)` retourne actuellement :
```graphql
type ClientTree {
  root: ClientTreeNode!
  nodes: [ClientTreeNode!]!
  totalNodes: Int!
  maxLevel: Int!
}

type ClientTreeNode {
  id: ID!
  clientId: String!
  name: String!
  phone: String
  parentId: ID
  level: Int!
  position: String # "left" or "right"
}
```

## Tâches à Effectuer

### Phase 1 : Extension du Schema GraphQL

#### 1.1 Ajouter de nouveaux champs au type `ClientTreeNode`
Étendre le type `ClientTreeNode` dans `graph/schema.graphqls` pour inclure :
```graphql
type ClientTreeNode {
  id: ID!
  clientId: String!
  name: String!
  phone: String
  parentId: ID
  level: Int!
  position: String # "left" or "right"
  
  # Nouveaux champs à ajouter
  networkVolumeLeft: Float!
  networkVolumeRight: Float!
  binaryPairs: Int!
  totalEarnings: Float!
  walletBalance: Float!
  isActive: Boolean! # Indique si le membre a fait au moins 1 vente
  leftActives: Int! # Nombre d'actifs dans la jambe gauche
  rightActives: Int! # Nombre d'actifs dans la jambe droite
  isQualified: Boolean! # Est qualifié pour recevoir des commissions
  cyclesAvailable: Int # Cycles disponibles (calculé)
  cyclesPaidToday: Int # Cycles payés aujourd'hui
}
```

#### 1.2 Mettre à jour le resolver `ClientTree`
Dans `graph/schema.resolvers.go`, modifier la fonction `ClientTree` pour :
1. Récupérer les informations supplémentaires de chaque client
2. Calculer le nombre d'actifs dans chaque jambe (utiliser `BinaryCommissionService.countActivesInLeg`)
3. Vérifier la qualification de chaque membre (utiliser `BinaryCommissionService.checkQualification`)
4. Récupérer les informations de capping pour chaque membre
5. Calculer les cycles disponibles pour chaque membre

**Note** : Pour des raisons de performance, considérer :
- Limiter le calcul des actifs aux 3-4 premiers niveaux
- Mettre en cache les résultats si possible
- Ajouter un paramètre `depth` à la query pour limiter la profondeur

### Phase 2 : Mise à jour de l'Interface Web

#### 2.1 Affichage des informations dans l'arbre
Pour chaque nœud de l'arbre, afficher :

**Informations de base (déjà présentes)** :
- Nom du client
- ClientID
- Téléphone
- Position (gauche/droite)
- Niveau dans l'arbre

**Nouvelles informations à ajouter** :
- **Volumes de réseau** : Afficher `networkVolumeLeft` et `networkVolumeRight` avec des indicateurs visuels
  - Utiliser des barres de progression ou des badges colorés
  - Couleur verte si volumes équilibrés, orange si déséquilibrés
  - Afficher le ratio d'équilibre : `min(left, right) / max(left, right) * 100%`

- **Statut d'activité** : Badge "Actif" ou "Inactif" basé sur `isActive`
  - Vert pour actif (a fait au moins 1 vente)
  - Gris pour inactif

- **Qualification** : Badge "Qualifié" ou "Non qualifié"
  - Vert si `isQualified = true`
  - Rouge si `isQualified = false`
  - Tooltip expliquant pourquoi (manque direct gauche/droite)

- **Actifs par jambe** : Afficher `leftActives` et `rightActives`
  - Format : "G: 15 | D: 23"
  - Utiliser des icônes ou badges distincts

- **Cycles disponibles** : Afficher `cyclesAvailable`
  - Format : "Cycles: 15"
  - Badge avec couleur selon la disponibilité

- **Cycles payés aujourd'hui** : Afficher `cyclesPaidToday`
  - Format : "Aujourd'hui: 2/4"
  - Barre de progression montrant la limite journalière

- **Gains** : Afficher `totalEarnings` et `walletBalance`
  - Format : "Gains: 500$ | Wallet: 300$"
  - Utiliser des icônes monétaires

#### 2.2 Amélioration visuelle de l'arbre
- **Couleurs selon le statut** :
  - Bordure verte pour les membres qualifiés
  - Bordure orange pour les membres non qualifiés mais actifs
  - Bordure grise pour les membres inactifs

- **Indicateurs visuels** :
  - Icône de feuille pour les membres actifs
  - Icône de cadenas pour les membres non qualifiés
  - Badge de volume déséquilibré si `|left - right| > 50%`

- **Tooltips informatifs** :
  - Au survol d'un nœud, afficher un tooltip avec toutes les informations
  - Inclure la raison de non-qualification si applicable
  - Afficher les statistiques détaillées

#### 2.3 Panneau de détails latéral
Ajouter un panneau latéral qui s'ouvre au clic sur un nœud, affichant :
- Informations complètes du client
- Statistiques détaillées (volumes, actifs, cycles)
- Historique des commissions binaires
- Graphique d'évolution des volumes
- Actions disponibles (calculer commission, voir détails, etc.)

#### 2.4 Filtres et options d'affichage
Ajouter des filtres pour :
- Afficher uniquement les membres qualifiés
- Afficher uniquement les membres actifs
- Filtrer par niveau de déséquilibre des volumes
- Filtrer par nombre de cycles disponibles

#### 2.5 Légende et aide
Ajouter une légende expliquant :
- Les différentes couleurs et badges
- La signification des volumes
- Comment fonctionne la qualification
- Comment sont calculés les cycles

### Phase 3 : Fonctionnalités Interactives

#### 3.1 Calcul de commission en temps réel
Ajouter un bouton "Calculer Commission" sur chaque nœud qui :
- Appelle la mutation `runBinaryCommissionCheck(clientId: ID!)`
- Affiche le résultat dans une modal ou notification
- Met à jour automatiquement les informations affichées

#### 3.2 Actualisation automatique
- Option pour actualiser automatiquement l'arbre toutes les X minutes
- Indicateur visuel lors de l'actualisation
- Notification si de nouveaux cycles sont disponibles

#### 3.3 Export et impression
- Bouton pour exporter l'arbre en PDF ou image
- Option pour imprimer l'arbre avec toutes les informations

### Phase 4 : Optimisations et Performance

#### 4.1 Chargement progressif
- Charger d'abord les 3 premiers niveaux
- Charger les niveaux suivants à la demande (lazy loading)
- Afficher un indicateur de chargement

#### 4.2 Mise en cache
- Mettre en cache les résultats de calculs coûteux (actifs, qualification)
- Invalider le cache lors des mises à jour
- Utiliser des timestamps pour déterminer la fraîcheur des données

#### 4.3 Pagination virtuelle
- Pour les grands arbres, utiliser la pagination virtuelle
- Ne charger que les nœuds visibles à l'écran

## Exemple de Requête GraphQL Étendue

```graphql
query GetClientTree($id: ID!, $depth: Int) {
  clientTree(id: $id) {
    root {
      id
      clientId
      name
      phone
      level
      position
      networkVolumeLeft
      networkVolumeRight
      binaryPairs
      totalEarnings
      walletBalance
      isActive
      leftActives
      rightActives
      isQualified
      cyclesAvailable
      cyclesPaidToday
    }
    nodes {
      id
      clientId
      name
      phone
      parentId
      level
      position
      networkVolumeLeft
      networkVolumeRight
      binaryPairs
      totalEarnings
      walletBalance
      isActive
      leftActives
      rightActives
      isQualified
      cyclesAvailable
      cyclesPaidToday
    }
    totalNodes
    maxLevel
  }
}
```

## Exemple de Mutation pour Calculer une Commission

```graphql
mutation CalculateBinaryCommission($clientId: ID!) {
  runBinaryCommissionCheck(clientId: $clientId) {
    commissionsCreated
    totalAmount
    message
  }
}
```

## Technologies Recommandées

### Frontend
- **React** ou **Vue.js** pour la structure
- **D3.js** ou **React Flow** pour la visualisation de l'arbre
- **Apollo Client** ou **urql** pour GraphQL
- **Tailwind CSS** ou **Material-UI** pour le styling
- **Recharts** ou **Chart.js** pour les graphiques

### Bibliothèques utiles
- `react-flow-renderer` ou `vis-network` pour les arbres interactifs
- `react-tooltip` pour les tooltips
- `react-modal` pour les modals
- `date-fns` pour le formatage des dates

## Priorités d'Implémentation

1. **Priorité Haute** :
   - Extension du schema GraphQL
   - Affichage des volumes de réseau
   - Affichage du statut de qualification
   - Calcul de commission interactif

2. **Priorité Moyenne** :
   - Panneau de détails latéral
   - Filtres et options d'affichage
   - Amélioration visuelle avec couleurs et badges

3. **Priorité Basse** :
   - Export et impression
   - Actualisation automatique
   - Optimisations de performance avancées

## Notes Importantes

- **Performance** : Le calcul des actifs peut être coûteux pour les grands arbres. Considérer :
  - Limiter la profondeur par défaut
  - Calculer les actifs uniquement pour les nœuds visibles
  - Utiliser des agrégations MongoDB côté backend

- **Sécurité** : S'assurer que seuls les utilisateurs autorisés peuvent :
  - Voir l'arbre complet
  - Calculer des commissions
  - Voir les informations financières détaillées

- **UX** : L'interface doit être :
  - Intuitive et facile à naviguer
  - Responsive (mobile-friendly)
  - Accessible (WCAG 2.1)
  - Performante même avec de grands arbres

## Tests à Effectuer

1. **Tests fonctionnels** :
   - Vérifier l'affichage correct de toutes les nouvelles informations
   - Tester le calcul de commission depuis l'interface
   - Vérifier les filtres et options d'affichage

2. **Tests de performance** :
   - Tester avec des arbres de différentes tailles (10, 100, 1000+ nœuds)
   - Mesurer le temps de chargement
   - Vérifier l'utilisation mémoire

3. **Tests d'interface** :
   - Tester sur différents navigateurs
   - Tester sur mobile et tablette
   - Vérifier l'accessibilité

## Documentation à Créer

1. Guide utilisateur pour l'interface de l'arbre
2. Documentation technique de l'implémentation
3. Guide de maintenance et de dépannage








