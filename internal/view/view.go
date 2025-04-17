package view

import (
	"net/http"

	"github.com/1ef7yy/medods_test_task/pkg/logger"
)

type View interface {
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

type view struct {
	log logger.Logger
}

func NewView(log logger.Logger) (View, error) {
	return &view{
		log: log,
	}, nil
}
