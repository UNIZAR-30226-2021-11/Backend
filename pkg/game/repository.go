package game

import "context"

// Repository handle the CRUD operations with Games.
type Repository interface {
	GetOne(ctx context.Context, gameID uint) (Game, error)
	GetAll(ctx context.Context) ([]Game, error)
	GetTournament(ctx context.Context) ([]Game, error)
	GetByName(ctx context.Context, gamename string) (Game, error)
	GetByUser(ctx context.Context, userID uint) ([]Game, error)
	Create(ctx context.Context, g *Game, userID uint) error
	CreateTournament(ctx context.Context, g *Game) error
	Update(ctx context.Context, id uint, game Game) error
	Delete(ctx context.Context, id uint) error
	End(ctx context.Context, game Game) error
}
