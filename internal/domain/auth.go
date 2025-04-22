package domain

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/1ef7yy/medods_test_task/internal/errors"
	"github.com/1ef7yy/medods_test_task/models"
	"github.com/1ef7yy/medods_test_task/pkg/jwt"
	"github.com/1ef7yy/medods_test_task/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func (d domain) Login(ctx context.Context, req models.GenerateTokenRequest) (models.Token, error) {
	tokens, err := jwt.GenerateTokenPair(req)

	d.log.Debugf("refresh token: %s\naccess token: %s", tokens.RefreshToken, tokens.AccessToken)

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

	err = d.db.StoreRefresh(ctx, bcryptRefresh, req.Guid)

	if err != nil {
		d.log.Errorf("error storing refresh: %s", err.Error())
		return models.Token{}, err
	}

	return tokens, nil
}

func (d domain) Refresh(ctx context.Context, req models.RefreshTokenRequest) (models.Token, error) {
	accessToken, err := jwt.DecodeAccess(req.Tokens.AccessToken)
	if err != nil {
		d.log.Errorf("error decoding access token: %s", err.Error())
		return models.Token{}, err
	}

	refreshToken, err := jwt.DecodeRefresh(req.Tokens.RefreshToken)
	if err != nil {
		d.log.Errorf("error decoding refresh token: %s", err.Error())
		return models.Token{}, err
	}

	if accessToken.Guid != refreshToken.Guid {
		return models.Token{}, errors.GuidIsDifferentErr
	}

	if refreshToken.IP != req.IP {
		userAddr, err := d.db.GetUserEmail(ctx, refreshToken.Guid)
		if err != nil {
			d.log.Errorf("error getting user's email: %s", err.Error())
		}
		err = d.smtp.SendMail(userAddr, fmt.Sprintf("warning, we noticed that your ip has changed, new address: %s", req.IP))
		if err != nil {
			d.log.Warnf("error sending mail: %s", err.Error())
		}
	}

	h := sha256.New()
	_, err = h.Write([]byte(req.Tokens.RefreshToken))
	if err != nil {
		d.log.Errorf("error writing to sha hash")
		return models.Token{}, err
	}
	shaRefresh := h.Sum(nil)

	refreshHash, err := d.db.GetHash(ctx, refreshToken.Guid)
	if err != nil {
		d.log.Errorf("error getting refresh_hash: %s", err.Error())
		return models.Token{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(refreshHash), shaRefresh)
	if err != nil {
		d.log.Errorf("sha256 (%s) and db bcrypt (%s) did not match", shaRefresh, refreshHash)
		return models.Token{}, errors.HashedRefreshDiffErr
	}

	generation, err := d.db.NewGeneration(ctx, refreshHash)

	if err != nil {
		d.log.Errorf("error incrementing generation: %s", err.Error())
		return models.Token{}, err
	}

	generateReq := models.GenerateTokenRequest{
		Guid:       refreshToken.Guid,
		IP:         refreshToken.IP,
		Generation: generation,
	}

	tokens, err := jwt.GenerateTokenPair(generateReq)
	if err != nil {
		d.log.Errorf("error generating jwt token pair: %s", err.Error())
		return models.Token{}, err
	}

	return tokens, nil
}
