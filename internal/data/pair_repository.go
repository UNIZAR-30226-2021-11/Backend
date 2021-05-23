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
	UPDATE pairs set winned=true
		WHERE id=$1;
	`

	stmt, err := pr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, id,
	)
	if err != nil {
		return err
	}

	return nil
}
