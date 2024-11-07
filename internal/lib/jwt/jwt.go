package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"grpc-auth/internal/models"
	"time"
)

// NewToken creates new JWT token for given user.
func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(user.Salt))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
