package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// TODO добавить обработку ошибок в получение куки

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
	fmt.Fprintf(w, strconv.Itoa(game.Board.WhoseTurn))
}

func handleGetSide(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	login, err := getLogin(r)
	if err != nil {
		return
	}
	game, ok := getGame(id)
	if !ok {
		return
	}
	if game.Players[0] == login {
		fmt.Fprintf(w, "0")
	} else {
		fmt.Fprintf(w, "1")
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
	winner := strconv.Itoa(game.whoWin())
	if winner == "-1" {
		if len(game.Turns) < 50 {
			fmt.Fprintf(w, winner)
		} else {
			fmt.Fprintf(w, "-2")
		}
	} else {
		fmt.Fprintf(w, winner)
	}
}

func handleMakeMove(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		redirectToIndex(w, r)
		return
	}
	id := r.URL.Query().Get("id")
	login, err := getLogin(r)
	if err != nil {
		return
	}
	fromX, _ := strconv.Atoi(r.URL.Query().Get("from_x"))
	fromY, _ := strconv.Atoi(r.URL.Query().Get("from_y"))
	toX, _ := strconv.Atoi(r.URL.Query().Get("to_x"))
	toY, _ := strconv.Atoi(r.URL.Query().Get("to_y"))
	game, ok := getGame(id)
	if !ok || game.Players[game.Board.WhoseTurn] != login || len(game.Turns) >= 50 {
		return
	}
	if game.makeMove([2]int{fromX, fromY}, [2]int{toX, toY}) {
		setGame(id, game)
		fmt.Fprintf(w, "1")
	} else {
		fmt.Fprintf(w, "0")
	}
}

func handleEndMove(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		redirectToIndex(w, r)
		return
	}
	id := r.URL.Query().Get("id")
	login, err := getLogin(r)
	if err != nil {
		return
	}
	game, ok := getGame(id)
	if !ok || game.Players[game.Board.WhoseTurn] != login || game.Board.LastPiece == [2]int{-1, -1} {
		//fmt.Fprintf(w, "0")
		return
	}
	game.endMove()
	if game.Players[game.Board.WhoseTurn] == "BOT" {
		BOT.makeMove(game)
		//game = BOT.findBestMove(game, game.Board.Whose_turn, (game.Board.Whose_turn+1)%2)
	}

	setGame(id, game)
	//fmt.Fprintf(w, "1")
}
