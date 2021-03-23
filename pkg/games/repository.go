package games

import "context"

// Repository handle the CRUD operations with Games.
type Repository interface {
	GetAllPublicStarted(ctx context.Context) ([]Game, error)
	GetOnePrivate(ctx context.Context, gamename string) (Game, error)
	GetEndedByUser(ctx context.Context, userID uint) ([]Game, error)
	Create(ctx context.Context, game *Game) error
	Update(ctx context.Context, id uint, game Game) error
	Delete(ctx context.Context, id uint) error
}