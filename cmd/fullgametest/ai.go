package main

import (
	"Backend/pkg/events"
	"Backend/pkg/simulation"
	"Backend/pkg/state"
	"fmt"
)

var (
	IdNewPlayer uint32
	PAIR        uint32
)

type AI struct {
	*state.Player
	out chan events.Event
	in  chan *simulation.GameState
}

func (ai *AI) controlAI() {
	// TODO react to player
	for {
		select {
		case newState := <-ai.in:
			if ai.CanPlay(newState) {
				ai.pickBestCard(newState)
			} else if ai.CanSing(newState) {
				ai.Sing()
			} else if ai.CanChange(newState) {
				ai.Change()
			}
		}
	}
}

func (ai *AI) CanPlay(newState *simulation.GameState) bool {
	for _, p := range newState.Players.All {
		if p.Id == ai.Id && p.CanPlay {
			return true
		}
	}
	return false
}

func (ai *AI) CanSing(newState *simulation.GameState) bool {
	for _, p := range newState.Players.All {
		if p.Id == ai.Id && p.CanSing {
			return true
		}
	}
	return false
}

func (ai *AI) CanChange(newState *simulation.GameState) bool {
	for _, p := range newState.Players.All {
		if p.Id == ai.Id && p.CanChange {
			return true
		}
	}
	return false
}

func (ai *AI) pickBestCard(newState *simulation.GameState) {

}

func (ai *AI) Sing() {

}

func (ai *AI) Change() {

}

func createAI(pair uint32) (ai *AI, out chan events.Event, in chan *simulation.GameState) {
	evtChan := make(chan events.Event)
	stateChan := make(chan *simulation.GameState)
	aiPlayer := state.CreatePlayer(
		IdNewPlayer, pair,
		fmt.Sprintf("IA_%d", IdNewPlayer))

	ai = &AI{
		Player: aiPlayer,
		out:    evtChan,
		in:     stateChan,
	}
	defer func() {
		IdNewPlayer++
		PAIR++
	}()
	return ai, evtChan, stateChan
}
