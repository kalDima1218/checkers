package main

import (
	"math"
	"strconv"
	"sync"
)

func _dist(a [2]int, b [2]int) [2]int {
	return [2]int{b[0] - a[0], b[1] - a[1]}
}

func _abs(a int) int {
	return int(math.Abs(float64(a)))
}

func _add(a [2]int, b [2]int) [2]int {
	return [2]int{a[0] + b[0], a[1] + b[1]}
}

func _sub(a [2]int, b [2]int) [2]int {
	return [2]int{a[0] - b[0], a[1] - b[1]}
}

func _div(a [2]int, b int) [2]int {
	return [2]int{a[0] / b, a[1] / b}
}

func _len(a [2]int, b [2]int) int {
	return int(math.Sqrt((math.Pow(float64(b[0]-a[0]), 2) + math.Pow(float64(b[1]-a[1]), 2)) / 2))
}

func _isBetwen(l [2]int, m [2]int, r [2]int) bool {
	return _len(l, m)+_len(m, r) == _len(l, r)
}

func _max(a int, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func _min(a int, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func _less(a Item, b Item) bool {
	return a.getVal(b.getSize()) < b.getVal(a.getSize())
}

func _less_equal(a Item, b Item) bool {
	return a.getVal(b.getSize()) <= b.getVal(a.getSize())
}

func _hash(s string) int {
	const b, m = 131, 100000000000031
	var h = 0
	for _, c := range s {
		h *= b
		h %= m
		h += int(c + 1)
		h %= m
	}
	return h
}

type Move struct {
	score float64
	game  Game
}

func newMove(score float64, game Game) Move {
	var tmp Move
	tmp.score = score
	tmp.game = game
	return tmp
}

type Item interface {
	hash() int
	getSize() int
	getVal(n int) string
	getFieldString(field string) string
}

type ItemWatingGame struct {
	val    int
	player string
}

func (i ItemWatingGame) hash() int {
	return _abs(i.val)
}

func (i ItemWatingGame) getSize() int {
	return len(strconv.Itoa(i.val))
}

func (i ItemWatingGame) getVal(n int) string {
	var val = strconv.Itoa(i.val)
	for len(val) < n {
		val = "0" + val
	}
	return val
}

func (i ItemWatingGame) getFieldString(field string) string {
	if field == "player" {
		return i.player
	}
	return ""
}

func newItemWatingGame(player string, _time int) Item {
	var tmp ItemWatingGame
	tmp.val = _time
	tmp.player = player
	return tmp
}

type ItemGame struct {
	val Game
}

func (it ItemGame) hash() int {
	const b, m = 5, 100000000000031
	var h = it.val.Whose_turn
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			h *= b
			h %= m
			h += it.val.Board[i][j]
			h %= m
		}
	}
	return h
}

func (it ItemGame) getSize() int {
	return 0
}

func (it ItemGame) getVal(n int) string {
	return ""
}

func (it ItemGame) getFieldString(field string) string {
	return ""
}

func newItemGame(x Game) ItemGame {
	var tmp ItemGame
	tmp.val = x
	return tmp
}

type GameKey struct {
	Board      [8][8]int
	Whose_turn int
	Last_piece [2]int
}

func newGameKey(game Game) GameKey {
	var tmp GameKey
	tmp.Board = game.Board
	tmp.Whose_turn = game.Whose_turn
	tmp.Last_piece = game.Last_piece
	return tmp
}

type MapGame struct {
	mx sync.Mutex
	mp map[GameKey]float64
}

func newMapGame() *MapGame {
	return &MapGame{
		mp: make(map[GameKey]float64),
	}
}

func (mp *MapGame) get(key GameKey) (float64, bool) {
	mp.mx.Lock()
	defer mp.mx.Unlock()
	val, ok := mp.mp[key]
	return val, ok
}

func (mp *MapGame) insert(key GameKey, value float64) {
	mp.mx.Lock()
	defer mp.mx.Unlock()
	mp.mp[key] = value
}

func (mp *MapGame) clear() {
	for k := range mp.mp {
		delete(mp.mp, k)
	}
}

func _init() {
	for i := 1; i <= 7; i++ {
		POSSIBLE_TURNS[i-1] = [2]int{i, i}
		POSSIBLE_TURNS[7+i-1] = [2]int{i, -i}
		POSSIBLE_TURNS[14+i-1] = [2]int{-i, i}
		POSSIBLE_TURNS[21+i-1] = [2]int{-i, -i}
	}
	players[BOT_PLAYER.Login] = BOT_PLAYER
	logins[BOT_PLAYER.Login] = true
	names[BOT_PLAYER.Name] = true
}
