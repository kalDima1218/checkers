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

var waiting_game = newTreap()
var waiting_for = make(map[string]string)
var last_seen = make(map[string]int)

var players = make(map[string]Player)

var logins = make(map[string]bool)
var names = make(map[string]bool)

// redirectToIndex redirects the request to the index page.
//
// It takes two parameters:
// - w: an http.ResponseWriter object used to write the response.
// - r: an *http.Request object representing the incoming request.
// It does not return any value.
func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://"+URL+":"+PORT+"/", http.StatusSeeOther)
}

// redirectTo redirects the user to the specified page.
//
// It takes in the http.ResponseWriter and *http.Request as parameters.
// It does not return anything.
func redirectTo(w http.ResponseWriter, r *http.Request, page string) {
	http.Redirect(w, r, "http://"+URL+":"+PORT+"/"+page, http.StatusSeeOther)
}

// resetCookie resets the login and password cookies.
//
// It takes a http.ResponseWriter as a parameter.
// It does not return anything.
func resetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "login", Value: "", MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: "password", Value: "", MaxAge: -1})
}

// getCookie returns the value of a cookie with the specified data key from the provided http.Request object.
//
// Parameters:
// - r: The http.Request object from which to retrieve the cookie.
// - data_key: The key of the cookie data to retrieve.
//
// Return:
// - string: The value of the cookie data.
func getCookie(r *http.Request, data_key string) string {
	data, _ := r.Cookie(data_key)
	return data.Value
}

// checkUser checks if the user is authenticated based on the provided request.
//
// The function takes a pointer to a http.Request as a parameter.
// It returns a boolean value indicating whether the user is authenticated or not.
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

// index is a Go function that handles the index route.
//
// It takes in two parameters:
//   - w: an http.ResponseWriter object used to write the response.
//   - r: a pointer to an http.Request object representing the incoming request.
//
// It does not return any value.
func index(w http.ResponseWriter, r *http.Request) {
	if checkUser(r) {
		page, _ := template.ParseFiles(path.Join("html", "index.html"))
		page.Execute(w, "")
	} else {
		resetCookie(w)
		redirectTo(w, r, "login")
	}
}

// reg handles the registration of a new player.
//
// It takes in two parameters: w, an http.ResponseWriter used to write the response,
// and r, an http.Request representing the incoming request.
//
// The function checks if the request method is not http.MethodPost and returns early if true.
// It then extracts the values of the "name", "login", and "password" query parameters from the request URL.
//
// Next, it checks if the name and login already exist in the names and logins maps respectively.
// If any of the conditions name == "", login == "", password == "", okName, or okLogin are true,
// the function writes "wrong" to the response and returns.
//
// If all the conditions are false, it creates a new player using the newPlayer function,
// adds the player's name, login, and password to the names and logins maps respectively,
// and sets the "login" and "password" cookies in the response.
// Finally, it writes "ok" to the response.
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

// login is a Go function that handles the login functionality.
//
// It takes in two parameters:
// - w: an http.ResponseWriter object that is used to write the response to the client.
// - r: an *http.Request object that represents the client request.
//
// The function does the following:
// - Checks if the request method is not POST and returns if true.
// - Retrieves the values of the "login" and "password" query parameters from the request URL.
// - Checks if the login exists in the "logins" map and if the password matches the stored password.
// - Writes "wrong" to the response writer if the login and password are invalid.
// - Sets two cookies ("login" and "password") with the respective values.
// - Writes "ok" to the response writer.
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

// logout logs out the user by resetting the cookie and redirecting them to the index page.
//
// Parameters:
// - w: the http.ResponseWriter used to write the response.
// - r: the *http.Request representing the incoming request.
//
// Returns: None.
func logout(w http.ResponseWriter, r *http.Request) {
	resetCookie(w)
	redirectToIndex(w, r)
}

// game handles the game logic for the HTTP server.
//
// It takes two parameters, w http.ResponseWriter and r *http.Request.
// It does not return anything.
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

