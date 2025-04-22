package db

import (
	"context"
	"sync"
	"time"

	"github.com/1ef7yy/medods_test_task/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	log logger.Logger
	DB  *pgxpool.Pool
}

func Config(dsn string, log logger.Logger) *pgxpool.Config {
	const defaultMaxConns = 10
	const defaultMinConns = 0
	const defaultMaxConnLifetime = time.Hour * 1
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		log.Fatalf("Failed to create a config: %s", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	log.Infof("postgres: defaultMaxConns: %d", dbConfig.MaxConns)
	log.Infof("postgres: defaultMinConns: %d", dbConfig.MinConns)
	log.Infof("postgres: defaultMaxConnLifetime: %s", dbConfig.MaxConnLifetime)
	log.Infof("postgres: defaultMaxConnIdleTime: %s", dbConfig.MaxConnIdleTime)
	log.Infof("postgres: defaultHealthCheckPeriod: %s", dbConfig.HealthCheckPeriod)
	log.Infof("postgres: defaultConnectTimeout: %s", dbConfig.ConnConfig.ConnectTimeout)

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Info("Closed the connection pool.")
	}

	return dbConfig
}

func NewPostgres(ctx context.Context, dsn string, log logger.Logger) (*Postgres, error) {
	var (
		pgInstance *Postgres
		pgOnce     sync.Once
		pgErr      error
	)

	pgOnce.Do(func() {
		db, err := pgxpool.NewWithConfig(ctx, Config(dsn, log))
		if err != nil {
			log.Fatalf("Unable to connect to database: %s", err.Error())
			pgErr = err
		}

		pgInstance = &Postgres{
			log: log,
			DB:  db,
		}
	})

	if pgErr != nil {
		return nil, pgErr
	}
	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.DB.Close()
}
