package data

import (
	"context"

	"Backend/pkg/game"
)

// GameRepository manages the operations with the database that
// correspond to the game model.
type GameRepository struct {
	Data *Data
}

// GetAll returns all public started games.
func (gr *GameRepository) GetAll(ctx context.Context) ([]game.Game, error) {
	q := `
	SELECT g.id, g.name, COUNT(pl.id)
	FROM
		games g
			INNER JOIN pairs pa
					   ON pa.game_id = g.id
			INNER JOIN players pl
					   ON pl.pair_id = pa.id
	WHERE end_date is NULL AND 
		  g.public = TRUE
	GROUP BY g.id, g.name, g.creation_date
	ORDER BY g.creation_date;
	`
	rows, err := gr.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []game.Game
	for rows.Next() {
		var g game.Game
		rows.Scan(&g.ID, &g.Name, &g.PlayersCount)
		games = append(games, g)
	}

	return games, nil
}

// GetOne returns one started game by name.
func (gr *GameRepository) GetOne(ctx context.Context, gamename string) (game.Game, error) {
	q := `
	SELECT g.id, g.name, COUNT(pl.id)
	FROM
		games g
			INNER JOIN pairs pa
					   ON pa.game_id = g.id
			INNER JOIN players pl
					   ON pl.pair_id = pa.id
	WHERE end_date is NULL AND
		  g.public = FALSE AND
		  g.name LIKE $1
	GROUP BY g.id, g.name, g.creation_date
	ORDER BY g.creation_date;
	`
	row := gr.Data.DB.QueryRowContext(ctx, q, gamename)

	var g game.Game
	err := row.Scan(&g.ID, &g.Name, &g.PlayersCount)
	if err != nil {
		return game.Game{}, err
	}

	return g, nil
}

// GetByUser returns all user ended games.
func (gr *GameRepository) GetByUser(ctx context.Context, userID uint) ([]game.Game, error) {
	q := `
	SELECT g.id, g.name, g.end_date, pa.winned
	FROM
		games g
			INNER JOIN pairs pa
					   ON pa.game_id = g.id
			INNER JOIN players pl
					   ON pl.pair_id = pa.id
			INNER JOIN users u
					   ON u.id = pl.user_id
	WHERE end_date is NOT NULL AND
		  u.id = $1
	ORDER BY g.end_date;
	`
	rows, err := gr.Data.DB.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var games []game.Game
	for rows.Next() {
		var g game.Game
		rows.Scan(&g.ID, &g.Name, &g.EndDate, &g.Winned)
		games = append(games, g)
	}

	return games, nil

}

// Create adds a new game.
func (gr *GameRepository) Create(ctx context.Context, game *game.Game) error {
	//TODO
	return nil
}

// Update updates a game by id.
func (gr *GameRepository) Update(ctx context.Context, id uint, game game.Game) error {
	//TODO
	return nil
}

// Delete removes a game by id.
func (gr *GameRepository) Delete(ctx context.Context, id uint) error {
	//TODO
	return nil
}