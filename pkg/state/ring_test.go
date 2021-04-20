package state

import (
	"encoding/json"
	"testing"
)

func TestSerializeRing(t *testing.T) {
	r := CreateTestRing()

	b, err := json.Marshal(r)
	if err != nil {
		t.Errorf("error marshaling")
	}

	t.Log(string(b))
}


func TestCreateRing(t *testing.T) {

	var players []*Player
	for i := 0; i < 4; i++{
		players = append(players, CreateTestPlayer())
	}
	r := NewPlayerRing(players)
	for _,p := range players {
		currentID := r.p.Id
		if !p.sameId(currentID){
			t.Errorf("got %v, want %v",currentID, p.Id )
		}
		r.Next()
	}
}



func CreateTestRing() *Ring {

	var players []*Player
	for i := 0; i < 4; i++{
		players = append(players, CreateTestPlayer())
	}

	return NewPlayerRing(players)
}
