package data

import (
	"Backend/pkg/player"
	"context"
)

type PlayerRepository struct {
	Data *Data
}

// Create adds a new player.
func (pr *PlayerRepository) Create(ctx context.Context, p *player.Player) error {
	q := `
	INSERT INTO players(user_id, pair_id)
		VALUES ($1, $2)
		RETURNING id;
	`
	stmt, err := pr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, p.UserID, p.PairID)

	err = row.Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

