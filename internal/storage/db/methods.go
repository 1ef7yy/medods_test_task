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
