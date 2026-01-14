package validation

import (
	"errors"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidObjectID     = errors.New("ID invalide (format ObjectID requis)")
	ErrNegativeAmount      = errors.New("le montant doit être positif")
	ErrZeroAmount          = errors.New("le montant doit être supérieur à zéro")
	ErrNegativePrice       = errors.New("le prix doit être positif")
	ErrNegativeStock        = errors.New("le stock doit être positif ou zéro")
	ErrNegativeQuantity     = errors.New("la quantité doit être positive")
	ErrZeroQuantity         = errors.New("la quantité doit être supérieure à zéro")
	ErrInvalidEmail         = errors.New("format d'email invalide")
	ErrInvalidStatus        = errors.New("statut invalide")
	ErrInvalidMethod        = errors.New("méthode de paiement invalide")
	ErrInvalidTransactionType = errors.New("type de transaction invalide (doit être 'entree' ou 'sortie')")
	ErrInvalidLevel          = errors.New("le niveau doit être positif")
	ErrInvalidCommissionType = errors.New("type de commission invalide")
	ErrInvalidPosition       = errors.New("position invalide (doit être 'left' ou 'right')")
	ErrEmptyName             = errors.New("le nom ne peut pas être vide")
	ErrEmptyString           = errors.New("ce champ ne peut pas être vide")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateObjectID validates that a string is a valid MongoDB ObjectID
func ValidateObjectID(id string) error {
	if id == "" {
		return ErrInvalidObjectID
	}
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidObjectID
	}
	return nil
}

// ValidateObjectIDPtr validates that a pointer to string is a valid MongoDB ObjectID (if not nil)
func ValidateObjectIDPtr(id *string) error {
	if id == nil {
		return nil // nil is valid for optional fields
	}
	return ValidateObjectID(*id)
}

// ValidateAmount validates that an amount is positive
func ValidateAmount(amount float64) error {
	if amount < 0 {
		return ErrNegativeAmount
	}
	return nil
}

// ValidateAmountPositive validates that an amount is strictly positive
func ValidateAmountPositive(amount float64) error {
	if amount <= 0 {
		return ErrZeroAmount
	}
	return nil
}

// ValidateAmountPtr validates that a pointer to amount is positive (if not nil)
func ValidateAmountPtr(amount *float64) error {
	if amount == nil {
		return nil // nil is valid for optional fields
	}
	return ValidateAmount(*amount)
}

// ValidatePrice validates that a price is positive
func ValidatePrice(price float64) error {
	if price < 0 {
		return ErrNegativePrice
	}
	return nil
}

// ValidateStock validates that stock is non-negative
func ValidateStock(stock int32) error {
	if stock < 0 {
		return ErrNegativeStock
	}
	return nil
}

// ValidateQuantity validates that a quantity is positive
func ValidateQuantity(quantity int32) error {
	if quantity <= 0 {
		return ErrZeroQuantity
	}
	return nil
}

// ValidateEmail validates that a string is a valid email format
func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmptyString
	}
	email = strings.TrimSpace(email)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateName validates that a name is not empty
func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}
	return nil
}

// ValidateSaleStatus validates that a sale status is valid
func ValidateSaleStatus(status string) error {
	validStatuses := map[string]bool{
		"pending": true,
		"paid":    true,
		"partial": true,
		"cancelled": true,
	}
	if !validStatuses[status] {
		return ErrInvalidStatus
	}
	return nil
}

// ValidatePaymentMethod validates that a payment method is valid
func ValidatePaymentMethod(method string) error {
	validMethods := map[string]bool{
		"cash":      true,
		"card":      true,
		"bank":      true,
		"mobile":    true,
		"transfer":  true,
		"other":     true,
	}
	if !validMethods[method] {
		return ErrInvalidMethod
	}
	return nil
}

// ValidateTransactionType validates that a transaction type is valid
func ValidateTransactionType(transactionType string) error {
	if transactionType != "entree" && transactionType != "sortie" {
		return ErrInvalidTransactionType
	}
	return nil
}

// ValidateCommissionType validates that a commission type is valid
func ValidateCommissionType(commissionType string) error {
	validTypes := map[string]bool{
		"binary":   true,
		"unilevel": true,
		"direct":   true,
		"bonus":    true,
	}
	if !validTypes[commissionType] {
		return ErrInvalidCommissionType
	}
	return nil
}

// ValidateLevel validates that a level is positive
func ValidateLevel(level int32) error {
	if level <= 0 {
		return ErrInvalidLevel
	}
	return nil
}

// ValidatePosition validates that a position is valid
func ValidatePosition(position string) error {
	if position != "left" && position != "right" {
		return ErrInvalidPosition
	}
	return nil
}

// ValidatePositionPtr validates that a pointer to position is valid (if not nil)
func ValidatePositionPtr(position *string) error {
	if position == nil {
		return nil // nil is valid for optional fields
	}
	return ValidatePosition(*position)
}



