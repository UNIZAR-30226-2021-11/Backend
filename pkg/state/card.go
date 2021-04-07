package state

// Card represents a guiñote Card
// val allowed [1-7] [10-12]
type Card struct {
	Suit        string
	val         int
	Points      int
	canBePlayed bool
}

// Creates a new card with the correct Points and value
func CreateCard(suit string, val int) *Card {

	return &Card{
		Suit:   suit,
		val:    val,
		Points: getPoints(val),
	}
}

// Changes the playability of the card
func (c *Card) AllowPlay(canBePlayed bool) {
	c.canBePlayed = canBePlayed
}

// Checks wether this card can be played or not
func (c *Card) CanBePlayed() bool {
	return c.canBePlayed
}

// Checks if the value of c2 is less than that of c
// They must be of the same Suit
func (c *Card) Wins(c2 *Card) bool {
	// If the have no value, check with order
	if c.Points == 0 && c2.Points == 0 {
		return c.val > c2.val
	}
	return c.Points > c2.Points
}

// Checks if the card is triumph
func (c *Card) IsTriumph(triumph string) bool {
	return c.Suit == triumph
}

// Checks whether the cards have the same Suit
func (c *Card) SameSuit(c2 *Card) bool {
	return c.Suit == c2.Suit
}

// Returns the Points that this Card gives
func (c *Card) value() int {
	return c.Points
}

// Checks if 2 cards are equal
func (c *Card) equals(c2 *Card) bool {
	return c.Suit == c2.Suit && c.val == c2.val
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
