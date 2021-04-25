package state

import "fmt"

// Card represents a guiñote Card
// Val allowed [1-7] [10-12]
type Card struct {
	Suit     string `json:"suit"`
	Val      int    `json:"val"`
	Points   int    `json:"points"`
	Playable bool   `json:"playable"`
}

// CreateCard Creates a new card with the correct Points and value
func CreateCard(suit string, val int) *Card {

	return &Card{
		Suit:   suit,
		Val:    val,
		Points: getPoints(val),
	}
}

// AllowPlay Changes the playability of the card
func (c *Card) AllowPlay(canBePlayed bool) {
	c.Playable = canBePlayed
}

// CanBePlayed Checks wether this card can be played or not
func (c *Card) CanBePlayed() bool {
	return c.Playable
}

// Wins Checks if the value of c2 is less than that of c
// They must be of the same Suit
func (c *Card) Wins(c2 *Card) bool {
	// If the have no value, check with order
	if c.Points == 0 && c2.Points == 0 {
		return c.Val > c2.Val
	}
	return c.Points > c2.Points
}

// IsTriumph Checks if the card is triumph
func (c *Card) IsTriumph(triumph string) bool {
	return c.Suit == triumph
}

// SameSuit Checks whether the cards have the same Suit
func (c *Card) SameSuit(c2 *Card) bool {
	return c.Suit == c2.Suit
}

// IsSingingPair Checks if th
func (c *Card) IsSingingPair(c2 *Card) bool {
	if c.SameSuit(c2) {
		return (c.Val == 10 || c.Val == 12) && (c2.Val == 10 || c2.Val == 12)

	} else {
		return false
	}
}

// Returns the Points that this Card gives
func (c *Card) value() int {
	return c.Points
}

// Equals Checks if 2 cards are equal
func (c *Card) Equals(c2 *Card) bool {
	return c.Suit == c2.Suit && c.Val == c2.Val
}

// Returns the Points that this Card gives
func getPoints(cardName int) int {
	switch cardName {
	case 1:
		return 11
	case 3:
		return 10
	case 12:
		return 4
	case 11:
		return 2
	case 10:
		return 3
	default:
		return 0
	}
}

func (c *Card) String() string {
	return fmt.Sprintf("%d de %s", c.Val, c.Suit)
}
