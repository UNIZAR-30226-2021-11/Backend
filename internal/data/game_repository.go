package data

import (
	"Backend/pkg/pair"
	"Backend/pkg/user"
	"context"
	"time"

	"Backend/pkg/game"
)

// GameRepository manages the operations with the database that
// correspond to the game model.
type GameRepository struct {
	Data *Data
}

// GetOne returns one started game.
func (gr *GameRepository) GetOne(ctx context.Context, gameID uint) (game.Game, error) {
	qPairs := `
	SELECT pa.id
	FROM
		games g
			INNER JOIN pairs pa
					   ON pa.game_id = g.id
	WHERE g.id = $1;
	`
	qUsers := `
	SELECT u.username
	FROM
		pairs pa
			INNER JOIN players pl
				ON pl.pair_id = pa.id
			INNER JOIN users u
				ON u.id = pl.user_id
	WHERE pa.id = $1;
	`
	var g game.Game
	//PAIRS
	rows,err := gr.Data.DB.QueryContext(ctx, qPairs, gameID)
	if err != nil {
		return game.Game{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var p pair.Pair
		rows.Scan(&p.ID)

		// USERS
		rows,err := gr.Data.DB.QueryContext(ctx, qUsers, p.ID)
		if err != nil {
			return game.Game{}, nil
		}

		defer rows.Close()

		for rows.Next() {
			var u user.User
			rows.Scan(&u.Username)
			p.Users = append(p.Users, u)
		}

		g.Pairs = append(g.Pairs, p)
	}

	return g, nil
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

// GetByName returns one started game by name.
func (gr *GameRepository) GetByName(ctx context.Context, gamename string) (game.Game, error) {
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
	SELECT g.id, g.name, g.end_date, pa.winned, pa.game_points
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
		rows.Scan(&g.ID, &g.Name, &g.EndDate, &g.Winned, &g.Points)
		games = append(games, g)
	}

	return games, nil

}

// Create adds and joins a new game.
func (gr *GameRepository) Create(ctx context.Context, g *game.Game, userID uint) error {
	// GAMES
	q := `
	INSERT INTO games(name, public, creation_date)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	stmt, err := gr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, g.Name, g.Public, time.Now())

	err = row.Scan(&g.ID)
	if err != nil {
		return err
	}
	g.PlayersCount = 1

	// PAIRS
	var pairID int
	q = `
	INSERT INTO pairs(game_id)
		VALUES ($1)
		RETURNING id;
	`
	stmt, err = gr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	//Create and join pair
	row = stmt.QueryRowContext(ctx, g.ID)
	//Create pair
	_ = stmt.QueryRowContext(ctx, g.ID)

	err = row.Scan(&pairID)
	if err != nil {
		return err
	}

	// PLAYERS
	var playerID uint
	q = `
	INSERT INTO players(user_id, pair_id)
		VALUES ($1, $2)
		RETURNING id;
	`
	stmt, err = gr.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row = stmt.QueryRowContext(ctx, userID, pairID)

	err = row.Scan(&playerID)
	if err != nil {
		return err
	}

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