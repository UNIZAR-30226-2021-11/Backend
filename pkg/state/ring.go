package state

import (
	"math/rand"
	"time"
)

type Ring struct {
	*ringNode
	All   []*Player `json:"players"`
	start *ringNode
	end   *ringNode
}

type ringNode struct {
	p    *Player
	next *ringNode
}

// NewPlayerRing creates a new 4 Player ring
func NewPlayerRing(players []*Player) *Ring {
	var r Ring
	var first ringNode
	for i, player := range players {
		if i == 0 {
			first = ringNode{
				p:    player,
				next: nil,
			}
			r = Ring{
				ringNode: &first,
				All:      players,
			}
			continue
		}
		rn := ringNode{
			p:    player,
			next: nil,
		}
		r.next = &rn
		r.ringNode = &rn
		if i == 3 {
			r.next = &first
		}
	}
	r.Next()
	return &r
}

// SetFirstPlayer sets the initial Player of the round
// The ID must be a correct identifier
func (r *Ring) SetFirstPlayer(p *Player) {
	//
	for !r.ringNode.p.sameId(p.Id) {
		r.ringNode = r.next
	}
}

// Sets a random first Player
func (r *Ring) SetRandomFirstPlayer() {

	rand.Seed(time.Now().UnixNano())
	for i := rand.Intn(4); i != 0; i-- {
		r.ringNode = r.next
	}
}

// Returns the current Player and advances the head to the next
func (r *Ring) Next() (p *Player) {
	p = r.p
	r.ringNode = r.next
	return p
}

// Gets the Player n positions ahead of the current Player
func (r *Ring) GetN(n int) (p *Player) {

	n = n % 4

	current := r.ringNode

	for i := 0; i < n; i++ {
		current = current.next
	}
	p = current.p
	return p
}

// DealCards deals a card to each Player
func (r *Ring) DealCards(cards []*Card) {

	for i := 0; i < 4; i++ {
		r.p.dealCard(cards[i])
	}
}
