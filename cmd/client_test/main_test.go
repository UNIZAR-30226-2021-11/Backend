package main

import (
	"testing"
	"time"
)

func TestClient_CreateSoloGame(t *testing.T) {

	c := Client{
		Id:     5,
		PairId: 2,
	}
	c.Start()
	c.CreateSoloGame(6)
	time.Sleep(time.Second * 1000)
}
