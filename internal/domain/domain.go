package domain

import (
	"github.com/1ef7yy/medods_test_task/models"
	"github.com/1ef7yy/medods_test_task/pkg/logger"
)

type Domain interface {
	Login(guid string) (models.Token, error)
	Refresh(models.Token) (models.Token, error)
}

type domain struct {
	log logger.Logger
}

func NewDomain(log logger.Logger) (Domain, error) {
	return domain{
		log: log,
	}, nil
}
