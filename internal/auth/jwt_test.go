package auth_test

import (
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/config"
)

// 6. Authentication / 6. JWTs
func TestValidate(t *testing.T) {
	config.MustLoad("../../.env")
	// userID := uuid.New()
	// validToken, _ := MakeJWT(userID, "secret", time.Hour)
	userID := int64(2)
	email := "some@example.com"
	jwtManager := auth.NewJWTManager([]byte(config.Cfg.App.JwtSecret))
	validToken, err := jwtManager.Generate(userID, email, auth.RolesUser)
	if err != nil {
		t.Errorf("jwtManager.Generate error = %v", err)
		return
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  int64
		wantErr     error
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: config.Cfg.App.JwtSecret,
			wantUserID:  userID,
			wantErr:     nil,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: config.Cfg.App.JwtSecret,
			wantUserID:  0,
			wantErr:     auth.ErrInvalidToken,
		},
		// {
		// 	name:        "Wrong secret",
		// 	tokenString: validToken,
		// 	tokenSecret: "wrong_secret",
		// 	wantUserID:  0,
		// 	wantErr:     nil,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id int64
			user, err := jwtManager.Validate(tt.tokenString)
			if err != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if user == nil {
				id = 0
			} else {
				id = user.UserID
			}
			if id != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", user.UserID, tt.wantUserID)
				return
			}
		})
	}
}
