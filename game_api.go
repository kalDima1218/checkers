package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func getBoardHist(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	turn, _ := strconv.Atoi(r.URL.Query().Get("turn"))
	game, ok := games[id]
	if !ok || turn < 0 || turn >= len(game.Turns) {
		return
	}
	board_json, _ := json.Marshal(game.Turns[turn])
	fmt.Fprintf(w, string(board_json))
}

func getLastMoveNumber(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	_, ok := games[id]
	if !ok {
		return
	}
	fmt.Fprintf(w, strconv.Itoa(len(games[id].Turns)-1))
}

func whoseMove(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	_, ok := games[id]
	if !ok {
		return
	}
	fmt.Fprintf(w, strconv.Itoa(games[id].Board.Whose_turn))
}

func getSide(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	var login = getCookie(r, "login")
	game, ok := games[id]
	if !ok {
		return
	}
	if game.Players[0] == login {
		fmt.Fprintf(w, "0")
		return
	} else {
		fmt.Fprintf(w, "1")
		return
	}
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := games[id]
	if !ok {
		return
	}
	players_json, _ := json.Marshal([2]string{players[game.Players[0]].Username, players[game.Players[1]].Username})
	fmt.Fprintf(w, string(players_json))
}

func whoWin(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := games[id]
	if !ok {
		return
	}
	fmt.Fprintf(w, strconv.Itoa(game.whoWin()))
}

func makeMove(w http.ResponseWriter, r *http.Request) {
	if !checkUser(r) {
		redirectToIndex(w, r)
		return
	}
	var id = r.URL.Query().Get("id")
	var login = getCookie(r, "login")
	from_x, _ := strconv.Atoi(r.URL.Query().Get("from_x"))
	from_y, _ := strconv.Atoi(r.URL.Query().Get("from_y"))
	to_x, _ := strconv.Atoi(r.URL.Query().Get("to_x"))
	to_y, _ := strconv.Atoi(r.URL.Query().Get("to_y"))
	game, ok := games[id]
	if !ok || game.Players[game.Board.Whose_turn] != login {
		return
	}
	if game.makeMove([2]int{from_x, from_y}, [2]int{to_x, to_y}) {
		games[id] = game
		fmt.Fprintf(w, "1")
		fmt.Println(from_x, from_y, to_x, to_y)
	} else {
		fmt.Fprintf(w, "0")
	}
}

func endMove(w http.ResponseWriter, r *http.Request) {
	if !checkUser(r) {
		redirectToIndex(w, r)
		return
	}
	var id = r.URL.Query().Get("id")
	var login = getCookie(r, "login")
	game, ok := games[id]
	if !ok || game.Players[game.Board.Whose_turn] != login || game.Board.Last_piece == [2]int{-1, -1} {
		fmt.Fprintf(w, "0")
		return
	}
	game.endMove()
	if game.Players[game.Board.Whose_turn] == "BOT" {
		BOT.makeMove(&game)
		//game = BOT.findBestMove(game, game.Board.Whose_turn, (game.Board.Whose_turn+1)%2)
	}
	games[id] = game
	fmt.Fprintf(w, "1")
}
