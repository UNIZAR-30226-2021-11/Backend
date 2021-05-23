package data

import (
	"Backend/pkg/pair"
	"context"
)

type PairRepository struct {
	Data *Data
}

func (pr *PairRepository) UpdateWinned(ctx context.Context, id uint, p pair.Pair) error {
	q := `
	UPDATE pairs set winned=$1, game_points=$2
		WHERE id=$3;
	`

	stmt, err := pr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, p.Winned, p.GamePoints, id,
	)
	if err != nil {
		return err
	}

	return nil
}
