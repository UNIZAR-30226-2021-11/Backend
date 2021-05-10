package data

import (
	"Backend/pkg/pair"
	"context"
)

type PairRepository struct {
	Data *Data
}

func (pr *PairRepository) Update(ctx context.Context, p *pair.Pair) error {
	//TODO: update when pair wins/loses game
	return nil
}