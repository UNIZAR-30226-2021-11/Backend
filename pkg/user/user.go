package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User of the system.
type User struct {
	ID           uint      `json:"id,omitempty"`
	Username     string    `json:"username,omitempty"`
	PlayerId     uint      `json:"player_id,omitempty"`
	Email        string    `json:"email,omitempty"`
	Location     string    `json:"location,omitempty"`
	GamesWon     int       `json:"games_won"`
	GamesLost    int       `json:"games_lost"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// HashPassword generates a hash of the password and places the result in PasswordHash.
func (u *User) HashPassword() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)

	return nil
}

// PasswordMatch compares HashPassword with the password and returns true if they match.
func (u User) PasswordMatch(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}
