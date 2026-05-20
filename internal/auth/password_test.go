package auth_test

import (
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
)

func TestCheckPassword(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := auth.HashPassword(password1)
	hash2, _ := auth.HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantRes  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantRes:  true,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantRes:  false,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantRes:  false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantRes:  false,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantRes:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := auth.CheckPassword(tt.password, tt.hash)
			if ok != tt.wantRes {
				t.Errorf("CheckPasswordHash() wantRes: %t get %t", ok, tt.wantRes)
			}
		})
	}
}
