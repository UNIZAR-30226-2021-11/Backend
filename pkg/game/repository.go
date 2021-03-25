package game

import "context"

// Repository handle the CRUD operations with Games.
type Repository interface {
	GetAll(ctx context.Context) ([]Game, error)
	GetOne(ctx context.Context, gamename string) (Game, error)
	GetByUser(ctx context.Context, userID uint) ([]Game, error)
	Create(ctx context.Context, g *Game, userID uint) error
	Update(ctx context.Context, id uint, game Game) error
	Delete(ctx context.Context, id uint) error
}