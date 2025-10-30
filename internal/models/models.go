package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product in the MLM system
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	ImageURL    string             `bson:"imageUrl" json:"imageUrl"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Client represents a client in the MLM system
type Client struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ClientID           string              `bson:"clientId" json:"clientId"` // 8-digit numeric ID
	Name               string              `bson:"name" json:"name"`
	PasswordHash       string              `bson:"passwordHash" json:"-"`
	SponsorID          *primitive.ObjectID `bson:"sponsorId,omitempty" json:"sponsorId"`
	Position           *string             `bson:"position,omitempty" json:"position"` // "left" or "right"
	LeftChildID        *primitive.ObjectID `bson:"leftChildId,omitempty" json:"leftChildId"`
	RightChildID       *primitive.ObjectID `bson:"rightChildId,omitempty" json:"rightChildId"`
	JoinDate           time.Time           `bson:"joinDate" json:"joinDate"`
	TotalEarnings      float64             `bson:"totalEarnings" json:"totalEarnings"`
	WalletBalance      float64             `bson:"walletBalance" json:"walletBalance"`
	NetworkVolumeLeft  float64             `bson:"networkVolumeLeft" json:"networkVolumeLeft"`
	NetworkVolumeRight float64             `bson:"networkVolumeRight" json:"networkVolumeRight"`
	BinaryPairs        int                 `bson:"binaryPairs" json:"binaryPairs"`
}

// Sale represents a sale in the MLM system
type Sale struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ClientID  primitive.ObjectID  `bson:"clientId" json:"clientId"`
	SponsorID primitive.ObjectID  `bson:"sponsorId" json:"sponsorId"`
	ProductID *primitive.ObjectID `bson:"productId,omitempty" json:"productId"`
	Amount    float64             `bson:"amount" json:"amount"`
	Quantity  int                 `bson:"quantity" json:"quantity"`
	Side      *string             `bson:"side,omitempty" json:"side"` // "left" or "right"
	Date      time.Time           `bson:"date" json:"date"`
	Status    string              `bson:"status" json:"status"` // "paid", "pending", "partial", "cancelled"
	Note      *string             `bson:"note,omitempty" json:"note"`
}

// Payment represents a payment in the MLM system
type Payment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID    primitive.ObjectID `bson:"clientId" json:"clientId"`
	Amount      float64            `bson:"amount" json:"amount"`
	Date        time.Time          `bson:"date" json:"date"`
	Method      string             `bson:"method" json:"method"` // 'mobile-money', 'cash', 'bank', etc.
	Status      string             `bson:"status" json:"status"` // "completed", "pending", "failed"
	Description *string            `bson:"description,omitempty" json:"description"`
}

// Commission represents a commission in the MLM system
type Commission struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID       primitive.ObjectID `bson:"clientId" json:"clientId"`
	SourceClientID primitive.ObjectID `bson:"sourceClientId" json:"sourceClientId"`
	Amount         float64            `bson:"amount" json:"amount"`
	Level          int                `bson:"level" json:"level"`
	Type           string             `bson:"type" json:"type"` // "binary-match", "override", etc.
	Date           time.Time          `bson:"date" json:"date"`
}

// Admin represents an admin user
type Admin struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	Role         string             `bson:"role" json:"role"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalClients     int     `json:"totalClients"`
	TotalSales       float64 `json:"totalSales"`
	TotalCommissions float64 `json:"totalCommissions"`
	TotalProducts    int     `json:"totalProducts"`
	ActiveClients    int     `json:"activeClients"`
}

// AuthPayload represents authentication response
type AuthPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Admin        *Admin `json:"admin"`
}

// CommissionResult represents the result of commission calculation
type CommissionResult struct {
	CommissionsCreated int     `json:"commissionsCreated"`
	TotalAmount        float64 `json:"totalAmount"`
	Message            string  `json:"message"`
}

// FilterInput represents filtering options for queries
type FilterInput struct {
	Search   *string    `json:"search,omitempty"`
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
	Status   *string    `json:"status,omitempty"`
}

// PagingInput represents pagination options for queries
type PagingInput struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// ProductInput represents input for creating/updating products
type ProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"imageUrl"`
}

// ClientInput represents input for creating/updating clients
type ClientInput struct {
	Name      string  `json:"name"`
	Password  string  `json:"password"`
	SponsorID *string `json:"sponsorId,omitempty"`
}

// SaleInput represents input for creating sales
type SaleInput struct {
	ClientID  string  `json:"clientId"`
	ProductID *string `json:"productId,omitempty"`
	Amount    float64 `json:"amount"`
	Note      *string `json:"note,omitempty"`
}

// PaymentInput represents input for creating payments
type PaymentInput struct {
	ClientID    string  `json:"clientId"`
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
	Description *string `json:"description,omitempty"`
}

// CommissionInput represents input for creating commissions
type CommissionInput struct {
	ClientID       string  `json:"clientId"`
	SourceClientID string  `json:"sourceClientId"`
	Amount         float64 `json:"amount"`
	Level          int     `json:"level"`
	Type           string  `json:"type"`
}

// LoginInput represents input for admin login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ClientLoginInput represents input for client login
type ClientLoginInput struct {
	ClientID string `json:"clientId"`
	Password string `json:"password"`
}

// RefreshTokenInput represents input for token refresh
type RefreshTokenInput struct {
	Token string `json:"token"`
}
