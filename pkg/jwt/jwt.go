package jwt

import (
	"fmt"
	"os"

	"github.com/1ef7yy/medods_test_task/internal/errors"
	"github.com/1ef7yy/medods_test_task/models"
	"github.com/golang-jwt/jwt/v4"
)

var (
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
)

func GenerateTokenPair(req models.GenerateTokenRequest) (models.Token, error) {

	if JWTSecret == nil {
		return models.Token{}, errors.CouldNotFindSecretErr
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
			"gen": req.Generation,
		}).SignedString(JWTSecret)

	if err != nil {
		return models.Token{}, err
	}

	return models.Token{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func DecodeRefresh(token string) (models.RefreshToken, error) {
	if JWTSecret == nil {
		return models.RefreshToken{}, errors.CouldNotFindSecretErr
	}
	refreshToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.TokenInvalidErr
		}
		return JWTSecret, nil
	})

	if err != nil {
		return models.RefreshToken{}, errors.TokenInvalidErr
	}

	if !refreshToken.Valid {
		return models.RefreshToken{}, errors.TokenInvalidErr
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)

	if !ok {
		return models.RefreshToken{}, fmt.Errorf("token %s could not be handled", token)
	}

	var ip, guid string
	if claims["sub"] != nil && claims["ip"] != nil {
		guid = claims["sub"].(string)
		ip = claims["ip"].(string)
	} else {
		return models.RefreshToken{}, errors.TokenInvalidErr
	}

	return models.RefreshToken{
		Guid: guid,
		IP:   ip,
	}, nil
}

func DecodeAccess(token string) (models.AccessToken, error) {
	if JWTSecret == nil {
		return models.AccessToken{}, errors.CouldNotFindSecretErr
	}
	accessToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.TokenInvalidErr
		}
		return JWTSecret, nil
	})

	if !accessToken.Valid {
		return models.AccessToken{}, errors.TokenInvalidErr
	}
	if err != nil {
		return models.AccessToken{}, errors.TokenInvalidErr
	}

	claims, ok := accessToken.Claims.(jwt.MapClaims)

	if !ok {
		return models.AccessToken{}, fmt.Errorf("token %s could not be handled", token)
	}
	var (
		guid string
		gen  int
	)

	if claims["sub"] != nil && claims["gen"] != nil {
		guid = claims["sub"].(string)
		gen = int(claims["gen"].(float64))
	} else {
		return models.AccessToken{}, errors.TokenInvalidErr
	}

	return models.AccessToken{
		Guid:       guid,
		Generation: gen,
	}, nil
}
