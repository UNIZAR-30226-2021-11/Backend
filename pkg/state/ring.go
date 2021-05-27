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
			}
			r.All = append(r.All, player)
			continue
		}
		r.All = append(r.All, player)
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
	r.ringNode = r.next
	return &r
}

// SetFirstPlayer sets the initial Player of the round
// The ID must be a correct identifier
func (r *Ring) SetFirstPlayer(p *Player) {
	//
	for !r.ringNode.p.sameId(p.Id) {
		r.p.SetPlay(false)
		r.ringNode = r.next
	}
	r.p.SetPlay(true)
}

// SetRandomFirstPlayer Sets a random first Player
func (r *Ring) SetRandomFirstPlayer() {

	rand.Seed(time.Now().UnixNano())
	for i := rand.Intn(4); i != 0; i-- {
		r.ringNode = r.next
	}
	r.ringNode.p.SetPlay(true)

}

// Next Returns the current Player and advances the head to the next
func (r *Ring) Next() (p *Player) {
	p = r.p
	p.SetPlay(false)
	r.ringNode = r.next
	r.ringNode.p.SetPlay(true)
	return p
}

// GetN Gets the Player n positions ahead of the current Player
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
func (r *Ring) DealCards(cards [4]*Card) {

	for i := 0; i < 4; i++ {
		r.GetN(i).dealCard(cards[i])
	}
}

func (r *Ring) InitialCardDealing(cards [4][6]*Card) {
	for i := 0; i < 4; i++ {
		p := r.GetN(i)
		p.DealCards(cards[i])
	}
}

func (r *Ring) Find(id uint32) (p *Player) {
	current := r.ringNode

	for i := 0; i < 4; i++ {
		if current.p.Id == id {
			return current.p
		}
		current = current.next
	}
	return nil
}

func (r *Ring) GetPlayersIds() (ids []uint32) {
	for i := 0; i < 4; i++ {
		ids = append(ids, r.GetN(i).Id)
	}
	return ids
}

func (r *Ring) Current() *Player {
	return r.p
}

func (r *Ring) String() string {
	return r.p.String()
}
