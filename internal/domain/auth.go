package domain

import (
	"context"
	"crypto/sha256"

	"github.com/1ef7yy/medods_test_task/models"
	"github.com/1ef7yy/medods_test_task/pkg/jwt"
	"github.com/1ef7yy/medods_test_task/pkg/utils"
)

func (d domain) Login(ctx context.Context, req models.GenerateTokenRequest) (models.Token, error) {
	tokens, err := jwt.GenerateTokenPair(req)

	d.log.Infof("refresh token: %s\naccess token: %s", tokens.RefreshToken, tokens.AccessToken)

	if err != nil {
		d.log.Errorf("error generating token pair: %s", err.Error())
		return models.Token{}, err
	}

	// bcrypt has upper limit of 72 bytes, base64 producers more
	// so we put the base64 into sha256 and then into bcrypt
	h := sha256.New()
	_, err = h.Write([]byte(tokens.RefreshToken))
	if err != nil {
		d.log.Errorf("error writing to sha hash")
		return models.Token{}, err
	}
	shaRefresh := h.Sum(nil)

	bcryptRefresh, err := utils.StringToBcrypt(string(shaRefresh))

	if err != nil {
		d.log.Errorf("error generating bcrypt from refresh (%s): %s", err.Error())
		return models.Token{}, err
	}

	err = d.db.StoreRefresh(ctx, bcryptRefresh)

	if err != nil {
		d.log.Errorf("error storing refresh: %s", err.Error())
		return models.Token{}, err
	}

	return tokens, nil
}

func (d domain) Refresh(tokens models.Token) (models.Token, error) {
	return models.Token{}, nil
}
