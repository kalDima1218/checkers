package main

import "sync"

type MapGame struct {
	mx sync.Mutex
	mp map[Board]float64
}

func newMapGame() *MapGame {
	return &MapGame{
		mp: make(map[Board]float64),
	}
}

func (mp *MapGame) get(key Board) (float64, bool) {
	mp.mx.Lock()
	defer mp.mx.Unlock()
	val, ok := mp.mp[key]
	return val, ok
}

func (mp *MapGame) insert(key Board, value float64) {
	mp.mx.Lock()
	defer mp.mx.Unlock()
	mp.mp[key] = value
}

func (mp *MapGame) clear() {
	for k := range mp.mp {
		delete(mp.mp, k)
	}
}
