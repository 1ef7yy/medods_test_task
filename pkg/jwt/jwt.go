package jwt

import (
	"fmt"
	"os"

	"github.com/1ef7yy/medods_test_task/models"
	"github.com/golang-jwt/jwt/v4"
)

var (
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
)

func GenerateAccessToken() (string, error) {
	return "", nil
}

func GenerateRefreshToken() (string, error) {
	return "", nil
}

func GenerateTokenPair(req models.GenerateTokenRequest) (models.Token, error) {

	if JWTSecret == nil {
		return models.Token{}, fmt.Errorf("could not find JWT_SECRET IN environment")
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"sub": req.Guid,
			"ip":  req.IP,
		},
	).SignedString(JWTSecret)

	if err != nil {
		return models.Token{}, err
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"sub": req.Guid,
		}).SignedString(JWTSecret)

	if err != nil {
		return models.Token{}, err
	}

	return models.Token{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
