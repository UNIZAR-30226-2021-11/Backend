package data

import "sync/atomic"

type IdManager struct {
	playerIds    map[uint32]bool

	lastId uint32
}

func NewIdManager() *IdManager {
	return &IdManager{
		playerIds:    make(map[uint32]bool),
		lastId:       0,
	}
}

func (manager *IdManager) NextPlayerId() uint32 {
	return manager.nextId(manager.playerIds)
}

func (manager *IdManager) IsPlayerId(id uint32) bool {
	_, ok := manager.playerIds[id]
	return ok
}

func (manager *IdManager) nextId(idsMap map[uint32]bool) uint32 {
	id := atomic.AddUint32(&manager.lastId, 1)
	idsMap[id] = true
	return id
}
