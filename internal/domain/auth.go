package domain

import (
	"github.com/1ef7yy/medods_test_task/models"
	"github.com/1ef7yy/medods_test_task/pkg/jwt"
)

func (d domain) Login(guid string) (models.Token, error) {
	tokens, err := jwt.GenerateTokenPair(guid)

	if err != nil {
		d.log.Errorf("error generating token pair: %s", err.Error())
		return models.Token{}, err
	}

	// add db

	return tokens, nil
}

func (d domain) Refresh(tokens models.Token) (models.Token, error) {
	return models.Token{}, nil
}
