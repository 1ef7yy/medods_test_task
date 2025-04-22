package view

import (
	"net/http"

	"github.com/1ef7yy/medods_test_task/internal/domain"
	"github.com/1ef7yy/medods_test_task/pkg/logger"
)

type View interface {
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

type view struct {
	log    logger.Logger
	domain domain.Domain
}

func NewView(log logger.Logger) (View, error) {
	domain, err := domain.NewDomain(log)

	if err != nil {
		return nil, err
	}

	return &view{
		log:    log,
		domain: domain,
	}, nil
}
