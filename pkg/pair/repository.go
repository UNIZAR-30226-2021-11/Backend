package pair

import "context"

type Repository interface {
	Update(ctx context.Context, id uint, pair Pair) error
}