package pair

import "context"

type Repository interface {
	UpdateWinned(ctx context.Context, id uint, p Pair) error
}
