package state


// Card represents a guiñote card
// value allowed [1-7] [10-12]
type card struct {
	suit string
	value string
	canBePlayed bool
}


// Checks if the value of c2 is less than that of c
func (c *card) wins(c2 *card) bool {
	return c.value > c2.value
}


// Checks whether the cards have the same suit
func (c *card) sameSuit(c2 *card) bool {
	return c.suit == c2.suit
}
