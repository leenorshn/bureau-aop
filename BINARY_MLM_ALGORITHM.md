# Algorithme MLM Binaire - Documentation Complète

## Vue d'ensemble

Cet algorithme calcule et paie les commissions binaires dans un système MLM basé sur un arbre binaire. Chaque membre possède une jambe gauche et une jambe droite, et les gains sont générés par des cycles (1 actif gauche + 1 actif droite).

## Pseudo-code de l'algorithme principal

```
FONCTION ComputeBinaryCommission(clientID):
    // 1. Vérifier que le client existe
    client = GetClient(clientID)
    SI client == NULL:
        RETOURNER erreur "Client introuvable"
    
    // 2. Vérifier la qualification
    qualification = CheckQualification(client)
    SI qualification.IsQualified == FALSE:
        RETOURNER {
            success: true,
            qualified: false,
            reason: "Membre non qualifié: doit avoir au moins 1 direct actif à gauche ET 1 direct actif à droite"
        }
    
    // 3. Lire les volumes des jambes
    legs = GetLegsVolumes(client)
    
    // 4. Vérifier les conditions de base
    SI legs.LeftActives == 0 OU legs.RightActives == 0:
        RETOURNER {
            success: true,
            qualified: true,
            reason: "Jambe gauche ou droite vide - aucun cycle possible"
        }
    
    // 5. Calculer les cycles possibles
    cyclesAvailable = MIN(legs.LeftActives, legs.RightActives)
    
    SI cyclesAvailable == 0:
        RETOURNER {
            success: true,
            qualified: true,
            reason: "Aucun cycle disponible - volumes insuffisants"
        }
    
    // 6. Appliquer la limite journalière
    cyclesToPay = ApplyDailyLimit(clientID, cyclesAvailable)
    
    SI cyclesToPay == 0:
        RETOURNER {
            success: true,
            qualified: true,
            reason: "Limite journalière atteinte"
        }
    
    // 7. Calculer le montant
    amount = cyclesToPay × CycleValue
    
    // 8. Enregistrer le paiement (avec mutex pour éviter double paiement)
    VERROUILLER mutex
        // Double vérification après verrouillage
        cyclesToPayFinal = ApplyDailyLimit(clientID, cyclesAvailable)
        
        SI cyclesToPayFinal == 0:
            DÉVERROUILLER mutex
            RETOURNER { reason: "Limite journalière atteinte (double vérification)" }
        
        amount = cyclesToPayFinal × CycleValue
        
        // 9. Créer la commission
        commission = RecordPayment(clientID, cyclesToPayFinal, amount)
        
        // 10. Déduire les volumes utilisés
        (leftRemaining, rightRemaining) = DeductVolume(clientID, legs, cyclesToPayFinal)
        
        // 11. Mettre à jour les gains du client
        UpdateClientEarnings(clientID, amount)
    DÉVERROUILLER mutex
    
    RETOURNER {
        success: true,
        qualified: true,
        cyclesAvailable: cyclesAvailable,
        cyclesPaid: cyclesToPayFinal,
        amount: amount,
        leftVolumeRemaining: leftRemaining,
        rightVolumeRemaining: rightRemaining,
        commissionId: commission.ID
    }
FIN FONCTION

FONCTION CheckQualification(client):
    qualification = {
        isQualified: false,
        hasDirectLeft: false,
        hasDirectRight: false
    }
    
    // Vérifier direct gauche
    SI client.LeftChildID != NULL:
        leftChild = GetClient(client.LeftChildID)
        SI leftChild != NULL ET IsClientActive(leftChild):
            qualification.hasDirectLeft = true
            qualification.directLeftCount = 1
    
    // Vérifier direct droite
    SI client.RightChildID != NULL:
        rightChild = GetClient(client.RightChildID)
        SI rightChild != NULL ET IsClientActive(rightChild):
            qualification.hasDirectRight = true
            qualification.directRightCount = 1
    
    // Qualification = avoir les deux
    qualification.isQualified = qualification.hasDirectLeft ET qualification.hasDirectRight
    
    RETOURNER qualification
FIN FONCTION

FONCTION CalculateCycles(legs):
    SI legs.LeftActives == 0 OU legs.RightActives == 0:
        RETOURNER 0
    
    RETOURNER MIN(legs.LeftActives, legs.RightActives)
FIN FONCTION

FONCTION ApplyDailyLimit(clientID, cyclesAvailable):
    SI DailyCycleLimit <= 0:
        RETOURNER cyclesAvailable // Pas de limite
    
    capping = GetOrCreateCapping(clientID, TODAY)
    
    SI capping.CyclesPaidToday >= DailyCycleLimit:
        RETOURNER 0
    
    remainingLimit = DailyCycleLimit - capping.CyclesPaidToday
    cyclesToPay = MIN(cyclesAvailable, remainingLimit)
    
    capping.CyclesPaidToday += cyclesToPay
    UpdateCapping(capping)
    
    RETOURNER cyclesToPay
FIN FONCTION
```

