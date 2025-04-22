package domain

import (
	"context"
	"fmt"
	"os"

	"github.com/1ef7yy/medods_test_task/internal/storage/db"
	"github.com/1ef7yy/medods_test_task/models"
	"github.com/1ef7yy/medods_test_task/pkg/logger"
	"github.com/1ef7yy/medods_test_task/pkg/mail"
)

type Domain interface {
	Login(ctx context.Context, req models.GenerateTokenRequest) (models.Token, error)
	Refresh(context.Context, models.RefreshTokenRequest) (models.Token, error)
}

type domain struct {
	log  logger.Logger
	db   db.Postgres
	smtp mail.SMTPService
}

func NewDomain(log logger.Logger) (Domain, error) {
	dsn, ok := os.LookupEnv("POSTGRES_CONN")
	if !ok {
		return nil, fmt.Errorf("error looking up postgres dsn in env")
	}
	pg, err := db.NewPostgres(context.Background(), dsn, log)
	if err != nil {
		log.Errorf("error creating postgres instance: %s", err.Error())
		return nil, err
	}
	addr, ok := os.LookupEnv("SMTP_ADDRESS")
	smtp := mail.NewSMTP(log, addr)
	if !ok {
		return nil, fmt.Errorf("error looking up mail address in env")
	}
	return domain{
		log:  log,
		db:   *pg,
		smtp: smtp,
	}, nil
}
