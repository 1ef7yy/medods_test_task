package domain

import (
	"context"
	"fmt"
	"os"

	"github.com/1ef7yy/medods_test_task/internal/storage/db"
	"github.com/1ef7yy/medods_test_task/models"
	"github.com/1ef7yy/medods_test_task/pkg/logger"
)

type Domain interface {
	Login(ctx context.Context, req models.GenerateTokenRequest) (models.Token, error)
	Refresh(models.Token) (models.Token, error)
}

type domain struct {
	log logger.Logger
	db  db.Postgres
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
	return domain{
		log: log,
		db:  *pg,
	}, nil
}
