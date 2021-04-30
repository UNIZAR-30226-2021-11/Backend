package events

import (
	"Backend/internal/data"
)

type StateChanged struct {
	ClientsID []uint32
	GameData  *data.GameData
}
