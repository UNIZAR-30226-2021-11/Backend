package player

import "context"

type Repository interface {
	Create(ctx context.Context, p *Player) error
}