// startGame is a function that handles the logic for starting a game.
//
// It takes in two parameters:
// - w: an http.ResponseWriter object
// - r: an http.Request object
//
// The function does the following:
// - Checks if the user is authenticated. If not, it redirects to the index page.
// - Retrieves the login cookie.
// - Checks if the user's last seen time exists in the `last_seen` map. If it does, it removes the user from the waiting game.
// - If there are players waiting in the waiting game, it retrieves the partner and generates a new game ID.
// - Creates a new game object and adds it to the `games` map.
// - Adds the partner and user IDs to the `waiting_for` map.
// - Redirects the user to the "waiting_game" page.
// - If there are no players waiting in the waiting game, it updates the user's last seen time and adds the user to the waiting game.
// - Redirects the user to the "waiting_game" page.
func startGame(w http.ResponseWriter, r *http.Request) {
	if !checkUser(r) {
		redirectToIndex(w, r)
		return
	}
	var login = getCookie(r, "login")
	_, ok := last_seen[login]
	if ok {
		waiting_game.erase(newItemWatingGame(login, last_seen[login]))
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
		waiting_game.insert(newItemWatingGame(login, last_seen[login]))
		redirectTo(w, r, "waiting_game")
	}
}

// getWaiting is a Go function that handles the "GET /waiting" endpoint.
//
// It takes in a http.ResponseWriter and a *http.Request as parameters.
// The function checks if the user is authenticated by calling the checkUser function.
// If the user is not authenticated, the function returns.
//
// The function then retrieves the login information from the cookie by calling the getCookie function.
// It checks if the login information exists in the waiting_for map.
// If the login information does not exist, it writes "wrong" to the http.ResponseWriter.
// Otherwise, it deletes the login information from the waiting_for map and writes the id to the http.ResponseWriter.
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

// startBotGame is a Go function that handles the start of a bot game.
//
// It takes in two parameters: w, which is an http.ResponseWriter object used to write the response, and r, which is an http.Request object representing the incoming request.
// Both parameters are required and cannot be nil.
//
// This function does the following:
// - Checks if the user is authenticated by calling checkUser(r). If the user is not authenticated, it redirects to the index page.
// - Retrieves the login value from the cookie by calling getCookie(r, "login").
// - Generates a random id using the strconv.Itoa and rand.Int functions.
// - Creates a new game using the newGame function, passing in the player associated with the login and the BOT_PLAYER constant.
// - Redirects the user to the game page with the generated id.
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

// waitingGame handles the HTTP request and response for the "waitingGame" function.
//
// It takes in a http.ResponseWriter and a *http.Request as parameters.
// It does not return any values.
func waitingGame(w http.ResponseWriter, r *http.Request) {
	page, _ := template.ParseFiles(path.Join("html", "waiting_game.html"))
	page.Execute(w, "")
}

// getBoardHist is a Go function that retrieves the board history for a game.
//
// It takes in two parameters:
//   - w: an http.ResponseWriter object used to write the response.
//   - r: a *http.Request object representing the incoming request.
//
// This function does not return any values.
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

// whoseMove returns the current player's turn for a given game ID.
//
// Parameters:
//   - w: The http.ResponseWriter used to write the response.
//   - r: The http.Request containing the game ID.
//
// Return:
//   - None.
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

// getPlayers retrieves the players' names and sends them as a JSON response.
//
// It takes a ResponseWriter and a Request as parameters.
// Returns nothing.
func getPlayers(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := games[id]
	if !ok {
		return
	}
	players_json, _ := json.Marshal([2]string{players[game.Players[0]].Name, players[game.Players[1]].Name})
	fmt.Fprintf(w, string(players_json))
}

// whoWin determines the winner of a game based on the provided ID.
//
// Parameters:
// - w: the http.ResponseWriter used to send the game result.
// - r: the *http.Request containing the ID of the game.
//
// Return:
// None.
func whoWin(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	game, ok := games[id]
	if !ok {
		return
	}
	fmt.Fprintf(w, strconv.Itoa(game.whoWin()))
}

// makeMove handles the logic for making a move in the game.
//
// It takes in the http.ResponseWriter and http.Request as parameters.
// It does not return anything.
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

// endMove handles the end of a player's move in the game.
//
// It takes in an http.ResponseWriter and an *http.Request as parameters.
// It does not return any values.
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

// setupRoutes sets up the routes for the HTTP server.
//
// No parameters.
// No return type.
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

// startSite initializes the site and starts the server.
//
// It calls _init_load to initialize the site.
// It spawns a goroutine to run _autosave in the background.
// It calls setupRoutes to set up the routes for the server.
// It listens for incoming requests on the specified port.
// The function does not take any parameters.
// It does not return any values.
func startSite() {
	_init_load()
	go _autosave()
	setupRoutes()
	http.ListenAndServe(":"+PORT, nil)
}