## Règles MLM implémentées

### 1. Calcul des cycles
- **1 cycle** = 1 actif gauche + 1 actif droite
- **Valeur d'un cycle** = 20$ (configurable via `BinaryConfig.CycleValue`)
- **cycles** = min(gaucheActifs, droiteActifs)
- **gain** = cycles × valeurCycle

### 2. Conditions pour être payé

#### A - Le parrain doit avoir :
- Au moins 1 actif à gauche
- Au moins 1 actif à droite

#### B - Qualification personnelle :
- 1 filleul direct actif à gauche
- 1 filleul direct actif à droite
- Sinon : non qualifié → gain = 0

#### C - Volume minimum :
- Si une jambe manque, gain = 0

#### D - Limitation journalière/hebdomadaire :
- Ex: 4 cycles/jour
- Si cycles > limite → cyclesPayés = limite

### 3. Cas où le gain = 0
- Jambe gauche vide
- Jambe droite vide
- Non qualifié (moins de 2 directs)
- Volume insuffisant
- Limite journalière atteinte

## Structure des données

### BinaryConfig
```go
type BinaryConfig struct {
    CycleValue        float64 // Valeur d'un cycle en $ (ex: 20$)
    DailyCycleLimit   int     // Limite de cycles par jour (ex: 4)
    WeeklyCycleLimit  int     // Limite de cycles par semaine (optionnel)
    MinVolumePerLeg   float64 // Volume minimum par jambe
    RequireDirectLeft bool    // Requiert 1 direct actif à gauche
    RequireDirectRight bool   // Requiert 1 direct actif à droite
}
```

### BinaryCommissionResult
```go
type BinaryCommissionResult struct {
    Success              bool    // Succès de l'opération
    Qualified            bool    // Est qualifié ou non
    CyclesAvailable      int     // Cycles possibles avant limite
    CyclesPaid           int     // Cycles effectivement payés
    Amount               float64 // Montant gagné
    LeftVolumeRemaining  float64 // Volume gauche restant
    RightVolumeRemaining float64 // Volume droite restant
    Reason               string  // Raison si gain = 0
    CommissionID         *string // ID de la commission créée
}
```

## Exemples de tests

### Cas 1: 50 gauche, 100 droite → cycles = 50 → gain = 1000$
```go
legs := BinaryLegs{
    LeftActives:  50,
    RightActives: 100,
}
cycles := min(50, 100) = 50
amount := 50 × 20$ = 1000$
```

### Cas 2: 3 gauche, 5 droite → cycles = 3 → gain = 60$
```go
legs := BinaryLegs{
    LeftActives:  3,
    RightActives: 5,
}
cycles := min(3, 5) = 3
amount := 3 × 20$ = 60$
```

### Cas 3: 0 gauche, 10 droite → gain = 0
```go
legs := BinaryLegs{
    LeftActives:  0,
    RightActives: 10,
}
// Jambe gauche vide → gain = 0
```

### Cas 4: Non qualifié → gain = 0
```go
// Pas de direct gauche OU pas de direct droite
qualification := {
    hasDirectLeft: false,
    hasDirectRight: true,
}
// Non qualifié → gain = 0
```

### Cas 5: Limite journalière 4 cycles → payer 4 cycles
```go
cyclesAvailable := 10
dailyLimit := 4
cyclesPaid := min(10, 4) = 4
amount := 4 × 20$ = 80$
```

## Utilisation

```go
// 1. Créer la configuration
config := models.BinaryConfig{
    CycleValue:        20.0,
    DailyCycleLimit:   4,
    MinVolumePerLeg:   1.0,
    RequireDirectLeft: true,
    RequireDirectRight: true,
}

// 2. Créer le service
service := NewBinaryCommissionService(
    clientRepo,
    commissionRepo,
    saleRepo,
    cappingRepo,
    logger,
    config,
)

// 3. Calculer la commission
result, err := service.ComputeBinaryCommission(ctx, clientID)
if err != nil {
    // Gérer l'erreur
}

// 4. Vérifier le résultat
if result.Success && result.Qualified {
    fmt.Printf("Cycles payés: %d, Montant: %.2f$\n", result.CyclesPaid, result.Amount)
} else {
    fmt.Printf("Raison: %s\n", result.Reason)
}
```

## Sécurité et concurrence

- **Mutex** : Utilisé pour éviter les doubles paiements
- **Double vérification** : Après verrouillage, on revérifie la limite
- **Transactions atomiques** : Les opérations DB sont atomiques
- **Validation** : Toutes les conditions sont vérifiées avant paiement

## Performance

- **Comptage récursif optimisé** : Utilise une queue au lieu de récursion
- **Cache des actifs** : Évite de recalculer les actifs à chaque fois
- **Requêtes parallèles** : Possibilité d'optimiser avec goroutines

## Maintenance

- Code modulaire et testé
- Logique explicite et commentée
- Facile à étendre (ajout de nouvelles règles)
- Compatible MongoDB avec agrégations optimisées










