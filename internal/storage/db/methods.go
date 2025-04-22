package db

import (
	"context"
)

func (p *Postgres) StoreRefresh(ctx context.Context, refresh_hash string) error {
	query := `
	INSERT INTO tokens(refresh_hash, generation)
	VALUES ($1, 1)
	`

	_, err := p.DB.Query(ctx, query, refresh_hash)

	if err != nil {
		p.log.Errorf("error storing refresh token: %s", err.Error())
	}

	return err
}

func (p *Postgres) GetUserEmail(ctx context.Context, guid string) (string, error) {
	query := `
	SELECT email
	FROM users
	WHERE guid=$1
	`

	val, err := p.DB.Query(ctx, query, guid)
	if err != nil {
		p.log.Errorf("error getting user's mail: %s", err.Error())
		return "", err
	}

	var addr string
	if val.Next() {
		err = val.Scan(&addr)
		if err != nil {
			p.log.Errorf("error scanning rows: %s", err.Error())
			return "", err
		}
	}

	return addr, nil
}
