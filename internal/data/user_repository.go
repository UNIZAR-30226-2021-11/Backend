package data

import (
	"context"
	"time"

	"Backend/pkg/user"
)

// UserRepository manages the operations with the database that
// correspond to the user model.
type UserRepository struct {
	Data *Data
}

// GetByUsername returns one user by username.
func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (user.User, error) {
	q := `
	SELECT id, username, email, location, games_won, games_lost,
		password, created_at, updated_at
		FROM users WHERE username = $1;
	`

	row := ur.Data.DB.QueryRowContext(ctx, q, username)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Location, &u.GamesWon, &u.GamesLost,
		&u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

// Create adds a new user.
func (ur *UserRepository) Create(ctx context.Context, u *user.User) error {
	q := `
	INSERT INTO users (username, password, email, location, games_won, games_lost, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`

	if err := u.HashPassword(); err != nil {
		return err
	}

	stmt, err := ur.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, u.Username, u.PasswordHash, u.Email, u.Location, u.GamesWon,
		u.GamesLost, time.Now(), time.Now(),
	)

	err = row.Scan(&u.ID)
	if err != nil {
		return err
	}

	send(u.Email, u.Username)

	return nil
}

// Update updates a user by id.
func (ur *UserRepository) Update(ctx context.Context, id uint, u user.User) error {
	q := `
	UPDATE users set email=$1, location=$2, games_won=$3, games_lost=$4, updated_at=$5
		WHERE id=$6;
	`

	stmt, err := ur.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx, u.Email, u.Location,
		u.GamesWon, u.GamesLost, time.Now(), id,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a user by id.
func (ur *UserRepository) Delete(ctx context.Context, id uint) error {
	q := `
	DELETE FROM users WHERE id=$1;
	`

	stmt, err := ur.Data.DB.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
