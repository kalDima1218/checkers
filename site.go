package main

import (
	"fmt"
	_ "github.com/golang-jwt/jwt"
	"html/template"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"time"
)

var URL = "127.0.0.1"
var PORT = "8080"

var BOT_PLAYER = newPlayer("BOT", "BOT", "")

var games = make(map[string]Game)

var waiting_game = newTreap()
var waiting_for = make(map[string]string)
var last_seen = make(map[string]int)

var players = make(map[string]Player)

var logins = make(map[string]bool)
var names = make(map[string]bool)

func checkUser(r *http.Request) bool {
	login, errLogin := r.Cookie("login")
	password, errPassword := r.Cookie("password")
	if errLogin != nil || errPassword != nil {
		return false
	}
	_, okLogin := logins[login.Value]
	if !okLogin || players[login.Value].Password != password.Value {
		return false
	} else {
		return true
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) {
		page, _ := template.ParseFiles(path.Join("html", "index.html"))
		page.Execute(w, "")
	} else {
		resetCookie(w)
		redirectTo(w, r, "login")
	}
}

func game(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) {
		var id = r.URL.Query().Get("id")
		_, ok := games[id]
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

func startGame(w http.ResponseWriter, r *http.Request) {
	if !checkUser(r) {
		redirectToIndex(w, r)
		return
	}
	var login = getCookie(r, "login")
	_, ok := last_seen[login]
	if ok {
		waiting_game.erase(newItemWaitingGame(login, last_seen[login]))
		delete(last_seen, login)
	}
	if !waiting_game.empty() {
		var partner = waiting_game.begin().i.getFieldString("player")
		var id = strconv.Itoa(rand.Int())
		games[id] = newGame(players[login], players[partner])
		waiting_for[partner] = id
		waiting_for[login] = id
		redirectTo(w, r, "waiting_game")
	} else {
		last_seen[login] = int(time.Now().Unix())
		waiting_game.insert(newItemWaitingGame(login, last_seen[login]))
		redirectTo(w, r, "waiting_game")
	}
}

func getWaiting(w http.ResponseWriter, r *http.Request) {
	if !checkUser(r) {
		return
	}
	var login = getCookie(r, "login")
	id, ok := waiting_for[login]
	if !ok {
		fmt.Fprintf(w, "wrong")
	} else {
		delete(waiting_for, login)
		fmt.Fprintf(w, id)
	}
}

func startBotGame(w http.ResponseWriter, r *http.Request) {
	if !checkUser(r) {
		redirectToIndex(w, r)
		return
	}
	var login = getCookie(r, "login")
	var id = strconv.Itoa(rand.Int())
	games[id] = newGame(players[login], BOT_PLAYER)
	redirectTo(w, r, "game?id="+id)
}

func waitingGame(w http.ResponseWriter, r *http.Request) {
	page, _ := template.ParseFiles(path.Join("html", "waiting_game.html"))
	page.Execute(w, "")
}

func setupRoutes() {
	// PAGES SECTION
	http.HandleFunc("/", index)
	http.HandleFunc("/game", game)
	http.HandleFunc("/reg", reg)
	http.HandleFunc("/login", login)
	// FUNCTION SECTION
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/make_move", makeMove)
	http.HandleFunc("/end_move", endMove)
	http.HandleFunc("/whose_move", whoseMove)
	http.HandleFunc("/get_side", getSide)
	http.HandleFunc("/get_players", getPlayers)
	http.HandleFunc("/who_win", whoWin)
	http.HandleFunc("/start_game", startGame)
	http.HandleFunc("/get_waiting", getWaiting)
	http.HandleFunc("/start_bot_game", startBotGame)
	http.HandleFunc("/waiting_game", waitingGame)
	http.HandleFunc("/get_last_move_number", getLastMoveNumber)
	http.HandleFunc("/get_board_hist", getBoardHist)
	// FILE SECTION
	http.HandleFunc("/game.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("js", "game.js"))
	})
	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("js", "index.js"))
	})
	http.HandleFunc("/reg.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("js", "reg.js"))
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
	http.HandleFunc("/login_reg.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("css", "login_reg.css"))
	})
}

func startSite() {
	_init_load()
	go _autosave()
	setupRoutes()
	http.ListenAndServe(":"+PORT, nil)
}
