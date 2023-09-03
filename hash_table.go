package main

import "time"

type HashTable struct {
	load, size int
	val        []Item
	t          []float64
	isLocked   bool
	accessing  int
}

func newHashTable() HashTable {
	return _newHashTableSized(3)
}

func _newHashTableSized(size int) HashTable {
	var tmp HashTable
	tmp.size, tmp.load = size, 0
	tmp.val = make([]Item, size)
	tmp.t = make([]float64, size)
	tmp.isLocked = false
	tmp.accessing = 0
	go tmp.Daemon()
	return tmp
}

func _h1(x, size int) int {
	return x % size
}

func _h2(x, size int) int {
	return 1 + x%(size-2)
}

func (mp *HashTable) Daemon() {
	for true {
		if float64(mp.load)/float64(mp.size) >= 0.5 {
			mp.isLocked = true
			for mp.accessing != 0 {
				<-time.After(time.Millisecond * 10)
			}
			mp._rehash()
			mp.isLocked = false
		}
		<-time.After(time.Millisecond * 100)
	}
}

func (mp *HashTable) insert(key Item, x float64) {
	mp.accessing++
	var i = _h1(key.hash(), mp.size)
	var d = _h2(key.hash(), mp.size)
	for mp.t[i] > 0 {
		if mp.val[i] == key {
			return
		}
		i = (i + d) % mp.size
	}
	mp.t[i] = x
	mp.val[i] = key
	mp.load++
}

func (mp *HashTable) get(x Item) (float64, bool) {
	var i = _h1(x.hash(), mp.size)
	var d = _h2(x.hash(), mp.size)
	for mp.t[i] > 0 {
		if mp.val[i] == x {
			return mp.t[i], true
		}
		i = (i + d) % mp.size
	}
	return 0, false
}

func (mp *HashTable) clear() {
	*mp = newHashTable()
}

func (mp *HashTable) _rehash() {
	var _mp HashTable
	_mp.size, _mp.load = mp.size*2, 0
	_mp.val = make([]Item, mp.size*2)
	_mp.t = make([]float64, _mp.size*2)
	for i := 0; i < mp.size; i++ {
		_mp.insert(mp.val[i], mp.t[i])
	}
	*mp = _mp
}
