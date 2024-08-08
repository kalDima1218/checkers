package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func handleGetBoardHist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	turn, isTurnConversed := strconv.Atoi(r.URL.Query().Get("turn"))
	if isTurnConversed != nil {
		return
	}
	game, ok := getGame(id)
	if !ok || turn < 0 || turn >= len(game.Turns) {
		return
	}
	boardJson, _ := json.Marshal(game.Turns[turn])
	fmt.Fprintf(w, string(boardJson))
}

func handleGetLastMoveNumber(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := getGame(id)
	if !ok {
		return
	}
	fmt.Fprintf(w, strconv.Itoa(len(game.Turns)-1))
}

func handleWhoseMove(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := getGame(id)
	if !ok {
		return
	}
	fmt.Fprintf(w, strconv.Itoa(game.Board.Whose_turn))
}

func handleGetSide(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	var login = getCookie(r, "login")
	game, ok := getGame(id)
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

func handleGetPlayersUsernames(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := getGame(id)
	if !ok {
		return
	}
	playersUsernamesJson, _ := json.Marshal([2]string{getUsername(game.Players[0]), getUsername(game.Players[1])})
	fmt.Fprintf(w, string(playersUsernamesJson))
}

func handleWhoWin(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := getGame(id)
	if !ok {
		return
	}
	fmt.Fprintf(w, strconv.Itoa(game.whoWin()))
}

func handleMakeMove(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		redirectToIndex(w, r)
		return
	}
	var id = r.URL.Query().Get("id")
	var login = getCookie(r, "login")
	from_x, _ := strconv.Atoi(r.URL.Query().Get("from_x"))
	from_y, _ := strconv.Atoi(r.URL.Query().Get("from_y"))
	to_x, _ := strconv.Atoi(r.URL.Query().Get("to_x"))
	to_y, _ := strconv.Atoi(r.URL.Query().Get("to_y"))
	game, ok := getGame(id)
	if !ok || game.Players[game.Board.Whose_turn] != login {
		return
	}
	if game.makeMove([2]int{from_x, from_y}, [2]int{to_x, to_y}) {
		setGame(id, game)
		fmt.Fprintf(w, "1")
		fmt.Println(from_x, from_y, to_x, to_y)
	} else {
		fmt.Fprintf(w, "0")
	}
}

func handleEndMove(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		redirectToIndex(w, r)
		return
	}
	var id = r.URL.Query().Get("id")
	var login = getCookie(r, "login")
	game, ok := getGame(id)
	if !ok || game.Players[game.Board.Whose_turn] != login || game.Board.Last_piece == [2]int{-1, -1} {
		fmt.Fprintf(w, "0")
		return
	}
	game.endMove()
	if game.Players[game.Board.Whose_turn] == "BOT" {
		BOT.makeMove(game)
		//game = BOT.findBestMove(game, game.Board.Whose_turn, (game.Board.Whose_turn+1)%2)
	}

	setGame(id, game)
	fmt.Fprintf(w, "1")
}
