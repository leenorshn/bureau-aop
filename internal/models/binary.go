package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BinaryConfig représente la configuration du système binaire MLM
type BinaryConfig struct {
	CycleValue         float64 `bson:"cycleValue" json:"cycleValue"`                 // Valeur d'un cycle en $ (ex: 20$)
	DailyCycleLimit    int     `bson:"dailyCycleLimit" json:"dailyCycleLimit"`       // Limite de cycles par jour (ex: 4)
	WeeklyCycleLimit   int     `bson:"weeklyCycleLimit" json:"weeklyCycleLimit"`     // Limite de cycles par semaine (optionnel)
	MinVolumePerLeg    float64 `bson:"minVolumePerLeg" json:"minVolumePerLeg"`       // Volume minimum par jambe pour être payé
	RequireDirectLeft  bool    `bson:"requireDirectLeft" json:"requireDirectLeft"`   // Requiert 1 direct actif à gauche
	RequireDirectRight bool    `bson:"requireDirectRight" json:"requireDirectRight"` // Requiert 1 direct actif à droite
}

// BinaryLegs représente les jambes gauche et droite d'un membre
type BinaryLegs struct {
	LeftVolume   float64 `bson:"leftVolume" json:"leftVolume"`     // Volume total de la jambe gauche
	RightVolume  float64 `bson:"rightVolume" json:"rightVolume"`   // Volume total de la jambe droite
	LeftActives  int     `bson:"leftActives" json:"leftActives"`   // Nombre d'actifs à gauche
	RightActives int     `bson:"rightActives" json:"rightActives"` // Nombre d'actifs à droite
}

// BinaryQualification représente la qualification d'un membre pour recevoir des commissions
type BinaryQualification struct {
	IsQualified      bool `bson:"isQualified" json:"isQualified"`           // Est qualifié ou non
	HasDirectLeft    bool `bson:"hasDirectLeft" json:"hasDirectLeft"`       // A au moins 1 direct actif à gauche
	HasDirectRight   bool `bson:"hasDirectRight" json:"hasDirectRight"`     // A au moins 1 direct actif à droite
	DirectLeftCount  int  `bson:"directLeftCount" json:"directLeftCount"`   // Nombre de directs actifs à gauche
	DirectRightCount int  `bson:"directRightCount" json:"directRightCount"` // Nombre de directs actifs à droite
}

// BinaryCycle représente un cycle calculé et payé
type BinaryCycle struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID        primitive.ObjectID `bson:"clientId" json:"clientId"`
	Cycles          int                `bson:"cycles" json:"cycles"`                   // Nombre de cycles payés
	Amount          float64            `bson:"amount" json:"amount"`                   // Montant gagné
	LeftVolumeUsed  float64            `bson:"leftVolumeUsed" json:"leftVolumeUsed"`   // Volume gauche utilisé
	RightVolumeUsed float64            `bson:"rightVolumeUsed" json:"rightVolumeUsed"` // Volume droite utilisé
	Date            time.Time          `bson:"date" json:"date"`                       // Date du calcul
	ProcessedAt     time.Time          `bson:"processedAt" json:"processedAt"`         // Date de traitement
}

// BinaryCapping représente les limites journalières/hebdomadaires d'un membre
type BinaryCapping struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID           primitive.ObjectID `bson:"clientId" json:"clientId"`
	Date               time.Time          `bson:"date" json:"date"`                             // Date (pour limite journalière)
	WeekStart          time.Time          `bson:"weekStart" json:"weekStart"`                   // Début de semaine (pour limite hebdomadaire)
	CyclesPaidToday    int                `bson:"cyclesPaidToday" json:"cyclesPaidToday"`       // Cycles payés aujourd'hui
	CyclesPaidThisWeek int                `bson:"cyclesPaidThisWeek" json:"cyclesPaidThisWeek"` // Cycles payés cette semaine
	LastResetDate      time.Time          `bson:"lastResetDate" json:"lastResetDate"`           // Dernière date de reset
}

// BinaryCommissionResult représente le résultat du calcul de commission binaire
type BinaryCommissionResult struct {
	Success              bool    `json:"success"`
	Qualified            bool    `json:"qualified"`
	CyclesAvailable      int     `json:"cyclesAvailable"`        // Cycles possibles avant limite
	CyclesPaid           int     `json:"cyclesPaid"`             // Cycles effectivement payés (après limite)
	Amount               float64 `json:"amount"`                 // Montant gagné
	LeftVolumeRemaining  float64 `json:"leftVolumeRemaining"`    // Volume gauche restant
	RightVolumeRemaining float64 `json:"rightVolumeRemaining"`   // Volume droite restant
	Reason               string  `json:"reason"`                 // Raison si gain = 0
	CommissionID         *string `json:"commissionId,omitempty"` // ID de la commission créée
}

// BinaryNode représente un nœud dans l'arbre binaire (pour calculs récursifs)
type BinaryNode struct {
	ClientID     primitive.ObjectID  `bson:"clientId" json:"clientId"`
	LeftChildID  *primitive.ObjectID `bson:"leftChildId,omitempty" json:"leftChildId,omitempty"`
	RightChildID *primitive.ObjectID `bson:"rightChildId,omitempty" json:"rightChildId,omitempty"`
	IsActive     bool                `bson:"isActive" json:"isActive"` // Un membre est actif s'il a fait au moins 1 vente
}












