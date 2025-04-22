package db

import (
	"context"
	"errors"

	customErrors "github.com/1ef7yy/medods_test_task/internal/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) StoreRefresh(ctx context.Context, refresh_hash, guid string) error {
	query := `
	INSERT INTO tokens(refresh_hash, guid, generation)
	VALUES ($1, $2, 1)
	RETURNING guid
	`

	var returnedGuid string
	err := p.DB.QueryRow(ctx, query, refresh_hash, guid).Scan(&returnedGuid)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				p.log.Errorf("duplicate guid: %s", guid)
				return customErrors.UserAlreadyLoggedIn
			case "23503":
				p.log.Errorf("invalid guid (user doesn't exist): %s", guid)
				return customErrors.CouldNotFindGuid
			}
		}
		p.log.Errorf("error storing refresh token: %s", err.Error())
		return err
	}

	if returnedGuid == "" {
		return customErrors.CouldNotFindGuid
	}

	return err
}

func (p *Postgres) GetHash(ctx context.Context, guid string) (string, error) {
	query := `
	SELECT refresh_hash FROM tokens
	WHERE guid=$1
	`
	var refreshHash string
	err := p.DB.QueryRow(ctx, query, guid).Scan(&refreshHash)

	if err != nil {
		p.log.Errorf("error getting hash: %s", err.Error())
		return "", err
	}

	if refreshHash == "" {
		return "", customErrors.CouldNotFindRefreshHash
	}

	return refreshHash, nil
}

func (p *Postgres) NewGeneration(ctx context.Context, refresh_hash string) (int, error) {
	query := `
	UPDATE tokens
	SET generation = generation + 1
	WHERE refresh_hash = $1
	RETURNING generation
	`
	var gen int
	err := p.DB.QueryRow(ctx, query, refresh_hash).Scan(&gen)
	if err != nil {
		p.log.Errorf("error incrementing generation: %s", err.Error())
		return 0, err
	}

	return gen, nil
}

func (p *Postgres) GetUserEmail(ctx context.Context, guid string) (string, error) {
	query := `
	SELECT email
	FROM users
	WHERE guid=$1
	`

	var addr string
	err := p.DB.QueryRow(ctx, query, guid).Scan(&addr)
	if err != nil {
		p.log.Errorf("error getting user's mail: %s", err.Error())
		return "", err
	}

	return addr, nil
}
