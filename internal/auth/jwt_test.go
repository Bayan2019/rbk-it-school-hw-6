package auth_test

import (
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/config"
)

// 6. Authentication / 6. JWTs
func TestValidate(t *testing.T) {
	// userID := uuid.New()
	// validToken, _ := MakeJWT(userID, "secret", time.Hour)
	userID := int64(2)
	email := "some@example.com"
	jwtManager := auth.NewJWTManager([]byte(config.Cfg.App.JwtSecret))
	validToken, _ := jwtManager.Generate(userID, email, auth.RolesUser)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  int64
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: config.Cfg.App.JwtSecret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: config.Cfg.App.JwtSecret,
			wantUserID:  0,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := jwtManager.Validate(tt.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if user.UserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", user.UserID, tt.wantUserID)
			}
		})
	}
}
