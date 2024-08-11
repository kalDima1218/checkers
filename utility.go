package main

import (
	"math"
	"net/http"
	"strconv"
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
	return _abs(b[0] - a[0])
}

func _isBetween(l [2]int, m [2]int, r [2]int) bool {
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

func hash(s string) int {
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
	game  Board
}

func newMove(score float64, game Board) Move {
	var tmp Move
	tmp.score = score
	tmp.game = game
	return tmp
}

type ItemWaitingGame struct {
	val    int64
	player string
}

func (i ItemWaitingGame) getSize() int {
	return len(strconv.FormatInt(i.val, 10))
}

func (i ItemWaitingGame) getVal(n int) string {
	val := strconv.FormatInt(i.val, 10)
	for len(val) < n {
		val = "0" + val
	}
	return val
}

func (i ItemWaitingGame) getFieldString(field string) string {
	if field == "player" {
		return i.player
	}
	return ""
}

func newItemWaitingGame(player string, _time int64) Item {
	var tmp ItemWaitingGame
	tmp.val = _time
	tmp.player = player
	return tmp
}

func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://"+URL+":"+PORT+"/", http.StatusSeeOther)
}

func redirectTo(w http.ResponseWriter, r *http.Request, page string) {
	http.Redirect(w, r, "http://"+URL+":"+PORT+"/"+page, http.StatusSeeOther)
}

func resetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "", MaxAge: -1})
}

func getCookie(r *http.Request, dataKey string) string {
	data, _ := r.Cookie(dataKey)
	return data.Value
}
