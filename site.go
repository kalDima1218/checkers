package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"time"
)

var URL = "127.0.0.1"
var PORT = "8080"

var waiting_game = newSet()
var waiting_for = make(map[string]string)
var last_seen = make(map[string]int)

// TODO добавить обработку ошибок в получение куки
// TODO протестирвать матчмейкинг

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if checkSession(r) {
		page, _ := template.ParseFiles(path.Join("html", "index.html"))
		page.Execute(w, "")
	} else {
		resetCookie(w)
		redirectTo(w, r, "login")
	}
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	if checkSession(r) {
		id := r.URL.Query().Get("id")
		ok := isGameExists(id)
		if !ok {
			redirectToIndex(w, r)
		} else {
			page, _ := template.ParseFiles(path.Join("html", "game.html"))
			page.Execute(w, "")
		}
	} else {
		redirectToIndex(w, r)
	}
}

func handleStartGame(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		redirectToIndex(w, r)
		return
	}
	login, err := getLogin(r)
	if err != nil {
		return
	}
	_, ok := last_seen[login]
	if ok {
		waiting_game.erase(newItemWaitingGame(login, last_seen[login]))
		delete(last_seen, login)
	}
	if !waiting_game.empty() {
		partner := waiting_game.begin().i.getFieldString("player")
		waiting_game.erase(waiting_game.begin().i)
		id := strconv.Itoa(rand.Int())
		game := newGame(getUsername(login), getUsername(partner))
		insertGame(id, &game)
		waiting_for[partner] = id
		waiting_for[login] = id
		redirectTo(w, r, "waiting_game")
	} else {
		last_seen[login] = int(time.Now().Unix())
		waiting_game.insert(newItemWaitingGame(login, last_seen[login]))
		redirectTo(w, r, "waiting_game")
	}
}

func handleGetWaiting(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		return
	}
	login, err := getLogin(r)
	if err != nil {
		return
	}
	id, ok := waiting_for[login]
	if !ok {
		fmt.Fprintf(w, "wrong")
	} else {
		delete(waiting_for, login)
		fmt.Fprintf(w, id)
	}
}

func handleStartBotGame(w http.ResponseWriter, r *http.Request) {
	if !checkSession(r) {
		redirectToIndex(w, r)
		return
	}
	login, err := getLogin(r)
	if err != nil {
		return
	}
	id := strconv.Itoa(rand.Int())
	game := newGame(getUsername(login), "BOT")
	err = insertGame(id, &game)
	if err != nil {
		return
	}
	redirectTo(w, r, "game?id="+id)
}

func handleWaitingGame(w http.ResponseWriter, r *http.Request) {
	page, _ := template.ParseFiles(path.Join("html", "waiting_game.html"))
	page.Execute(w, "")
}

func startSite() {
	loadDB()
	// PAGES
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/game", handleGame)
	http.HandleFunc("/registration", handleRegistration)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/waiting_game", handleWaitingGame)

	// AUTH
	http.HandleFunc("/logout", handleLogout)

	// GAME API
	http.HandleFunc("/make_move", handleMakeMove)
	http.HandleFunc("/end_move", handleEndMove)
	http.HandleFunc("/whose_move", handleWhoseMove)
	http.HandleFunc("/get_side", handleGetSide)
	http.HandleFunc("/get_players", handleGetPlayersUsernames)
	http.HandleFunc("/who_win", handleWhoWin)
	http.HandleFunc("/get_last_move_number", handleGetLastMoveNumber)
	http.HandleFunc("/get_board_hist", handleGetBoardHist)

	// FUNCTIONS
	http.HandleFunc("/start_bot_game", handleStartBotGame)
	http.HandleFunc("/start_game", handleStartGame)
	http.HandleFunc("/get_waiting", handleGetWaiting)

	// FILES
	http.HandleFunc("/game.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("js", "game.js"))
	})
	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("js", "index.js"))
	})
	http.HandleFunc("/registration.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("js", "registration.js"))
	})
	http.HandleFunc("/login.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("js", "login.js"))
	})
	http.HandleFunc("/waiting_game.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("js", "waiting_game.js"))
	})
	http.HandleFunc("/main.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("css", "main.css"))
	})
	http.HandleFunc("/game.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("css", "game.css"))
	})
	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("css", "index.css"))
	})
	http.HandleFunc("/login_registration.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("css", "login_registration.css"))
	})

	http.ListenAndServe(":"+PORT, nil)
}
