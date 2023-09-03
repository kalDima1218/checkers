package main

import (
	"encoding/json"
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

var BOT_PLAYER = newPlayer("BOT", "BOT", "")

var games = make(map[string]Game)

var wating_game = newTreap()
var waiting_for = make(map[string]string)
var last_seen = make(map[string]int)

var players = make(map[string]Player)

var logins = make(map[string]bool)
var names = make(map[string]bool)

func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://"+URL+":"+PORT+"/", http.StatusSeeOther)
}

func redirectTo(w http.ResponseWriter, r *http.Request, page string) {
	http.Redirect(w, r, "http://"+URL+":"+PORT+"/"+page, http.StatusSeeOther)
}

func resetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "login", Value: "", MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "password", Value: "", MaxAge: -1})
}

func getCookie(r *http.Request, data_key string) string {
	data, _ := r.Cookie(data_key)
	return data.Value
}

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

func reg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var name = r.URL.Query().Get("name")
	var login = r.URL.Query().Get("login")
	var password = r.URL.Query().Get("password")
	_, okName := names[name]
	_, okLogin := logins[login]
	if name == "" || login == "" || password == "" || okName || okLogin {
		fmt.Fprintf(w, "wrong")
		return
	}
	players[login] = newPlayer(name, login, password)
	names[name] = true
	logins[login] = true
	http.SetCookie(w, &http.Cookie{Name: "login", Value: login})
	http.SetCookie(w, &http.Cookie{Name: "password", Value: password})
	fmt.Fprintf(w, "ok")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var login = r.URL.Query().Get("login")
	var password = r.URL.Query().Get("password")
	_, ok := logins[login]
	if !ok || players[login].Password != password {
		fmt.Fprintf(w, "wrong")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "login", Value: login})
	http.SetCookie(w, &http.Cookie{Name: "password", Value: password})
	fmt.Fprint(w, "ok")
}

func logout(w http.ResponseWriter, r *http.Request) {
	resetCookie(w)
	redirectToIndex(w, r)
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
		wating_game.erase(newItemWatingGame(login, last_seen[login]))
		delete(last_seen, login)
	}
	if !wating_game.empty() {
		var partner = wating_game.begin().i.getFieldString("player")
		var id = strconv.Itoa(rand.Int())
		games[id] = newGame(players[login], players[partner])
		waiting_for[partner] = id
		waiting_for[login] = id
		redirectTo(w, r, "waiting_game")
	} else {
		last_seen[login] = int(time.Now().Unix())
		wating_game.insert(newItemWatingGame(login, last_seen[login]))
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

func getBoard(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := games[id]
	if !ok {
		return
	}
	fmt.Fprintf(w, game.getJsonBoard())
}

func whoseMove(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	_, ok := games[id]
	if !ok {
		return
	}
	fmt.Fprintf(w, strconv.Itoa(games[id].Turn))
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := games[id]
	if !ok {
		return
	}
	players_json, _ := json.Marshal([2]string{players[game.Players[0]].Name, players[game.Players[1]].Name})
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
	if !ok || game.Players[game.Turn] != login {
		return
	}
	game.makeMove([2]int{from_x, from_y}, [2]int{to_x, to_y})
	games[id] = game
}

func endMove(w http.ResponseWriter, r *http.Request) {
	if !checkUser(r) {
		redirectToIndex(w, r)
		return
	}
	var id = r.URL.Query().Get("id")
	var login = getCookie(r, "login")
	game, ok := games[id]
	if !ok || game.Players[game.Turn] != login {
		return
	}
	game.endMove()
	if game.Players[game.Turn] == "BOT" {
		game = BOT.findBestMove(game, game.Turn, (game.Turn+1)%2)
	}
	games[id] = game
}

func setupRoutes() {
	// PAGES SECTION
	http.HandleFunc("/", index)
	http.HandleFunc("/game", game)
	http.HandleFunc("/reg", func(w http.ResponseWriter, r *http.Request) {
		page, _ := template.ParseFiles(path.Join("html", "reg.html"))
		page.Execute(w, "")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		page, _ := template.ParseFiles(path.Join("html", "login.html"))
		page.Execute(w, "")
	})
	// FUNCTION SECTION
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/_reg", reg)
	http.HandleFunc("/_login", login)
	http.HandleFunc("/get_board", getBoard)
	http.HandleFunc("/make_move", makeMove)
	http.HandleFunc("/end_move", endMove)
	http.HandleFunc("/whose_move", whoseMove)
	http.HandleFunc("/get_players", getPlayers)
	http.HandleFunc("/who_win", whoWin)
	http.HandleFunc("/start_game", startGame)
	http.HandleFunc("/get_waiting", getWaiting)
	http.HandleFunc("/start_bot_game", startBotGame)
	http.HandleFunc("/waiting_game", waitingGame)
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
