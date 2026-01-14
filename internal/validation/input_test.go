package validation

import "testing"

func TestValidateObjectID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		wantError bool
	}{
		{"Valid ObjectID", "507f1f77bcf86cd799439011", false},
		{"Empty string", "", true},
		{"Invalid format", "invalid-id", true},
		{"Too short", "507f1f77bcf86cd79943901", true},
		{"Too long", "507f1f77bcf86cd7994390111", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateObjectID(tt.id)
			if tt.wantError && err == nil {
				t.Errorf("ValidateObjectID() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateObjectID() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		wantError bool
	}{
		{"Positive amount", 100.0, false},
		{"Zero amount", 0.0, false},
		{"Negative amount", -10.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAmount(tt.amount)
			if tt.wantError && err == nil {
				t.Errorf("ValidateAmount() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateAmount() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAmountPositive(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		wantError bool
	}{
		{"Positive amount", 100.0, false},
		{"Zero amount", 0.0, true},
		{"Negative amount", -10.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAmountPositive(tt.amount)
			if tt.wantError && err == nil {
				t.Errorf("ValidateAmountPositive() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateAmountPositive() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantError bool
	}{
		{"Valid email", "test@example.com", false},
		{"Valid email with subdomain", "user@mail.example.com", false},
		{"Invalid email - no @", "testexample.com", true},
		{"Invalid email - no domain", "test@", true},
		{"Invalid email - no TLD", "test@example", true},
		{"Empty email", "", true},
		{"Email with spaces", "test @example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantError && err == nil {
				t.Errorf("ValidateEmail() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateEmail() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateSaleStatus(t *testing.T) {
	tests := []struct {
		name      string
		status    string
		wantError bool
	}{
		{"Valid status - pending", "pending", false},
		{"Valid status - paid", "paid", false},
		{"Valid status - partial", "partial", false},
		{"Valid status - cancelled", "cancelled", false},
		{"Invalid status", "invalid", true},
		{"Empty status", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSaleStatus(tt.status)
			if tt.wantError && err == nil {
				t.Errorf("ValidateSaleStatus() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateSaleStatus() unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePaymentMethod(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		wantError bool
	}{
		{"Valid method - cash", "cash", false},
		{"Valid method - card", "card", false},
		{"Valid method - bank", "bank", false},
		{"Valid method - mobile", "mobile", false},
		{"Valid method - transfer", "transfer", false},
		{"Valid method - other", "other", false},
		{"Invalid method", "invalid", true},
		{"Empty method", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePaymentMethod(tt.method)
			if tt.wantError && err == nil {
				t.Errorf("ValidatePaymentMethod() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidatePaymentMethod() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateTransactionType(t *testing.T) {
	tests := []struct {
		name      string
		transactionType string
		wantError bool
	}{
		{"Valid type - entree", "entree", false},
		{"Valid type - sortie", "sortie", false},
		{"Invalid type", "invalid", true},
		{"Empty type", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransactionType(tt.transactionType)
			if tt.wantError && err == nil {
				t.Errorf("ValidateTransactionType() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateTransactionType() unexpected error: %v", err)
			}
		})
	}
}

func TestValidatePosition(t *testing.T) {
	tests := []struct {
		name      string
		position  string
		wantError bool
	}{
		{"Valid position - left", "left", false},
		{"Valid position - right", "right", false},
		{"Invalid position", "invalid", true},
		{"Empty position", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePosition(tt.position)
			if tt.wantError && err == nil {
				t.Errorf("ValidatePosition() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidatePosition() unexpected error: %v", err)
			}
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		name      string
		quantity  int32
		wantError bool
	}{
		{"Valid quantity", 10, false},
		{"Valid quantity - 1", 1, false},
		{"Zero quantity", 0, true},
		{"Negative quantity", -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuantity(tt.quantity)
			if tt.wantError && err == nil {
				t.Errorf("ValidateQuantity() expected error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateQuantity() unexpected error: %v", err)
			}
		})
	}
}



