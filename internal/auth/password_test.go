package auth

import "testing"

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
		errorType error
	}{
		{
			name:      "Valid password",
			password:  "Password123@",
			wantError: false,
		},
		{
			name:      "Password too short",
			password:  "Pass1@",
			wantError: true,
			errorType: ErrPasswordTooShort,
		},
		{
			name:      "Password without uppercase",
			password:  "password123@",
			wantError: true,
			errorType: ErrPasswordNoUpper,
		},
		{
			name:      "Password without lowercase",
			password:  "PASSWORD123@",
			wantError: true,
			errorType: ErrPasswordNoLower,
		},
		{
			name:      "Password without digit",
			password:  "Password@",
			wantError: true,
			errorType: ErrPasswordNoDigit,
		},
		{
			name:      "Password without special character",
			password:  "Password123",
			wantError: true,
			errorType: ErrPasswordNoSpecialChar,
		},
		{
			name:      "Password with all requirements",
			password:  "MyP@ssw0rd",
			wantError: false,
		},
		{
			name:      "Password with different special chars",
			password:  "Test123$",
			wantError: false,
		},
		{
			name:      "Password with multiple special chars",
			password:  "Test123@$!",
			wantError: false,
		},
		{
			name:      "Minimum length password",
			password:  "Pass1@",
			wantError: true,
			errorType: ErrPasswordTooShort,
		},
		{
			name:      "Exactly 8 characters valid",
			password:  "Pass1@ab",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidatePassword() expected error but got nil")
					return
				}
				if tt.errorType != nil && err != tt.errorType {
					t.Errorf("ValidatePassword() expected error %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePassword() unexpected error: %v", err)
				}
			}
		})
	}
}



